// Package a2a fournit la signature et la validation des Agent Cards A2A :
// canonicalisation JCS (RFC 8785) et signature/vérification JWS ES256 (stdlib),
// par-dessus les types du SDK a2a-go. FR-01 (card conforme au schéma), FR-03
// (signatures vérifiables).
package a2a

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
)

var b64 = base64.RawURLEncoding

// Canonicalize produit une forme JSON canonique (clés d'objet triées, compact,
// nombres préservés) — JCS (RFC 8785) pragmatique, suffisant pour les Agent
// Cards (données majoritairement textuelles). Un aller-retour Marshal→Decode
// (UseNumber, pour préserver les nombres verbatim)→Marshal suffit : le
// marshaling standard de Go trie déjà les clés de map et compacte la sortie,
// donc encoder à nouveau la structure décodée EST la forme canonique (B1 :
// remplace un walker maison équivalent, verrouillé octet-à-octet par
// TestCardPayloadCanonicalBytes).
//
// ponytail: JCS suffisant pour nos cards (signature ET vérification internes,
// sur le struct parsé) ; les cas limites RFC 8785 (nombres ES6, échappement
// HTML de <>&, normalisation Unicode NFC) ne sont pas implémentés → la
// vérification croisée avec des signataires A2A **externes** conformes RFC 8785
// n'est pas garantie (hors périmètre du démonstrateur).
func Canonicalize(payload any) ([]byte, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

// cardPayload retourne la forme canonique de la card **sans ses signatures**
// (on signe la card, puis on y attache la signature).
func cardPayload(card *a2asdk.AgentCard) ([]byte, error) {
	c := *card
	c.Signatures = nil
	return Canonicalize(&c)
}

// Sign calcule une signature JWS ES256 (détachée) de la card, avec le kid donné.
func Sign(card *a2asdk.AgentCard, key *ecdsa.PrivateKey, kid string) (a2asdk.AgentCardSignature, error) {
	payload, err := cardPayload(card)
	if err != nil {
		return a2asdk.AgentCardSignature{}, err
	}
	ph, err := json.Marshal(map[string]any{"alg": SigAlgES256, "kid": kid})
	if err != nil {
		return a2asdk.AgentCardSignature{}, err
	}
	phB64 := b64.EncodeToString(ph)
	signingInput := phB64 + "." + b64.EncodeToString(payload)
	h := sha256.Sum256([]byte(signingInput))
	r, s, err := ecdsa.Sign(rand.Reader, key, h[:])
	if err != nil {
		return a2asdk.AgentCardSignature{}, err
	}
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])
	return a2asdk.AgentCardSignature{Protected: phB64, Signature: b64.EncodeToString(sig)}, nil
}

// checkProtectedHeader authentifie l'en-tête protégé (RFC 7515 §4.1.1) : exige
// alg == "ES256" (rejet explicite de "none" et de toute autre valeur) et
// l'absence d'extension 'crit' non supportée. Ferme la confusion d'algorithme.
func checkProtectedHeader(protected string) error {
	raw, err := b64.DecodeString(protected)
	if err != nil {
		return fmt.Errorf("a2a: en-tête protégé illisible: %w", err)
	}
	var h struct {
		Alg  string   `json:"alg"`
		Crit []string `json:"crit"`
	}
	if err := json.Unmarshal(raw, &h); err != nil {
		return fmt.Errorf("a2a: en-tête protégé invalide: %w", err)
	}
	if h.Alg != SigAlgES256 {
		return fmt.Errorf("a2a: alg %q non supporté (ES256 requis ; 'none' et le placeholder PQ rejetés)", h.Alg)
	}
	if len(h.Crit) > 0 {
		return errors.New("a2a: extension 'crit' non supportée")
	}
	return nil
}

// Verify vérifie une signature JWS ES256 de la card contre la clé publique.
// L'en-tête protégé est authentifié (alg==ES256) avant toute vérification.
func Verify(card *a2asdk.AgentCard, sig a2asdk.AgentCardSignature, pub *ecdsa.PublicKey) error {
	if err := checkProtectedHeader(sig.Protected); err != nil {
		return err
	}
	payload, err := cardPayload(card)
	if err != nil {
		return err
	}
	signingInput := sig.Protected + "." + b64.EncodeToString(payload)
	h := sha256.Sum256([]byte(signingInput))
	raw, err := b64.DecodeString(sig.Signature)
	if err != nil {
		return fmt.Errorf("a2a: signature illisible: %w", err)
	}
	if len(raw) != 64 {
		return errors.New("a2a: longueur de signature invalide")
	}
	r := new(big.Int).SetBytes(raw[:32])
	s := new(big.Int).SetBytes(raw[32:])
	if !ecdsa.Verify(pub, h[:], r, s) {
		return errors.New("a2a: signature invalide")
	}
	return nil
}

// KID extrait le key id (kid) de l'en-tête protégé JWS d'une signature, pour
// permettre à un vérificateur de choisir la clé publique correspondante.
func KID(sig a2asdk.AgentCardSignature) (string, error) {
	raw, err := b64.DecodeString(sig.Protected)
	if err != nil {
		return "", fmt.Errorf("a2a: en-tête protégé illisible: %w", err)
	}
	var h struct {
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(raw, &h); err != nil {
		return "", fmt.Errorf("a2a: en-tête protégé invalide: %w", err)
	}
	return h.Kid, nil
}

// Validate vérifie qu'une AgentCard porte les champs requis par le schéma A2A
// v1.0 (FR-01).
func Validate(card *a2asdk.AgentCard) error {
	switch {
	case card.Name == "":
		return errors.New("a2a: card sans name")
	case card.Version == "":
		return errors.New("a2a: card sans version")
	case card.Description == "":
		return errors.New("a2a: card sans description")
	case len(card.SupportedInterfaces) == 0:
		return errors.New("a2a: card sans supportedInterfaces")
	case len(card.Skills) == 0:
		return errors.New("a2a: card sans skill")
	}
	for i, sk := range card.Skills {
		if sk.ID == "" || sk.Name == "" {
			return fmt.Errorf("a2a: skill %d sans id ou name", i)
		}
	}
	return nil
}
