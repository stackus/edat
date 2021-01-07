package retry

import (
	"context"
)

type Retryer interface {
	Retry(ctx context.Context, fn func() error) error
}
