package core_test

import (
	"reflect"
	"testing"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/core/coretest"
)

type (
	testSnapshot         struct{ Value string }
	unregisteredSnapshot struct{ Value string }
)

func (testSnapshot) SnapshotName() string         { return "core_test.testSnapshot" }
func (unregisteredSnapshot) SnapshotName() string { return "core_test.unregisteredSnapshot" }

var (
	testSs         = &testSnapshot{"snapshot"}
	unregisteredSs = &unregisteredSnapshot{"snapshot"}
)

func TestDeserializeSnapshot(t *testing.T) {
	type args struct {
		snapshotName string
		data         []byte
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterSnapshots(testSnapshot{})
	core.RegisterEvents(testEvent{})

	tests := map[string]struct {
		args    args
		want    core.Snapshot
		wantErr bool
	}{
		"Success": {
			args: args{
				snapshotName: testSnapshot{}.SnapshotName(),
				data:         getGoldenFileData(t, testSnapshot{}.SnapshotName()),
			},
			want:    testSs,
			wantErr: false,
		},
		"SuccessEmpty": {
			args: args{
				snapshotName: testSnapshot{}.SnapshotName(),
				data:         []byte("{}"),
			},
			want:    &testSnapshot{},
			wantErr: false,
		},
		"FailureNoData": {
			args: args{
				snapshotName: testSnapshot{}.SnapshotName(),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureWrongType": {
			args: args{
				snapshotName: testEvent{}.EventName(),
				data:         []byte("{}"),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureUnregistered": {
			args: args{
				snapshotName: unregisteredSnapshot{}.SnapshotName(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.DeserializeSnapshot(tt.args.snapshotName, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeSnapshot() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeSnapshot(t *testing.T) {
	type args struct {
		v core.Snapshot
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterSnapshots(testSnapshot{})

	tests := map[string]struct {
		args    args
		want    []byte
		wantErr bool
	}{
		"Success": {
			args:    args{testSs},
			want:    getGoldenFileData(t, testSnapshot{}.SnapshotName()),
			wantErr: false,
		},
		"SuccessEmpty": {
			args:    args{testSnapshot{}},
			want:    []byte(`{"Value":""}`),
			wantErr: false,
		},
		"FailureUnregistered": {
			args:    args{unregisteredSs},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.SerializeSnapshot(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeSnapshot() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
