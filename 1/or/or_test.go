package or

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	t.Run("single channel", func(t *testing.T) {
		start := time.Now()
		<-Or(sig(100 * time.Millisecond))
		if d := time.Since(start); d < 100*time.Millisecond {
			t.Fatalf("expected ~100ms, got %v", d)
		}
	})

	t.Run("closes on shortest channel", func(t *testing.T) {
		start := time.Now()
		<-Or(
			sig(2*time.Hour),
			sig(5*time.Minute),
			sig(100*time.Millisecond),
			sig(1*time.Hour),
			sig(1*time.Minute),
		)
		if d := time.Since(start); d > 200*time.Millisecond {
			t.Fatalf("expected ~100ms, got %v", d)
		}
	})

	t.Run("no channels returns nil", func(t *testing.T) {
		ch := Or()
		if ch != nil {
			t.Fatal("expected nil for no channels")
		}
	})

	t.Run("single channel returns same channel", func(t *testing.T) {
		ch := make(chan interface{})
		defer close(ch)
		result := Or(ch)
		if result != ch {
			t.Fatal("expected same channel back")
		}
	})

	t.Run("already closed channel closes immediately", func(t *testing.T) {
		closed := make(chan interface{})
		close(closed)

		start := time.Now()
		<-Or(closed, sig(1*time.Hour))
		if d := time.Since(start); d > 50*time.Millisecond {
			t.Fatalf("expected immediate close, got %v", d)
		}
	})

	t.Run("multiple channels with same duration", func(t *testing.T) {
		start := time.Now()
		<-Or(sig(50*time.Millisecond), sig(50*time.Millisecond))
		if d := time.Since(start); d > 100*time.Millisecond {
			t.Fatalf("expected ~50ms, got %v", d)
		}
	})
}

func ExampleOr() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// Output will be close to 1 second
	_ = start
}
