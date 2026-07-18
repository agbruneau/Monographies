package idpmock

import (
	"testing"
	"time"
)

func fixedClock(ts int64) func() time.Time {
	tm := time.Unix(ts, 0)
	return func() time.Time { return tm }
}

func TestIssueVerifyRoundtrip(t *testing.T) {
	idp := New([]byte("dev"), WithClock(fixedClock(1000)))
	tok, err := idp.Issue("orchestrator", "borealis-kyc")
	if err != nil {
		t.Fatal(err)
	}
	c, err := idp.Verify(tok, "borealis-kyc")
	if err != nil {
		t.Fatalf("verify : %v", err)
	}
	if c.Sub != "orchestrator" || c.Aud != "borealis-kyc" {
		t.Errorf("claims = %+v", c)
	}
}

func TestVerifyWrongAudience(t *testing.T) {
	idp := New([]byte("dev"))
	tok, _ := idp.Issue("orch", "agent-a")
	if _, err := idp.Verify(tok, "agent-b"); err == nil {
		t.Error("audience non conforme devrait être rejetée (RFC 8707)")
	}
}

func TestVerifyExpired(t *testing.T) {
	idp := New([]byte("dev"), WithClock(fixedClock(1000)), WithTTL(time.Second))
	tok, _ := idp.Issue("orch", "a")
	idp.now = fixedClock(2000) // au-delà de l'expiration
	if _, err := idp.Verify(tok, "a"); err == nil {
		t.Error("jeton expiré devrait être rejeté")
	}
}

func TestVerifyBadSignature(t *testing.T) {
	idp := New([]byte("dev"))
	tok, _ := idp.Issue("orch", "a")
	if _, err := idp.Verify(tok+"tamper", "a"); err == nil {
		t.Error("signature altérée devrait être rejetée")
	}
	if _, err := New([]byte("autre")).Verify(tok, "a"); err == nil {
		t.Error("secret différent devrait rejeter")
	}
}

func TestVerifyMalformed(t *testing.T) {
	idp := New([]byte("dev"))
	if _, err := idp.Verify("pasunjwt", "a"); err == nil {
		t.Error("format invalide devrait être rejeté")
	}
	// Payload illisible (base64url invalide), mais signature cohérente.
	badPayload := jwtHeader + "." + "!!!"
	if _, err := idp.Verify(badPayload+"."+idp.sign(badPayload), "a"); err == nil {
		t.Error("payload illisible devrait être rejeté")
	}
	// Payload valide base64 mais non-JSON.
	notJSON := jwtHeader + "." + b64.EncodeToString([]byte("pas du json"))
	if _, err := idp.Verify(notJSON+"."+idp.sign(notJSON), "a"); err == nil {
		t.Error("payload non-JSON devrait être rejeté")
	}
}
