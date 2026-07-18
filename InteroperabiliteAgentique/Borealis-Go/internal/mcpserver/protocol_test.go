package mcpserver

import (
	"context"
	"testing"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestProtocolVsBusinessError distingue une erreur PROTOCOLE (outil inconnu →
// erreur JSON-RPC remontée comme erreur Go) d'une erreur MÉTIER (id absent →
// résultat IsError, sans erreur protocole). FR-15 / ADR-0001.
func TestProtocolVsBusinessError(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 10)
	cs := connectServer(t, NewBureauServer(ds))
	ctx := context.Background()

	// Erreur protocole : outil inexistant.
	if _, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "does_not_exist"}); err == nil {
		t.Error("outil inconnu : erreur protocole attendue")
	}

	// Erreur métier : demandeur absent → IsError, sans erreur protocole.
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "get_credit_report", Arguments: map[string]any{"applicantId": "ZZZ"}})
	if err != nil {
		t.Fatalf("erreur métier ne doit pas être une erreur protocole : %v", err)
	}
	if !res.IsError {
		t.Error("id absent : IsError attendu")
	}
}
