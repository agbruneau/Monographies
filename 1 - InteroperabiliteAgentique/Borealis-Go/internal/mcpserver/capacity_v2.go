package mcpserver

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CapacityInV2 fait évoluer l'entrée de compute_capacity de façon **BACKWARD**
// (M4.5, G-4) : elle ajoute un champ **optionnel** (omitempty → non requis dans
// le schéma inféré) `stressRate` (choc de taux). v2 accepte donc toute entrée
// v1 valide (stressRate absent = 0 = comportement v1).
type CapacityInV2 struct {
	Income          float64 `json:"income" jsonschema:"revenu annuel brut"`
	Debts           float64 `json:"debts" jsonschema:"total des autres dettes mensuelles"`
	RequestedAmount float64 `json:"requestedAmount" jsonschema:"montant demande"`
	TermMonths      int     `json:"termMonths" jsonschema:"duree du pret en mois"`
	StressRate      float64 `json:"stressRate,omitempty" jsonschema:"choc de taux additionnel (0 = v1)"`
}

// CapacityOutV2 enrichit la sortie d'un champ **ajouté** `riskBand` : un
// consommateur v1 lit ses champs (monthlyPayment/dtiRatio/capacityOk) et ignore
// simplement le nouveau — la sortie v2 reste lisible en v1 (BACKWARD).
type CapacityOutV2 struct {
	MonthlyPayment float64 `json:"monthlyPayment"`
	DTIRatio       float64 `json:"dtiRatio"`
	CapacityOk     bool    `json:"capacityOk"`
	RiskBand       string  `json:"riskBand"` // nouveau en v2 : "low" | "medium" | "high"
}

// errStressRateFloor signale un stressRate qui annule ou inverse le taux
// choqué (0,09 + stressRate <= 0) : au taux net nul, la formule de mensualité
// divise par (1-1) = 0 → NaN, non sérialisable en JSON (M1). Un taux net
// négatif serait mathématiquement absurde (mensualité négative).
var errStressRateFloor = errors.New("stressRate: le taux choqué (0.09 + stressRate) doit être strictement positif")

// computeCapacityV2 réutilise la formule v1 au taux choqué (9 % + stressRate) et
// ajoute la bande de risque dérivée du DTI. À stressRate = 0 (toute entrée v1),
// les champs v1 sont **identiques** à ceux de computeCapacity : la substitution
// ne change pas la décision d'un consommateur v1.
func computeCapacityV2(in CapacityInV2) (CapacityOutV2, error) {
	rate := 0.09 + in.StressRate
	if rate <= 0 {
		return CapacityOutV2{}, errStressRateFloor
	}
	base := capacityAt(CapacityIn{
		Income:          in.Income,
		Debts:           in.Debts,
		RequestedAmount: in.RequestedAmount,
		TermMonths:      in.TermMonths,
	}, rate)
	return CapacityOutV2{
		MonthlyPayment: base.MonthlyPayment,
		DTIRatio:       base.DTIRatio,
		CapacityOk:     base.CapacityOk,
		RiskBand:       riskBand(base),
	}, nil
}

// riskBand classe une demande en bande de risque (illustratif, démonstrateur).
// La borne « high » dérive de **capacityOk** (hors capacité : dégénérée OU DTI
// trop élevé), non de `dti <= 0` : une demande valide au DTI arrondi à 0,00 est
// « low », pas « high » — sinon un consommateur v2 déclinerait un dossier que le
// champ v1 `capacityOk` approuve (incohérence de contrat).
func riskBand(c CapacityOut) string {
	switch {
	case !c.CapacityOk:
		return "high" // hors capacité (dégénérée ou DTI > 0,40)
	case c.DTIRatio < 0.30:
		return "low"
	default:
		return "medium" // capacityOk ⇒ DTI ≤ 0,40
	}
}

// NewCapacityServerV2 construit la v2 du Capacity Calculator (compute_capacity),
// substituable à la v1 sans modifier le consommateur (M4.5, US-9).
func NewCapacityServerV2() *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "capacity-calculator", Version: "v2.0.0"}, nil)
	mcp.AddTool(s, &mcp.Tool{
		Name:        "compute_capacity",
		Description: "Calcule le versement mensuel, le ratio DTI et la bande de risque d'une demande.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in CapacityInV2) (*mcp.CallToolResult, CapacityOutV2, error) {
		out, err := computeCapacityV2(in)
		if err != nil {
			return toolError(err.Error()), CapacityOutV2{}, nil
		}
		return nil, out, nil
	})
	return s
}
