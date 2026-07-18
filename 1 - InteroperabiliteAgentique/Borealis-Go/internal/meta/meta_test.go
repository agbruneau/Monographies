package meta

import "testing"

func TestInfo(t *testing.T) {
	// Sous `go test` (sans ldflags), les valeurs par défaut s'appliquent :
	// l'assertion porte sur le FORMAT exact (ordre des champs, libellés
	// « commit »/« built », ponctuation), pas seulement sur la présence des vars.
	const want = "borealis dev (commit unknown, built unknown)"
	if got := Info(); got != want {
		t.Errorf("Info() = %q, want %q", got, want)
	}
}
