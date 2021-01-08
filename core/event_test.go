package core_test

import (
	"reflect"
	"testing"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/core/coretest"
)

type (
	testEvent         struct{ Value string }
	unregisteredEvent struct{ Value string }
)

func (testEvent) EventName() string         { return "core_test.testEvent" }
func (unregisteredEvent) EventName() string { return "core_test.unregisteredEvent" }

var (
	testEvt         = &testEvent{"event"}
	unregisteredEvt = &unregisteredEvent{"event"}
)

func TestDeserializeEvent(t *testing.T) {
	type args struct {
		eventName string
		data      []byte
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterEvents(testEvent{})
	core.RegisterCommands(testCommand{})

	tests := map[string]struct {
		args    args
		want    core.Event
		wantErr bool
	}{
		"Success": {
			args: args{
				eventName: testEvent{}.EventName(),
				data:      getGoldenFileData(t, testEvent{}.EventName()),
			},
			want:    testEvt,
			wantErr: false,
		},
		"SuccessEmpty": {
			args: args{
				eventName: testEvent{}.EventName(),
				data:      []byte("{}"),
			},
			want:    &testEvent{},
			wantErr: false,
		},
		"FailureNoData": {
			args: args{
				eventName: testEvent{}.EventName(),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureWrongType": {
			args: args{
				eventName: testCommand{}.CommandName(),
				data:      []byte("{}"),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureUnregistered": {
			args: args{
				eventName: unregisteredEvent{}.EventName(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.DeserializeEvent(tt.args.eventName, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeEvent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeEvent(t *testing.T) {
	type args struct {
		v core.Event
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterEvents(testEvent{})

	tests := map[string]struct {
		args    args
		want    []byte
		wantErr bool
	}{
		"Success": {
			args:    args{testEvt},
			want:    getGoldenFileData(t, testEvent{}.EventName()),
			wantErr: false,
		},
		"SuccessEmpty": {
			args:    args{testEvent{}},
			want:    []byte(`{"Value":""}`),
			wantErr: false,
		},
		"FailureUnregistered": {
			args:    args{unregisteredEvt},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.SerializeEvent(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeEvent() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
