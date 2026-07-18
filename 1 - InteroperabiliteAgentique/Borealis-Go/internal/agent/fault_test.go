package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/agbruneau/borealis/internal/resilience"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestFaultInjectionServerStopped injecte une panne (serveur MCP **arrêté**) et
// vérifie la réponse résiliente (FR-17, NFR-04) : chaque appel échoue proprement
// (erreur, jamais de panic), le **disjoncteur s'ouvre** au seuil (fast-fail), et
// la détection reste bien en deçà de 30 s. La borne des tentatives par appel est
// prouvée par les tests unitaires du breaker (TestRetryExhausted) ; le cas d'un
// serveur qui **bloque** est couvert par TestFaultInjectionServerHangs.
func TestFaultInjectionServerStopped(t *testing.T) {
	sess := newIdentitySession(t)
	caller := NewResilientCaller(sess)
	ctx := context.Background()
	args := map[string]any{"applicantId": "A00001"}

	// Serveur vivant : l'appel réussit (le disjoncteur reste fermé).
	if _, err := caller.CallTool(ctx, "verify_identity", args); err != nil {
		t.Fatalf("appel initial devrait réussir : %v", err)
	}

	// Injection de panne : le serveur devient inaccessible.
	if err := sess.Close(); err != nil {
		t.Fatalf("fermeture de la session : %v", err)
	}

	// Détection bornée : chaque appel échoue proprement. Le seuil du disjoncteur
	// (5) atteint, l'appel suivant doit *fast-fail* en ErrOpen sans toucher le
	// serveur mort.
	start := time.Now()
	for i := 0; i < 5; i++ {
		if _, err := caller.CallTool(ctx, "verify_identity", args); err == nil {
			t.Fatalf("appel %d sur serveur arrêté devrait échouer", i+1)
		}
	}
	_, err := caller.CallTool(ctx, "verify_identity", args)
	if !errors.Is(err, resilience.ErrOpen) {
		t.Errorf("après 5 échecs, le disjoncteur devrait être ouvert (ErrOpen), got %v", err)
	}
	if elapsed := time.Since(start); elapsed > 30*time.Second {
		t.Errorf("détection en %v, want < 30s (NFR-04)", elapsed)
	}
}

// TestFaultInjectionServerHangs : un serveur qui **accepte mais ne répond jamais**
// est borné par le timeout PAR TENTATIVE (500 ms) — l'appel échoue en quelques
// secondes, pas indéfiniment. Prouve que CallTool honore l'annulation du contexte
// (sinon un hang violerait NFR-04 sans être détecté).
func TestFaultInjectionServerHangs(t *testing.T) {
	block := make(chan struct{})
	t.Cleanup(func() { close(block) })

	server := mcp.NewServer(&mcp.Implementation{Name: "hang", Version: "v0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "slow", Description: "bloque"},
		func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, struct{}, error) {
			select {
			case <-ctx.Done(): // annulé par le timeout par tentative
			case <-block: // repli : fin du test
			}
			return nil, struct{}{}, ctx.Err()
		})

	ctx := context.Background()
	sess := connectInMemory(t, server)

	caller := NewResilientCaller(sess)
	start := time.Now()
	if _, err := caller.CallTool(ctx, "slow", nil); err == nil {
		t.Error("un outil qui bloque devrait échouer par timeout")
	}
	if el := time.Since(start); el > 10*time.Second {
		t.Errorf("hang non borné : %v — le timeout par tentative devrait borner", el)
	}
}

// TestFaultInjectionHangOpensBreaker : le timeout PAR TENTATIVE est une panne du
// service distant (il n'a pas répondu dans le délai), pas une annulation de
// l'appelant — il doit compter comme échec et OUVRIR le disjoncteur au seuil
// (NFR-04, FR-17). Avant correctif, CallTool remontait le ctx.Err() nu du
// timeout par tentative, que onResult neutralisait (m5) : un serveur MCP
// suspendu ne déclenchait jamais le fast-fail ErrOpen.
func TestFaultInjectionHangOpensBreaker(t *testing.T) {
	block := make(chan struct{})
	t.Cleanup(func() { close(block) })

	server := mcp.NewServer(&mcp.Implementation{Name: "hang", Version: "v0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "slow", Description: "bloque"},
		func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, struct{}, error) {
			select {
			case <-ctx.Done():
			case <-block:
			}
			return nil, struct{}{}, ctx.Err()
		})

	caller := NewResilientCaller(connectInMemory(t, server))
	// Réglages resserrés pour un test rapide : 1 tentative de 20 ms, seuil 2.
	caller.timeout = 20 * time.Millisecond
	caller.retry = resilience.RetryPolicy{MaxAttempts: 1}
	caller.breaker = resilience.NewBreaker(2, time.Minute)

	ctx := context.Background()
	for i := 0; i < 2; i++ {
		_, err := caller.CallTool(ctx, "slow", nil)
		if err == nil {
			t.Fatalf("appel %d sur serveur suspendu devrait échouer", i+1)
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			t.Errorf("appel %d : le timeout par tentative ne doit pas être classé annulation d'appelant : %v", i+1, err)
		}
	}
	if _, err := caller.CallTool(ctx, "slow", nil); !errors.Is(err, resilience.ErrOpen) {
		t.Errorf("après 2 timeouts de tentative, disjoncteur ouvert attendu (ErrOpen), got %v", err)
	}
}
