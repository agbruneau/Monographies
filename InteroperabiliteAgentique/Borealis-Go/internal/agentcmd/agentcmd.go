// Package agentcmd factorise le lancement d'un agent A2A en binaire : construit
// sa card, sert le point d'invocation + la card well-known, et s'arrête
// proprement à l'annulation du contexte. La logique est ici (testable) ; les
// main() des cmd/agent-* et cmd/orchestrator se réduisent à un appel à Run.
package agentcmd

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/agbruneau/borealis/internal/a2aserver"
	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/idpmock"
)

// env retourne la variable d'environnement key, ou def si absente/vide.
func env(key, def string) string { return cmp.Or(os.Getenv(key), def) }

// maxRequestBody borne la taille du corps des requêtes HTTP entrantes (m7) :
// sans cette limite, un corps JSON arbitraire est lu en entier en mémoire par
// le handler JSON-RPC (DoS trivial sur un service exposé).
const maxRequestBody = 1 << 20 // 1 MiB

// limitBody rejette (413) tout corps annoncé (Content-Length) au-delà de n, et
// borne aussi la lecture effective via MaxBytesReader (défense en profondeur
// pour le transfert par blocs, sans Content-Length connu à l'avance).
func limitBody(next http.Handler, n int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > n {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, n)
		next.ServeHTTP(w, r)
	})
}

// handlerFor construit le routeur HTTP **protégé** d'un agent (SecureMux) : la
// card reste publique (FR-01) mais le point d'invocation exige un bearer JWT à
// audience restreinte = nom de l'agent (FR-26) — jeton absent/invalide/aud
// étranger → 401 + audit. L'exécuteur est un stub de réponse (la logique métier
// — appels MCP — reste au niveau bibliothèque).
//
// Le secret IdP est un PARAMÈTRE (fourni par Run depuis IDP_SECRET, garanti non
// vide par la garde m9) : aucun repli committé — un appelant qui contournerait
// Run ne peut plus redémarrer silencieusement avec un secret lisible dans le
// dépôt (jetons forgeables). La délivrance/attache du jeton côté appelant
// (orchestrateur) reste la couture attendue de FR-26.
func handlerFor(spec a2aserver.Spec, baseURL string, idpSecret []byte) http.Handler {
	card := a2aserver.BuildCard(spec, baseURL)
	exec := &a2aserver.ReplyExecutor{Reply: spec.Name + " ready"}
	idp := idpmock.New(idpSecret)
	return limitBody(a2aserver.SecureMux(card, exec, idp, spec.Name, audit.New()), maxRequestBody)
}

// Run sert l'agent nommé sur l'adresse et la base d'URL lues de l'environnement
// (ADDR, BASE_URL), jusqu'à l'annulation de ctx (arrêt gracieux). Retourne un
// code de sortie de processus. Nom d'agent inconnu → erreur (code 2).
// IDP_SECRET absent → erreur (code 3, m9) : sans cette garde, un déploiement
// oubliant la variable démarre silencieusement avec le secret de démonstration
// committé (jetons forgeables par quiconque lit le dépôt).
func Run(ctx context.Context, agentName string, out io.Writer) int {
	spec, ok := a2aserver.SpecByName(agentName)
	if !ok {
		fmt.Fprintf(out, "error: agent inconnu %q\n", agentName)
		return 2
	}
	secret := os.Getenv("IDP_SECRET")
	if secret == "" {
		fmt.Fprintln(out, "error: IDP_SECRET requis en mode réseau (aucun secret de démonstration par défaut)")
		return 3
	}
	addr := env("ADDR", ":8080")
	baseURL := env("BASE_URL", "http://localhost"+addr)
	if err := Serve(ctx, addr, handlerFor(spec, baseURL, []byte(secret)), out); err != nil {
		fmt.Fprintln(out, "error:", err)
		return 1
	}
	return 0
}

// Serve écoute sur addr et sert h jusqu'à l'annulation de ctx (arrêt gracieux).
func Serve(ctx context.Context, addr string, h http.Handler, out io.Writer) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return serve(ctx, ln, h, out)
}

// shutdownGrace borne le drain gracieux à l'arrêt (surchargeable en test).
var shutdownGrace = 5 * time.Second

// newServer construit le http.Server durci commun aux agents A2A et aux
// serveurs MCP HTTP (via Serve) : en-têtes 5 s, lecture de la requête 30 s,
// connexions oisives 120 s. limitBody borne la TAILLE du corps, pas sa DURÉE
// de lecture — sans ReadTimeout, un client égrenant son corps octet par octet
// retient goroutine et connexion indéfiniment (slowloris). Pas de WriteTimeout
// (il romprait les flux SSE A2A/MCP).
// ponytail: ReadTimeout borne aussi la durée TOTALE d'une réponse en flux
// (lecture d'arrière-plan) — plafond 30 s, large pour des flux qui se terminent
// en quelques secondes ; à relever si un flux longue durée apparaît.
func newServer(h http.Handler) *http.Server {
	return &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

// serve sert h sur un listener déjà ouvert jusqu'à l'annulation de ctx, puis
// s'arrête proprement (drain borné). Cœur testable (listener injecté).
func serve(ctx context.Context, ln net.Listener, h http.Handler, out io.Writer) error {
	srv := newServer(h)
	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()
	fmt.Fprintf(out, "agent en écoute sur %s\n", ln.Addr())

	select {
	case err := <-errc:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		sctx, cancel := context.WithTimeout(context.Background(), shutdownGrace)
		defer cancel()
		if err := srv.Shutdown(sctx); err != nil {
			// Drain gracieux expiré (flux SSE longue durée encore ouvert) : forcer
			// la fermeture. L'arrêt étant demandé (ctx annulé), c'est un arrêt propre,
			// pas un échec — ne pas remonter le timeout comme code d'erreur.
			_ = srv.Close()
		}
		return nil
	}
}
