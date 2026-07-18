// Commande mcp-capacity : serveur MCP Capacity Calculator (compute_capacity).
// Transport stdio par défaut ; HTTP streamable si MCP_HTTP_ADDR est défini
// (déploiement compose). L'aiguillage est dans internal/servercmd (testable).
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/agbruneau/borealis/internal/mcpserver"
	"github.com/agbruneau/borealis/internal/servercmd"
)

func main() { os.Exit(run()) }

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return servercmd.RunMCP(ctx, os.Args, os.Stdout, mcpserver.NewCapacityServer)
}
