// Commande agent-kyc : agent A2A KYC (verify-identity), servi en HTTP. La
// logique de service est dans internal/agentcmd (testable) ; adresse et base
// d'URL lues de l'environnement (ADDR, BASE_URL).
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
	return agentcmd.Run(ctx, "borealis-kyc", os.Stdout)
}
