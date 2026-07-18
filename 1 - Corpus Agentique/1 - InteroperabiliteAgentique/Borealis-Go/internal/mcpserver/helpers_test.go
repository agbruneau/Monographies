package mcpserver

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/agbruneau/borealis/internal/testfakes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// connectServer connecte un client à un serveur MCP via transport mémoire.
func connectServer(t *testing.T, s *mcp.Server) *mcp.ClientSession {
	return testfakes.ConnectInMemory(t, s)
}

// callStructured appelle un outil et, si le résultat n'est pas en erreur,
// désérialise la sortie structurée dans out (out peut être nil).
func callStructured(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any, out any) *mcp.CallToolResult {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("CallTool %s : %v", name, err)
	}
	if !res.IsError && out != nil {
		b, _ := json.Marshal(res.StructuredContent)
		if err := json.Unmarshal(b, out); err != nil {
			t.Fatalf("unmarshal %s : %v", name, err)
		}
	}
	return res
}
