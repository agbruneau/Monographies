package main

import (
	"os"
	"testing"
)

// TestRunVersionFlag exécute run() (T5 : jamais exercé jusqu'ici) via la
// branche --version — retour immédiat, sans écoute réseau ni lecture stdio.
func TestRunVersionFlag(t *testing.T) {
	old := os.Args
	os.Args = []string{"mcp-identity", "--version"}
	defer func() { os.Args = old }()
	if code := run(); code != 0 {
		t.Errorf("run() --version = %d, want 0", code)
	}
}
