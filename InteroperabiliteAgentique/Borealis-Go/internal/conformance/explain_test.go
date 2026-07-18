package conformance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/agbruneau/borealis/internal/creditdomain"
)

func goldenDecision() creditdomain.Decision {
	return creditdomain.Decision{
		ApplicantID: "A00042",
		Outcome:     creditdomain.OutcomeDeclined,
		DTI:         0.52,
		Reasons:     []string{"credit_score_below_minimum", "dti_above_maximum"},
	}
}

// goldenPreQualified est le cas le plus sensible pour l'invariant 1 (jamais un
// octroi ferme) : le libellé exact et l'avertissement du chemin pré-qualifié
// sont figés par golden, pour qu'aucune révision ne les fasse dériver vers un
// vocabulaire d'octroi.
func goldenPreQualified() creditdomain.Decision {
	return creditdomain.Decision{
		ApplicantID: "A00042",
		Outcome:     creditdomain.OutcomePreQualified,
		DTI:         0.28,
		Reasons:     []string{"meets_prequalification_criteria"},
	}
}

func compareGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil { // #nosec G306 -- golden de test
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path) // #nosec G304 -- chemin de golden contrôlé
	if err != nil {
		t.Fatalf("golden %s absent (lancer UPDATE_GOLDEN=1) : %v", name, err)
	}
	if string(got) != string(want) {
		t.Errorf("%s diffère du golden :\n got=%s\nwant=%s", name, got, want)
	}
}

func TestExplainGolden(t *testing.T) {
	cases := []struct {
		prefix string
		d      creditdomain.Decision
	}{
		{"explain", goldenDecision()},                  // refusé, 2 motifs
		{"explain-prequalified", goldenPreQualified()}, // pré-qualifié (invariant 1)
	}
	for _, c := range cases {
		for _, lang := range []string{"fr", "en"} {
			e := Explain(c.d, lang)
			got, err := json.MarshalIndent(e, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			compareGolden(t, c.prefix+"-"+lang+".json", got)
		}
	}
}

func TestExplainAllOutcomes(t *testing.T) {
	cases := []creditdomain.Decision{
		{ApplicantID: "A1", Outcome: creditdomain.OutcomeDeclined, Reasons: []string{"credit_score_below_minimum"}},
		{ApplicantID: "A2", Outcome: creditdomain.OutcomeRefer, Reasons: []string{"borderline_requires_human_review"}},
		{ApplicantID: "A3", Outcome: creditdomain.OutcomePreQualified, Reasons: []string{"meets_prequalification_criteria"}},
		{ApplicantID: "A4", Outcome: creditdomain.OutcomeDeclined, Reasons: []string{"motif_non_traduit"}}, // repli
	}
	for _, lang := range []string{"fr", "en", "de"} { // 'de' → repli fr
		for _, d := range cases {
			e := Explain(d, lang)
			if e.Summary == "" || len(e.Factors) == 0 || e.Disclaimer == "" || e.Outcome == "" {
				t.Errorf("explication incomplète pour %s/%s : %+v", lang, d.ApplicantID, e)
			}
		}
	}
}

func TestExplainFast(t *testing.T) {
	d := goldenDecision()
	start := time.Now()
	for i := 0; i < 1000; i++ {
		_ = Explain(d, "fr")
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Errorf("1000 explications en %v, want < 500ms", elapsed)
	}
}

func TestLangSelection(t *testing.T) {
	d := goldenDecision()
	// Étiquettes de locale (casse, région) : "en", "EN", "en-US", "en_CA" →
	// anglais (US-11) ; défaut/inconnu → français.
	for lang, want := range map[string]string{
		"en": "en", "EN": "en", "en-US": "en", "en_CA": "en", "fr": "fr",
		"FR": "fr", "fr-CA": "fr", "": "fr", "zz": "fr", "de": "fr", "english": "fr",
	} {
		if got := Explain(d, lang).Lang; got != want {
			t.Errorf("Explain(lang=%q).Lang = %q, want %q", lang, got, want)
		}
	}
}
