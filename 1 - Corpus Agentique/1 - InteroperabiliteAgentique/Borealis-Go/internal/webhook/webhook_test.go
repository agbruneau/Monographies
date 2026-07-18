package webhook

import (
	"errors"
	"testing"
	"time"
)

func fixedClock(ts time.Time) func() time.Time { return func() time.Time { return ts } }

func signerVerifier(t *testing.T, now time.Time) (*Signer, *Verifier) {
	t.Helper()
	secret := []byte("secret-webhook-demonstrateur-32b")
	s := NewSigner(secret)
	s.now = fixedClock(now)
	v := NewVerifier(secret, 5*time.Minute)
	v.now = fixedClock(now)
	return s, v
}

// TestSignVerifyValid : une signature fraîche et intègre est acceptée.
func TestSignVerifyValid(t *testing.T) {
	now := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	s, v := signerVerifier(t, now)
	body := []byte(`{"event":"prequal.completed","id":"A00042"}`)
	sig := s.Sign("nonce-1", body)
	if err := v.Verify(sig, body); err != nil {
		t.Fatalf("signature valide rejetée : %v", err)
	}
}

// TestVerifyTamperedBody : un corps modifié rompt le MAC.
func TestVerifyTamperedBody(t *testing.T) {
	now := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	s, v := signerVerifier(t, now)
	sig := s.Sign("nonce-1", []byte("original"))
	if err := v.Verify(sig, []byte("falsifie")); !errors.Is(err, ErrBadMAC) {
		t.Errorf("corps falsifié : %v, want ErrBadMAC", err)
	}
}

// TestVerifyTamperedMeta : modifier l'horodatage ou le nonce sans re-signer rompt
// le MAC (ils sont couverts par la signature).
func TestVerifyTamperedMeta(t *testing.T) {
	now := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	s, v := signerVerifier(t, now)
	body := []byte("payload")
	sig := s.Sign("nonce-1", body)

	bad := sig
	bad.Nonce = "nonce-2"
	if err := v.Verify(bad, body); !errors.Is(err, ErrBadMAC) {
		t.Errorf("nonce falsifié : %v, want ErrBadMAC", err)
	}
	bad2 := sig
	bad2.Timestamp += 1
	if err := v.Verify(bad2, body); !errors.Is(err, ErrBadMAC) {
		t.Errorf("horodatage falsifié : %v, want ErrBadMAC", err)
	}
}

// TestVerifyWrongSecret : un vérifieur à secret différent rejette.
func TestVerifyWrongSecret(t *testing.T) {
	now := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	s := NewSigner([]byte("secret-emetteur-aaaaaaaaaaaaaaaaa"))
	s.now = fixedClock(now)
	v := NewVerifier([]byte("secret-different-bbbbbbbbbbbbbbbbb"), 5*time.Minute)
	v.now = fixedClock(now)
	sig := s.Sign("n1", []byte("x"))
	if err := v.Verify(sig, []byte("x")); !errors.Is(err, ErrBadMAC) {
		t.Errorf("secret erroné : %v, want ErrBadMAC", err)
	}
}

// TestVerifyStale : une signature hors de la fenêtre de fraîcheur est rejetée
// (anti-rejeu par horodatage).
func TestVerifyStale(t *testing.T) {
	t0 := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	secret := []byte("secret-webhook-demonstrateur-32b")
	s := NewSigner(secret)
	s.now = fixedClock(t0)
	sig := s.Sign("n1", []byte("x"))
	// Vérifieur 10 min plus tard, fenêtre 5 min → périmé.
	v := NewVerifier(secret, 5*time.Minute)
	v.now = fixedClock(t0.Add(10 * time.Minute))
	if err := v.Verify(sig, []byte("x")); !errors.Is(err, ErrStale) {
		t.Errorf("horodatage périmé : %v, want ErrStale", err)
	}
	// Un horodatage dans le futur au-delà de la fenêtre est aussi rejeté.
	vpast := NewVerifier(secret, 5*time.Minute)
	vpast.now = fixedClock(t0.Add(-10 * time.Minute))
	if err := vpast.Verify(sig, []byte("x")); !errors.Is(err, ErrStale) {
		t.Errorf("horodatage trop futur : %v, want ErrStale", err)
	}
}

// TestVerifyReplayAcrossPurgeBoundary : reproduit le trou W1 (rétention indexée
// sur l'heure murale au lieu de l'horodatage du message). Horloge du vérifieur en
// retard sur l'émetteur → le nonce doit rester mémorisé tant que le message est
// frais. Avec le correctif, le rejeu est rejeté ; sans lui, le nonce serait purgé
// prématurément et le rejeu accepté.
func TestVerifyReplayAcrossPurgeBoundary(t *testing.T) {
	secret := []byte("secret-webhook-demonstrateur-32b")
	ts := time.Date(2026, 7, 7, 12, 5, 0, 0, time.UTC) // horodatage du message
	s := NewSigner(secret)
	s.now = fixedClock(ts)
	sig := s.Sign("n1", []byte("x"))

	now := ts.Add(-5 * time.Minute) // vérifieur en retard de 5 min (12:00:00)
	v := NewVerifier(secret, 5*time.Minute)
	v.now = func() time.Time { return now }
	if err := v.Verify(sig, []byte("x")); err != nil {
		t.Fatalf("1re vérif (horloge en retard, message frais) devrait passer : %v", err)
	}
	now = ts.Add(1 * time.Second) // 12:05:01 : 5 min 1 s depuis la 1re vue murale, mais 1 s depuis ts
	if err := v.Verify(sig, []byte("x")); !errors.Is(err, ErrReplay) {
		t.Errorf("rejeu à travers la frontière de purge : %v, want ErrReplay (trou W1)", err)
	}
}

// TestMACNoNULCollision : le préfixe de longueur rend le MAC injectif — un octet
// NUL déplacé entre nonce et corps ne produit plus la même signature (trou W2).
func TestMACNoNULCollision(t *testing.T) {
	secret := []byte("k")
	if mac(secret, 1, "a\x00b", []byte("X")) == mac(secret, 1, "a", []byte("b\x00X")) {
		t.Error("collision de canonicalisation : un NUL déplacé produit le même MAC")
	}
}

// TestVerifyReplayNonce : rejouer la même signature (même nonce) dans la fenêtre
// est rejeté (anti-rejeu par nonce) ; un nonce distinct passe.
func TestVerifyReplayNonce(t *testing.T) {
	now := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	s, v := signerVerifier(t, now)
	body := []byte("x")
	sig := s.Sign("n1", body)
	if err := v.Verify(sig, body); err != nil {
		t.Fatalf("premier envoi rejeté : %v", err)
	}
	if err := v.Verify(sig, body); !errors.Is(err, ErrReplay) {
		t.Errorf("rejeu du même nonce : %v, want ErrReplay", err)
	}
	// Un nonce différent (même instant) passe.
	if err := v.Verify(s.Sign("n2", body), body); err != nil {
		t.Errorf("nonce distinct rejeté : %v", err)
	}
}
