package core_test

import (
	"reflect"
	"testing"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/core/coretest"
)

type (
	testReply         struct{ Value string }
	unregisteredReply struct{ Value string }
)

func (testReply) ReplyName() string         { return "core_test.testReply" }
func (unregisteredReply) ReplyName() string { return "core_test.unregisteredReply" }

var (
	testRp         = &testReply{"reply"}
	unregisteredRp = &unregisteredReply{"reply"}
)

func TestDeserializeReply(t *testing.T) {
	type args struct {
		replyName string
		data      []byte
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterReplies(testReply{})
	core.RegisterEvents(testEvent{})

	tests := map[string]struct {
		args    args
		want    core.Reply
		wantErr bool
	}{
		"Success": {
			args: args{
				replyName: testReply{}.ReplyName(),
				data:      getGoldenFileData(t, testReply{}.ReplyName()),
			},
			want:    testRp,
			wantErr: false,
		},
		"SuccessEmpty": {
			args: args{
				replyName: testReply{}.ReplyName(),
				data:      []byte("{}"),
			},
			want:    &testReply{},
			wantErr: false,
		},
		"FailureNoData": {
			args: args{
				replyName: testReply{}.ReplyName(),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureWrongType": {
			args: args{
				replyName: testEvent{}.EventName(),
				data:      []byte("{}"),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureUnregistered": {
			args: args{
				replyName: unregisteredReply{}.ReplyName(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.DeserializeReply(tt.args.replyName, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeReply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeReply() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeReply(t *testing.T) {
	type args struct {
		v core.Reply
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterReplies(testReply{})

	tests := map[string]struct {
		args    args
		want    []byte
		wantErr bool
	}{
		"Success": {
			args:    args{testRp},
			want:    getGoldenFileData(t, testReply{}.ReplyName()),
			wantErr: false,
		},
		"SuccessEmpty": {
			args:    args{testReply{}},
			want:    []byte(`{"Value":""}`),
			wantErr: false,
		},
		"FailureUnregistered": {
			args:    args{unregisteredRp},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.SerializeReply(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeReply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeReply() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
