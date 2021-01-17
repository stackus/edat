package msg_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackus/edat/msg"
)

func TestReceiveMessageFunc_ReceiveMessage(t *testing.T) {
	type args struct {
		ctx     context.Context
		message msg.Message
	}
	tests := map[string]struct {
		f       msg.ReceiveMessageFunc
		args    args
		wantErr bool
	}{
		"Success": {
			f: func(ctx context.Context, m msg.Message) error {
				return nil
			},
			args: args{
				ctx:     context.Background(),
				message: msg.NewMessage([]byte(`{}`)),
			},
			wantErr: false,
		},
		"ReceiverError": {
			f: func(ctx context.Context, m msg.Message) error {
				return fmt.Errorf("receiver-error")
			},
			args: args{
				ctx:     context.Background(),
				message: msg.NewMessage([]byte(`{}`)),
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if err := tt.f.ReceiveMessage(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("ReceiveMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
