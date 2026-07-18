package mcpserver

import (
	"context"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BureauIn est l'entrée de get_credit_report (PRD §6.2).
type BureauIn struct {
	ApplicantID string `json:"applicantId" jsonschema:"identifiant du demandeur"`
}

// BureauOut est la sortie de get_credit_report.
type BureauOut struct {
	Score    int     `json:"score"`
	ABDRatio float64 `json:"abdRatio"`
	ATDRatio float64 `json:"atdRatio"`
	Defaults int     `json:"defaults"`
}

// ratioNonComputable est un ratio sentinelle (revenu non calculable) : élevé,
// afin que la politique rejette plutôt que d'interpréter 0 comme un profil sain.
const ratioNonComputable = 999.0

// creditReport dérive le rapport de crédit d'un emprunteur (logique pure). Un
// revenu <= 0 rend les ratios non calculables (sentinelle élevée) — jamais 0,
// qui ferait passer un dossier à revenu nul pour le meilleur profil de dette.
func creditReport(b fixtures.Borrower) BureauOut {
	monthlyIncome := b.IncomeAnnual / 12.0
	if monthlyIncome <= 0 {
		return BureauOut{Score: b.Score, ABDRatio: ratioNonComputable, ATDRatio: ratioNonComputable, Defaults: b.Defaults}
	}
	abd := round2(b.MonthlyHousing / monthlyIncome)
	atd := round2((b.MonthlyHousing + b.MonthlyDebts) / monthlyIncome)
	return BureauOut{Score: b.Score, ABDRatio: abd, ATDRatio: atd, Defaults: b.Defaults}
}

// NewBureauServer construit le serveur MCP Credit Bureau Sim sur le jeu fourni.
func NewBureauServer(ds *fixtures.Dataset) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "credit-bureau-sim", Version: "v1.0.0"}, nil)
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_credit_report",
		Description: "Retourne le rapport de crédit synthétique d'un demandeur.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in BureauIn) (*mcp.CallToolResult, BureauOut, error) {
		if token := req.Params.GetProgressToken(); token != nil {
			_ = req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
				ProgressToken: token, Message: "looking up credit report", Progress: 1, Total: 1,
			}) // FR-14 : notification de progression si un jeton est demandé
		}
		b, ok := ds.Borrower(in.ApplicantID)
		if !ok {
			return toolError("NotFound: applicant " + in.ApplicantID), BureauOut{}, nil
		}
		return nil, creditReport(b), nil
	})
	addAssessmentResource(s, ds) // FR-12 : ressource d'évaluation
	addAdvicePrompt(s)           // FR-13 : gabarit d'avis
	return s
}
