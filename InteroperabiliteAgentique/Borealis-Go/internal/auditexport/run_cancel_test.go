package auditexport

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

// TestRunPropagatesCancelledContext : un contexte déjà annulé fait échouer
// l'appel d'outil MCP sous-jacent (Prequalify) — Run propage l'erreur au lieu
// d'écrire un export partiel (T1, branches d'erreur de Run non exercées par
// le chemin nominal).
func TestRunPropagatesCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var buf bytes.Buffer
	if err := Run(ctx, &buf); !errors.Is(err, context.Canceled) {
		t.Errorf("Run(ctx annulé) = %v, want context.Canceled", err)
	}
	if buf.Len() != 0 {
		t.Errorf("export partiel écrit malgré l'erreur : %d octets", buf.Len())
	}
}
