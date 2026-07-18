// Commande orchestrator : agent A2A Orchestrateur (orchestrate-prequal), servi
// en HTTP. Service dans internal/agentcmd ; ADDR/BASE_URL par l'environnement.
// La découverte des pairs et le routage vivent dans internal/orchestrator,
// consommés par l'exécuteur de l'agent (câblage métier complet hors périmètre du
// démonstrateur — exécuteur de réponse par défaut).
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/agbruneau/borealis/internal/agentcmd"
	"github.com/agbruneau/borealis/internal/auditexport"
)

func main() { os.Exit(run()) }

func run() int {
	// --audit-export : exporte le journal WORM d'une pré-qualification canonique
	// en JSONL rejouable (FR-23, M5.2b) puis sort — pas de service.
	if len(os.Args) > 1 && os.Args[1] == "--audit-export" {
		if err := auditexport.Run(context.Background(), os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, "audit-export:", err)
			return 1
		}
		return 0
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return agentcmd.Run(ctx, "borealis-orchestrator", os.Stdout)
}
