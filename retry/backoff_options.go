package retry

import (
	"time"
)

// BackoffOption options for Backoff
type BackoffOption func(*Backoff)

// WithBackoffInitialInterval sets the initialInterval for Backoff
func WithBackoffInitialInterval(initialInterval time.Duration) BackoffOption {
	return func(backoff *Backoff) {
		backoff.initialInterval = initialInterval
	}
}

// WithBackoffMaxRetries sets the maximum number of retries for Backoff
func WithBackoffMaxRetries(maxRetries int) BackoffOption {
	return func(backoff *Backoff) {
		backoff.maxRetries = maxRetries
	}
}

// WithBackoffMaxInterval sets the maximum interval duration for Backoff
func WithBackoffMaxInterval(maxInterval time.Duration) BackoffOption {
	return func(backoff *Backoff) {
		backoff.maxInterval = maxInterval
	}
}

// WithBackoffMultiplier sets the interval duration multipler for Backoff
func WithBackoffMultiplier(multiplier float64) BackoffOption {
	return func(backoff *Backoff) {
		backoff.multiplier = multiplier
	}
}

// WithBackoffRandomizationFactor sets the randomization factor (min and max jigger) for Backoff
func WithBackoffRandomizationFactor(randomizationFactor float64) BackoffOption {
	return func(backoff *Backoff) {
		backoff.randomization = randomizationFactor
	}
}
