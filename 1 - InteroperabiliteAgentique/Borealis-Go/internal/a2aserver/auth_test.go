package a2aserver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a2aproject/a2a-go/v2/a2asrv"
	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/idpmock"
	"go.opentelemetry.io/otel/trace"
)

func httpGetCode(t *testing.T, url, authz string) int {
	t.Helper()
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if authz != "" {
		req.Header.Set("Authorization", authz)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}

// TestSecureMux : le montage sécurisé protège le point d'invocation par
// construction (401 sans jeton, non-401 avec jeton valide) et laisse la card
// publique (FR-26, FR-01).
func TestSecureMux(t *testing.T) {
	journal := audit.New()
	idp := idpmock.New([]byte("dev"))
	card := BuildCard(Specs[1], "http://x")
	srv := httptest.NewServer(SecureMux(card, &ReplyExecutor{Reply: "ok"}, idp, "borealis-kyc", journal))
	defer srv.Close()

	if code := httpGetCode(t, srv.URL+InvokePath, ""); code != http.StatusUnauthorized {
		t.Errorf("invoke sans jeton : code %d, want 401", code)
	}
	tok, _ := idp.Issue("orch", "borealis-kyc")
	if code := httpGetCode(t, srv.URL+InvokePath, "Bearer "+tok); code == http.StatusUnauthorized {
		t.Errorf("invoke avec jeton valide : rejeté (401) à tort")
	}
	if code := httpGetCode(t, srv.URL+a2asrv.WellKnownAgentCardPath, ""); code != http.StatusOK {
		t.Errorf("card publique : code %d, want 200", code)
	}
}

// TestTraceExtract : un traceparent entrant est extrait vers le contexte
// (contexte de span distant valide) — trace ininterrompue A2A (NFR-07).
func TestTraceExtract(t *testing.T) {
	var remoteValid bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sc := trace.SpanContextFromContext(r.Context())
		remoteValid = sc.IsValid() && sc.IsRemote()
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(TraceExtract(next))
	defer srv.Close()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	req.Header.Set("traceparent", "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if !remoteValid {
		t.Error("traceparent entrant non extrait vers le contexte (NFR-07)")
	}
}

func TestBearerAuth(t *testing.T) {
	journal := audit.New()
	idp := idpmock.New([]byte("dev-secret"))
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	srv := httptest.NewServer(BearerAuth(next, idp, "borealis-kyc", journal))
	defer srv.Close()

	get := func(path, authz string) int {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL+path, nil)
		if authz != "" {
			req.Header.Set("Authorization", authz)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = resp.Body.Close() }()
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp.StatusCode
	}

	tok, _ := idp.Issue("orchestrator", "borealis-kyc")
	if code := get(InvokePath, "Bearer "+tok); code != http.StatusOK {
		t.Errorf("jeton valide : code %d, want 200", code)
	}
	if code := get(InvokePath, ""); code != http.StatusUnauthorized {
		t.Errorf("jeton absent : code %d, want 401", code)
	}
	tokWrong, _ := idp.Issue("orch", "autre-agent")
	if code := get(InvokePath, "Bearer "+tokWrong); code != http.StatusUnauthorized {
		t.Errorf("aud non conforme : code %d, want 401 (RFC 8707)", code)
	}
	if code := get(InvokePath, "Bearer "+tok+"tamper"); code != http.StatusUnauthorized {
		t.Errorf("signature invalide : code %d, want 401", code)
	}
	// Agent Card publique accessible sans jeton.
	if code := get(a2asrv.WellKnownAgentCardPath, ""); code != http.StatusOK {
		t.Errorf("card publique : code %d, want 200", code)
	}

	// Trois rejets journalisés (absent, aud, signature).
	if n := journal.Len(); n != 3 {
		t.Errorf("entrées d'audit = %d, want 3 (rejets journalisés)", n)
	}
}

// TestBearerAuthRejectAuditCapped : la journalisation des rejets d'auth est
// plafonnée — elle est pilotable AVANT authentification, et sans plafond un
// flot de requêtes sans jeton fait croître le journal en mémoire sans borne
// (DoS mémoire non authentifié). Au plafond : une entrée de coupure, puis plus
// aucune écriture ; le 401 reste servi.
func TestBearerAuthRejectAuditCapped(t *testing.T) {
	old := maxAuthAuditEntries
	maxAuthAuditEntries = 2
	defer func() { maxAuthAuditEntries = old }()

	journal := audit.New()
	idp := idpmock.New([]byte("dev-secret"))
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	srv := httptest.NewServer(BearerAuth(next, idp, "borealis-kyc", journal))
	defer srv.Close()

	for i := 0; i < 5; i++ {
		if code := httpGetCode(t, srv.URL+InvokePath, ""); code != http.StatusUnauthorized {
			t.Fatalf("requête %d sans jeton : code %d, want 401", i+1, code)
		}
	}
	if n := journal.Len(); n != 3 {
		t.Errorf("entrées d'audit = %d, want 3 (2 rejets + 1 entrée de coupure)", n)
	}
}
