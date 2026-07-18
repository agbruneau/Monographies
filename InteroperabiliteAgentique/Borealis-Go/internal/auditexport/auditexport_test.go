package auditexport

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/agbruneau/borealis/internal/audit"
)

// TestRunProducesReplayableJSONL : l'export d'une pré-qualification canonique est
// un JSONL non vide, relisible et intègre (FR-23, M5.2b).
func TestRunProducesReplayableJSONL(t *testing.T) {
	var buf bytes.Buffer
	if err := Run(context.Background(), &buf); err != nil {
		t.Fatalf("audit-export : %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("export vide")
	}
	var entries []audit.Entry
	dec := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	for dec.More() {
		var e audit.Entry
		if err := dec.Decode(&e); err != nil {
			t.Fatalf("ligne JSONL invalide : %v", err)
		}
		entries = append(entries, e)
	}
	if len(entries) < 3 {
		t.Errorf("entrées = %d, want >= 3 (verify_identity, screen_watchlists, …, decision)", len(entries))
	}
	if err := audit.VerifyEntries(entries); err != nil {
		t.Errorf("chaîne exportée non intègre : %v", err)
	}
	if last := entries[len(entries)-1]; last.Action != "decision" {
		t.Errorf("dernière entrée = %q, want decision (piste close)", last.Action)
	}
}
