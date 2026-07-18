// Package webhook signe et vérifie des notifications sortantes (push webhook,
// FR-09, US-10). La signature HMAC-SHA256 couvre horodatage + nonce + corps ;
// la vérification rejette les rejeux par **double garde** : fraîcheur de
// l'horodatage (fenêtre) ET unicité du nonce dans la fenêtre. Illustratif :
// données synthétiques, aucun secret réel.
package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Erreurs de vérification.
var (
	ErrBadMAC = errors.New("webhook: signature invalide")      // MAC incohérent (corps/méta/secret)
	ErrStale  = errors.New("webhook: horodatage hors fenêtre") // trop ancien ou trop futur
	ErrReplay = errors.New("webhook: nonce déjà vu (rejeu)")   // nonce réutilisé dans la fenêtre
)

// Signature accompagne un corps de webhook : horodatage Unix, nonce unique, et
// MAC hexadécimal. Tous trois sont transmis avec la notification.
type Signature struct {
	Timestamp int64
	Nonce     string
	MAC       string
}

// mac calcule le HMAC-SHA256 hex sur (ts, longueur du nonce, nonce, body). Le
// **préfixe de longueur** du nonce fixe la frontière nonce|corps : l'encodage
// est injectif même si le nonce ou le corps contient un octet NUL (un simple
// séparateur \x00 serait ambigu — un NUL déplacé du corps vers le nonce
// produirait le même MAC, permettant d'altérer le corps sous signature valide).
func mac(secret []byte, ts int64, nonce string, body []byte) string {
	h := hmac.New(sha256.New, secret)
	fmt.Fprintf(h, "%d\x00%d\x00", ts, len(nonce))
	h.Write([]byte(nonce))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

// Signer produit des signatures. Le nonce est fourni par l'appelant (identifiant
// de requête unique) : la bibliothèque reste déterministe (pas d'aléa interne).
type Signer struct {
	secret []byte
	now    func() time.Time
}

// NewSigner construit un signeur. L'horloge (champ now, déterminisme des
// tests) se règle en écrivant directement le champ depuis le paquet (B8 :
// aucun appelant externe n'a jamais eu besoin de l'option).
func NewSigner(secret []byte) *Signer {
	return &Signer{secret: secret, now: time.Now}
}

// Sign signe body avec le nonce donné à l'instant courant.
func (s *Signer) Sign(nonce string, body []byte) Signature {
	ts := s.now().Unix()
	return Signature{Timestamp: ts, Nonce: nonce, MAC: mac(s.secret, ts, nonce, body)}
}

// Verifier vérifie les signatures et refuse les rejeux. Il mémorise les nonces
// vus dans la fenêtre (purge des périmés à chaque vérification).
type Verifier struct {
	secret []byte
	now    func() time.Time
	window time.Duration
	mu     sync.Mutex
	seen   map[string]time.Time // nonce → instant de première vue
}

// NewVerifier construit un vérifieur à fenêtre de fraîcheur donnée. L'horloge
// (champ now, déterminisme des tests) se règle en écrivant directement le
// champ depuis le paquet (B8 : aucun appelant externe n'a jamais eu besoin de
// l'option).
func NewVerifier(secret []byte, window time.Duration) *Verifier {
	return &Verifier{secret: secret, now: time.Now, window: window, seen: make(map[string]time.Time)}
}

// Verify vérifie le MAC, la fraîcheur puis l'unicité du nonce. Ordre voulu : un
// MAC invalide est rejeté AVANT toute mémorisation (un attaquant ne peut pas
// polluer le cache de nonces avec des signatures forgées).
func (v *Verifier) Verify(sig Signature, body []byte) error {
	want := mac(v.secret, sig.Timestamp, sig.Nonce, body)
	if !hmac.Equal([]byte(want), []byte(sig.MAC)) {
		return ErrBadMAC
	}
	now := v.now()
	age := now.Sub(time.Unix(sig.Timestamp, 0))
	if age < 0 {
		age = -age // horodatage futur : même fenêtre, symétrique
	}
	if age > v.window {
		return ErrStale
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	// Rétention et purge indexées sur l'HORODATAGE DU MESSAGE, pas sur l'heure
	// murale de première vue : un nonce reste mémorisé exactement aussi longtemps
	// que son message peut être jugé frais (|now − ts| ≤ window). Aligner les deux
	// références ferme le trou de rejeu quand l'horloge du vérifieur est en retard
	// sur celle de l'émetteur (sinon le nonce serait purgé avant péremption).
	//
	// ponytail: purge O(n) sur tout v.seen à CHAQUE Verify — acceptable au débit
	// du démonstrateur (n borné par le volume de webhooks dans une fenêtre de
	// quelques minutes). Si le débit devient significatif : une file/deque
	// triée par première-vue évite le balayage complet à chaque appel.
	msgTime := time.Unix(sig.Timestamp, 0)
	for n, seenAt := range v.seen {
		if now.Sub(seenAt) > v.window {
			delete(v.seen, n)
		}
	}
	if _, dup := v.seen[sig.Nonce]; dup {
		return ErrReplay
	}
	v.seen[sig.Nonce] = msgTime
	return nil
}
