package agentcmd

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/a2aproject/a2a-go/v2/a2asrv"
	"github.com/agbruneau/borealis/internal/a2aserver"
	"github.com/agbruneau/borealis/internal/idpmock"
)

// getWithRetry tolère la fenêtre entre l'ouverture du listener et l'acceptation.
func getWithRetry(t *testing.T, url string) *http.Response {
	t.Helper()
	client := &http.Client{Timeout: time.Second}
	for i := 0; i < 50; i++ {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		if resp, err := client.Do(req); err == nil {
			return resp
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("aucune réponse de %s", url)
	return nil
}

// TestServeCardAndShutdown : l'agent sert sa card well-known, puis s'arrête
// proprement à l'annulation du contexte (drain gracieux, aucune erreur).
func TestServeCardAndShutdown(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	spec, ok := a2aserver.SpecByName("borealis-kyc")
	if !ok {
		t.Fatal("spec borealis-kyc absente")
	}
	h := handlerFor(spec, "http://"+ln.Addr().String(), []byte("test-idp-secret-not-for-prod"))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- serve(ctx, ln, h, io.Discard) }()

	resp := getWithRetry(t, "http://"+ln.Addr().String()+a2asrv.WellKnownAgentCardPath)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("card well-known : status %d, want 200", resp.StatusCode)
	}
	_ = resp.Body.Close()

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("arrêt non propre : %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Serve ne s'est pas arrêté après annulation du contexte")
	}
}

// TestServeForcesShutdownOnHang : une requête active qui ne se draine pas (flux
// SSE longue durée) ne doit pas faire échouer l'arrêt : après le drain gracieux
// expiré, serve force la fermeture et retourne nil (W3).
func TestServeForcesShutdownOnHang(t *testing.T) {
	old := shutdownGrace
	shutdownGrace = 50 * time.Millisecond
	defer func() { shutdownGrace = old }()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	reached := make(chan struct{}, 1)
	block := make(chan struct{})
	h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		select {
		case reached <- struct{}{}:
		default:
		}
		<-block // requête active jusqu'à la fin du test
	})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- serve(ctx, ln, h, io.Discard) }()

	go func() {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://"+ln.Addr().String()+"/", nil)
		if resp, e := (&http.Client{}).Do(req); e == nil {
			_ = resp.Body.Close()
		}
	}()
	<-reached // la requête bloquante est arrivée au handler
	cancel()
	select {
	case e := <-done:
		if e != nil {
			t.Errorf("arrêt forcé devrait retourner nil, got %v", e)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("serve ne s'est pas arrêté (repli Close manquant ?)")
	}
	close(block)
}

// TestSecureInvokeRequiresAuth : dans le chemin DÉPLOYÉ (handlerFor), la card
// well-known reste publique (FR-01) mais le point d'invocation exige un bearer —
// une requête non authentifiée reçoit 401 (FR-26). Ferme l'écart « auth non
// câblée dans le binaire » relevé en vérification finale.
func TestSecureInvokeRequiresAuth(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	spec, ok := a2aserver.SpecByName("borealis-kyc")
	if !ok {
		t.Fatal("spec absente")
	}
	base := "http://" + ln.Addr().String()
	h := handlerFor(spec, base, []byte("test-idp-secret-not-for-prod"))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- serve(ctx, ln, h, io.Discard) }()
	defer func() { cancel(); <-done }()

	// Card publique (FR-01).
	cardResp := getWithRetry(t, base+a2asrv.WellKnownAgentCardPath)
	if cardResp.StatusCode != http.StatusOK {
		t.Errorf("card well-known = %d, want 200 (publique)", cardResp.StatusCode)
	}
	_ = cardResp.Body.Close()

	// Invocation SANS jeton → 401 (FR-26).
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, base+a2aserver.InvokePath, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := (&http.Client{Timeout: time.Second}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("invocation non authentifiée = %d, want 401 (FR-26)", resp.StatusCode)
	}
}

// TestRunListenError : un ADDR déjà occupé fait échouer net.Listen — Run
// retourne 1 et signale l'erreur, sans jamais démarrer de service (T1).
func TestRunListenError(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	t.Setenv("IDP_SECRET", "test-idp-secret-not-for-prod")
	t.Setenv("ADDR", ln.Addr().String()) // déjà occupé
	var buf bytes.Buffer
	if code := Run(context.Background(), "borealis-kyc", &buf); code != 1 {
		t.Errorf("Run (port occupé) = %d, want 1", code)
	}
	if !strings.Contains(buf.String(), "error") {
		t.Errorf("sortie = %q, devrait mentionner l'erreur", buf.String())
	}
}

// TestServeListenerErrorPropagates : une erreur de service autre que la
// fermeture gracieuse (ex. listener déjà fermé) doit être propagée par serve,
// pas absorbée en nil (T1).
func TestServeListenerErrorPropagates(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	_ = ln.Close() // fermé AVANT que serve() ne l'utilise
	if err := serve(context.Background(), ln, http.NewServeMux(), io.Discard); err == nil {
		t.Error("listener déjà fermé : erreur attendue")
	}
}

// TestSecureInvokeRejectsOversizedBody : le point d'invocation limite la
// taille du corps (m7) — un corps authentifié dépassant la limite est rejeté
// (413), jamais lu en entier en mémoire (DoS trivial sinon).
func TestSecureInvokeRejectsOversizedBody(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	spec, ok := a2aserver.SpecByName("borealis-kyc")
	if !ok {
		t.Fatal("spec absente")
	}
	base := "http://" + ln.Addr().String()
	h := handlerFor(spec, base, []byte("test-idp-secret-not-for-prod"))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- serve(ctx, ln, h, io.Discard) }()
	defer func() { cancel(); <-done }()

	idp := idpmock.New([]byte("test-idp-secret-not-for-prod")) // même secret que handlerFor ci-dessus
	tok, err := idp.Issue("test", spec.Name)
	if err != nil {
		t.Fatal(err)
	}

	body := bytes.NewReader(make([]byte, 2<<20)) // 2 MiB > limite (1 MiB)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, base+a2aserver.InvokePath, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("corps surdimensionné = %d, want 413", resp.StatusCode)
	}
}

// TestNewServerTimeouts : le serveur HTTP commun est durci — en-têtes, lecture
// du corps et connexions oisives bornés (anti-slowloris : limitBody borne la
// TAILLE du corps, pas sa DURÉE de lecture) ; PAS de WriteTimeout (il romprait
// les flux SSE A2A/MCP).
func TestNewServerTimeouts(t *testing.T) {
	srv := newServer(http.NotFoundHandler())
	if srv.ReadHeaderTimeout <= 0 || srv.ReadTimeout <= 0 || srv.IdleTimeout <= 0 {
		t.Errorf("timeouts manquants : ReadHeader=%v Read=%v Idle=%v",
			srv.ReadHeaderTimeout, srv.ReadTimeout, srv.IdleTimeout)
	}
	if srv.WriteTimeout != 0 {
		t.Errorf("WriteTimeout=%v, want 0 (romprait les flux SSE)", srv.WriteTimeout)
	}
}

// TestRunUnknownAgent : un nom d'agent inconnu retourne le code 2.
func TestRunUnknownAgent(t *testing.T) {
	if code := Run(context.Background(), "agent-fantome", io.Discard); code != 2 {
		t.Errorf("agent inconnu → code %d, want 2", code)
	}
}

// TestRunServesThenShutsDown : Run sur un agent connu, contexte annulé d'emblée,
// se termine proprement (code 0). ADDR sur un port éphémère.
func TestRunServesThenShutsDown(t *testing.T) {
	t.Setenv("ADDR", "127.0.0.1:0")
	t.Setenv("IDP_SECRET", "test-idp-secret-not-for-prod")
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // annulé d'emblée → arrêt gracieux immédiat
	if code := Run(ctx, "borealis-orchestrator", io.Discard); code != 0 {
		t.Errorf("Run (ctx annulé) → code %d, want 0", code)
	}
}

// TestRunRequiresIDPSecret : Run refuse de démarrer en mode réseau sans
// IDP_SECRET explicite (m9) — sans cette garde, tout déploiement oubliant la
// variable acceptait silencieusement des jetons forgeables avec le secret de
// démonstration committé.
func TestRunRequiresIDPSecret(t *testing.T) {
	t.Setenv("IDP_SECRET", "")
	t.Setenv("ADDR", "127.0.0.1:0")
	var buf bytes.Buffer
	if code := Run(context.Background(), "borealis-orchestrator", &buf); code == 0 {
		t.Error("Run sans IDP_SECRET devrait refuser de démarrer (code non nul)")
	}
	if !strings.Contains(buf.String(), "IDP_SECRET") {
		t.Errorf("message d'erreur = %q, devrait mentionner IDP_SECRET", buf.String())
	}
}
