package outbox

import (
	"context"
	"time"
)

type MessageStore interface {
	Fetch(ctx context.Context, limit int) ([]Message, error)
	Save(ctx context.Context, message Message) error
	MarkPublished(ctx context.Context, messageIDs []string) error
	PurgePublished(ctx context.Context, olderThan time.Duration) error
}
