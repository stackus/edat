package retry

import (
	"time"
)

// Retry constants
const (
	DefaultInitialInterval = 500 * time.Millisecond
	DefaultMaxRetries      = 100
	DefaultMaxInterval     = 60 * time.Second

	DefaultMultiplier          float64 = 1.5
	DefaultRandomizationFactor float64 = 0.5
)
