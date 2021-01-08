package core_test

import (
	"reflect"
	"testing"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/core/coretest"
)

type (
	testCommand         struct{ Value string }
	unregisteredCommand struct{ Value string }
)

func (testCommand) CommandName() string         { return "core_test.testCommand" }
func (unregisteredCommand) CommandName() string { return "core_test.unregisteredCommand" }

var (
	testCmd         = &testCommand{"command"}
	unregisteredCmd = &unregisteredCommand{"command"}
)

func TestDeserializeCommand(t *testing.T) {
	type args struct {
		commandName string
		data        []byte
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterCommands(testCommand{})
	core.RegisterEvents(testEvent{})

	tests := map[string]struct {
		args    args
		want    core.Command
		wantErr bool
	}{
		"Success": {
			args: args{
				commandName: testCommand{}.CommandName(),
				data:        getGoldenFileData(t, testCommand{}.CommandName()),
			},
			want:    testCmd,
			wantErr: false,
		},
		"SuccessEmpty": {
			args: args{
				commandName: testCommand{}.CommandName(),
				data:        []byte("{}"),
			},
			want:    &testCommand{},
			wantErr: false,
		},
		"FailureNoData": {
			args: args{
				commandName: testCommand{}.CommandName(),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureWrongType": {
			args: args{
				commandName: testEvent{}.EventName(),
				data:        []byte("{}"),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureUnregistered": {
			args: args{
				commandName: unregisteredCommand{}.CommandName(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.DeserializeCommand(tt.args.commandName, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeCommand(t *testing.T) {
	type args struct {
		v core.Command
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterCommands(testCommand{})

	tests := map[string]struct {
		args    args
		want    []byte
		wantErr bool
	}{
		"Success": {
			args:    args{testCmd},
			want:    getGoldenFileData(t, testCommand{}.CommandName()),
			wantErr: false,
		},
		"SuccessEmpty": {
			args:    args{testCommand{}},
			want:    []byte(`{"Value":""}`),
			wantErr: false,
		},
		"FailureUnregistered": {
			args:    args{unregisteredCmd},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.SerializeCommand(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeCommand() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
