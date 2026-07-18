package a2a

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
)

func testCard() *a2asdk.AgentCard {
	return &a2asdk.AgentCard{
		Name:                "borealis-kyc",
		Version:             "1.0.0",
		Description:         "Agent KYC de pré-qualification",
		SupportedInterfaces: []*a2asdk.AgentInterface{{URL: "http://kyc.local/a2a", ProtocolBinding: a2asdk.TransportProtocolJSONRPC, ProtocolVersion: "1.0"}},
		Capabilities:        a2asdk.AgentCapabilities{Streaming: true},
		DefaultInputModes:   []string{"application/json"},
		DefaultOutputModes:  []string{"application/json"},
		Skills:              []a2asdk.AgentSkill{{ID: "verify-identity", Name: "Verify Identity", Description: "vérifie l'identité", Tags: []string{"kyc"}}},
	}
}

func newKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func TestSignVerify(t *testing.T) {
	key := newKey(t)
	card := testCard()
	sig, err := Sign(card, key, "kid-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(card, sig, &key.PublicKey); err != nil {
		t.Fatalf("vérification d'une signature valide : %v", err)
	}
}

func TestVerifyRejectsTampered(t *testing.T) {
	key := newKey(t)
	card := testCard()
	sig, _ := Sign(card, key, "kid-1")
	card.Name = "borealis-evil" // altération après signature
	if err := Verify(card, sig, &key.PublicKey); err == nil {
		t.Error("card altérée : la vérification devrait échouer")
	}
}

func TestVerifyRejectsWrongKey(t *testing.T) {
	key, other := newKey(t), newKey(t)
	card := testCard()
	sig, _ := Sign(card, key, "kid-1")
	if err := Verify(card, sig, &other.PublicKey); err == nil {
		t.Error("mauvaise clé : la vérification devrait échouer")
	}
}

func TestVerifyRejectsMalformedSignature(t *testing.T) {
	key := newKey(t)
	card := testCard()
	if err := Verify(card, a2asdk.AgentCardSignature{Protected: "x", Signature: "!!!"}, &key.PublicKey); err == nil {
		t.Error("signature illisible : erreur attendue")
	}
	if err := Verify(card, a2asdk.AgentCardSignature{Protected: "x", Signature: "AAAA"}, &key.PublicKey); err == nil {
		t.Error("longueur invalide : erreur attendue")
	}
}

// TestVerifyRejectsAlgConfusion : un en-tête protégé avec alg="none" (ou un
// alg autre qu'ES256, ou une extension crit) est rejeté avant toute vérif.
func TestVerifyRejectsAlgConfusion(t *testing.T) {
	key := newKey(t)
	card := testCard()
	sig, _ := Sign(card, key, "kid-1")

	for _, protectedJSON := range []string{
		`{"alg":"none"}`,
		`{"alg":"HS256","kid":"kid-1"}`,
		`{"alg":"ES256","crit":["exp"]}`,
	} {
		tampered := sig
		tampered.Protected = b64.EncodeToString([]byte(protectedJSON))
		if err := Verify(card, tampered, &key.PublicKey); err == nil {
			t.Errorf("en-tête %s devrait être rejeté", protectedJSON)
		}
	}
	// En-tête protégé illisible (base64 invalide).
	if err := Verify(card, a2asdk.AgentCardSignature{Protected: "!!!", Signature: sig.Signature}, &key.PublicKey); err == nil {
		t.Error("en-tête illisible devrait être rejeté")
	}
}

func TestCanonicalize(t *testing.T) {
	out, err := Canonicalize(map[string]any{"b": 1, "a": 2, "c": []any{3, "x"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"a":2,"b":1,"c":[3,"x"]}`; string(out) != want {
		t.Errorf("Canonicalize = %s, want %s", out, want)
	}
	// Déterminisme (ordre d'insertion différent → même sortie).
	out2, _ := Canonicalize(map[string]any{"c": []any{3, "x"}, "a": 2, "b": 1})
	if string(out) != string(out2) {
		t.Error("Canonicalize non déterministe")
	}
	// bool + null.
	if out3, _ := Canonicalize(map[string]any{"t": true, "f": false, "n": nil}); string(out3) != `{"f":false,"n":null,"t":true}` {
		t.Errorf("bool/null : %s", out3)
	}
}

func TestCanonicalizeRejectsUnmarshalable(t *testing.T) {
	// Frontière de l'API publique : une valeur non sérialisable → erreur.
	if _, err := Canonicalize(make(chan int)); err == nil {
		t.Error("valeur non sérialisable : erreur attendue")
	}
}

// TestCardPayloadCanonicalBytes verrouille la forme canonique octet-à-octet
// d'une card de référence (B1) : golden capturé sur l'implémentation avant
// refactorisation (canonicalObject/canonicalArray maison), pour prouver que
// le remplacement par un second json.Marshal post-UseNumber ne change pas un
// seul octet — sensible car la signature JWS porte sur ces octets exacts.
func TestCardPayloadCanonicalBytes(t *testing.T) {
	b, err := cardPayload(testCard())
	if err != nil {
		t.Fatal(err)
	}
	want := `{"capabilities":{"streaming":true},"defaultInputModes":["application/json"],"defaultOutputModes":["application/json"],"description":"Agent KYC de pré-qualification","name":"borealis-kyc","skills":[{"description":"vérifie l'identité","id":"verify-identity","name":"Verify Identity","tags":["kyc"]}],"supportedInterfaces":[{"protocolBinding":"JSONRPC","protocolVersion":"1.0","url":"http://kyc.local/a2a"}],"version":"1.0.0"}`
	if string(b) != want {
		t.Errorf("cardPayload a changé d'un octet :\n got=%s\nwant=%s", b, want)
	}
}

func TestKID(t *testing.T) {
	key := newKey(t)
	sig, _ := Sign(testCard(), key, "kid-42")
	if kid, err := KID(sig); err != nil || kid != "kid-42" {
		t.Errorf("KID = %q, err=%v ; want kid-42", kid, err)
	}
	// En-tête protégé illisible.
	if _, err := KID(a2asdk.AgentCardSignature{Protected: "!!!"}); err == nil {
		t.Error("protected illisible : erreur attendue")
	}
	// En-tête protégé décodable mais non-JSON.
	if _, err := KID(a2asdk.AgentCardSignature{Protected: b64.EncodeToString([]byte("pas du json"))}); err == nil {
		t.Error("protected non-JSON : erreur attendue")
	}
}

func TestValidate(t *testing.T) {
	if err := Validate(testCard()); err != nil {
		t.Fatalf("card valide rejetée : %v", err)
	}
	tests := []struct {
		name string
		mut  func(*a2asdk.AgentCard)
	}{
		{"sans name", func(c *a2asdk.AgentCard) { c.Name = "" }},
		{"sans version", func(c *a2asdk.AgentCard) { c.Version = "" }},
		{"sans description", func(c *a2asdk.AgentCard) { c.Description = "" }},
		{"sans interface", func(c *a2asdk.AgentCard) { c.SupportedInterfaces = nil }},
		{"sans skill", func(c *a2asdk.AgentCard) { c.Skills = nil }},
		{"skill sans id", func(c *a2asdk.AgentCard) { c.Skills[0].ID = "" }},
	}
	for _, tt := range tests {
		c := testCard()
		tt.mut(c)
		if err := Validate(c); err == nil {
			t.Errorf("%s : erreur attendue", tt.name)
		}
	}
}
