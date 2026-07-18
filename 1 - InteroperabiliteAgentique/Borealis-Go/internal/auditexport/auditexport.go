// Package auditexport compose une pré-qualification canonique (happy path,
// données 100 % synthétiques) et exporte son journal WORM en JSONL rejouable
// (FR-23, M5.2b). Sert la cible `make audit-export`. Horloge figée → sortie
// stable d'une exécution à l'autre (G-8).
package auditexport

import (
	"context"
	"io"
	"time"

	"github.com/agbruneau/borealis/internal/agent"
	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/creditdomain"
	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/agbruneau/borealis/internal/mcpserver"
	"github.com/agbruneau/borealis/internal/testfakes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Run exécute une pré-qualification canonique et écrit son journal d'audit en
// JSONL sur w. La chaîne exportée est intègre et re-vérifiable (audit.VerifyEntries).
func Run(ctx context.Context, w io.Writer) error {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	server := mcpserver.NewIdentityServer(ds)
	client := mcp.NewClient(&mcp.Implementation{Name: "audit-export", Version: "v0"}, nil)
	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		return err
	}
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		return err
	}
	defer func() { _ = cs.Close() }()

	clock := func() time.Time { return time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC) }
	j := audit.New(audit.WithClock(clock))
	o := agent.NewOrchestrator(cs, testfakes.NewStubReasoner(), j, policy())
	applicant := creditdomain.Applicant{
		ID: "A00001", MonthlyIncome: 1000000, MonthlyDebt: 50000,
		RequestedAmount: 1000000, AnnualRate: 0.079, TermMonths: 60, CreditScore: 720,
	}
	if _, err := o.Prequalify(ctx, applicant); err != nil {
		return err
	}
	return j.ExportJSONL(w)
}

func policy() creditdomain.Policy {
	return creditdomain.Policy{MinScore: 600, MaxDTI: 0.43, ReferScore: 660, ReferDTI: 0.36}
}
