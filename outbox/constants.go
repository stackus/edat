package outbox

import (
	"time"

	"github.com/stackus/edat/retry"
)

// Package defaults
const (
	DefaultMessagesPerPolling = 500
	DefaultPollingInterval    = 500 * time.Millisecond
	DefaultPurgeOlderThan     = 60 * time.Second
	DefaultPurgeInterval      = 30 * time.Second

	DefaultMaxRetries               = 100
	DefaultRetryMultiplier          = 1.25
	DefaultRetryRandomizationFactor = 0.33
)

// DefaultRetryer with exponential backoff strategy
var DefaultRetryer = retry.NewExponentialBackoff(
	retry.WithBackoffInitialInterval(DefaultPollingInterval),
	retry.WithBackoffMaxRetries(DefaultMaxRetries),
	retry.WithBackoffMultiplier(DefaultRetryMultiplier),
	retry.WithBackoffRandomizationFactor(DefaultRetryRandomizationFactor),
)
