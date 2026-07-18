package testfakes

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestConnectInMemory : le client se connecte au serveur et peut appeler un
// outil (couverture directe du paquet, T3 — sinon exercé seulement
// cross-paquet, non compté sans -coverpkg).
func TestConnectInMemory(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "t", Version: "v0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "ping"}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, struct{}, error) {
		return nil, struct{}{}, nil
	})

	cs := ConnectInMemory(t, server)
	if _, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "ping"}); err != nil {
		t.Fatalf("appel d'outil : %v", err)
	}
}
