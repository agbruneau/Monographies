// Package servercmd factorise le lancement d'un serveur MCP en binaire :
// drapeau --version, service sur un transport **stdio** (défaut) ou **HTTP
// streamable** (déploiement compose). La logique est ici (testable) ; les main()
// des cmd/mcp-* se réduisent à un aiguillage.
package servercmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/agbruneau/borealis/internal/agentcmd"
	"github.com/agbruneau/borealis/internal/meta"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// serve sert le serveur sur le transport fourni jusqu'à sa fermeture ou
// l'annulation de ctx (arrêt gracieux, M2 : SIGINT/SIGTERM en mode stdio).
// Transport injecté pour le test (en mémoire, échec simulé).
func serve(ctx context.Context, out io.Writer, newServer func() *mcp.Server, t mcp.Transport) int {
	err := newServer().Run(ctx, t)
	if err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return 0
	}
	fmt.Fprintln(out, "error:", err)
	return 1
}

// maxRequestBody borne la taille du corps des requêtes HTTP entrantes (m7) :
// sans cette limite, un corps JSON arbitraire est lu en entier en mémoire —
// d'autant plus sensible ici que ce mode n'a aucune authentification (m8).
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

// httpHandler construit un handler MCP **HTTP streamable** (SDK) servant une
// unique instance de serveur à toutes les sessions. Facteur testable de ServeHTTP.
//
// Pas de garde bearer ici (contrairement à agentcmd.handlerFor, FR-26) : la
// frontière de confiance assumée est le réseau compose `mesh` (internal: true,
// aucun port publié) — voir ADR-0011 pour la justification et la condition de
// réversion.
func httpHandler(newServer func() *mcp.Server) http.Handler {
	srv := newServer()
	return limitBody(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return srv }, nil), maxRequestBody)
}

// ServeHTTP sert le serveur MCP en HTTP streamable sur addr jusqu'à l'annulation
// de ctx (arrêt gracieux). C'est le mode de déploiement en compose, où chaque
// serveur MCP est un service réseau (§14, M4.3). Retourne un code de sortie.
func ServeHTTP(ctx context.Context, addr string, newServer func() *mcp.Server, out io.Writer) int {
	if err := agentcmd.Serve(ctx, addr, httpHandler(newServer), out); err != nil {
		fmt.Fprintln(out, "error:", err)
		return 1
	}
	return 0
}

// RunMCP aiguille le lancement d'un serveur MCP : --version, puis HTTP streamable
// si MCP_HTTP_ADDR est défini (déploiement compose), sinon stdio (défaut). Facteur
// **testable** partagé par les 4 binaires cmd/mcp-* (leur run() n'ajoute que la
// gestion du signal). ctx pilote l'arrêt gracieux du mode HTTP.
func RunMCP(ctx context.Context, args []string, out io.Writer, newServer func() *mcp.Server) int {
	if len(args) > 1 && args[1] == "--version" {
		fmt.Fprintln(out, meta.Info())
		return 0
	}
	if addr := os.Getenv("MCP_HTTP_ADDR"); addr != "" {
		return ServeHTTP(ctx, addr, newServer, out)
	}
	return serve(ctx, out, newServer, &mcp.StdioTransport{})
}
