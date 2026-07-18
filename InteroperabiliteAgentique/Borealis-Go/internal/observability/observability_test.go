package observability

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// tracerFor construit un traceur OTel minimal (span recorder), sans dépendre
// d'un fournisseur applicatif — seul TraceRoundTripper/InjectTraceparent sont
// du code de production ici (B2 : observability.Provider supprimé, jamais
// instancié par les binaires).
func tracerFor(t *testing.T, rec sdktrace.SpanProcessor) trace.Tracer {
	t.Helper()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	t.Cleanup(func() { _ = tp.Shutdown(context.Background()) })
	return tp.Tracer("x")
}

func TestTraceRoundTripper(t *testing.T) {
	rec := tracetest.NewSpanRecorder()
	ctx, span := tracerFor(t, rec).Start(context.Background(), "op")
	defer span.End()

	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("traceparent")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &http.Client{Transport: &TraceRoundTripper{}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	if !strings.HasPrefix(got, "00-") {
		t.Errorf("traceparent reçu = %q, want format W3C (00-...)", got)
	}
	// Contrat http.RoundTripper : la requête ORIGINALE ne doit pas être mutée
	// (l'injection se fait sur un clone) — une réémission ou un partage de la
	// même *http.Request par un transport supérieur muterait sinon sa map Header.
	if h := req.Header.Get("traceparent"); h != "" {
		t.Errorf("requête originale mutée (traceparent=%q) — RoundTrip doit cloner", h)
	}
}

// TestInjectTraceparent : le contexte d'un span actif produit un traceparent
// W3C via la carte retournée (canal non-HTTP, ex. champ _meta MCP, T5).
func TestInjectTraceparent(t *testing.T) {
	rec := tracetest.NewSpanRecorder()
	ctx, span := tracerFor(t, rec).Start(context.Background(), "op")
	defer span.End()

	got := InjectTraceparent(ctx)
	if tp, ok := got["traceparent"]; !ok || !strings.HasPrefix(tp, "00-") {
		t.Errorf("traceparent = %q (présent=%v), want format W3C (00-...)", got["traceparent"], ok)
	}
}
