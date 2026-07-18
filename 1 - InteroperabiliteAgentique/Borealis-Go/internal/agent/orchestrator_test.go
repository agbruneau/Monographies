package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/creditdomain"
	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/agbruneau/borealis/internal/hitl"
	"github.com/agbruneau/borealis/internal/mcpserver"
	"github.com/agbruneau/borealis/internal/pep"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// stubReasoner est un double local (les tests internes du paquet agent ne
// peuvent pas importer testfakes, qui importe agent — cycle). Il mappe le
// scénario à une recommandation fixe.
type stubReasoner struct{}

func (stubReasoner) Reason(_ context.Context, req ReasonRequest) (ReasonResult, error) {
	switch req.Signals["scenario"] {
	case "approved":
		return ReasonResult{Recommendation: RecommendPreQualify}, nil
	case "denied":
		return ReasonResult{Recommendation: RecommendDecline}, nil
	default:
		return ReasonResult{Recommendation: RecommendEscalate}, nil
	}
}

// connectInMemory connecte un client à un serveur MCP via transport mémoire
// (T3). Helper LOCAL au paquet, pas dans internal/testfakes : testfakes
// importe agent (stubReasoner y vit conceptuellement), donc un test interne
// de agent important testfakes créerait un cycle (vérifié : "import cycle
// not allowed in test").
func connectInMemory(t testing.TB, s *mcp.Server) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "v0"}, nil)
	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := s.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("connexion serveur : %v", err)
	}
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("connexion client : %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func newIdentitySession(t testing.TB) *mcp.ClientSession {
	return connectInMemory(t, mcpserver.NewIdentityServer(fixtures.Generate(fixtures.DefaultSeed, 100)))
}

func testPolicy() creditdomain.Policy {
	return creditdomain.Policy{MinScore: 600, MaxDTI: 0.43, ReferScore: 660, ReferDTI: 0.36}
}

func fixedClock() func() time.Time {
	ts := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	return func() time.Time { return ts }
}

func goodApplicant() creditdomain.Applicant {
	return creditdomain.Applicant{
		ID: "A00001", MonthlyIncome: 1000000, MonthlyDebt: 50000,
		RequestedAmount: 1000000, AnnualRate: 0.079, TermMonths: 60, CreditScore: 720,
	}
}

func TestPrequalifyVerifiedPreQualified(t *testing.T) {
	sess := newIdentitySession(t)
	j := audit.New(audit.WithClock(fixedClock()))
	o := NewOrchestrator(sess, stubReasoner{}, j, testPolicy())

	res, err := o.Prequalify(context.Background(), goodApplicant())
	if err != nil {
		t.Fatal(err)
	}
	if res.Decision.Outcome != creditdomain.OutcomePreQualified {
		t.Errorf("outcome = %s, want pre_qualified", res.Decision.Outcome)
	}
	if res.Recommendation != RecommendPreQualify {
		t.Errorf("reco = %s, want pre_qualify", res.Recommendation)
	}
	if !res.Identity.Match {
		t.Error("identité non vérifiée alors que l'ID existe dans le jeu")
	}
	// 5 entrées : verify_identity, screen_watchlists, assess, reason, decision.
	if j.Len() != 5 {
		t.Errorf("audit = %d entrées, want 5", j.Len())
	}
}

func TestPrequalifyUnverified(t *testing.T) {
	sess := newIdentitySession(t)
	j := audit.New(audit.WithClock(fixedClock()))
	o := NewOrchestrator(sess, stubReasoner{}, j, testPolicy())

	a := goodApplicant()
	a.ID = "unknown-9" // non reconnu → identité non vérifiée
	res, err := o.Prequalify(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if res.Decision.Outcome != creditdomain.OutcomeDeclined {
		t.Errorf("outcome = %s, want declined", res.Decision.Outcome)
	}
	if res.Recommendation != RecommendDecline {
		t.Errorf("reco = %s, want decline", res.Recommendation)
	}
	// 2 entrées : verify_identity, decision (pas d'assess ni de reason).
	if j.Len() != 2 {
		t.Errorf("audit = %d entrées, want 2", j.Len())
	}
}

func TestPrequalifyDeniedAndRefer(t *testing.T) {
	sess := newIdentitySession(t)
	tests := []struct {
		name     string
		mut      func(*creditdomain.Applicant)
		want     creditdomain.Outcome
		wantReco Recommendation
	}{
		{"denied", func(a *creditdomain.Applicant) { a.CreditScore = 500 }, creditdomain.OutcomeDeclined, RecommendDecline},
		{"refer", func(a *creditdomain.Applicant) { a.CreditScore = 640 }, creditdomain.OutcomeRefer, RecommendEscalate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := audit.New(audit.WithClock(fixedClock()))
			o := NewOrchestrator(sess, stubReasoner{}, j, testPolicy())
			a := goodApplicant()
			tt.mut(&a)
			res, err := o.Prequalify(context.Background(), a)
			if err != nil {
				t.Fatal(err)
			}
			if res.Decision.Outcome != tt.want {
				t.Errorf("outcome = %s, want %s", res.Decision.Outcome, tt.want)
			}
			// Verrouille le mapping issue → scénario → recommandation consultative.
			if res.Recommendation != tt.wantReco {
				t.Errorf("recommendation = %s, want %s", res.Recommendation, tt.wantReco)
			}
		})
	}
}

// TestAMLGateThresholds isole la logique de seuils (M3.6) : la bande médiane
// « refus sans escalade » (0,7 < score <= 0,8) n'est pas produite par un vrai
// dossier (scores dans [0,0.5) ou [0.85,1.0)) ; on l'exerce donc en synthétique.
func TestAMLGateThresholds(t *testing.T) {
	cases := []struct {
		name        string
		score       float64
		wantDone    bool
		wantPending int // propositions HITL attendues
	}{
		{"sous le seuil → passe", 0.50, false, 0},
		{"pile au seuil de refus → passe", 0.70, false, 0},
		{"refus sans escalade", 0.75, true, 0},
		{"pile au seuil d'escalade → refus sans escalade", 0.80, true, 0},
		{"refus + escalade", 0.85, true, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			j := audit.New(audit.WithClock(fixedClock()))
			broker := hitl.New(j)
			o := NewOrchestrator(newIdentitySession(t), stubReasoner{}, j, testPolicy()).WithHITL(broker)
			dec, done := o.amlGate("cust:x", "A00042", Screening{AMLScore: c.score})
			if done != c.wantDone {
				t.Fatalf("done = %v, want %v", done, c.wantDone)
			}
			if done && dec.Outcome != creditdomain.OutcomeDeclined {
				t.Errorf("outcome = %s, want declined", dec.Outcome)
			}
			if got := broker.PendingProposals(); got != c.wantPending {
				t.Errorf("propositions HITL = %d, want %d", got, c.wantPending)
			}
		})
	}
}

// TestAMLGateWatchlistHitOverridesScore : un hit de liste de surveillance doit
// refuser indépendamment du score AML (m12). Avant correctif, amlGate ne
// regardait que AMLScore ; l'invariant « hit ⇒ refus » ne tenait que parce que
// le serveur d'identité simulé couple hit→score élevé — un serveur substitué
// renvoyant hit=true avec un score bas passerait à l'évaluation métier.
func TestAMLGateWatchlistHitOverridesScore(t *testing.T) {
	j := audit.New(audit.WithClock(fixedClock()))
	o := NewOrchestrator(newIdentitySession(t), stubReasoner{}, j, testPolicy())
	dec, done := o.amlGate("cust:x", "A00042", Screening{WatchlistHit: true, AMLScore: 0.1})
	if !done {
		t.Fatal("watchlistHit=true devrait refuser indépendamment du score")
	}
	if dec.Outcome != creditdomain.OutcomeDeclined {
		t.Errorf("outcome = %s, want declined", dec.Outcome)
	}
	// Motif fidèle au déclencheur : « score au-dessus du seuil » serait ici
	// factuellement faux (score 0.1) — piste d'audit et explication trompeuses.
	if len(dec.Reasons) == 0 || dec.Reasons[0] != "watchlist_hit" {
		t.Errorf("reasons = %v, want [watchlist_hit] (hit avec score sous le seuil)", dec.Reasons)
	}
}

// TestAMLGateEscalateWithoutBroker : sans broker, un score > 0,8 refuse sans
// paniquer (l'escalade est optionnelle).
func TestAMLGateEscalateWithoutBroker(t *testing.T) {
	j := audit.New(audit.WithClock(fixedClock()))
	o := NewOrchestrator(newIdentitySession(t), stubReasoner{}, j, testPolicy())
	if _, done := o.amlGate("cust:x", "A00042", Screening{AMLScore: 0.9}); !done {
		t.Error("score 0.9 devrait refuser")
	}
}

// TestPrequalifyAMLEscalation : un dossier sur liste (A00007, Gita Diallo →
// score >= 0.85) est refusé ET escaladé en HITL, l'entrée « hitl.propose »
// précédant le scellement « decision » (escalade AVANT scellement, FR-16).
func TestPrequalifyAMLEscalation(t *testing.T) {
	sess := newIdentitySession(t)
	j := audit.New(audit.WithClock(fixedClock()))
	broker := hitl.New(j)
	o := NewOrchestrator(sess, stubReasoner{}, j, testPolicy()).WithHITL(broker)

	a := goodApplicant()
	a.ID = "A00007" // Gita Diallo → liste OFAC simulée
	res, err := o.Prequalify(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if res.Decision.Outcome != creditdomain.OutcomeDeclined {
		t.Errorf("outcome = %s, want declined (AML)", res.Decision.Outcome)
	}
	if len(res.Decision.Reasons) == 0 || res.Decision.Reasons[0] != "aml_score_above_threshold" {
		t.Errorf("reasons = %v, want [aml_score_above_threshold]", res.Decision.Reasons)
	}
	if broker.PendingProposals() != 1 {
		t.Errorf("propositions HITL = %d, want 1", broker.PendingProposals())
	}
	// L'escalade précède le scellement : index(hitl.propose) < index(decision).
	es := j.Entries()
	iPropose, iDecision := -1, -1
	for i, e := range es {
		switch e.Action {
		case "hitl.propose":
			iPropose = i
		case "decision":
			iDecision = i
		}
	}
	if iPropose < 0 || iDecision < 0 || iPropose >= iDecision {
		t.Errorf("ordre d'audit : propose=%d decision=%d, want propose < decision", iPropose, iDecision)
	}
	// La file de commande reste vide sans release (deux files distinctes, FR-16).
	if len(broker.Commands()) != 0 {
		t.Errorf("file de commande = %d, want 0 (aucune release)", len(broker.Commands()))
	}
}

// assertClosedTrail vérifie qu'une exécution laisse une piste d'audit close :
// nombre d'entrées attendu et **entrée terminale « decision »** (jamais une
// piste tronquée, même sur erreur).
func assertClosedTrail(t *testing.T, j *audit.Journal, wantLen int) {
	t.Helper()
	es := j.Entries()
	if len(es) != wantLen {
		t.Fatalf("audit = %d entrées, want %d : %+v", len(es), wantLen, es)
	}
	if last := es[len(es)-1]; last.Action != "decision" {
		t.Errorf("dernière entrée Action=%q, want decision (piste non close)", last.Action)
	}
}

func TestPrequalifyAssessError(t *testing.T) {
	sess := newIdentitySession(t)
	j := audit.New(audit.WithClock(fixedClock()))
	o := NewOrchestrator(sess, stubReasoner{}, j, testPolicy())
	a := goodApplicant()
	a.ID = "A00002"     // existe dans le jeu → identité OK
	a.MonthlyIncome = 0 // mais Assess échoue (revenu nul)
	if _, err := o.Prequalify(context.Background(), a); err == nil {
		t.Fatal("erreur attendue sur demande invalide")
	}
	assertClosedTrail(t, j, 4) // verify_identity, screen_watchlists, assess(error), decision(error)
}

type failingReasoner struct{}

func (failingReasoner) Reason(context.Context, ReasonRequest) (ReasonResult, error) {
	return ReasonResult{}, errors.New("reasoner down")
}

func TestPrequalifyReasonerError(t *testing.T) {
	sess := newIdentitySession(t)
	j := audit.New()
	o := NewOrchestrator(sess, failingReasoner{}, j, testPolicy())
	if _, err := o.Prequalify(context.Background(), goodApplicant()); err == nil {
		t.Fatal("erreur attendue quand le reasoner échoue")
	}
	assertClosedTrail(t, j, 5) // verify_identity, screen_watchlists, assess, reason(error), decision(error)
}

// TestPrequalifyFailClosed : PEP indisponible → la garde non contournable retombe
// fail-closed (L1→L0), le premier appel d'outil échoue, et la retombée est
// journalisée — aucune action autonome (invariant 3, §12.4 ; mécanisme M3.2 câblé
// dans la boucle via WithPEP).
func TestPrequalifyFailClosed(t *testing.T) {
	sess := newIdentitySession(t)
	j := audit.New(audit.WithClock(fixedClock()))
	p := pep.New(j, pep.L1)
	p.SetAvailable(false)
	o := NewOrchestrator(sess, stubReasoner{}, j, testPolicy()).WithPEP(p, "agent:orchestrator")
	if _, err := o.Prequalify(context.Background(), goodApplicant()); err == nil {
		t.Fatal("PEP indisponible : la boucle devrait échouer fail-closed")
	}
	found := false
	for _, e := range j.Entries() {
		if e.Action == "pep.fail_closed" {
			found = true
		}
	}
	if !found {
		t.Errorf("retombée fail-closed non journalisée : %+v", j.Entries())
	}
}

func TestPrequalifyIdentityCallError(t *testing.T) {
	sess := newIdentitySession(t)
	if err := sess.Close(); err != nil { // fermer → CallTool échoue
		t.Fatalf("fermeture : %v", err)
	}
	j := audit.New()
	o := NewOrchestrator(sess, stubReasoner{}, j, testPolicy())
	if _, err := o.Prequalify(context.Background(), goodApplicant()); err == nil {
		t.Fatal("erreur attendue quand la session MCP est fermée")
	}
	assertClosedTrail(t, j, 2) // verify_identity(error), decision(error)
}
