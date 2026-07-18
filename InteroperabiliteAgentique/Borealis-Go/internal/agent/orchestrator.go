package agent

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/creditdomain"
	"github.com/agbruneau/borealis/internal/hitl"
	"github.com/agbruneau/borealis/internal/pep"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Seuils AML cumulatifs (FINTRAC simulé, §11.4, FR-16) : au-delà de 0,7, la
// demande est refusée ; au-delà de 0,8, elle est refusée **et** escaladée en
// HITL avant scellement de la décision. Valeurs à calibrer (hypothèses).
const (
	amlDeclineThreshold  = 0.7
	amlEscalateThreshold = 0.8
)

// Result agrège la sortie de la pré-qualification. Jamais un octroi ferme
// (FR-25, invariant 1) : Decision.Outcome ∈ {declined, refer, pre_qualified}.
//
// Decision est la sortie métier **autoritative**. Recommendation est l'avis
// **consultatif** du reasoner (fake en M0) : aucune autorité décisionnelle, et
// avec un reasoner réel elle pourrait diverger de Decision. Un consommateur
// doit agir sur Decision.Outcome, pas sur Recommendation.
type Result struct {
	Applicant      string
	Identity       IdentityCheck
	Decision       creditdomain.Decision
	Recommendation Recommendation
}

// Orchestrator exécute une boucle L0 minimale : identité (appel MCP) → règles
// métier → raisonnement (fake déterministe) → journal d'audit à chaque étape.
// Aucun chemin d'octroi.
type Orchestrator struct {
	caller   *ResilientCaller
	reasoner Reasoner
	journal  *audit.Journal
	policy   creditdomain.Policy
	hitl     *hitl.Broker // optionnel : escalade AML élevée (nil → refus sans escalade)
	kya      string
	ownerKYH string
}

// WithHITL branche un broker HITL : un score AML au-dessus du seuil d'escalade
// (> 0,8) soumet une proposition **avant** le scellement de la décision (FR-16).
// Sans broker, le dépassement reste un refus motivé, sans escalade.
func (o *Orchestrator) WithHITL(b *hitl.Broker) *Orchestrator {
	o.hitl = b
	return o
}

// WithPEP branche un Policy Enforcement Point sur le caller : toute la boucle
// passe par la garde non contournable (invariant 3). Si le PEP est indisponible,
// la retombée fail-closed L1→L0 fait échouer le premier appel d'outil — aucune
// action autonome (scénario M5.1 fail-closed).
func (o *Orchestrator) WithPEP(p *pep.PEP, kya string) *Orchestrator {
	o.caller.WithPEP(p, kya)
	return o
}

// NewOrchestrator construit l'orchestrateur avec ses dépendances injectées. La
// session MCP est enveloppée d'un caller résilient (FR-17, NFR-04).
func NewOrchestrator(identity *mcp.ClientSession, reasoner Reasoner, journal *audit.Journal, policy creditdomain.Policy) *Orchestrator {
	return &Orchestrator{
		caller:   NewResilientCaller(identity),
		reasoner: reasoner,
		journal:  journal,
		policy:   policy,
		kya:      "agent:orchestrator",
		ownerKYH: "human:credit-lead@borealis.example",
	}
}

// scenarioFor associe une issue métier au scénario du reasoner (fake).
func scenarioFor(o creditdomain.Outcome) string {
	switch o {
	case creditdomain.OutcomePreQualified:
		return "approved"
	case creditdomain.OutcomeDeclined:
		return "denied"
	default:
		return "escalate"
	}
}

func (o *Orchestrator) audit(kyc, action, result string) {
	o.journal.Append(audit.Record{
		KYA: o.kya, OwnerKYH: o.ownerKYH, KYC: kyc,
		Action: action, Result: result, Version: "m0",
	})
}

// fail journalise l'échec d'une étape puis une entrée « decision » terminale
// (outcome=error), garantissant qu'une exécution en erreur laisse une piste
// d'audit **close** — jamais tronquée ni indétectable.
func (o *Orchestrator) fail(kyc, action string, err error) (Result, error) {
	o.audit(kyc, action, "error: "+err.Error())
	o.audit(kyc, "decision", "error")
	return Result{}, err
}

// subjectRef pseudonymise le sujet au point d'écriture du journal : le champ
// KYC ne doit jamais porter de PII (invariant appliqué ici, pas seulement
// documenté ; données 100 % synthétiques).
func subjectRef(applicantID string) string {
	sum := sha256.Sum256([]byte(applicantID))
	return "cust:" + hex.EncodeToString(sum[:8])
}

// amlGate applique les seuils AML cumulatifs. Il retourne done=true (avec une
// décision de refus scellée) quand le score dépasse le seuil de refus **ou**
// qu'un hit de liste de surveillance est signalé (m12 : l'invariant « hit ⇒
// refus » doit être porté ici, pas seulement supposé du couplage hit→score du
// serveur d'identité simulé) ; au-delà du seuil d'escalade, il soumet d'abord
// une proposition HITL — l'entrée d'audit « hitl.propose » précède ainsi le
// scellement « decision », prouvant que l'escalade a lieu **avant** que la
// décision soit close (FR-16).
func (o *Orchestrator) amlGate(kyc, applicantID string, scr Screening) (creditdomain.Decision, bool) {
	if !scr.WatchlistHit && scr.AMLScore <= amlDeclineThreshold {
		return creditdomain.Decision{}, false
	}
	// Motif fidèle au déclencheur : quand le refus vient d'un hit avec un score
	// SOUS le seuil (serveur découplant hit et score), « score au-dessus du
	// seuil » serait factuellement faux — piste d'audit et explication aval
	// trompeuses. Le cas couplé (score > seuil) garde son motif historique.
	reason := "aml_score_above_threshold"
	if scr.WatchlistHit && scr.AMLScore <= amlDeclineThreshold {
		reason = "watchlist_hit"
	}
	dec := creditdomain.Decision{
		ApplicantID: applicantID,
		Outcome:     creditdomain.OutcomeDeclined,
		Reasons:     []string{reason},
	}
	if scr.AMLScore > amlEscalateThreshold && o.hitl != nil {
		// ponytail: ID de proposition dérivé du sujet — unique par demandeur pour
		// le démonstrateur (chaque Prequalify est mono-passe). Plafond connu : si un
		// même dossier était rejoué APRÈS release, la 2e proposition réutiliserait
		// cet ID et serait avalée (released[id] reste vrai → PendingProposals ne la
		// voit pas, Release renvoie ErrAlreadyReleased). Correctif à ce moment :
		// nonce de requête dans l'ID.
		o.hitl.Propose(hitl.Proposal{
			ID:      "aml-" + kyc,
			Subject: kyc,
			Summary: "AML élevé (score > 0.8) — escalade avant scellement",
		})
	}
	o.audit(kyc, "decision", dec.Outcome.String())
	return dec, true
}

// Prequalify pré-qualifie une demande. Ne produit jamais d'octroi ferme (FR-25).
func (o *Orchestrator) Prequalify(ctx context.Context, a creditdomain.Applicant) (Result, error) {
	kyc := subjectRef(a.ID)

	id, err := VerifyIdentity(ctx, o.caller, a.ID)
	if err != nil {
		return o.fail(kyc, "verify_identity", err)
	}
	o.audit(kyc, "verify_identity", id.Reason)

	if !id.Match {
		dec := creditdomain.Decision{
			ApplicantID: a.ID,
			Outcome:     creditdomain.OutcomeDeclined,
			Reasons:     []string{"identity_not_verified"},
		}
		o.audit(kyc, "decision", dec.Outcome.String())
		return Result{Applicant: a.ID, Identity: id, Decision: dec, Recommendation: RecommendDecline}, nil
	}

	// Criblage AML : arrêt de conformité qui court-circuite l'évaluation métier.
	scr, err := ScreenWatchlists(ctx, o.caller, a.ID)
	if err != nil {
		return o.fail(kyc, "screen_watchlists", err)
	}
	o.audit(kyc, "screen_watchlists", scr.Reason)
	if dec, done := o.amlGate(kyc, a.ID, scr); done {
		return Result{Applicant: a.ID, Identity: id, Decision: dec, Recommendation: RecommendDecline}, nil
	}

	dec, err := creditdomain.Assess(a, o.policy)
	if err != nil {
		return o.fail(kyc, "assess", err)
	}
	o.audit(kyc, "assess", dec.Outcome.String())

	rr, err := o.reasoner.Reason(ctx, ReasonRequest{
		ApplicantID: a.ID,
		Signals:     map[string]string{"scenario": scenarioFor(dec.Outcome)},
	})
	if err != nil {
		return o.fail(kyc, "reason", err)
	}
	o.audit(kyc, "reason", string(rr.Recommendation))
	o.audit(kyc, "decision", dec.Outcome.String())

	return Result{Applicant: a.ID, Identity: id, Decision: dec, Recommendation: rr.Recommendation}, nil
}
