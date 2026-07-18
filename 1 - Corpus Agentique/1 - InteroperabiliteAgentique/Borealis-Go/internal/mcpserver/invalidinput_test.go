package mcpserver

import (
	"testing"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestAllToolsRejectInvalidInput couvre le critère de sortie M1 « CHAQUE outil
// validé par test de contrat » : une entrée de type invalide est rejetée
// (IsError, ADR-0001) pour les 5 outils exposés.
func TestAllToolsRejectInvalidInput(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 10)
	cases := []struct {
		server *mcp.Server
		tool   string
		args   map[string]any
	}{
		{NewIdentityServer(ds), "verify_identity", map[string]any{"applicantId": 123}},
		{NewIdentityServer(ds), "check_watchlists", map[string]any{"applicantId": 123}},
		{NewBureauServer(ds), "get_credit_report", map[string]any{"applicantId": 123}},
		{NewCapacityServer(), "compute_capacity", map[string]any{"income": "x", "debts": 0, "requestedAmount": 1000, "termMonths": 12}},
		{NewPolicyServer(ds), "evaluate_policy", map[string]any{"segment": 123, "age": 30, "abdRatio": 0.1, "atdRatio": 0.1, "amount": 1000}},
	}
	for _, c := range cases {
		cs := connectServer(t, c.server)
		if res := callStructured(t, cs, c.tool, c.args, nil); !res.IsError {
			t.Errorf("%s : entrée de type invalide acceptée (IsError attendu)", c.tool)
		}
	}
}
