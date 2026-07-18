// Package testfakes fournit des doubles déterministes (fakes) pour le gate :
// aucun LLM réel, sorties reproductibles (G-8, NFR-10).
package testfakes

import (
	"context"

	"github.com/agbruneau/borealis/internal/agent"
)

// StubReasoner est un Reasoner déterministe piloté par un scénario nommé
// (signal "scenario" : approved | denied | escalate). Scénario inconnu →
// repli prudent sur l'escalade humaine (fail-safe).
type StubReasoner struct {
	byScenario map[string]agent.ReasonResult
	fallback   agent.ReasonResult
}

var _ agent.Reasoner = (*StubReasoner)(nil)

// NewStubReasoner construit le fake avec les trois scénarios canoniques.
func NewStubReasoner() *StubReasoner {
	return &StubReasoner{
		byScenario: map[string]agent.ReasonResult{
			"approved": {Recommendation: agent.RecommendPreQualify, Rationale: "signals within pre-qualification bounds", Confidence: 0.92},
			"denied":   {Recommendation: agent.RecommendDecline, Rationale: "hard policy breach", Confidence: 0.97},
			"escalate": {Recommendation: agent.RecommendEscalate, Rationale: "conflicting signals", Confidence: 0.55},
		},
		fallback: agent.ReasonResult{Recommendation: agent.RecommendEscalate, Rationale: "unknown scenario: defer to human", Confidence: 0},
	}
}

// Reason retourne le résultat fixe du scénario, ou le repli prudent.
func (s *StubReasoner) Reason(_ context.Context, req agent.ReasonRequest) (agent.ReasonResult, error) {
	if r, ok := s.byScenario[req.Signals["scenario"]]; ok {
		return r, nil
	}
	return s.fallback, nil
}
