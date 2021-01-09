package retry

// NewExponentialBackoff constructs a Backoff strategy with a exponential backoff retry rate
func NewExponentialBackoff(options ...BackoffOption) *Backoff {
	// Set some defaults for exponential backoff
	expOptions := append([]BackoffOption{WithBackoffMultiplier(DefaultMultiplier)}, options...)

	b := NewBackoff(expOptions...)

	if b.multiplier <= 1.0 {
		panic("backoff: multiplier must be greater than 1 for exponential backoffs")
	}

	return b
}
