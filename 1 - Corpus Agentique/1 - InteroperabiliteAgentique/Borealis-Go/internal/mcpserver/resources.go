package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	assessmentURIPrefix = "credit:///application/"
	assessmentURISuffix = "/assessment"
)

// applicantIDFromURI extrait l'ID de credit:///application/{id}/assessment.
func applicantIDFromURI(uri string) (string, bool) {
	if !strings.HasPrefix(uri, assessmentURIPrefix) || !strings.HasSuffix(uri, assessmentURISuffix) {
		return "", false
	}
	id := strings.TrimSuffix(strings.TrimPrefix(uri, assessmentURIPrefix), assessmentURISuffix)
	if id == "" {
		return "", false
	}
	return id, true
}

// addAssessmentResource enregistre la ressource d'évaluation (FR-12).
func addAssessmentResource(s *mcp.Server, ds *fixtures.Dataset) {
	s.AddResourceTemplate(&mcp.ResourceTemplate{
		Name:        "credit-assessment",
		URITemplate: assessmentURIPrefix + "{id}" + assessmentURISuffix,
		MIMEType:    "application/json",
		Description: "Évaluation de crédit synthétique d'un demandeur.",
	}, func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		uri := req.Params.URI
		id, ok := applicantIDFromURI(uri)
		if !ok {
			return nil, mcp.ResourceNotFoundError(uri)
		}
		b, ok := ds.Borrower(id)
		if !ok {
			return nil, mcp.ResourceNotFoundError(uri)
		}
		body, err := json.Marshal(creditReport(b))
		if err != nil {
			return nil, err
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{URI: uri, MIMEType: "application/json", Text: string(body)}},
		}, nil
	})
}

// addAdvicePrompt enregistre le gabarit d'avis de pré-qualification (FR-13).
func addAdvicePrompt(s *mcp.Server) {
	s.AddPrompt(&mcp.Prompt{
		Name:        "credit_assessment_advice",
		Description: "Gabarit d'avis de pré-qualification (fr/en).",
		Arguments: []*mcp.PromptArgument{
			{Name: "applicantId", Description: "identifiant du demandeur", Required: true},
			{Name: "lang", Description: "langue de l'avis (fr|en)"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		id := req.Params.Arguments["applicantId"]
		var text string
		if req.Params.Arguments["lang"] == "en" {
			text = fmt.Sprintf("Draft a concise pre-qualification advice for applicant %s. State clearly it is a pre-qualification, never a firm grant.", id)
		} else {
			text = fmt.Sprintf("Rédige un avis de pré-qualification concis pour le demandeur %s. Précise qu'il s'agit d'une pré-qualification, jamais d'un octroi ferme.", id)
		}
		return &mcp.GetPromptResult{
			Description: "Avis de pré-qualification",
			Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
		}, nil
	})
}
