package mcpserver

import (
	"slices"
	"testing"

	"github.com/agbruneau/borealis/internal/fixtures"
)

func TestEvaluatePolicy(t *testing.T) {
	p := fixtures.Policy{Segment: "retail", AgeMin: 25, AgeMax: 34, ABDMax: 0.35, ATDMax: 0.43, MaxAmount: 20000}

	if o := evaluatePolicy(p, PolicyIn{ABDRatio: 0.30, ATDRatio: 0.40, Amount: 15000}); !o.Eligible || !slices.Contains(o.Reasons, "meets_policy") {
		t.Errorf("attendu éligible : %+v", o)
	}
	cases := []struct {
		name   string
		in     PolicyIn
		reason string
	}{
		{"abd", PolicyIn{ABDRatio: 0.40, ATDRatio: 0.40, Amount: 15000}, "abd_above_max"},
		{"atd", PolicyIn{ABDRatio: 0.30, ATDRatio: 0.50, Amount: 15000}, "atd_above_max"},
		{"amount", PolicyIn{ABDRatio: 0.30, ATDRatio: 0.40, Amount: 30000}, "amount_above_ceiling"},
	}
	for _, c := range cases {
		o := evaluatePolicy(p, c.in)
		if o.Eligible || !slices.Contains(o.Reasons, c.reason) {
			t.Errorf("%s : %+v, want inéligible + %s", c.name, o, c.reason)
		}
	}
}

func TestPolicyRoundtripAndNotFound(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 10)
	cs := connectServer(t, NewPolicyServer(ds))

	var out PolicyOut
	res := callStructured(t, cs, "evaluate_policy", map[string]any{
		"segment": "retail", "age": 30, "abdRatio": 0.30, "atdRatio": 0.40, "amount": 10000.0,
	}, &out)
	if res.IsError {
		t.Fatalf("retail/30 en erreur : %+v", res.Content)
	}
	if !out.Eligible {
		t.Errorf("attendu éligible : %+v", out)
	}
	// Segment inconnu → IsError (NotFound).
	if res2 := callStructured(t, cs, "evaluate_policy", map[string]any{
		"segment": "unknown", "age": 30, "abdRatio": 0.1, "atdRatio": 0.1, "amount": 1000.0,
	}, nil); !res2.IsError {
		t.Error("segment inconnu : IsError attendu (NotFound)")
	}
}

func BenchmarkEvaluatePolicy(b *testing.B) {
	p := fixtures.Policy{Segment: "retail", ABDMax: 0.35, ATDMax: 0.43, MaxAmount: 20000}
	in := PolicyIn{ABDRatio: 0.30, ATDRatio: 0.40, Amount: 15000}
	for i := 0; i < b.N; i++ {
		_ = evaluatePolicy(p, in)
	}
}
