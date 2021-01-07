package retry

func NewConstantBackoff(options ...BackoffOption) *Backoff {
	// Set some defaults for constant backoff
	conOptions := append(options, WithBackoffMultiplier(1), WithBackoffRandomizationFactor(0))

	b := NewBackoff(conOptions...)

	return b
}
