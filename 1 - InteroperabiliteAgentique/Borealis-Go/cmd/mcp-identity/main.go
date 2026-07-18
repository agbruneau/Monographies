// Commande mcp-identity : serveur MCP Identity Hub (verify_identity +
// check_watchlists). Transport stdio par défaut ; HTTP streamable si
// MCP_HTTP_ADDR est défini (déploiement compose). Charge le jeu synthétique
// seedé en mémoire. L'aiguillage est dans internal/servercmd (testable).
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/agbruneau/borealis/internal/mcpserver"
	"github.com/agbruneau/borealis/internal/servercmd"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() { os.Exit(run()) }

func run() int {
	newServer := func() *mcp.Server {
		return mcpserver.NewIdentityServer(fixtures.Generate(fixtures.DefaultSeed, 10000))
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return servercmd.RunMCP(ctx, os.Args, os.Stdout, newServer)
}
