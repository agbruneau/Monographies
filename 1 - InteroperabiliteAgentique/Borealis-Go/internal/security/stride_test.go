package security_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/idpmock"
	"github.com/agbruneau/borealis/pkg/a2a"
)

// TestSTRIDESpoofing : une AgentCard falsifiée après signature est rejetée par
// la vérification (ES256/JCS). Un pair ne peut pas usurper l'identité d'un agent.
func TestSTRIDESpoofing(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	card := &a2asdk.AgentCard{
		Name:                "borealis-kyc",
		Version:             "1.0.0",
		Description:         "Agent KYC de pré-qualification",
		SupportedInterfaces: []*a2asdk.AgentInterface{{URL: "http://kyc.local/a2a", ProtocolBinding: a2asdk.TransportProtocolJSONRPC, ProtocolVersion: "1.0"}},
		DefaultInputModes:   []string{"application/json"},
		DefaultOutputModes:  []string{"application/json"},
	}
	sig, err := a2a.Sign(card, key, "kid-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := a2a.Verify(card, sig, &key.PublicKey); err != nil {
		t.Fatalf("signature intègre rejetée : %v", err)
	}
	card.Name = "borealis-evil" // usurpation après signature
	if err := a2a.Verify(card, sig, &key.PublicKey); err == nil {
		t.Error("spoofing : une card falsifiée devrait être rejetée")
	}
}

// TestSTRIDETampering : la réécriture d'un log d'audit **exporté** (hors du
// journal) rompt la chaîne de hachage et reste détectable par un auditeur tiers.
func TestSTRIDETampering(t *testing.T) {
	ts := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	j := audit.New(audit.WithClock(func() time.Time { return ts }))
	for i := 0; i < 3; i++ {
		j.Append(audit.Record{KYA: "agent:kyc", KYC: "cust:h", Action: "decision", Result: "declined", Version: "m3"})
	}
	exported := j.Entries() // copie exportée vers un stockage WORM
	// Ancrage engagé à part au moment de l'export (compte + hash de tête).
	wantLen, wantTip := j.Len(), j.Tip()
	if err := audit.VerifyExport(exported, wantLen, wantTip); err != nil {
		t.Fatalf("chaîne exportée intègre rejetée : %v", err)
	}

	// Réécriture d'un champ en milieu de chaîne → rompue (relecture nue suffit).
	tampered := append([]audit.Entry(nil), exported...)
	tampered[1].Result = "octroi_falsifie"
	if err := audit.VerifyEntries(tampered); err == nil {
		t.Error("tampering : la réécriture d'un log exporté devrait rompre la chaîne")
	}

	// Troncature de fin : la relecture nue NE la détecte pas (chaîne plus courte
	// mais valide), l'ancrage OUI (compte/tête ne correspondent plus).
	truncated := exported[:len(exported)-1]
	if err := audit.VerifyEntries(truncated); err != nil {
		t.Error("relecture nue : une troncature de fin reste une chaîne valide (limite connue)")
	}
	if err := audit.VerifyExport(truncated, wantLen, wantTip); err == nil {
		t.Error("tampering : l'ancrage devrait détecter la troncature de fin")
	}

	// Effacement total : détecté par l'ancrage (compte 0 != attendu).
	if err := audit.VerifyExport(nil, wantLen, wantTip); err == nil {
		t.Error("tampering : l'ancrage devrait détecter un log entièrement effacé")
	}

	// Re-forge cohérent (même longueur, contenu réécrit puis re-scellé de bout en
	// bout) : la chaîne est interne-valide, seul l'ancrage de tête la démasque.
	reforged := audit.New(audit.WithClock(func() time.Time { return ts }))
	for i := 0; i < 3; i++ {
		reforged.Append(audit.Record{KYA: "agent:kyc", KYC: "cust:h", Action: "decision", Result: "octroi_falsifie", Version: "m3"})
	}
	ref := reforged.Entries()
	if err := audit.VerifyEntries(ref); err != nil {
		t.Fatalf("re-forge doit être interne-valide : %v", err)
	}
	if err := audit.VerifyExport(ref, wantLen, wantTip); err == nil {
		t.Error("tampering : l'ancrage de tête devrait démasquer un re-forge cohérent de même longueur")
	}
}

// TestSTRIDEElevation : un jeton émis pour une audience donnée est refusé par
// tout autre destinataire (aud étranger → rejet, équivalent 401 ; RFC 8707).
func TestSTRIDEElevation(t *testing.T) {
	idp := idpmock.New([]byte("secret-demo-32-bytes-minimum-xx!"))
	tok, err := idp.Issue("agent:orchestrator", "identity-hub")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := idp.Verify(tok, "identity-hub"); err != nil {
		t.Fatalf("jeton d'audience correcte rejeté : %v", err)
	}
	if _, err := idp.Verify(tok, "scoring-engine"); err == nil {
		t.Error("elevation : un jeton d'audience étrangère devrait être rejeté (RFC 8707)")
	}
}
