package retry

import (
	"context"
)

// Retryer interface
type Retryer interface {
	Retry(ctx context.Context, fn func() error) error
}
