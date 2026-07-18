package testfakes

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ConnectInMemory connecte un client à un serveur MCP via transport en
// mémoire, et ferme la session client à la fin du test. Factorise le
// dialogue de connexion dupliqué dans plusieurs paquets de test (T3).
func ConnectInMemory(t testing.TB, s *mcp.Server) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "v0"}, nil)
	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := s.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("connexion serveur : %v", err)
	}
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("connexion client : %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}
