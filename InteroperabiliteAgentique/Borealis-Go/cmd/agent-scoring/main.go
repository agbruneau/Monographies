// Commande agent-scoring : agent A2A Scoring (compute-credit-score), servi en
// HTTP. Service dans internal/agentcmd ; ADDR/BASE_URL par l'environnement.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/agbruneau/borealis/internal/agentcmd"
)

func main() { os.Exit(run()) }

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return agentcmd.Run(ctx, "borealis-scoring", os.Stdout)
}
