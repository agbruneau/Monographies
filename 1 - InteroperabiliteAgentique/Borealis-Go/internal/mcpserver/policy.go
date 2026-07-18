package mcpserver

import (
	"context"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// PolicyIn est l'entrée de evaluate_policy (PRD §6.2).
type PolicyIn struct {
	Segment  string  `json:"segment" jsonschema:"segment de clientele"`
	Age      int     `json:"age" jsonschema:"age du demandeur"`
	ABDRatio float64 `json:"abdRatio" jsonschema:"ratio d'amortissement brut (housing)"`
	ATDRatio float64 `json:"atdRatio" jsonschema:"ratio d'endettement total"`
	Amount   float64 `json:"amount" jsonschema:"montant demande"`
}

// PolicyOut est la sortie de evaluate_policy.
type PolicyOut struct {
	Eligible  bool     `json:"eligible"`
	MaxAmount float64  `json:"maxAmount"`
	Reasons   []string `json:"reasons"`
}

// evaluatePolicy applique une ligne d'octroi à une demande (logique pure).
func evaluatePolicy(p fixtures.Policy, in PolicyIn) PolicyOut {
	out := PolicyOut{Eligible: true, MaxAmount: p.MaxAmount}
	if in.ABDRatio > p.ABDMax {
		out.Eligible = false
		out.Reasons = append(out.Reasons, "abd_above_max")
	}
	if in.ATDRatio > p.ATDMax {
		out.Eligible = false
		out.Reasons = append(out.Reasons, "atd_above_max")
	}
	if in.Amount > p.MaxAmount {
		out.Eligible = false
		out.Reasons = append(out.Reasons, "amount_above_ceiling")
	}
	if out.Eligible {
		out.Reasons = append(out.Reasons, "meets_policy")
	}
	return out
}

// NewPolicyServer construit le serveur MCP Policy Engine sur le jeu fourni.
func NewPolicyServer(ds *fixtures.Dataset) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "policy-engine", Version: "v1.0.0"}, nil)
	mcp.AddTool(s, &mcp.Tool{
		Name:        "evaluate_policy",
		Description: "Évalue l'éligibilité d'une demande selon la table d'octroi.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in PolicyIn) (*mcp.CallToolResult, PolicyOut, error) {
		p, ok := ds.PolicyFor(in.Segment, in.Age)
		if !ok {
			return toolError("NotFound: no policy for segment=" + in.Segment), PolicyOut{}, nil
		}
		return nil, evaluatePolicy(p, in), nil
	})
	return s
}
