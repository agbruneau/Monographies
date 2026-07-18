package testfakes

import (
	"context"
	"testing"

	"github.com/agbruneau/borealis/internal/agent"
)

func TestStubReasonerScenarios(t *testing.T) {
	r := NewStubReasoner()
	cases := map[string]agent.Recommendation{
		"approved": agent.RecommendPreQualify,
		"denied":   agent.RecommendDecline,
		"escalate": agent.RecommendEscalate,
		"unknown":  agent.RecommendEscalate, // repli prudent
		"":         agent.RecommendEscalate, // signal absent
	}
	for scenario, want := range cases {
		got, err := r.Reason(context.Background(), agent.ReasonRequest{
			ApplicantID: "a1",
			Signals:     map[string]string{"scenario": scenario},
		})
		if err != nil {
			t.Fatalf("scénario %q : erreur %v", scenario, err)
		}
		if got.Recommendation != want {
			t.Errorf("scénario %q : Recommendation = %q, want %q", scenario, got.Recommendation, want)
		}
	}
}

// TestStubDeterminism : deux appels identiques → résultat identique.
func TestStubDeterminism(t *testing.T) {
	r := NewStubReasoner()
	req := agent.ReasonRequest{ApplicantID: "a1", Signals: map[string]string{"scenario": "approved"}}
	a, _ := r.Reason(context.Background(), req)
	b, _ := r.Reason(context.Background(), req)
	if a != b {
		t.Errorf("non déterministe : %+v != %+v", a, b)
	}
}
