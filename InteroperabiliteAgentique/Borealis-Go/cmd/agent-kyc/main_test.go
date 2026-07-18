package main

import "testing"

// TestRunRequiresIDPSecret exécute run() (T5 : jamais exercé jusqu'ici) via le
// repli fail-fast sans IDP_SECRET (m9) — rapide et sans écoute réseau, tout en
// prouvant que le nom d'agent câblé (borealis-kyc) est valide (SpecByName
// vérifié avant le repli).
func TestRunRequiresIDPSecret(t *testing.T) {
	t.Setenv("IDP_SECRET", "")
	if code := run(); code != 3 {
		t.Errorf("run() sans IDP_SECRET = %d, want 3", code)
	}
}
