package retry

import (
	"time"
)

type BackoffOption func(*Backoff)

func WithBackoffInitialInterval(initialInterval time.Duration) BackoffOption {
	return func(backoff *Backoff) {
		backoff.initialInterval = initialInterval
	}
}

func WithBackoffMaxRetries(maxRetries int) BackoffOption {
	return func(backoff *Backoff) {
		backoff.maxRetries = maxRetries
	}
}

func WithBackoffMaxInterval(maxInterval time.Duration) BackoffOption {
	return func(backoff *Backoff) {
		backoff.maxInterval = maxInterval
	}
}

func WithBackoffMultiplier(multiplier float64) BackoffOption {
	return func(backoff *Backoff) {
		backoff.multiplier = multiplier
	}
}

func WithBackoffRandomizationFactor(randomizationFactor float64) BackoffOption {
	return func(backoff *Backoff) {
		backoff.randomization = randomizationFactor
	}
}
