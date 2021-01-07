package saga

import (
	"context"
)

type InstanceStore interface {
	Find(ctx context.Context, sagaName, sagaID string) (*Instance, error)
	Save(ctx context.Context, sagaInstance *Instance) error
	Update(ctx context.Context, sagaInstance *Instance) error
}
