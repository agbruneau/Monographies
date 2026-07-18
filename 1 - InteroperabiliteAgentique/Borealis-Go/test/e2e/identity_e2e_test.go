//go:build e2e

// Package e2e porte les scénarios de bout en bout (tag `e2e`, hors gate par
// défaut ; lancés par `make e2e`). Ils exercent les **binaires réels** et les
// transports de déploiement (stdio), pas seulement les transports en mémoire.
package e2e

import (
	"context"
	"encoding/json"
	"os/exec"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestIdentityStdioRoundtrip lance le binaire réel cmd/mcp-identity et effectue
// un aller-retour verify_identity sur le transport **stdio** — démontre FR-11
// (« appel MCP réel de bout en bout ») sur la voie de déploiement effective,
// au-delà des transports en mémoire des tests unitaires.
func TestIdentityStdioRoundtrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.Command("go", "run", "github.com/agbruneau/borealis/cmd/mcp-identity")
	client := mcp.NewClient(&mcp.Implementation{Name: "e2e", Version: "v0"}, nil)
	cs, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connexion au binaire stdio : %v", err)
	}
	defer func() { _ = cs.Close() }()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      "verify_identity",
		Arguments: map[string]any{"applicantId": "A00001"}, // présent dans le jeu seedé
	})
	if err != nil {
		t.Fatalf("CallTool : %v", err)
	}
	if res.IsError {
		t.Fatalf("résultat en erreur : %+v", res.Content)
	}
	var out struct {
		Match bool `json:"match"`
	}
	b, _ := json.Marshal(res.StructuredContent)
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal : %v", err)
	}
	if !out.Match {
		t.Errorf("sortie = %+v, want match=true", out)
	}
}
