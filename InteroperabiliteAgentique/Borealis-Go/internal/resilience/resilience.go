// Package resilience fournit un disjoncteur (circuit breaker) et une politique
// de retry bornée pour les appels sortants (MCP/A2A). NFR-04 / FR-17 : détection
// de panne, retries bornés, ouverture du disjoncteur après un seuil d'échecs.
//
// L'horloge est injectable (WithClock) pour un comportement déterministe en test.
package resilience

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State est l'état du disjoncteur.
type State int

const (
	// StateClosed : les appels passent normalement.
	StateClosed State = iota
	// StateOpen : les appels sont rejetés (ErrOpen) jusqu'à la fin du cooldown.
	StateOpen
	// StateHalfOpen : un appel d'essai est autorisé après le cooldown.
	StateHalfOpen
)

// String rend l'état lisible.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// ErrOpen est retournée quand le disjoncteur est ouvert.
var ErrOpen = errors.New("resilience: circuit breaker open")

// RetryPolicy borne les tentatives et l'attente entre elles.
type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
}

// Breaker est un disjoncteur à seuil d'échecs consécutifs.
type Breaker struct {
	mu               sync.Mutex
	threshold        int
	cooldown         time.Duration
	now              func() time.Time
	failures         int
	state            State
	openedAt         time.Time
	halfOpenInFlight bool // un seul appel d'essai autorisé en half-open
}

// NewBreaker construit un disjoncteur (ouvre après threshold échecs ; se réarme
// après cooldown). L'horloge (champ now, déterminisme des tests) se règle en
// écrivant directement le champ depuis le paquet (B8 : aucun appelant externe
// n'a jamais eu besoin de l'option, retirée au profit du champ non exporté).
func NewBreaker(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{threshold: threshold, cooldown: cooldown, now: time.Now, state: StateClosed}
}

// allow indique si un appel est permis, en transitionnant Open→HalfOpen après
// le cooldown.
func (b *Breaker) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == StateOpen && b.now().Sub(b.openedAt) >= b.cooldown {
		b.state = StateHalfOpen
		b.halfOpenInFlight = false
	}
	switch b.state {
	case StateOpen:
		return false
	case StateHalfOpen:
		if b.halfOpenInFlight {
			return false // un essai est déjà en cours
		}
		b.halfOpenInFlight = true
		return true
	default: // StateClosed
		return true
	}
}

// onResult met à jour l'état du disjoncteur selon le résultat de l'appel.
// Une annulation de contexte côté appelant (ctx.Canceled/DeadlineExceeded)
// n'est ni un succès ni un échec du **service distant** : elle ne compte pas
// comme échec, sous peine d'ouvrir le circuit d'un service sain à cause de
// timeouts purement côté appelant (m5).
func (b *Breaker) onResult(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.halfOpenInFlight = false
	switch {
	case err == nil:
		b.failures = 0
		b.state = StateClosed
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		// ni échec ni succès : état inchangé.
	default:
		b.failures++
		if b.failures >= b.threshold {
			b.state = StateOpen
			b.openedAt = b.now()
		}
	}
}

// State retourne l'état courant (en tenant compte de la transition HalfOpen).
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == StateOpen && b.now().Sub(b.openedAt) >= b.cooldown {
		return StateHalfOpen
	}
	return b.state
}

// Execute exécute fn sous le disjoncteur, avec une politique de retry bornée.
// Rejette immédiatement (ErrOpen) si le disjoncteur est ouvert ; un cycle de
// retries épuisé compte pour un seul échec du disjoncteur.
func (b *Breaker) Execute(ctx context.Context, r RetryPolicy, fn func() error) error {
	if !b.allow() {
		return ErrOpen
	}
	err := retry(ctx, r, fn)
	b.onResult(err)
	return err
}

func retry(ctx context.Context, r RetryPolicy, fn func() error) error {
	attempts := r.MaxAttempts
	if attempts < 1 {
		attempts = 1
	}
	var err error
	for i := 1; i <= attempts; i++ {
		if err = ctx.Err(); err != nil {
			return err
		}
		if err = fn(); err == nil {
			return nil
		}
		if i < attempts {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(r.Backoff):
			}
		}
	}
	return err
}
