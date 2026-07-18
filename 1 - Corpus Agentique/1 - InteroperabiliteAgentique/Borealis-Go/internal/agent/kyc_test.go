package agent

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestParseIdentityResult(t *testing.T) {
	// Cas valide.
	ok, err := parseIdentityResult(&mcp.CallToolResult{
		StructuredContent: map[string]any{"match": true, "reason": "verified"},
	})
	if err != nil || !ok.Match || ok.Reason != "verified" {
		t.Fatalf("cas valide : %+v err=%v", ok, err)
	}

	// IsError avec texte → l'erreur reprend le texte.
	if _, err := parseIdentityResult(&mcp.CallToolResult{
		IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: "bad input"}},
	}); err == nil {
		t.Error("IsError avec texte : erreur attendue")
	}

	// IsError sans contenu texte → repli "unknown tool error".
	if _, err := parseIdentityResult(&mcp.CallToolResult{IsError: true}); err == nil {
		t.Error("IsError sans texte : erreur attendue")
	}

	// Erreur de sérialisation (canal non marshalable).
	if _, err := parseIdentityResult(&mcp.CallToolResult{StructuredContent: make(chan int)}); err == nil {
		t.Error("marshal : erreur attendue")
	}

	// Erreur de désérialisation (match de type chaîne).
	if _, err := parseIdentityResult(&mcp.CallToolResult{
		StructuredContent: map[string]any{"match": "notabool"},
	}); err == nil {
		t.Error("unmarshal : erreur attendue")
	}
}
