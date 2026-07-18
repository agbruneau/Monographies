package mcpserver

import (
	"context"
	"math"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CapacityIn est l'entrée de compute_capacity (PRD §9.2).
type CapacityIn struct {
	Income          float64 `json:"income" jsonschema:"revenu annuel brut"`
	Debts           float64 `json:"debts" jsonschema:"total des autres dettes mensuelles"`
	RequestedAmount float64 `json:"requestedAmount" jsonschema:"montant demande"`
	TermMonths      int     `json:"termMonths" jsonschema:"duree du pret en mois"`
}

// CapacityOut est la sortie de compute_capacity.
type CapacityOut struct {
	MonthlyPayment float64 `json:"monthlyPayment"`
	DTIRatio       float64 `json:"dtiRatio"`
	CapacityOk     bool    `json:"capacityOk"`
}

// computeCapacity calcule la mensualité et le ratio DTI (taux fixe 9 %/an,
// démonstrateur). Les entrées dégénérées (durée, montant, revenu <= 0, ou dette
// < 0) sont **hors capacité** : jamais capacityOk sur une demande non-sens (une
// durée nulle ou un montant négatif ne doit pas être neutralisé en pmt=0 puis
// déclaré dans les capacités).
func computeCapacity(in CapacityIn) CapacityOut {
	return capacityAt(in, 0.09)
}

// capacityAt calcule mensualité, DTI et capacité au taux annuel donné (facteur
// commun réutilisé par le contrat v2, qui applique un choc de taux — M4.5).
func capacityAt(in CapacityIn, annualRate float64) CapacityOut {
	if in.TermMonths <= 0 || in.RequestedAmount <= 0 || in.Income <= 0 || in.Debts < 0 {
		return CapacityOut{MonthlyPayment: 0, DTIRatio: 0, CapacityOk: false}
	}
	r := annualRate / 12.0
	n := float64(in.TermMonths)
	pmt := in.RequestedAmount * r / (1 - math.Pow(1+r, -n))
	dti := round2((pmt + in.Debts) / (in.Income / 12.0))
	return CapacityOut{
		MonthlyPayment: round2(pmt),
		DTIRatio:       dti,
		CapacityOk:     dti <= 0.40, // décision sur la valeur PUBLIÉE (arrondie) : cohérente avec dtiRatio
	}
}

// NewCapacityServer construit le serveur MCP Capacity Calculator.
func NewCapacityServer() *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "capacity-calculator", Version: "v1.0.0"}, nil)
	mcp.AddTool(s, &mcp.Tool{
		Name:        "compute_capacity",
		Description: "Calcule le versement mensuel et le ratio DTI d'une demande.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in CapacityIn) (*mcp.CallToolResult, CapacityOut, error) {
		return nil, computeCapacity(in), nil
	})
	return s
}

func round2(x float64) float64 { return math.Round(x*100) / 100 }
