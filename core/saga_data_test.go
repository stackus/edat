package core_test

import (
	"reflect"
	"testing"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/core/coretest"
)

type (
	testSagaData         struct{ Value string }
	unregisteredSagaData struct{ Value string }
)

func (testSagaData) SagaDataName() string         { return "core_test.testSagaData" }
func (unregisteredSagaData) SagaDataName() string { return "core_test.unregisteredSagaData" }

var (
	testSd         = &testSagaData{"sagaData"}
	unregisteredSd = &unregisteredSagaData{"sagaData"}
)

func TestDeserializeSagaData(t *testing.T) {
	type args struct {
		sagaDataName string
		data         []byte
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterSagaData(testSagaData{})
	core.RegisterEvents(testEvent{})

	tests := map[string]struct {
		args    args
		want    core.SagaData
		wantErr bool
	}{
		"Success": {
			args: args{
				sagaDataName: testSagaData{}.SagaDataName(),
				data:         getGoldenFileData(t, testSagaData{}.SagaDataName()),
			},
			want:    testSd,
			wantErr: false,
		},
		"SuccessEmpty": {
			args: args{
				sagaDataName: testSagaData{}.SagaDataName(),
				data:         []byte("{}"),
			},
			want:    &testSagaData{},
			wantErr: false,
		},
		"FailureNoData": {
			args: args{
				sagaDataName: testSagaData{}.SagaDataName(),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureWrongType": {
			args: args{
				sagaDataName: testEvent{}.EventName(),
				data:         []byte("{}"),
			},
			want:    nil,
			wantErr: true,
		},
		"FailureUnregistered": {
			args: args{
				sagaDataName: unregisteredSagaData{}.SagaDataName(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.DeserializeSagaData(tt.args.sagaDataName, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeSagaData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeSagaData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeSagaData(t *testing.T) {
	type args struct {
		v core.SagaData
	}

	testMarshaller := coretest.NewTestMarshaller()
	core.RegisterDefaultMarshaller(testMarshaller)
	core.RegisterSagaData(testSagaData{})

	tests := map[string]struct {
		args    args
		want    []byte
		wantErr bool
	}{
		"Success": {
			args:    args{testSd},
			want:    getGoldenFileData(t, testSagaData{}.SagaDataName()),
			wantErr: false,
		},
		"SuccessEmpty": {
			args:    args{testSagaData{}},
			want:    []byte(`{"Value":""}`),
			wantErr: false,
		},
		"FailureUnregistered": {
			args:    args{unregisteredSd},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := core.SerializeSagaData(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeSagaData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeSagaData() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
