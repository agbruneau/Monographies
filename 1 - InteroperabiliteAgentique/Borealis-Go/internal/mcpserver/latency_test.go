package mcpserver

import (
	"context"
	"testing"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestToolRoundtripSmoke vérifie que chaque outil répond via un aller-retour
// réel (transport mémoire). Le SLO mesuré NFR-02 (P95 < 500 ms) n'est pas
// falsifiable ici : un aller-retour en mémoire se compte en dizaines de µs,
// plusieurs ordres de grandeur sous le seuil — la preuve mesurée revient à la
// composition déployée (T4 ; l'ancienne version bouclait 50x par outil pour un
// résultat jamais autrement que vert).
func TestToolRoundtripSmoke(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	calls := []struct {
		server *mcp.Server
		tool   string
		args   map[string]any
	}{
		{NewIdentityServer(ds), "verify_identity", map[string]any{"applicantId": "A00001"}},
		{NewIdentityServer(ds), "check_watchlists", map[string]any{"applicantId": "A00001", "name": "Alix"}},
		{NewBureauServer(ds), "get_credit_report", map[string]any{"applicantId": "A00001"}},
		{NewCapacityServer(), "compute_capacity", map[string]any{"income": 120000.0, "debts": 500.0, "requestedAmount": 20000.0, "termMonths": 60}},
		{NewPolicyServer(ds), "evaluate_policy", map[string]any{"segment": "retail", "age": 30, "abdRatio": 0.3, "atdRatio": 0.4, "amount": 10000.0}},
	}
	for _, c := range calls {
		cs := connectServer(t, c.server)
		if _, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: c.tool, Arguments: c.args}); err != nil {
			t.Errorf("%s : %v", c.tool, err)
		}
	}
}
