package main

import (
	"os"
	"testing"
)

// TestRunAuditExport exécute la branche --audit-export de run() (T5 : jamais
// exercée jusqu'ici) — auditexport.Run est autonome (transport MCP en mémoire,
// horloge figée), sans écoute réseau ni signal.
func TestRunAuditExport(t *testing.T) {
	old := os.Args
	os.Args = []string{"orchestrator", "--audit-export"}
	defer func() { os.Args = old }()
	if code := run(); code != 0 {
		t.Errorf("run() --audit-export = %d, want 0", code)
	}
}

// TestRunRequiresIDPSecret exécute la branche par défaut de run() via le repli
// fail-fast sans IDP_SECRET (m9) — rapide et sans écoute réseau, tout en
// prouvant que le nom d'agent câblé (borealis-orchestrator) est valide.
func TestRunRequiresIDPSecret(t *testing.T) {
	old := os.Args
	os.Args = []string{"orchestrator"}
	defer func() { os.Args = old }()
	t.Setenv("IDP_SECRET", "")
	if code := run(); code != 3 {
		t.Errorf("run() sans IDP_SECRET = %d, want 3", code)
	}
}
