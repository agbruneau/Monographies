//go:build e2e

// Les 6 scénarios e2e du PRD §12.4 (M5.1) composent le système réel — serveur
// MCP Identity (transport en mémoire), orchestrateur, fake Reasoner déterministe,
// journal WORM, PEP, HITL, explicabilité — et produisent une **sortie canonique**
// (champs volatils exclus : horodatages, hachages, task id) déterministe d'une
// exécution à l'autre (G-8, NFR-10). Le happy path est figé par un golden.
package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/agbruneau/borealis/internal/agent"
	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/conformance"
	"github.com/agbruneau/borealis/internal/creditdomain"
	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/agbruneau/borealis/internal/hitl"
	"github.com/agbruneau/borealis/internal/mcpserver"
	"github.com/agbruneau/borealis/internal/pep"
	"github.com/agbruneau/borealis/internal/testfakes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func fixedClock() func() time.Time {
	ts := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	return func() time.Time { return ts }
}

func identitySession(t *testing.T, ds *fixtures.Dataset) *mcp.ClientSession {
	return testfakes.ConnectInMemory(t, mcpserver.NewIdentityServer(ds))
}

// fakeReasoner mappe l'issue métier à une recommandation fixe (déterministe, G-8).
type fakeReasoner struct{}

func (fakeReasoner) Reason(_ context.Context, req agent.ReasonRequest) (agent.ReasonResult, error) {
	switch req.Signals["scenario"] {
	case "approved":
		return agent.ReasonResult{Recommendation: agent.RecommendPreQualify}, nil
	case "denied":
		return agent.ReasonResult{Recommendation: agent.RecommendDecline}, nil
	default:
		return agent.ReasonResult{Recommendation: agent.RecommendEscalate}, nil
	}
}

// timeoutReasoner simule un délai modèle dépassé (scénario 3).
type timeoutReasoner struct{}

func (timeoutReasoner) Reason(context.Context, agent.ReasonRequest) (agent.ReasonResult, error) {
	return agent.ReasonResult{}, context.DeadlineExceeded
}

func policy() creditdomain.Policy {
	return creditdomain.Policy{MinScore: 600, MaxDTI: 0.43, ReferScore: 660, ReferDTI: 0.36}
}

func goodApplicant(id string) creditdomain.Applicant {
	return creditdomain.Applicant{
		ID: id, MonthlyIncome: 1000000, MonthlyDebt: 50000,
		RequestedAmount: 1000000, AnnualRate: 0.079, TermMonths: 60, CreditScore: 720,
	}
}

// --- Sortie canonique (exclut les champs volatils) ---

type canonExplain struct {
	Outcome    string   `json:"outcome"`
	Summary    string   `json:"summary"`
	Factors    []string `json:"factors"`
	Disclaimer string   `json:"disclaimer"`
}

type canonical struct {
	Scenario       string        `json:"scenario"`
	Outcome        string        `json:"outcome"`
	Recommendation string        `json:"recommendation,omitempty"`
	Reasons        []string      `json:"reasons,omitempty"`
	HITLPending    int           `json:"hitl_pending"`
	AuditActions   []string      `json:"audit_actions"`
	Explanation    *canonExplain `json:"explanation,omitempty"`
	ExplainError   string        `json:"explain_error,omitempty"`
	Error          string        `json:"error,omitempty"`
}

func auditActions(j *audit.Journal) []string {
	es := j.Entries()
	out := make([]string, 0, len(es))
	for _, e := range es {
		out = append(out, e.Action)
	}
	return out
}

func explainOf(d creditdomain.Decision) *canonExplain {
	e := conformance.Explain(d, "fr")
	return &canonExplain{Outcome: e.Outcome, Summary: e.Summary, Factors: e.Factors, Disclaimer: e.Disclaimer}
}

func canonJSON(t *testing.T, c canonical) []byte {
	t.Helper()
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// compareGolden compare la sortie canonique au golden (UPDATE_GOLDEN=1 régénère).
//
// L'égalité est une comparaison de CHAÎNE stricte (octet à octet), pas un
// diff structurel après désérialisation : c'est délibéré, pas une fragilité
// non examinée (T7). Le déterminisme exigé par le PRDPlan §2 est « make e2e
// x3 → sorties canoniques identiques » — l'égalité de chaîne EST cette
// exigence (y compris la mise en forme MarshalIndent), pas une approximation
// de plus haut niveau. Un diff structurel masquerait une régression de
// formatage qui romprait pourtant la reproductibilité byte-à-byte visée.
func compareGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("..", "golden", name)
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden %s absent (lancer UPDATE_GOLDEN=1) : %v", name, err)
	}
	if string(got) != string(want) {
		t.Errorf("%s diffère du golden :\n got=%s\nwant=%s", name, got, want)
	}
}

func newOrch(t *testing.T, ds *fixtures.Dataset, j *audit.Journal, r agent.Reasoner) *agent.Orchestrator {
	return agent.NewOrchestrator(identitySession(t, ds), r, j, policy())
}

// 1. Happy path : identité vérifiée, AML propre, bon dossier → pré-qualifié.
// Trace canonique figée par golden (trace-approved.json, PRD §12.3).
func TestScenarioHappyPath(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	j := audit.New(audit.WithClock(fixedClock()))
	res, err := newOrch(t, ds, j, fakeReasoner{}).Prequalify(context.Background(), goodApplicant("A00001"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Decision.Outcome != creditdomain.OutcomePreQualified {
		t.Errorf("outcome = %s, want pre_qualified (jamais un octroi ferme, FR-25)", res.Decision.Outcome)
	}
	c := canonical{
		Scenario: "happy_path", Outcome: res.Decision.Outcome.String(),
		Recommendation: string(res.Recommendation), Reasons: res.Decision.Reasons,
		AuditActions: auditActions(j), Explanation: explainOf(res.Decision),
	}
	compareGolden(t, "trace-approved.json", canonJSON(t, c))
}

// 2. Refus AML : dossier sur liste (A00007, Diallo) → refus + escalade HITL.
func TestScenarioAMLRefusal(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	j := audit.New(audit.WithClock(fixedClock()))
	broker := hitl.New(j)
	o := newOrch(t, ds, j, fakeReasoner{}).WithHITL(broker)
	res, err := o.Prequalify(context.Background(), goodApplicant("A00007"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Decision.Outcome != creditdomain.OutcomeDeclined {
		t.Errorf("outcome = %s, want declined (AML)", res.Decision.Outcome)
	}
	if broker.PendingProposals() != 1 {
		t.Errorf("propositions HITL = %d, want 1", broker.PendingProposals())
	}
	if len(broker.Commands()) != 0 {
		t.Errorf("file de commande = %d, want 0 (aucune release)", len(broker.Commands()))
	}
	// L'escalade précède le scellement (FR-16).
	acts := auditActions(j)
	iP, iD := indexOf(acts, "hitl.propose"), indexOf(acts, "decision")
	if iP < 0 || iD < 0 || iP >= iD {
		t.Errorf("ordre d'audit propose=%d decision=%d, want propose < decision : %v", iP, iD, acts)
	}
}

// 3. Timeout modèle : le Reasoner dépasse le délai → échec propre, piste close.
func TestScenarioModelTimeout(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	j := audit.New(audit.WithClock(fixedClock()))
	_, err := newOrch(t, ds, j, timeoutReasoner{}).Prequalify(context.Background(), goodApplicant("A00001"))
	if err == nil {
		t.Fatal("délai modèle dépassé : la boucle devrait échouer proprement")
	}
	es := j.Entries()
	if last := es[len(es)-1]; last.Action != "decision" {
		t.Errorf("dernière entrée = %q, want decision (piste non close)", last.Action)
	}
}

// 4. Signaux conflictuels : dossier limite (score 640) → renvoi humain (refer).
func TestScenarioConflictingSignals(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	j := audit.New(audit.WithClock(fixedClock()))
	a := goodApplicant("A00001")
	a.CreditScore = 640 // >= MinScore 600 mais < ReferScore 660 → refer
	res, err := newOrch(t, ds, j, fakeReasoner{}).Prequalify(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if res.Decision.Outcome != creditdomain.OutcomeRefer {
		t.Errorf("outcome = %s, want refer", res.Decision.Outcome)
	}
	if res.Recommendation != agent.RecommendEscalate {
		t.Errorf("recommandation = %s, want escalate", res.Recommendation)
	}
}

// explainStep est l'étage aval d'explicabilité (agent Explicabilité), injectable
// pour simuler sa défaillance. L'explicabilité n'est PAS dans le chemin
// décisionnel de Prequalify : elle s'applique après, sur la décision produite.
type explainStep func(creditdomain.Decision) (conformance.Explanation, error)

// composeWithExplain applique l'étage d'explicabilité à un résultat : en cas
// d'échec, l'avis est absent (dégradation gracieuse), la décision demeure.
func composeWithExplain(res agent.Result, j *audit.Journal, explain explainStep) canonical {
	c := canonical{
		Scenario: "explainability_failure", Outcome: res.Decision.Outcome.String(),
		Recommendation: string(res.Recommendation), AuditActions: auditActions(j),
	}
	if exp, err := explain(res.Decision); err != nil {
		c.ExplainError = err.Error()
	} else {
		c.Explanation = &canonExplain{Outcome: exp.Outcome, Summary: exp.Summary, Factors: exp.Factors, Disclaimer: exp.Disclaimer}
	}
	return c
}

// 5. Échec Explicabilité : l'agent d'explicabilité est indisponible → dégradation
// gracieuse RÉELLEMENT exercée (le chemin d'échec est pris : aucun avis produit),
// mais la décision et la piste d'audit tiennent (G-2). Contraste avec un
// explicateur sain pour prouver que l'étage produit bien un avis quand il marche.
func TestScenarioExplainabilityFailure(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	j := audit.New(audit.WithClock(fixedClock()))
	res, err := newOrch(t, ds, j, fakeReasoner{}).Prequalify(context.Background(), goodApplicant("A00001"))
	if err != nil {
		t.Fatal(err)
	}

	// Contraste : un explicateur SAIN produit bien un avis (non vacuité).
	healthy := composeWithExplain(res, j, func(d creditdomain.Decision) (conformance.Explanation, error) {
		return conformance.Explain(d, "fr"), nil
	})
	if healthy.Explanation == nil || healthy.Explanation.Summary == "" {
		t.Fatal("prérequis : un explicateur sain devrait produire un avis")
	}

	// Défaillant : le chemin d'échec est pris.
	failing := composeWithExplain(res, j, func(creditdomain.Decision) (conformance.Explanation, error) {
		return conformance.Explanation{}, errors.New("explainer unavailable")
	})
	if failing.Explanation != nil {
		t.Error("aucun avis ne devrait être produit quand l'explicateur échoue")
	}
	if failing.ExplainError == "" {
		t.Error("l'échec de l'explicateur devrait être capturé")
	}
	// Dégradation gracieuse : la décision et l'audit tiennent malgré l'échec.
	if failing.Outcome != creditdomain.OutcomePreQualified.String() {
		t.Errorf("la décision doit tenir malgré l'échec de l'avis : %s", failing.Outcome)
	}
	if len(failing.AuditActions) == 0 {
		t.Error("la piste d'audit doit subsister malgré l'échec de l'avis")
	}
	_ = canonJSON(t, failing)
}

// 6. Fail-closed : PEP indisponible → retombée L1→L0 journalisée, aucune action
// autonome (le premier appel d'outil échoue). Mécanisme livré en M3.2.
func TestScenarioFailClosed(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	j := audit.New(audit.WithClock(fixedClock()))
	p := pep.New(j, pep.L1)
	p.SetAvailable(false) // PEP/orchestrateur indisponible
	o := newOrch(t, ds, j, fakeReasoner{}).WithPEP(p, "agent:orchestrator")
	_, err := o.Prequalify(context.Background(), goodApplicant("A00001"))
	if err == nil {
		t.Fatal("PEP indisponible : aucune action autonome — la boucle devrait échouer")
	}
	if indexOf(auditActions(j), "pep.fail_closed") < 0 {
		t.Errorf("retombée fail-closed non journalisée : %v", auditActions(j))
	}
}

func indexOf(ss []string, target string) int {
	for i, s := range ss {
		if s == target {
			return i
		}
	}
	return -1
}
