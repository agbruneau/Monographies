package a2aserver

import (
	"net/http"
	"strings"
	"sync/atomic"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/a2aproject/a2a-go/v2/a2asrv"
	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/idpmock"
	"github.com/agbruneau/borealis/internal/observability"
	"go.opentelemetry.io/otel/propagation"
)

// SecureMux assemble un agent A2A **protégé par construction** : le traceparent
// entrant est extrait vers le contexte (NFR-07), puis le point d'invocation est
// gardé par BearerAuth (FR-26). L'Agent Card reste publique. C'est le montage à
// utiliser en déploiement (le Mux nu ne sert qu'aux tests unitaires).
func SecureMux(card *a2asdk.AgentCard, exec a2asrv.AgentExecutor, idp *idpmock.IdP, audience string, journal *audit.Journal) http.Handler {
	return TraceExtract(BearerAuth(Mux(card, exec), idp, audience, journal))
}

// TraceExtract extrait le traceparent entrant vers le contexte de la requête,
// pour une trace ininterrompue à travers l'agent (NFR-07).
func TraceExtract(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := observability.Propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// maxAuthAuditEntries plafonne les entrées d'audit produites par les rejets
// d'authentification d'un même handler : la journalisation 401 est pilotable
// AVANT authentification, et le journal (en mémoire, append-only) croîtrait
// sans borne sous un flot de requêtes sans jeton (DoS mémoire non
// authentifié). Au plafond : une dernière entrée marque la coupure, puis les
// rejets suivants ne sont plus journalisés — ils restent refusés en 401.
// Variable (non const) pour être abaissée en test.
var maxAuthAuditEntries uint64 = 1000

// BearerAuth enveloppe un handler d'une vérification de bearer JWT à **audience
// restreinte** : jeton absent, invalide, ou dont l'aud ne correspond pas à cet
// agent → **401 + entrée d'audit** (FR-26, §11.1), journalisation plafonnée
// (maxAuthAuditEntries). Le chemin well-known (Agent Card publique) reste
// accessible sans jeton (FR-01 : GET sans auth).
func BearerAuth(next http.Handler, idp *idpmock.IdP, audience string, journal *audit.Journal) http.Handler {
	var rejects atomic.Uint64
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == a2asrv.WellKnownAgentCardPath {
			next.ServeHTTP(w, r)
			return
		}
		tok := bearerToken(r.Header.Get("Authorization"))
		if tok == "" {
			reject(w, journal, audience, "missing_bearer", &rejects)
			return
		}
		if _, err := idp.Verify(tok, audience); err != nil {
			reject(w, journal, audience, "invalid_token", &rejects)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if strings.HasPrefix(header, prefix) {
		return strings.TrimPrefix(header, prefix)
	}
	return ""
}

func reject(w http.ResponseWriter, journal *audit.Journal, audience, reason string, rejects *atomic.Uint64) {
	if journal != nil {
		switch n := rejects.Add(1); {
		case n <= maxAuthAuditEntries:
			journal.Append(audit.Record{
				KYA: "agent:" + audience, Action: "a2a.auth", Result: "401 " + reason, Version: "m2",
			})
		case n == maxAuthAuditEntries+1:
			journal.Append(audit.Record{
				KYA: "agent:" + audience, Action: "a2a.auth",
				Result: "401 auth_audit_capped (rejets suivants non journalisés)", Version: "m2",
			})
		}
	}
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}
