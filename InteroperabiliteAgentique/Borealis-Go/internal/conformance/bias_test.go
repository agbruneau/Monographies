package conformance

import "testing"

func groupRuns(group string, n, approved int) []GroupDecision {
	out := make([]GroupDecision, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, GroupDecision{Group: group, Approved: i < approved})
	}
	return out
}

func TestAnalyzeBiasAlert(t *testing.T) {
	// Groupe A à 90 %, groupe B à 40 % → divergence 0,50 > 0,15 → alerte.
	ds := append(groupRuns("A", 10, 9), groupRuns("B", 10, 4)...)
	r := AnalyzeBias(ds, 0.15)
	if !r.Alert {
		t.Errorf("divergence %.2f devrait alerter (> 0.15)", r.MaxDivergence)
	}
	// Valeurs exactes verrouillées (une erreur de dénominateur donnant un résultat
	// faux-mais-positif passerait une simple assertion « > 0 ») :
	// MaxDivergence = 0,9 − 0,4 = 0,50 ; Gini des taux [0,9 ; 0,4] = 1,0/(2·2·1,3)
	// = 0,192308 (moyenne des écarts absolus normalisée, définition standard).
	if diff := r.MaxDivergence - 0.5; diff < -1e-9 || diff > 1e-9 {
		t.Errorf("MaxDivergence = %.6f, want 0.500000", r.MaxDivergence)
	}
	if diff := r.Gini - 0.192308; diff < -1e-4 || diff > 1e-4 {
		t.Errorf("Gini = %.6f, want 0.192308", r.Gini)
	}
	if len(r.Groups) != 2 {
		t.Errorf("groupes = %d, want 2", len(r.Groups))
	}
}

func TestAnalyzeBiasNoAlert(t *testing.T) {
	ds := append(groupRuns("A", 10, 7), groupRuns("B", 10, 8)...)
	if r := AnalyzeBias(ds, 0.15); r.Alert {
		t.Errorf("divergence %.2f ne devrait pas alerter", r.MaxDivergence)
	}
}

func TestAnalyzeBiasEmpty(t *testing.T) {
	if r := AnalyzeBias(nil, 0.15); len(r.Groups) != 0 || r.Alert {
		t.Errorf("rapport vide attendu : %+v", r)
	}
}

func TestGiniEdgeCases(t *testing.T) {
	if giniOf(nil) != 0 {
		t.Error("gini(nil) devrait être 0")
	}
	if giniOf([]GroupStat{{Group: "A"}, {Group: "B"}}) != 0 {
		t.Error("gini de taux nuls devrait être 0")
	}
}

// TestAnalyzeBiasSingletonGroupAlert : un groupe singleton (n=1, 0 % approuvé)
// face à un grand groupe majoritairement approuvé doit déclencher l'alerte
// (T6, replié depuis zz_repro_test.go).
func TestAnalyzeBiasSingletonGroupAlert(t *testing.T) {
	ds := []GroupDecision{{Group: "Z", Approved: false}}
	for i := 0; i < 500; i++ {
		ds = append(ds, GroupDecision{Group: "A", Approved: i < 100}) // 100/500 = 0,20
	}
	r := AnalyzeBias(ds, 0.15)
	if !r.Alert {
		t.Fatalf("alerte attendue (groupe singleton) : MaxDivergence=%.4f groups=%+v", r.MaxDivergence, r.Groups)
	}
}
