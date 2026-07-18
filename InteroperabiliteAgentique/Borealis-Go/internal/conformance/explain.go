// Package conformance porte les dispositifs de conformité : explicabilité fr/en
// (avis lisible) et détection de biais (divergence par groupe). Déterministe
// (golden), sans dépendance à un LLM.
package conformance

import (
	"fmt"
	"strings"

	"github.com/agbruneau/borealis/internal/creditdomain"
)

// Explanation est un avis d'explicabilité **hybride** : règles (motifs codés) +
// importance (facteurs quantitatifs) + contexte (issue + avertissement).
type Explanation struct {
	Lang       string   `json:"lang"`
	Outcome    string   `json:"outcome"`
	Summary    string   `json:"summary"`
	Factors    []string `json:"factors"`
	Disclaimer string   `json:"disclaimer"`
}

var outcomeText = map[string]map[creditdomain.Outcome]string{
	"fr": {
		creditdomain.OutcomeDeclined:     "Refusé",
		creditdomain.OutcomeRefer:        "Renvoi à un conseiller",
		creditdomain.OutcomePreQualified: "Pré-qualifié",
	},
	"en": {
		creditdomain.OutcomeDeclined:     "Declined",
		creditdomain.OutcomeRefer:        "Referred to an advisor",
		creditdomain.OutcomePreQualified: "Pre-qualified",
	},
}

var reasonText = map[string]map[string]string{
	"fr": {
		"credit_score_below_minimum":       "cote de crédit sous le minimum requis",
		"dti_above_maximum":                "ratio d'endettement supérieur au maximum",
		"borderline_requires_human_review": "dossier limite nécessitant une revue humaine",
		"meets_prequalification_criteria":  "critères de pré-qualification satisfaits",
		"identity_not_verified":            "identité non vérifiée",
		"aml_score_above_threshold":        "score AML supérieur au seuil de refus",
		"watchlist_hit":                    "correspondance sur liste de surveillance",
	},
	"en": {
		"credit_score_below_minimum":       "credit score below the required minimum",
		"dti_above_maximum":                "debt-to-income ratio above the maximum",
		"borderline_requires_human_review": "borderline file requiring human review",
		"meets_prequalification_criteria":  "pre-qualification criteria met",
		"identity_not_verified":            "identity not verified",
		"aml_score_above_threshold":        "AML score above the decline threshold",
		"watchlist_hit":                    "watchlist match",
	},
}

var disclaimerText = map[string]string{
	"fr": "Avis de pré-qualification à titre indicatif — jamais un octroi ferme.",
	"en": "Indicative pre-qualification advice — never a firm grant.",
}

// normLang normalise une étiquette de langue en "fr" ou "en" (fr par défaut ;
// US-11). Insensible à la casse et à la sous-étiquette de région, pour accepter
// une valeur relayée d'un en-tête Accept-Language ("EN", "en-US", "en_CA" →
// anglais) plutôt que replier silencieusement en français sur tout sauf "en".
func normLang(lang string) string {
	primary, _, _ := strings.Cut(lang, "-")
	primary, _, _ = strings.Cut(primary, "_")
	if strings.EqualFold(primary, "en") {
		return "en"
	}
	return "fr"
}

// Explain produit l'avis d'explicabilité d'une décision, dans la langue lang
// (fr par défaut). Déterministe (golden fr/en) ; ne produit jamais un octroi.
func Explain(d creditdomain.Decision, lang string) Explanation {
	l := normLang(lang)
	factors := make([]string, 0, len(d.Reasons)+1)
	for _, r := range d.Reasons {
		if txt, ok := reasonText[l][r]; ok {
			factors = append(factors, txt)
		} else {
			factors = append(factors, r) // repli : motif brut si non traduit
		}
	}
	if l == "en" {
		factors = append(factors, fmt.Sprintf("debt-to-income ratio: %.2f", d.DTI))
	} else {
		factors = append(factors, fmt.Sprintf("ratio d'endettement : %.2f", d.DTI))
	}

	outcome := outcomeText[l][d.Outcome]
	var summary string
	if l == "en" {
		summary = fmt.Sprintf("Application %s: %s.", d.ApplicantID, outcome)
	} else {
		summary = fmt.Sprintf("Demande %s : %s.", d.ApplicantID, outcome)
	}
	return Explanation{
		Lang:       l,
		Outcome:    outcome,
		Summary:    summary,
		Factors:    factors,
		Disclaimer: disclaimerText[l],
	}
}
