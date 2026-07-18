// Package idpmock est un fournisseur d'identité **en processus** (dev) émettant
// des bearer JWT courts à **audience restreinte** (claim aud, RFC 8707) et les
// vérifiant. Clés de dev uniquement, aucun secret réel (FR-26, §11.1).
//
// Limite assumée (mock de dev) : **pas de protection anti-rejeu** (ni jti, ni
// nonce, ni liaison de possession DPoP/mTLS) — un bearer capturé reste
// rejouable pendant sa TTL. Le durcissement anti-rejeu relève du candidat #2
// (module d'identité complet, cf. ADR-0008).
package idpmock

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var b64 = base64.RawURLEncoding

// jwtHeader est l'en-tête JWT fixe {"alg":"HS256","typ":"JWT"} en base64url.
const jwtHeader = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

// Claims est le contenu vérifié d'un jeton.
type Claims struct {
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

// IdP émet et vérifie des bearer JWT HS256 (dev).
type IdP struct {
	secret []byte
	now    func() time.Time
	ttl    time.Duration
}

// Option configure un IdP.
type Option func(*IdP)

// WithClock injecte une horloge (déterminisme des tests).
func WithClock(now func() time.Time) Option { return func(i *IdP) { i.now = now } }

// WithTTL fixe la durée de vie des jetons.
func WithTTL(ttl time.Duration) Option { return func(i *IdP) { i.ttl = ttl } }

// New construit un IdP avec le secret dev fourni (TTL 5 min par défaut).
func New(secret []byte, opts ...Option) *IdP {
	i := &IdP{secret: secret, now: time.Now, ttl: 5 * time.Minute}
	for _, o := range opts {
		o(i)
	}
	return i
}

func (i *IdP) sign(signingInput string) string {
	mac := hmac.New(sha256.New, i.secret)
	mac.Write([]byte(signingInput))
	return b64.EncodeToString(mac.Sum(nil))
}

// Issue émet un jeton pour subject, à destination de audience (claim aud).
func (i *IdP) Issue(subject, audience string) (string, error) {
	now := i.now()
	cb, err := json.Marshal(Claims{Sub: subject, Aud: audience, Iat: now.Unix(), Exp: now.Add(i.ttl).Unix()})
	if err != nil {
		return "", err
	}
	signingInput := jwtHeader + "." + b64.EncodeToString(cb)
	return signingInput + "." + i.sign(signingInput), nil
}

// Verify vérifie la signature, l'expiration et l'audience attendue (RFC 8707) :
// un jeton dont l'aud ne correspond pas au destinataire est rejeté.
func (i *IdP) Verify(token, expectedAudience string) (Claims, error) {
	parts := strings.SplitN(token, ".", 3)
	if len(parts) != 3 {
		return Claims{}, errors.New("idpmock: format de jeton invalide")
	}
	signingInput := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(i.sign(signingInput)), []byte(parts[2])) {
		return Claims{}, errors.New("idpmock: signature invalide")
	}
	cb, err := b64.DecodeString(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("idpmock: payload illisible: %w", err)
	}
	var c Claims
	if err := json.Unmarshal(cb, &c); err != nil {
		return Claims{}, fmt.Errorf("idpmock: payload invalide: %w", err)
	}
	if i.now().Unix() >= c.Exp {
		return Claims{}, errors.New("idpmock: jeton expiré")
	}
	if c.Aud != expectedAudience {
		return Claims{}, fmt.Errorf("idpmock: audience %q != attendue %q (RFC 8707)", c.Aud, expectedAudience)
	}
	return c, nil
}
