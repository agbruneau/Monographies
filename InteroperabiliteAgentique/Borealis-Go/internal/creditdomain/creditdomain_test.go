package creditdomain

import (
	"errors"
	"testing"
)

func TestOutcomeString(t *testing.T) {
	cases := map[Outcome]string{
		OutcomeDeclined:     "declined",
		OutcomeRefer:        "refer",
		OutcomePreQualified: "pre_qualified",
		Outcome(99):         "unknown",
	}
	for o, want := range cases {
		if got := o.String(); got != want {
			t.Errorf("Outcome(%d).String() = %q, want %q", o, got, want)
		}
	}
}

func TestMonthlyPayment(t *testing.T) {
	tests := []struct {
		name       string
		principal  Money
		annualRate float64
		term       int
		want       Money
	}{
		{"taux nul, linéaire", 120000, 0, 12, 10000},
		{"durée nulle → 0", 120000, 0.05, 0, 0},
		{"durée négative → 0", 120000, 0.05, -1, 0},
		// 10 000,00 $ à 12 % sur 12 mois → 888,49 $ (88849 cents).
		{"taux positif amortissement", 1000000, 0.12, 12, 88849},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MonthlyPayment(tt.principal, tt.annualRate, tt.term); got != tt.want {
				t.Errorf("MonthlyPayment(%d,%v,%d) = %d, want %d", tt.principal, tt.annualRate, tt.term, got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	base := Applicant{MonthlyIncome: 500000, RequestedAmount: 1000000, TermMonths: 60}
	if err := base.Validate(); err != nil {
		t.Fatalf("demande valide rejetée : %v", err)
	}
	tests := []struct {
		name string
		mut  func(*Applicant)
		want error
	}{
		{"revenu nul", func(a *Applicant) { a.MonthlyIncome = 0 }, ErrIncome},
		{"durée nulle", func(a *Applicant) { a.TermMonths = 0 }, ErrTerm},
		{"montant nul", func(a *Applicant) { a.RequestedAmount = 0 }, ErrAmount},
		{"dette négative", func(a *Applicant) { a.MonthlyDebt = -1 }, ErrDebt},
		{"taux négatif", func(a *Applicant) { a.AnnualRate = -0.01 }, ErrRate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := base
			tt.mut(&a)
			if err := a.Validate(); !errors.Is(err, tt.want) {
				t.Errorf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestAssess(t *testing.T) {
	pol := Policy{MinScore: 600, MaxDTI: 0.43, ReferScore: 660, ReferDTI: 0.36}
	// Base : bon dossier → pré-qualifié.
	good := Applicant{ID: "a1", MonthlyIncome: 1000000, MonthlyDebt: 50000, RequestedAmount: 1000000, AnnualRate: 0.079, TermMonths: 60, CreditScore: 720}

	tests := []struct {
		name    string
		mut     func(*Applicant)
		want    Outcome
		wantErr bool
	}{
		{"pré-qualifié", nil, OutcomePreQualified, false},
		{"refus score bas", func(a *Applicant) { a.CreditScore = 550 }, OutcomeDeclined, false},
		{"refus DTI haut", func(a *Applicant) { a.MonthlyDebt = 450000 }, OutcomeDeclined, false},
		{"renvoi humain score limite", func(a *Applicant) { a.CreditScore = 640 }, OutcomeRefer, false},
		{"renvoi humain DTI limite", func(a *Applicant) { a.MonthlyDebt = 360000 }, OutcomeRefer, false},
		{"entrée invalide", func(a *Applicant) { a.MonthlyIncome = 0 }, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := good
			if tt.mut != nil {
				tt.mut(&a)
			}
			d, err := Assess(a, pol)
			if tt.wantErr {
				if err == nil {
					t.Fatal("erreur attendue, obtenu nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("erreur inattendue : %v", err)
			}
			if d.Outcome != tt.want {
				t.Errorf("Outcome = %s, want %s (DTI=%.3f, reasons=%v)", d.Outcome, tt.want, d.DTI, d.Reasons)
			}
			if len(d.Reasons) == 0 {
				t.Error("Decision sans motif")
			}
		})
	}
}

// TestPreQualificationCeiling : le meilleur dossier possible plafonne à
// « pré-qualifié » (jamais un octroi ferme, FR-25) ; un mauvais dossier n'est
// jamais pré-qualifié (les portes de refus ne sont pas contournables). Oracle
// portant sur le comportement, non sur l'énumération des valeurs possibles.
func TestPreQualificationCeiling(t *testing.T) {
	pol := Policy{MinScore: 600, MaxDTI: 0.43, ReferScore: 660, ReferDTI: 0.36}

	// Dossier maximal : le plafond atteignable est la pré-qualification.
	best := Applicant{ID: "best", MonthlyIncome: 5000000, MonthlyDebt: 0, RequestedAmount: 100000, AnnualRate: 0.05, TermMonths: 120, CreditScore: 900}
	d, err := Assess(best, pol)
	if err != nil {
		t.Fatal(err)
	}
	if d.Outcome != OutcomePreQualified {
		t.Errorf("meilleur dossier → %s, attendu pre_qualified (plafond)", d.Outcome)
	}

	// Dossiers clairement mauvais : jamais pré-qualifiés.
	bad := []Applicant{
		{ID: "lowscore", MonthlyIncome: 1000000, RequestedAmount: 100000, AnnualRate: 0.05, TermMonths: 60, CreditScore: 400},
		{ID: "overleveraged", MonthlyIncome: 300000, MonthlyDebt: 250000, RequestedAmount: 2000000, AnnualRate: 0.12, TermMonths: 24, CreditScore: 800},
		// Anciennes attaques (dette/taux négatifs) désormais rejetées à la validation.
		{ID: "negdebt", MonthlyIncome: 500000, MonthlyDebt: -500000, RequestedAmount: 1000000, AnnualRate: 0.12, TermMonths: 24, CreditScore: 800},
		{ID: "negrate", MonthlyIncome: 500000, RequestedAmount: 1000000, AnnualRate: -12.0, TermMonths: 60, CreditScore: 800},
	}
	for _, a := range bad {
		d, err := Assess(a, pol)
		if err != nil {
			continue // rejeté à la validation : acceptable (jamais pré-qualifié)
		}
		if d.Outcome == OutcomePreQualified {
			t.Errorf("dossier %q pré-qualifié à tort (DTI=%.3f) — porte de refus contournée", a.ID, d.DTI)
		}
	}
}
