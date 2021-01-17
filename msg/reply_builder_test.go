package msg_test

import (
	"reflect"
	"testing"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/msg"
)

func TestReplyBuilder_Failure(t *testing.T) {
	type fields struct {
		reply   core.Reply
		headers map[string]string
	}
	tests := map[string]struct {
		fields fields
		want   msg.Reply
	}{
		"Success": {
			fields: fields{
				reply: msg.Failure{},
				headers: map[string]string{
					"custom": "value",
				},
			},
			want: msg.NewReply(msg.Failure{}, map[string]string{
				msg.MessageReplyName:    msg.Failure{}.ReplyName(),
				msg.MessageReplyOutcome: msg.ReplyOutcomeFailure,
				"custom":                "value",
			}),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := msg.WithReply(tt.fields.reply).Headers(tt.fields.headers)
			if got := b.Failure(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Failure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplyBuilder_Success(t *testing.T) {
	type fields struct {
		reply   core.Reply
		headers map[string]string
	}
	tests := map[string]struct {
		fields fields
		want   msg.Reply
	}{
		"Success": {
			fields: fields{
				reply: msg.Success{},
				headers: map[string]string{
					"custom": "value",
				},
			},
			want: msg.NewReply(msg.Success{}, map[string]string{
				msg.MessageReplyName:    msg.Success{}.ReplyName(),
				msg.MessageReplyOutcome: msg.ReplyOutcomeSuccess,
				"custom":                "value",
			}),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b := msg.WithReply(tt.fields.reply).Headers(tt.fields.headers)
			if got := b.Success(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Success() = %v, want %v", got, tt.want)
			}
		})
	}
}
