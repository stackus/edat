package retry

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Backoff is a configurable retry backoff strategy
type Backoff struct {
	maxRetries      int
	initialInterval time.Duration
	maxInterval     time.Duration
	maxElapsed      time.Duration
	multiplier      float64
	randomization   float64
}

// NewBackoff constructs a new Backoff
func NewBackoff(options ...BackoffOption) *Backoff {
	b := &Backoff{
		maxRetries:      DefaultMaxRetries,
		initialInterval: DefaultInitialInterval,
		maxInterval:     DefaultMaxInterval,
		maxElapsed:      0,
		multiplier:      1,
		randomization:   0,
	}

	for _, option := range options {
		option(b)
	}

	// Simple infinite retries check
	if b.maxRetries == 0 && b.maxElapsed == 0 {
		panic("backoff: cannot set both maxRetries and maxElapsed to zero")
	}

	return b
}

// Retry executes a command until it succeeds, encounters an ErrDoNotRetry error, or reaches the backoff strategy limits
func (b Backoff) Retry(ctx context.Context, fn func() error) error {
	tries := 0
	started := time.Now()
	interval := b.initialInterval

	sleepTimer := time.NewTimer(0)
	defer sleepTimer.Stop()

	for {
		err := fn()
		if err == nil {
			return nil
		}

		var doNotRetry *ErrDoNotRetry
		if errors.As(err, &doNotRetry) {
			return doNotRetry.err
		}

		tries++
		if b.maxRetries != 0 && tries >= b.maxRetries {
			return fmt.Errorf("%v: %w", MaxRetriesExceeded, err)
		}

		if b.maxElapsed != 0 && started.Add(b.maxElapsed).After(time.Now()) {
			return fmt.Errorf("%v: %w", MaxElapsedExceeded, err)
		}

		// ensure the timer is stopped and drained
		if !sleepTimer.Stop() {
			select {
			case <-sleepTimer.C:
			default:
			}
		}

		sleepTimer.Reset(interval)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sleepTimer.C:
			interval = b.nextInterval(interval)
		}
	}
}

func (b Backoff) nextInterval(lastInterval time.Duration) time.Duration {
	// Either there isn't any delay or there is not growth
	if b.initialInterval == 0 || lastInterval == b.maxInterval {
		return lastInterval
	}

	next := float64(lastInterval) * b.multiplier

	// RandomizationFactor / Jitter
	if b.randomization != 0 {
		min := next - b.randomization*next
		max := next + b.randomization*next

		next = min + (rand.Float64() * (max - min)) // nolint:gosec
	}

	nextInterval := time.Duration(next)

	// Use the initial interval as the lower bounds
	if nextInterval < b.initialInterval {
		return b.initialInterval
	}

	// limit the upper bounds
	if nextInterval > b.maxInterval {
		return b.maxInterval
	}

	return nextInterval
}
