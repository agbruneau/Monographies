package mcpserver

import (
	"testing"

	"github.com/agbruneau/borealis/internal/fixtures"
)

func TestCreditReport(t *testing.T) {
	b := fixtures.Borrower{IncomeAnnual: 120000, MonthlyHousing: 2000, MonthlyDebts: 500, Score: 720, Defaults: 1}
	r := creditReport(b)
	if r.Score != 720 || r.Defaults != 1 || r.ABDRatio <= 0 || r.ATDRatio <= r.ABDRatio {
		t.Errorf("creditReport : %+v", r)
	}
	// Revenu <= 0 → ratios non calculables (sentinelle élevée), jamais 0, afin
	// que la politique rejette en aval.
	if r0 := creditReport(fixtures.Borrower{IncomeAnnual: 0, Score: 500}); r0.ABDRatio < 1 || r0.ATDRatio < 1 {
		t.Errorf("revenu 0 : %+v, want ratios élevés", r0)
	}
}

func TestBureauRoundtripAndNotFound(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	cs := connectServer(t, NewBureauServer(ds))

	var out BureauOut
	res := callStructured(t, cs, "get_credit_report", map[string]any{"applicantId": "A00001"}, &out)
	if res.IsError {
		t.Fatalf("A00001 en erreur : %+v", res.Content)
	}
	if out.Score < 300 || out.Score > 900 {
		t.Errorf("score hors plage : %d", out.Score)
	}
	// Inconnu → IsError (NotFound).
	if res2 := callStructured(t, cs, "get_credit_report", map[string]any{"applicantId": "ZZZ"}, nil); !res2.IsError {
		t.Error("ZZZ : IsError attendu (NotFound)")
	}
}

func BenchmarkCreditReport(b *testing.B) {
	bor := fixtures.Borrower{IncomeAnnual: 120000, MonthlyHousing: 2000, MonthlyDebts: 500, Score: 720}
	for i := 0; i < b.N; i++ {
		_ = creditReport(bor)
	}
}
