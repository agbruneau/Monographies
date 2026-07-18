package pep

import (
	"errors"
	"sync"
	"testing"

	"github.com/agbruneau/borealis/internal/audit"
)

func TestGuardAllowsAtL1(t *testing.T) {
	j := audit.New()
	p := New(j, L1)
	if err := p.Guard("agent:kyc", "cust:h", "mcp.CallTool:verify_identity"); err != nil {
		t.Fatalf("L1 devrait autoriser : %v", err)
	}
	if j.Len() != 1 {
		t.Errorf("décision non journalisée (Len=%d)", j.Len())
	}
}

func TestGuardDeniesAtL0(t *testing.T) {
	j := audit.New()
	p := New(j, L0)
	if err := p.Guard("agent:kyc", "cust:h", "x"); !errors.Is(err, ErrDenied) {
		t.Errorf("L0 devrait refuser (humain requis) : %v", err)
	}
	if j.Len() != 1 {
		t.Error("refus non journalisé")
	}
}

func TestGuardFailClosed(t *testing.T) {
	j := audit.New()
	p := New(j, L1)
	p.SetAvailable(false) // PEP/orchestrateur indisponible
	if err := p.Guard("agent:kyc", "cust:h", "x"); !errors.Is(err, ErrFailClosed) {
		t.Errorf("indisponible devrait retomber fail-closed : %v", err)
	}
	// La retombée L1->L0 est journalisée.
	es := j.Entries()
	if len(es) != 1 || es[0].Action != "pep.fail_closed" {
		t.Errorf("retombée fail-closed non journalisée : %+v", es)
	}
}

func TestGuardNilJournal(t *testing.T) {
	// Sans journal, la garde fonctionne sans paniquer.
	if err := New(nil, L1).Guard("a", "c", "x"); err != nil {
		t.Fatal(err)
	}
}

// TestSetAvailableConcurrent : SetAvailable et Guard concurrents ne doivent
// jamais paniquer ni être signalés par le détecteur de courses (P2 ; available
// est un atomic.Bool). Preuve complète : `go test -race`, indisponible sur cet
// hôte (CGO absent).
func TestSetAvailableConcurrent(t *testing.T) {
	p := New(nil, L1)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(v bool) { defer wg.Done(); p.SetAvailable(v) }(i%2 == 0)
		go func() { defer wg.Done(); _ = p.Guard("a", "c", "x") }()
	}
	wg.Wait()
}
