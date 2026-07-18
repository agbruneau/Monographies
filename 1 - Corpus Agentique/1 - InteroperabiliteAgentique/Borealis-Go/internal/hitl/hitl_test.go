package hitl

import (
	"errors"
	"testing"

	"github.com/agbruneau/borealis/internal/audit"
)

func TestProposeThenReleaseFlow(t *testing.T) {
	j := audit.New()
	b := New(j)

	b.Propose(Proposal{ID: "p1", Subject: "cust:h", Summary: "AML élevé"})
	// Sans release, AUCUN message en file de commande (FR-16).
	if len(b.Commands()) != 0 {
		t.Fatal("file de commande non vide sans release")
	}
	if b.PendingProposals() != 1 {
		t.Errorf("propositions en attente = %d, want 1", b.PendingProposals())
	}

	// Release humaine → une commande, journalisée.
	if err := b.Release("p1", "approve", "human:credit-lead"); err != nil {
		t.Fatalf("release : %v", err)
	}
	cmds := b.Commands()
	if len(cmds) != 1 || cmds[0].Decision != "approve" || cmds[0].Human != "human:credit-lead" {
		t.Errorf("file de commande = %+v", cmds)
	}
	if b.PendingProposals() != 0 {
		t.Error("proposition toujours en attente après release")
	}
	// Journal : propose + release.
	if j.Len() != 2 {
		t.Errorf("entrées d'audit = %d, want 2 (propose + release)", j.Len())
	}
	// La release doit être présente et attribuée à l'humain.
	if last := j.Entries()[1]; last.Action != "hitl.release" || last.OwnerKYH != "human:credit-lead" {
		t.Errorf("release non journalisée correctement : %+v", last)
	}
}

func TestReleaseErrors(t *testing.T) {
	b := New(audit.New())
	if err := b.Release("inconnue", "approve", "h"); !errors.Is(err, ErrUnknownProposal) {
		t.Errorf("release d'une proposition inconnue : %v", err)
	}
	b.Propose(Proposal{ID: "p1", Subject: "s"})
	if err := b.Release("p1", "reject", "h"); err != nil {
		t.Fatal(err)
	}
	// Double release refusée (pas de commande en double).
	if err := b.Release("p1", "approve", "h"); !errors.Is(err, ErrAlreadyReleased) {
		t.Errorf("double release : %v", err)
	}
	if len(b.Commands()) != 1 {
		t.Errorf("commandes = %d, want 1 (pas de double)", len(b.Commands()))
	}
}

// TestReleaseRejectsInvalidDecision : une décision hors {approve, reject} (ex.
// faute de frappe « aprove ») doit être refusée — jamais entrer en file de
// commande ni au journal comme release valide (m11).
func TestReleaseRejectsInvalidDecision(t *testing.T) {
	j := audit.New()
	b := New(j)
	b.Propose(Proposal{ID: "p1", Subject: "s"})
	if err := b.Release("p1", "aprove", "h"); !errors.Is(err, ErrInvalidDecision) {
		t.Errorf("décision invalide : ErrInvalidDecision attendu, got %v", err)
	}
	if len(b.Commands()) != 0 {
		t.Errorf("aucune commande ne devrait être journalisée : %v", b.Commands())
	}
	if b.PendingProposals() != 1 {
		t.Error("la proposition devrait rester en attente après un rejet de validation")
	}
	if j.Len() != 1 { // seul le propose initial, pas de release
		t.Errorf("entrées d'audit = %d, want 1 (propose seul)", j.Len())
	}
}
