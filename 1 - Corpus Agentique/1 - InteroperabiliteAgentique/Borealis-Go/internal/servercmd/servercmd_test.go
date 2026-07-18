package servercmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/agbruneau/borealis/internal/mcpserver"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestServeHTTPRoundtrip : le mode HTTP streamable sert un serveur MCP réel ; un
// client s'y connecte et appelle un outil (mode de déploiement compose, M4.3).
func TestServeHTTPRoundtrip(t *testing.T) {
	hs := httptest.NewServer(httpHandler(mcpserver.NewCapacityServer))
	defer hs.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "c", Version: "v0"}, nil)
	cs, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: hs.URL}, nil)
	if err != nil {
		t.Fatalf("connexion HTTP : %v", err)
	}
	defer func() { _ = cs.Close() }()

	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "compute_capacity",
		Arguments: map[string]any{"income": 120000.0, "debts": 1000.0, "requestedAmount": 30000.0, "termMonths": 60},
	})
	if err != nil {
		t.Fatalf("appel d'outil via HTTP : %v", err)
	}
	if res.IsError {
		t.Errorf("outil en erreur : %+v", res)
	}
}

// TestHTTPHandlerRejectsOversizedBody : le handler MCP HTTP limite la taille
// du corps (m7) — ce mode n'a par ailleurs aucune authentification (m8), donc
// c'est la seule garde contre un corps arbitraire lu en entier en mémoire.
func TestHTTPHandlerRejectsOversizedBody(t *testing.T) {
	hs := httptest.NewServer(httpHandler(mcpserver.NewCapacityServer))
	defer hs.Close()

	body := bytes.NewReader(make([]byte, 2<<20)) // 2 MiB > limite (1 MiB)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, hs.URL, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("corps surdimensionné = %d, want 413", resp.StatusCode)
	}
}

// TestServeHTTPListenError : un addr déjà occupé fait échouer le bind —
// ServeHTTP retourne 1 et signale l'erreur (T1).
func TestServeHTTPListenError(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	var buf bytes.Buffer
	if code := ServeHTTP(context.Background(), ln.Addr().String(), mcpserver.NewCapacityServer, &buf); code != 1 {
		t.Errorf("ServeHTTP (port occupé) = %d, want 1", code)
	}
}

// TestRunMCPHTTPListenError : même scénario via RunMCP (MCP_HTTP_ADDR).
func TestRunMCPHTTPListenError(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	t.Setenv("MCP_HTTP_ADDR", ln.Addr().String())
	var buf bytes.Buffer
	if code := RunMCP(context.Background(), []string{"srv"}, &buf, mcpserver.NewCapacityServer); code != 1 {
		t.Errorf("RunMCP (port occupé) = %d, want 1", code)
	}
}

// TestServeHTTPShutsDown : ServeHTTP s'arrête proprement (code 0) à l'annulation.
func TestServeHTTPShutsDown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if code := ServeHTTP(ctx, "127.0.0.1:0", mcpserver.NewCapacityServer, io.Discard); code != 0 {
		t.Errorf("ServeHTTP (ctx annulé) → code %d, want 0", code)
	}
}

// TestRunMCPVersion : l'aiguillage traite --version (code 0, sans servir).
func TestRunMCPVersion(t *testing.T) {
	var buf bytes.Buffer
	if code := RunMCP(context.Background(), []string{"srv", "--version"}, &buf, mcpserver.NewCapacityServer); code != 0 {
		t.Errorf("--version → code %d, want 0", code)
	}
	if !strings.Contains(buf.String(), "borealis") {
		t.Errorf("version = %q", buf.String())
	}
}

// TestRunMCPHTTP : MCP_HTTP_ADDR défini → mode HTTP ; ctx annulé → arrêt propre.
func TestRunMCPHTTP(t *testing.T) {
	t.Setenv("MCP_HTTP_ADDR", "127.0.0.1:0")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if code := RunMCP(ctx, []string{"srv"}, io.Discard, mcpserver.NewCapacityServer); code != 0 {
		t.Errorf("mode HTTP (ctx annulé) → code %d, want 0", code)
	}
}

// TestRunMCPStdio : sans MCP_HTTP_ADDR, aiguille vers stdio (serve). Se termine
// à la fermeture du client (transport en mémoire injecté via l'absence d'env).
func TestRunMCPStdio(t *testing.T) {
	t.Setenv("MCP_HTTP_ADDR", "")
	t1, t2 := mcp.NewInMemoryTransports()
	done := make(chan int, 1)
	// RunMCP en stdio utilise StdioTransport ; on couvre ici la branche via serve
	// directement avec un transport en mémoire (même chemin de service).
	go func() { done <- serve(context.Background(), io.Discard, mcpserver.NewCapacityServer, t1) }()
	client := mcp.NewClient(&mcp.Implementation{Name: "c", Version: "v0"}, nil)
	cs, err := client.Connect(context.Background(), t2, nil)
	if err != nil {
		t.Fatalf("connexion : %v", err)
	}
	_ = cs.Close()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("service stdio ne s'est pas terminé")
	}
}

// TestRunMCPStdioCtxCancelled : la branche stdio doit s'arrêter proprement (code
// 0) à l'annulation de ctx, même sans fermeture du client (M2 : avant ce
// correctif, `context.Background()` était codé en dur et SIGINT/SIGTERM
// n'atteignaient jamais le transport, forçant un SIGKILL).
func TestRunMCPStdioCtxCancelled(t *testing.T) {
	t1, _ := mcp.NewInMemoryTransports()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() { done <- serve(ctx, io.Discard, mcpserver.NewCapacityServer, t1) }()
	cancel()
	select {
	case code := <-done:
		if code != 0 {
			t.Errorf("stdio (ctx annulé) → code %d, want 0", code)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("le service stdio n'a pas observé l'annulation du contexte")
	}
}

type failTransport struct{}

func (failTransport) Connect(context.Context) (mcp.Connection, error) {
	return nil, errors.New("boom")
}

func TestServeError(t *testing.T) {
	var buf bytes.Buffer
	code := serve(context.Background(), &buf, mcpserver.NewCapacityServer, failTransport{})
	if code != 1 {
		t.Errorf("code = %d, want 1", code)
	}
	if !strings.Contains(buf.String(), "error") {
		t.Errorf("sortie erreur = %q", buf.String())
	}
}

func TestServeEndsOnClientClose(t *testing.T) {
	t1, t2 := mcp.NewInMemoryTransports()
	done := make(chan int, 1)
	go func() {
		done <- serve(context.Background(), io.Discard, mcpserver.NewCapacityServer, t1)
	}()
	client := mcp.NewClient(&mcp.Implementation{Name: "c", Version: "v0"}, nil)
	cs, err := client.Connect(context.Background(), t2, nil)
	if err != nil {
		t.Fatalf("connexion client : %v", err)
	}
	if err := cs.Close(); err != nil {
		t.Fatalf("fermeture client : %v", err)
	}
	select {
	case code := <-done:
		if code != 0 {
			t.Errorf("code = %d, want 0", code)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serve n'a pas terminé après fermeture du client")
	}
}
