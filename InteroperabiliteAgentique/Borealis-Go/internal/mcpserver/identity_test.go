package mcpserver

import (
	"testing"

	"github.com/agbruneau/borealis/internal/fixtures"
)

func TestVerifyIdentityPure(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	b := ds.Borrowers[0]
	cases := []struct {
		name   string
		in     VerifyIdentityInput
		want   bool
		reason string
	}{
		{"vide", VerifyIdentityInput{}, false, "missing_applicant_id"},
		{"inconnu", VerifyIdentityInput{ApplicantID: "A99999"}, false, "not_found"},
		{"existence seule", VerifyIdentityInput{ApplicantID: b.ID}, true, "verified"},
		{"nom correct", VerifyIdentityInput{ApplicantID: b.ID, Name: b.FirstName + " " + b.LastName}, true, "verified"},
		{"nom faux", VerifyIdentityInput{ApplicantID: b.ID, Name: "Faux Nom"}, false, "name_mismatch"},
		{"sin faux", VerifyIdentityInput{ApplicantID: b.ID, SIN: "000-000-000"}, false, "sin_mismatch"},
		{"sin correct", VerifyIdentityInput{ApplicantID: b.ID, SIN: b.SIN}, true, "verified"},
	}
	for _, c := range cases {
		got := verifyIdentity(ds, c.in)
		if got.Match != c.want || got.Reason != c.reason {
			t.Errorf("%s : %+v, want match=%v reason=%q", c.name, got, c.want, c.reason)
		}
	}
}

func TestCheckWatchlists(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	// A00001 (Farah Kaur) : propre → pas de hit, score SOUS les seuils FINTRAC
	// (< 0.5). Le nom criblé est résolu depuis le dossier, pas depuis l'entrée :
	// un nom falsifié est ignoré.
	clear := checkWatchlists(ds, CheckWatchlistsInput{ApplicantID: "A00001", Name: "Faux Nom"})
	if clear.WatchlistHit || clear.AMLScore >= 0.5 {
		t.Errorf("propre : %+v, want hit=false et score<0.5", clear)
	}
	// Déterminisme (NFR-10).
	if again := checkWatchlists(ds, CheckWatchlistsInput{ApplicantID: "A00001"}); again != clear {
		t.Errorf("non déterministe : %+v != %+v", again, clear)
	}
	// A00007 (Gita Diallo) : nom sur liste OFAC simulée → hit + score élevé
	// (>= 0.85, au-dessus des seuils) — atteint via la seule résolution serveur.
	hit := checkWatchlists(ds, CheckWatchlistsInput{ApplicantID: "A00007"})
	if !hit.WatchlistHit || hit.AMLScore < 0.85 {
		t.Errorf("hit : %+v, want hit=true et score>=0.85", hit)
	}
}

func TestIdentityRoundtrip(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	cs := connectServer(t, NewIdentityServer(ds))

	var vo VerifyIdentityOutput
	if res := callStructured(t, cs, "verify_identity", map[string]any{"applicantId": ds.Borrowers[0].ID}, &vo); res.IsError || !vo.Match {
		t.Errorf("verify_identity : isErr=%v out=%+v", res.IsError, vo)
	}
	var wo CheckWatchlistsOutput
	if res := callStructured(t, cs, "check_watchlists", map[string]any{"applicantId": "A00001", "name": "Alix Tremblay"}, &wo); res.IsError {
		t.Error("check_watchlists en erreur")
	}
}

// TestVerifyIdentityInvalidInput : le SDK valide l'entrée contre le schéma
// inféré ; violation → IsError (ADR-0001), pas -32602.
func TestVerifyIdentityInvalidInput(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 10)
	cs := connectServer(t, NewIdentityServer(ds))
	for _, args := range []map[string]any{
		{},                   // applicantId requis manquant
		{"applicantId": 123}, // type invalide
	} {
		if res := callStructured(t, cs, "verify_identity", args, nil); !res.IsError {
			t.Errorf("args=%v : IsError attendu (rejet de schéma)", args)
		}
	}
}
