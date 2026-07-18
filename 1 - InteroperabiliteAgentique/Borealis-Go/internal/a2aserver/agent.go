// Package a2aserver expose les agents Borealis comme serveurs A2A : Agent Card
// signée servie sur /.well-known/agent-card.json et point d'invocation JSON-RPC,
// via le SDK a2a-go. FR-01 (card conforme), §6.1 (inventaire des agents).
package a2aserver

import (
	"context"
	"iter"
	"net/http"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/a2aproject/a2a-go/v2/a2asrv"
	borealisa2a "github.com/agbruneau/borealis/pkg/a2a"
)

// InvokePath est le chemin d'invocation JSON-RPC des agents.
const InvokePath = "/a2a"

// Spec décrit un agent (§6.1) : nom et compétence unique.
type Spec struct {
	Name        string
	SkillID     string
	SkillName   string
	Description string
}

// Specs est l'inventaire des 5 agents A2A (PRD §6.1).
var Specs = []Spec{
	{"borealis-orchestrator", "orchestrate-prequal", "Orchestrate Pre-qualification", "Orchestrateur de pré-qualification de crédit"},
	{"borealis-kyc", "verify-identity", "Verify Identity", "Vérification d'identité synthétique (KYC)"},
	{"borealis-scoring", "compute-credit-score", "Compute Credit Score", "Évaluation de solvabilité (score, ratios)"},
	{"borealis-policy", "apply-lending-policy", "Apply Lending Policy", "Application des règles d'octroi par segment"},
	{"borealis-explainer", "explain-decision", "Explain Decision", "Génération de l'avis lisible fr/en"},
}

// SpecByName retourne la spec de l'agent nommé (nom de card, ex. "borealis-kyc")
// et true si elle existe. Sert de câblage aux binaires cmd/agent-*/orchestrator.
func SpecByName(name string) (Spec, bool) {
	for _, s := range Specs {
		if s.Name == name {
			return s, true
		}
	}
	return Spec{}, false
}

// BuildCard construit l'Agent Card d'un agent, avec son interface JSON-RPC à
// baseURL+InvokePath.
func BuildCard(s Spec, baseURL string) *a2asdk.AgentCard {
	return &a2asdk.AgentCard{
		Name:                s.Name,
		Version:             "1.0.0",
		Description:         s.Description,
		SupportedInterfaces: []*a2asdk.AgentInterface{a2asdk.NewAgentInterface(baseURL+InvokePath, a2asdk.TransportProtocolJSONRPC)},
		Capabilities:        a2asdk.AgentCapabilities{Streaming: true},
		DefaultInputModes:   []string{"application/json"},
		DefaultOutputModes:  []string{"application/json"},
		Skills:              []a2asdk.AgentSkill{{ID: s.SkillID, Name: s.SkillName, Description: s.Description, Tags: []string{"borealis"}}},
	}
}

// Mux assemble le routeur HTTP d'un agent : point JSON-RPC + Agent Card servie
// sur le chemin well-known.
func Mux(card *a2asdk.AgentCard, exec a2asrv.AgentExecutor) *http.ServeMux {
	rh := a2asrv.NewHandler(exec)
	mux := http.NewServeMux()
	mux.Handle(InvokePath, a2asrv.NewJSONRPCHandler(rh))
	mux.Handle(a2asrv.WellKnownAgentCardPath, a2asrv.NewStaticAgentCardHandler(card))
	return mux
}

// ReplyExecutor est un exécuteur A2A minimal qui répond par un message fixe
// (stub M2 ; la logique métier des agents — appels MCP — est câblée en M2.4).
type ReplyExecutor struct {
	Reply string
}

var _ a2asrv.AgentExecutor = (*ReplyExecutor)(nil)

// Execute émet un unique message de réponse.
func (e *ReplyExecutor) Execute(_ context.Context, _ *a2asrv.ExecutorContext) iter.Seq2[a2asdk.Event, error] {
	return func(yield func(a2asdk.Event, error) bool) {
		yield(a2asdk.NewMessage(a2asdk.MessageRoleAgent, a2asdk.NewTextPart(e.Reply)), nil)
	}
}

// Cancel n'a rien à annuler (exécution synchrone immédiate).
func (e *ReplyExecutor) Cancel(_ context.Context, _ *a2asrv.ExecutorContext) iter.Seq2[a2asdk.Event, error] {
	return func(_ func(a2asdk.Event, error) bool) {}
}

// LifecycleExecutor émet un cycle de vie de Task : working, puis un état final
// (completed par défaut, ou input-required) portant un message de réponse. Les
// émissions passent par une **machine à états** qui refuse toute transition
// depuis un état terminal (FR-07 appliqué sur le chemin réel, pas un helper isolé).
type LifecycleExecutor struct {
	Reply             string
	Final             a2asdk.TaskState // vide → completed
	IllegalAfterFinal bool             // tente une transition post-terminale (bloquée) — pour les tests FR-07
}

var _ a2asrv.AgentExecutor = (*LifecycleExecutor)(nil)

// Execute émet submitted (Task), working, puis l'état final, sous garde d'états.
func (e *LifecycleExecutor) Execute(_ context.Context, execCtx *a2asrv.ExecutorContext) iter.Seq2[a2asdk.Event, error] {
	return func(yield func(a2asdk.Event, error) bool) {
		m := borealisa2a.NewTaskMachine() // submitted
		if !yield(a2asdk.NewSubmittedTask(execCtx, execCtx.Message), nil) {
			return
		}
		if !guardedEmit(yield, execCtx, m, a2asdk.TaskStateWorking, nil) {
			return
		}
		final := e.Final
		if final == "" {
			final = a2asdk.TaskStateCompleted
		}
		msg := a2asdk.NewMessage(a2asdk.MessageRoleAgent, a2asdk.NewTextPart(e.Reply))
		if !guardedEmit(yield, execCtx, m, final, msg) {
			return
		}
		// Transition post-terminale : la garde (FR-07) la bloque — aucun
		// événement n'est émis.
		if e.IllegalAfterFinal {
			guardedEmit(yield, execCtx, m, a2asdk.TaskStateWorking, nil)
		}
	}
}

// guardedEmit n'émet une mise à jour de statut que si la transition est permise
// par la machine à états (états terminaux irréversibles, FR-07).
func guardedEmit(yield func(a2asdk.Event, error) bool, execCtx *a2asrv.ExecutorContext, m *borealisa2a.TaskMachine, state a2asdk.TaskState, msg *a2asdk.Message) bool {
	if err := m.To(state); err != nil {
		return true // transition interdite → non émise (la garde tient)
	}
	return yield(a2asdk.NewStatusUpdateEvent(execCtx, state, msg), nil)
}

// Cancel n'a rien à annuler.
func (e *LifecycleExecutor) Cancel(_ context.Context, _ *a2asrv.ExecutorContext) iter.Seq2[a2asdk.Event, error] {
	return func(_ func(a2asdk.Event, error) bool) {}
}
