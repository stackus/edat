package core

import (
	"fmt"
	"reflect"
)

// Snapshot interface
type Snapshot interface {
	SnapshotName() string
}

// SerializeSnapshot serializes snapshots with a registered marshaller
func SerializeSnapshot(v Snapshot) ([]byte, error) {
	return marshal(v.SnapshotName(), v)
}

// DeserializeSnapshot deserializes the snapshot data using a registered marshaller returning a *Snapshot
func DeserializeSnapshot(snapshotName string, data []byte) (Snapshot, error) {
	snapshot, err := unmarshal(snapshotName, data)
	if err != nil {
		return nil, err
	}

	if snapshot != nil {
		if _, ok := snapshot.(Snapshot); !ok {
			return nil, fmt.Errorf("`%s` was registered but not registered as a snapshot", snapshotName)
		}
	}

	return snapshot.(Snapshot), nil
}

// RegisterSnapshots registers one or more snapshots with a registered marshaller
//
// Register snapshots using any form desired "&MySnapshot{}", "MySnapshot{}", "(*MySnapshot)(nil)"
//
// Snapshots must be registered after first registering a marshaller you wish to use
func RegisterSnapshots(snapshots ...Snapshot) {
	for _, snapshot := range snapshots {
		if v := reflect.ValueOf(snapshot); v.Kind() == reflect.Ptr && v.Pointer() == 0 {
			snapshotName := reflect.Zero(reflect.TypeOf(snapshot).Elem()).Interface().(Snapshot).SnapshotName()
			registerType(snapshotName, snapshot)
		} else {
			registerType(snapshot.SnapshotName(), snapshot)
		}
	}
}
