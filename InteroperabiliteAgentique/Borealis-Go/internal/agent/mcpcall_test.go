package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/pep"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// TestResilientCallerPEPGuard : le PEP branché bloque l'appel d'outil quand il
// retombe fail-closed (garde non contournable, M3.2).
func TestResilientCallerPEPGuard(t *testing.T) {
	sess := newIdentitySession(t)
	p := pep.New(audit.New(), pep.L1)
	caller := NewResilientCaller(sess).WithPEP(p, "agent:kyc")

	if _, err := caller.CallTool(context.Background(), "verify_identity", map[string]any{"applicantId": "A00001"}); err != nil {
		t.Fatalf("PEP disponible (L1) : l'appel devrait passer : %v", err)
	}
	p.SetAvailable(false) // PEP/orchestrateur indisponible
	if _, err := caller.CallTool(context.Background(), "verify_identity", map[string]any{"applicantId": "A00001"}); err == nil {
		t.Error("PEP indisponible : l'appel d'outil devrait être bloqué (fail-closed)")
	}
}

// TestMCPTraceparent : le caller propage le traceparent W3C via le champ _meta
// d'un appel MCP (trace ininterrompue A2A→MCP, NFR-07).
func TestMCPTraceparent(t *testing.T) {
	var got string
	type in struct{}
	type out struct {
		TP string `json:"tp"`
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "tp", Version: "v0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "echo_tp"}, func(_ context.Context, req *mcp.CallToolRequest, _ in) (*mcp.CallToolResult, out, error) {
		if req.Params != nil {
			if v, ok := req.Params.Meta["traceparent"].(string); ok {
				got = v
			}
		}
		return nil, out{TP: got}, nil
	})

	ctx0 := context.Background()
	cs := connectInMemory(t, server)

	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(tracetest.NewSpanRecorder()))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })
	ctx, span := tp.Tracer("x").Start(ctx0, "op")
	defer span.End()

	if _, err := NewResilientCaller(cs).CallTool(ctx, "echo_tp", map[string]any{}); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "00-") {
		t.Errorf("traceparent MCP = %q, want format W3C (00-...)", got)
	}
}
