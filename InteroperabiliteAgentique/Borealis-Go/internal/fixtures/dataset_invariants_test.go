package fixtures

import "testing"

// watchlistSurnames reflète la liste OFAC/PPE simulée du serveur Identity
// (mcpserver.ofacSimNames). Dupliquée ici pour éviter un cycle d'import : ce
// test garde l'invariant, pas la valeur.
var watchlistSurnames = map[string]bool{"Rossi": true, "Diallo": true}

// TestDefaultDatasetWatchlistReachability garantit que le jeu synthétique par
// défaut contient **à la fois** des dossiers sur liste (score AML > 0,8 →
// escalade HITL de M3.6 atteignable) et des dossiers propres (happy path non
// faussement escaladé), ET verrouille les **identifiants exacts** dont dépendent
// d'autres tests : A00001 (propre, happy path) et A00007 (sur liste, escalade).
// Un changement de graine ou de la liste des noms qui déplacerait ces dossiers
// doit échouer ICI — pas seulement, plus loin, dans les tests qui les codent en
// dur (identity_test, orchestrator_test).
func TestDefaultDatasetWatchlistReachability(t *testing.T) {
	ds := Generate(DefaultSeed, 100)
	var hit, clean int
	for _, b := range ds.Borrowers {
		if watchlistSurnames[b.LastName] {
			hit++
		} else {
			clean++
		}
	}
	if hit == 0 {
		t.Error("aucun dossier sur liste : le chemin d'escalade AML (M3.6) est inatteignable")
	}
	if clean == 0 {
		t.Error("aucun dossier propre : le happy path serait faussement escaladé")
	}

	// Identifiants figés par d'autres tests : la propriété doit tenir au dossier près.
	if b, ok := ds.Borrower("A00001"); !ok || watchlistSurnames[b.LastName] {
		t.Errorf("A00001 attendu propre (happy path), a %q", b.LastName)
	}
	if b, ok := ds.Borrower("A00007"); !ok || !watchlistSurnames[b.LastName] {
		t.Errorf("A00007 attendu sur liste (escalade AML), a %q", b.LastName)
	}
}
