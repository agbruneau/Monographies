package mcpserver

import "testing"

func TestComputeCapacity(t *testing.T) {
	out := computeCapacity(CapacityIn{Income: 120000, Debts: 500, RequestedAmount: 20000, TermMonths: 60})
	if out.MonthlyPayment <= 0 || out.DTIRatio <= 0 {
		t.Fatalf("nominal : %+v", out)
	}
	// Entrées dégénérées → hors capacité (jamais capacityOk sur un non-sens).
	if o := computeCapacity(CapacityIn{Income: 120000, RequestedAmount: 20000, TermMonths: 0}); o.MonthlyPayment != 0 || o.CapacityOk {
		t.Errorf("term 0 : %+v, want pmt=0 capacityOk=false", o)
	}
	if o := computeCapacity(CapacityIn{Income: 120000, RequestedAmount: -500000, TermMonths: 60}); o.CapacityOk {
		t.Errorf("montant négatif : capacityOk true inattendu (%+v)", o)
	}
	if o := computeCapacity(CapacityIn{Income: 0, RequestedAmount: 20000, TermMonths: 60}); o.CapacityOk {
		t.Errorf("revenu 0 : capacityOk true inattendu (%+v)", o)
	}
	if o := computeCapacity(CapacityIn{Income: 120000, Debts: -100, RequestedAmount: 20000, TermMonths: 60}); o.CapacityOk {
		t.Errorf("dette négative : capacityOk true inattendu (%+v)", o)
	}
	// DTI élevé → capacityOk false.
	if o := computeCapacity(CapacityIn{Income: 12000, Debts: 5000, RequestedAmount: 50000, TermMonths: 12}); o.CapacityOk {
		t.Errorf("DTI élevé : capacityOk true inattendu (%+v)", o)
	}
}

func TestCapacityRoundtrip(t *testing.T) {
	cs := connectServer(t, NewCapacityServer())
	var out CapacityOut
	res := callStructured(t, cs, "compute_capacity", map[string]any{
		"income": 120000.0, "debts": 500.0, "requestedAmount": 20000.0, "termMonths": 60,
	}, &out)
	if res.IsError {
		t.Fatalf("erreur : %+v", res.Content)
	}
	if !out.CapacityOk {
		t.Errorf("attendu capacityOk pour un bon dossier : %+v", out)
	}
}

func BenchmarkComputeCapacity(b *testing.B) {
	in := CapacityIn{Income: 120000, Debts: 500, RequestedAmount: 20000, TermMonths: 60}
	for i := 0; i < b.N; i++ {
		_ = computeCapacity(in)
	}
}
