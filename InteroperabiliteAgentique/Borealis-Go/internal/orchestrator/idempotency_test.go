package orchestrator

import (
	"errors"
	"fmt"
	"testing"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/agbruneau/borealis/internal/audit"
)

func ctxMessage(contextID, text string) *a2asdk.Message {
	m := a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart(text))
	m.ContextID = contextID
	return m
}

// TestIdempotentDelegation : un rejeu (même contextId + params) ne produit qu'un
// seul effet (deliver appelé une fois), retourne le résultat mémorisé, et
// journalise le doublon (FR-18, NFR-05).
func TestIdempotentDelegation(t *testing.T) {
	j := audit.New()
	o := New(nil, nil).WithAudit(j)
	calls := 0
	deliver := func() (a2asdk.SendMessageResult, error) {
		calls++
		return &a2asdk.Message{ID: fmt.Sprintf("r%d", calls)}, nil
	}
	msg := ctxMessage("ctx-42", "dossier-42")

	r1, err := o.idempotent(msg.ContextID, "kyc", msg, deliver)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := o.idempotent(msg.ContextID, "kyc", msg, deliver) // rejeu identique
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Errorf("effet répété : deliver appelé %d fois, want 1", calls)
	}
	if r1 != r2 {
		t.Error("le rejeu ne retourne pas le résultat mémorisé")
	}
	es := j.Entries()
	if len(es) != 2 || es[0].Action != "delegation" || es[1].Action != "delegation.duplicate" {
		t.Errorf("audit = %+v, want [delegation, delegation.duplicate]", es)
	}
	if es[1].KYC != "ctx-42" || es[1].Result != "kyc" {
		t.Errorf("entrée doublon = %+v, want KYC=ctx-42 Result=kyc", es[1])
	}
}

// TestIdempotentDistinctParams : des params différents dans le même dossier ne
// sont PAS dédupliqués (deux effets distincts).
func TestIdempotentDistinctParams(t *testing.T) {
	o := New(nil, nil).WithAudit(audit.New())
	calls := 0
	deliver := func() (a2asdk.SendMessageResult, error) { calls++; return &a2asdk.Message{}, nil }
	o.idempotent("ctx", "kyc", ctxMessage("ctx", "A"), deliver)
	o.idempotent("ctx", "kyc", ctxMessage("ctx", "B"), deliver)
	if calls != 2 {
		t.Errorf("params distincts dédupliqués à tort : calls=%d, want 2", calls)
	}
}

// TestIdempotentNoContextExecutes : sans contextId, aucune identité de dossier
// pour asseoir l'idempotence → exécution à chaque fois, aucune entrée d'audit.
func TestIdempotentNoContextExecutes(t *testing.T) {
	j := audit.New()
	o := New(nil, nil).WithAudit(j)
	calls := 0
	deliver := func() (a2asdk.SendMessageResult, error) { calls++; return &a2asdk.Message{}, nil }
	m := ctxMessage("", "x")
	o.idempotent("", "kyc", m, deliver)
	o.idempotent("", "kyc", m, deliver)
	if calls != 2 {
		t.Errorf("sans contextId, devrait exécuter à chaque fois : calls=%d", calls)
	}
	if j.Len() != 0 {
		t.Errorf("sans contextId, aucune entrée d'audit attendue : %d", j.Len())
	}
}

// TestIdempotentErrorNotCached : un échec n'est pas mémorisé — un rejeu peut
// réussir (idempotence des EFFETS, pas des tentatives ratées).
func TestIdempotentErrorNotCached(t *testing.T) {
	o := New(nil, nil)
	calls := 0
	deliver := func() (a2asdk.SendMessageResult, error) {
		calls++
		if calls == 1 {
			return nil, errors.New("boom")
		}
		return &a2asdk.Message{}, nil
	}
	m := ctxMessage("ctx", "x")
	if _, err := o.idempotent("ctx", "kyc", m, deliver); err == nil {
		t.Fatal("premier appel : échec attendu")
	}
	if _, err := o.idempotent("ctx", "kyc", m, deliver); err != nil {
		t.Fatalf("rejeu après échec devrait réussir : %v", err)
	}
	if calls != 2 {
		t.Errorf("échec mémorisé à tort : calls=%d, want 2", calls)
	}
}

// TestIdempotencyKeyIgnoresMessageID : la clé dépend du contenu et du dossier,
// PAS du messageId (aléatoire par appel) — sinon un rejeu ne serait jamais
// reconnu. Elle distingue en revanche des contextId différents.
func TestIdempotencyKeyIgnoresMessageID(t *testing.T) {
	m1 := ctxMessage("c", "x")
	m2 := ctxMessage("c", "x")
	if m1.ID == m2.ID {
		t.Fatal("prérequis : les messageId devraient être distincts (aléatoires)")
	}
	k1, err := idempotencyKey("s", m1)
	if err != nil {
		t.Fatal(err)
	}
	k2, _ := idempotencyKey("s", m2)
	if k1 != k2 {
		t.Error("la clé dépend du messageId — un rejeu ne serait jamais reconnu")
	}
	k3, _ := idempotencyKey("s", ctxMessage("d", "x")) // autre dossier
	if k1 == k3 {
		t.Error("la clé ignore le contextId — deux dossiers seraient confondus")
	}
	k4, _ := idempotencyKey("autre-skill", m1) // autre compétence
	if k1 == k4 {
		t.Error("la clé ignore la compétence")
	}
}
