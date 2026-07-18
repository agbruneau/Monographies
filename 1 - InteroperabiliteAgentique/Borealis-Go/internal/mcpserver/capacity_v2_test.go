package mcpserver

import "testing"

// TestComputeCapacityV2StressRateFloor : un stressRate qui annule ou inverse le
// taux choqué (0,09 + stressRate <= 0) doit être rejeté, jamais calculé — sinon
// pmt = amount*0/(1-1) = NaN (non sérialisable en JSON) ou une mensualité
// sous-évaluée qui approuve un dossier que v1 refuserait (M1).
func TestComputeCapacityV2StressRateFloor(t *testing.T) {
	in := CapacityInV2{Income: 120000, Debts: 500, RequestedAmount: 20000, TermMonths: 60, StressRate: -0.09}
	if _, err := computeCapacityV2(in); err == nil {
		t.Fatalf("stressRate=-0.09 (taux net = 0) : erreur attendue")
	}
	in.StressRate = -0.19
	if _, err := computeCapacityV2(in); err == nil {
		t.Fatalf("stressRate=-0.19 (taux net < 0) : erreur attendue")
	}
}

// TestCapacityV2ToolRejectsInvalidStressRate : rejet en IsError (ADR-0001), pas
// en erreur protocole -32602 — cf. toolError.
func TestCapacityV2ToolRejectsInvalidStressRate(t *testing.T) {
	cs := connectServer(t, NewCapacityServerV2())
	res := callStructured(t, cs, "compute_capacity", map[string]any{
		"income": 120000.0, "debts": 500.0, "requestedAmount": 20000.0, "termMonths": 60, "stressRate": -0.09,
	}, nil)
	if !res.IsError {
		t.Errorf("stressRate=-0.09 : IsError attendu, eu %+v", res)
	}
}
