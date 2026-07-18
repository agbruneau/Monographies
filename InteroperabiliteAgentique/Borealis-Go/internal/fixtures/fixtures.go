// Package fixtures génère des jeux de données 100 % synthétiques, seedés et
// **déterministes** (PRD §9.5 ; aucun PII réel). Les serveurs MCP consomment
// le jeu en mémoire via Generate.
//
// ponytail: générateur seedé déterministe ; aucun PII réel.
package fixtures

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// DefaultSeed est la graine par défaut : même graine → mêmes données.
const DefaultSeed int64 = 20260707

// Segments de clientèle (coopérative fictive).
var segments = []string{"retail", "newcomer", "student", "business"}

// Borrower est un emprunteur synthétique.
type Borrower struct {
	ID             string  `json:"id"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	SIN            string  `json:"sin"` // jeton synthétique (SYN-…), volontairement hors du format d'un NAS
	IncomeAnnual   float64 `json:"incomeAnnual"`
	MonthlyDebts   float64 `json:"monthlyDebts"`
	MonthlyHousing float64 `json:"monthlyHousing"`
	Score          int     `json:"score"`
	Defaults       int     `json:"defaults"`
	Segment        string  `json:"segment"`
	Age            int     `json:"age"`
}

// Policy est une ligne de table d'octroi.
type Policy struct {
	Segment   string  `json:"segment" csv:"segment"`
	AgeMin    int     `json:"ageMin" csv:"age_min"`
	AgeMax    int     `json:"ageMax" csv:"age_max"`
	ABDMax    float64 `json:"abdMax" csv:"abd_max"`
	ATDMax    float64 `json:"atdMax" csv:"atd_max"`
	MaxAmount float64 `json:"maxAmount" csv:"plafond"`
}

// Dataset regroupe le jeu généré. Borrowers est trié par ID (déterminisme).
type Dataset struct {
	Borrowers []Borrower
	Policies  []Policy
	byID      map[string]Borrower
}

// Borrower retourne l'emprunteur d'ID donné (ok=false si absent).
func (d *Dataset) Borrower(id string) (Borrower, bool) {
	b, ok := d.byID[id]
	return b, ok
}

// Generate produit un Dataset déterministe : nBorrowers emprunteurs et 50
// politiques (une par (segment × 12 tranches d'âge, complétées à 50)).
func Generate(seed int64, nBorrowers int) *Dataset {
	r := rand.New(rand.NewSource(seed)) //nolint:gosec // G404 : RNG seedé volontaire pour le déterminisme (NFR-10) ; aucun usage sécuritaire
	firsts := []string{"Alix", "Bao", "Chen", "Dara", "Elif", "Farah", "Gita", "Hugo", "Ines", "Jules"}
	lasts := []string{"Tremblay", "Nguyen", "Singh", "Cote", "Diallo", "Rossi", "Kaur", "Lopez", "Wong", "Roy"}

	ds := &Dataset{byID: make(map[string]Borrower, nBorrowers)}
	for i := 1; i <= nBorrowers; i++ {
		income := 25000 + r.Float64()*175000 // 25k–200k
		b := Borrower{
			ID:             fmt.Sprintf("A%05d", i),
			FirstName:      firsts[r.Intn(len(firsts))],
			LastName:       lasts[r.Intn(len(lasts))],
			SIN:            fmt.Sprintf("SYN-%09d", r.Intn(1000000000)), // pas un NAS : jeton synthétique
			IncomeAnnual:   round2(income),
			MonthlyDebts:   round2(r.Float64() * (income / 12) * 0.25),
			MonthlyHousing: round2(r.Float64() * (income / 12) * 0.30),
			Score:          300 + r.Intn(601), // 300–900
			Defaults:       r.Intn(4),         // 0–3
			Segment:        segments[r.Intn(len(segments))],
			Age:            18 + r.Intn(63), // 18–80
		}
		ds.Borrowers = append(ds.Borrowers, b)
		ds.byID[b.ID] = b
	}
	sort.Slice(ds.Borrowers, func(i, j int) bool { return ds.Borrowers[i].ID < ds.Borrowers[j].ID })
	ds.Policies = generatePolicies()
	return ds
}

// generatePolicies produit 50 lignes d'octroi déterministes (segments × tranches).
func generatePolicies() []Policy {
	var pols []Policy
	ages := [][2]int{{18, 24}, {25, 34}, {35, 44}, {45, 54}, {55, 64}, {65, 80}}
	for _, seg := range segments {
		for _, a := range ages {
			pols = append(pols, Policy{
				Segment: seg, AgeMin: a[0], AgeMax: a[1],
				ABDMax: 0.35, ATDMax: 0.43, MaxAmount: 15000 + float64(a[0])*500,
			})
		}
	}
	// Complète à 50 lignes par des variantes de plafond (déterministe).
	for len(pols) < 50 {
		base := pols[len(pols)%(len(segments)*len(ages))]
		base.MaxAmount += 1000
		pols = append(pols, base)
	}
	return pols[:50]
}

// PolicyFor retourne la première politique couvrant (segment, age).
func (d *Dataset) PolicyFor(segment string, age int) (Policy, bool) {
	for _, p := range d.Policies {
		if p.Segment == segment && age >= p.AgeMin && age <= p.AgeMax {
			return p, true
		}
	}
	return Policy{}, false
}

func round2(x float64) float64 { return math.Round(x*100) / 100 }
