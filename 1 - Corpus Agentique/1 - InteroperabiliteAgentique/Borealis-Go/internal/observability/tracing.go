package observability

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/propagation"
)

// Propagator est le propagateur W3C **Trace Context** (en-tête traceparent) :
// il assure une trace ininterrompue à travers A2A et MCP (NFR-07).
var Propagator = propagation.TraceContext{}

// TraceRoundTripper injecte l'en-tête traceparent (issu du contexte de la
// requête) dans chaque appel HTTP sortant — propagation A2A.
type TraceRoundTripper struct {
	Base http.RoundTripper
}

// RoundTrip injecte traceparent puis délègue au transport sous-jacent.
// L'injection se fait sur un CLONE : le contrat http.RoundTripper interdit de
// modifier la requête reçue (une réémission ou un partage de la même
// *http.Request muterait sinon sa map Header — course de données possible).
func (t *TraceRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	Propagator.Inject(req.Context(), propagation.HeaderCarrier(req.Header))
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

// InjectTraceparent retourne la carte des en-têtes de propagation (traceparent)
// du contexte, pour les canaux non-HTTP comme le champ _meta d'un appel MCP.
func InjectTraceparent(ctx context.Context) map[string]string {
	carrier := propagation.MapCarrier{}
	Propagator.Inject(ctx, carrier)
	return carrier
}
