package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestStateString(t *testing.T) {
	for s, want := range map[State]string{StateClosed: "closed", StateOpen: "open", StateHalfOpen: "half-open", State(9): "unknown"} {
		if s.String() != want {
			t.Errorf("State(%d) = %q, want %q", s, s.String(), want)
		}
	}
}

func TestBreakerLifecycle(t *testing.T) {
	cur := time.Unix(1000, 0)
	b := NewBreaker(3, time.Minute)
	b.now = func() time.Time { return cur }
	fail := func() error { return errors.New("x") }
	ok := func() error { return nil }
	rp := RetryPolicy{MaxAttempts: 1}
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		if err := b.Execute(ctx, rp, fail); err == nil {
			t.Fatalf("appel %d : échec attendu", i)
		}
	}
	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open", b.State())
	}
	// Ouvert : appel rejeté sans exécuter fn.
	if err := b.Execute(ctx, rp, ok); !errors.Is(err, ErrOpen) {
		t.Errorf("attendu ErrOpen, got %v", err)
	}
	// Après cooldown → half-open.
	cur = cur.Add(2 * time.Minute)
	if b.State() != StateHalfOpen {
		t.Errorf("state = %s, want half-open", b.State())
	}
	// Un succès en half-open → closed.
	if err := b.Execute(ctx, rp, ok); err != nil {
		t.Fatalf("succès attendu : %v", err)
	}
	if b.State() != StateClosed {
		t.Errorf("state = %s, want closed", b.State())
	}
}

// TestBreakerIgnoresContextCancellation : une annulation de contexte côté
// appelant n'est ni un succès ni un échec du service distant — elle ne doit
// jamais ouvrir le disjoncteur (m5). Avant correctif, 5 annulations client
// (timeouts appelant) ouvraient le circuit d'un service par ailleurs sain.
func TestBreakerIgnoresContextCancellation(t *testing.T) {
	b := NewBreaker(3, time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rp := RetryPolicy{MaxAttempts: 1}
	fn := func() error { t.Fatal("fn ne devrait pas s'exécuter (ctx déjà annulé)"); return nil }
	for i := 0; i < 5; i++ {
		if err := b.Execute(ctx, rp, fn); !errors.Is(err, context.Canceled) {
			t.Fatalf("appel %d : context.Canceled attendu, got %v", i, err)
		}
	}
	if b.State() != StateClosed {
		t.Errorf("state = %s, want closed (annulations non comptées comme échecs)", b.State())
	}
}

func TestBreakerHalfOpenReopen(t *testing.T) {
	cur := time.Unix(1000, 0)
	b := NewBreaker(1, time.Minute)
	b.now = func() time.Time { return cur }
	fail := func() error { return errors.New("x") }
	ctx := context.Background()

	b.Execute(ctx, RetryPolicy{MaxAttempts: 1}, fail) // ouvre (seuil 1)
	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open", b.State())
	}
	cur = cur.Add(2 * time.Minute) // → half-open
	if b.State() != StateHalfOpen {
		t.Fatalf("state = %s, want half-open", b.State())
	}
	// Échec en half-open → ré-ouvre.
	if err := b.Execute(ctx, RetryPolicy{MaxAttempts: 1}, fail); err == nil {
		t.Fatal("échec attendu")
	}
	if b.State() != StateOpen {
		t.Errorf("state = %s, want open (ré-ouvert)", b.State())
	}
}

func TestBreakerHalfOpenSingleProbe(t *testing.T) {
	cur := time.Unix(1000, 0)
	b := NewBreaker(1, time.Minute)
	b.now = func() time.Time { return cur }
	ctx := context.Background()
	b.Execute(ctx, RetryPolicy{MaxAttempts: 1}, func() error { return errors.New("x") }) // ouvre
	cur = cur.Add(2 * time.Minute)                                                       // cooldown → half-open

	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- b.Execute(ctx, RetryPolicy{MaxAttempts: 1}, func() error {
			close(started)
			<-release
			return nil
		})
	}()
	<-started // l'essai half-open est en cours
	// Un appel concurrent doit être rejeté : un seul essai à la fois.
	if err := b.Execute(ctx, RetryPolicy{MaxAttempts: 1}, func() error { return nil }); !errors.Is(err, ErrOpen) {
		t.Errorf("appel concurrent en half-open : %v, want ErrOpen", err)
	}
	close(release)
	if err := <-done; err != nil {
		t.Fatalf("essai : %v", err)
	}
	// L'essai réussi a refermé le disjoncteur.
	if b.State() != StateClosed {
		t.Errorf("state = %s, want closed", b.State())
	}
}

func TestRetrySucceedsAfterFailures(t *testing.T) {
	b := NewBreaker(3, time.Minute)
	calls := 0
	err := b.Execute(context.Background(), RetryPolicy{MaxAttempts: 3, Backoff: time.Millisecond}, func() error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil || calls != 3 {
		t.Fatalf("err=%v calls=%d, want nil/3", err, calls)
	}
	if b.State() != StateClosed {
		t.Errorf("state = %s, want closed", b.State())
	}
}

func TestRetryExhausted(t *testing.T) {
	b := NewBreaker(5, time.Minute)
	calls := 0
	err := b.Execute(context.Background(), RetryPolicy{MaxAttempts: 3, Backoff: time.Millisecond}, func() error {
		calls++
		return errors.New("perm")
	})
	if err == nil || calls != 3 {
		t.Fatalf("err=%v calls=%d, want err/3", err, calls)
	}
}

func TestRetryZeroAttemptsRunsOnce(t *testing.T) {
	b := NewBreaker(5, time.Minute)
	calls := 0
	_ = b.Execute(context.Background(), RetryPolicy{MaxAttempts: 0}, func() error { calls++; return nil })
	if calls != 1 {
		t.Errorf("calls=%d, want 1", calls)
	}
}

func TestRetryContextCanceledUpfront(t *testing.T) {
	b := NewBreaker(5, time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := b.Execute(ctx, RetryPolicy{MaxAttempts: 3}, func() error { calls++; return nil })
	if !errors.Is(err, context.Canceled) || calls != 0 {
		t.Fatalf("err=%v calls=%d, want Canceled/0", err, calls)
	}
}

func TestRetryContextCanceledDuringBackoff(t *testing.T) {
	b := NewBreaker(5, time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	err := b.Execute(ctx, RetryPolicy{MaxAttempts: 3, Backoff: time.Hour}, func() error {
		calls++
		cancel() // annule pendant l'attente du backoff
		return errors.New("boom")
	})
	if !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatalf("err=%v calls=%d, want Canceled/1", err, calls)
	}
}
