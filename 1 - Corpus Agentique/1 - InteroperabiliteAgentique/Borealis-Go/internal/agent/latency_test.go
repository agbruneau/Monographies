package agent

import (
	"context"
	"testing"

	"github.com/agbruneau/borealis/internal/audit"
)

// BenchmarkPrequalifyLoop mesure la latence de la boucle de pré-qualification
// (NFR-01, `make bench`) : identité → criblage AML → évaluation → raisonnement,
// journal WORM à chaque étape. Serveur MCP en mémoire, fake Reasoner.
func BenchmarkPrequalifyLoop(b *testing.B) {
	sess := newIdentitySession(b)
	a := goodApplicant()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := NewOrchestrator(sess, stubReasoner{}, audit.New(), testPolicy())
		if _, err := o.Prequalify(context.Background(), a); err != nil {
			b.Fatal(err)
		}
	}
}

// TestPrequalifyLoopSmoke vérifie que la boucle complète sans erreur. Le SLO
// mesuré NFR-01 (P99 < 2s) n'est pas falsifiable ici : la boucle en mémoire
// (fake Reasoner) se compte en dizaines de µs, plusieurs ordres de grandeur
// sous le seuil — la preuve mesurée revient à la composition déployée (T4 ;
// l'ancienne version bouclait 300x pour un résultat jamais autrement que vert).
func TestPrequalifyLoopSmoke(t *testing.T) {
	sess := newIdentitySession(t)
	a := goodApplicant()
	o := NewOrchestrator(sess, stubReasoner{}, audit.New(), testPolicy())
	if _, err := o.Prequalify(context.Background(), a); err != nil {
		t.Fatal(err)
	}
}
