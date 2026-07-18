// Package creditdomain porte le modèle métier pur de la pré-qualification de
// crédit : montants, ratios (PMT, DTI) et règle d'évaluation. Aucune E/S,
// aucune dépendance externe — logique déterministe.
//
// Invariant (FR-25, invariant 1 du PRD) : jamais d'octroi ferme. Les seules
// issues sont pré-qualifié, refusé, ou renvoi à un humain (HITL).
package creditdomain

import (
	"errors"
	"math"
)

// Money représente un montant en cents (entier) pour éviter toute dérive de
// virgule flottante sur les montants.
type Money int64

// Outcome est l'issue d'une évaluation. Aucune valeur « octroyé » n'existe.
type Outcome int

const (
	// OutcomeDeclined : demande refusée.
	OutcomeDeclined Outcome = iota
	// OutcomeRefer : renvoi à un humain (HITL).
	OutcomeRefer
	// OutcomePreQualified : pré-qualifié (jamais un octroi ferme).
	OutcomePreQualified
)

// String rend l'issue sous forme stable (utilisée dans les golden).
func (o Outcome) String() string {
	switch o {
	case OutcomeDeclined:
		return "declined"
	case OutcomeRefer:
		return "refer"
	case OutcomePreQualified:
		return "pre_qualified"
	default:
		return "unknown"
	}
}

// Applicant est une demande de pré-qualification (données 100 % synthétiques).
type Applicant struct {
	ID              string
	MonthlyIncome   Money   // revenu mensuel, cents (> 0)
	MonthlyDebt     Money   // charges de dette mensuelles existantes, cents (>= 0)
	RequestedAmount Money   // principal demandé, cents (> 0)
	TermMonths      int     // durée en mois (> 0)
	AnnualRate      float64 // taux annuel nominal (0.079 = 7,9 %)
	CreditScore     int     // score de crédit (ex. 300–900)
}

// Policy porte les seuils d'évaluation (calibrables — hypothèses à ajuster).
type Policy struct {
	MinScore   int     // en deçà → refus d'emblée
	MaxDTI     float64 // au-delà → refus
	ReferScore int     // score >= MinScore mais < ReferScore → renvoi humain
	ReferDTI   float64 // DTI <= MaxDTI mais > ReferDTI → renvoi humain
}

// Decision est le résultat d'une évaluation.
type Decision struct {
	ApplicantID    string
	Outcome        Outcome
	DTI            float64
	MonthlyPayment Money
	Reasons        []string
}

// Erreurs de frontière de confiance (entrées invalides).
var (
	ErrIncome = errors.New("creditdomain: monthly income must be > 0")
	ErrTerm   = errors.New("creditdomain: term months must be > 0")
	ErrAmount = errors.New("creditdomain: requested amount must be > 0")
	ErrDebt   = errors.New("creditdomain: monthly debt must be >= 0")
	ErrRate   = errors.New("creditdomain: annual rate must be >= 0")
)

// MonthlyPayment calcule la mensualité d'amortissement (arrondie au cent).
// Taux nul → remboursement linéaire. Durée <= 0 → 0.
func MonthlyPayment(principal Money, annualRate float64, termMonths int) Money {
	if termMonths <= 0 {
		return 0
	}
	p := float64(principal)
	if annualRate == 0 {
		return Money(math.Round(p / float64(termMonths)))
	}
	r := annualRate / 12.0
	n := float64(termMonths)
	pmt := p * r / (1 - math.Pow(1+r, -n))
	return Money(math.Round(pmt))
}

// Validate vérifie les frontières de confiance de la demande. Une dette ou un
// taux négatifs sont rejetés : sinon ils sous-évalueraient le DTI (dette < 0)
// ou produiraient un DTI non fini (taux <= -12 %/an → NaN) qui traverserait
// silencieusement les portes de refus.
func (a Applicant) Validate() error {
	switch {
	case a.MonthlyIncome <= 0:
		return ErrIncome
	case a.TermMonths <= 0:
		return ErrTerm
	case a.RequestedAmount <= 0:
		return ErrAmount
	case a.MonthlyDebt < 0:
		return ErrDebt
	case a.AnnualRate < 0:
		return ErrRate
	default:
		return nil
	}
}

// dtiValue calcule le ratio dette/revenu en supposant un revenu > 0 (garanti
// par Validate). Facteur commun, sans branche morte.
func dtiValue(a Applicant) float64 {
	pmt := MonthlyPayment(a.RequestedAmount, a.AnnualRate, a.TermMonths)
	return float64(a.MonthlyDebt+pmt) / float64(a.MonthlyIncome)
}

// Assess évalue une demande selon une politique. Retourne une erreur en cas
// d'entrée invalide (frontière de confiance) ; sinon une Decision métier.
// Ne produit jamais d'octroi ferme (FR-25).
func Assess(a Applicant, p Policy) (Decision, error) {
	if err := a.Validate(); err != nil {
		return Decision{}, err
	}
	dti := dtiValue(a)
	d := Decision{
		ApplicantID:    a.ID,
		DTI:            dti,
		MonthlyPayment: MonthlyPayment(a.RequestedAmount, a.AnnualRate, a.TermMonths),
	}
	switch {
	case a.CreditScore < p.MinScore:
		d.Outcome = OutcomeDeclined
		d.Reasons = []string{"credit_score_below_minimum"}
	case dti > p.MaxDTI:
		d.Outcome = OutcomeDeclined
		d.Reasons = []string{"dti_above_maximum"}
	case a.CreditScore < p.ReferScore || dti > p.ReferDTI:
		d.Outcome = OutcomeRefer
		d.Reasons = []string{"borderline_requires_human_review"}
	default:
		d.Outcome = OutcomePreQualified
		d.Reasons = []string{"meets_prequalification_criteria"}
	}
	return d, nil
}
