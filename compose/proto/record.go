package proto

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func ToRecord(i *Record) *types.Record {
	if i == nil {
		return nil
	}

	var t = &types.Record{
		ID:          i.RecordID,
		ModuleID:    i.ModuleID,
		NamespaceID: i.NamespaceID,
		OwnedBy:     i.OwnedBy,
		CreatedBy:   i.CreatedBy,
		UpdatedBy:   i.UpdatedBy,
		DeletedBy:   i.DeletedBy,
		CreatedAt:   toTime(i.CreatedAt),
		UpdatedAt:   toTimePtr(i.UpdatedAt),
		DeletedAt:   toTimePtr(i.DeletedAt),
		Values:      make([]*types.RecordValue, len(i.Values)),
	}

	for v := range i.Values {
		t.Values[v] = &types.RecordValue{
			Value: i.Values[v].Value,
			Name:  i.Values[v].Name,
		}
	}

	return t
}

// Converts time.Time (ptr AND value) to *timestamp.Timestamp
//
// Intentionally ignoring
func toTime(ts *timestamp.Timestamp) time.Time {
	var t = time.Time{}
	return t.Add(time.Duration(ts.GetNanos()) + (time.Duration(ts.GetSeconds()) * time.Second))
}

func toTimePtr(ts *timestamp.Timestamp) *time.Time {
	var t = toTime(ts)
	return &t
}

func FromRecord(i *types.Record) *Record {
	if i == nil {
		return nil
	}

	var p = &Record{
		RecordID:    i.ID,
		ModuleID:    i.ModuleID,
		NamespaceID: i.NamespaceID,
		OwnedBy:     i.OwnedBy,
		CreatedBy:   i.CreatedBy,
		UpdatedBy:   i.UpdatedBy,
		DeletedBy:   i.DeletedBy,
		CreatedAt:   fromTime(i.CreatedAt),
		UpdatedAt:   fromTime(i.UpdatedAt),
		DeletedAt:   fromTime(i.DeletedAt),
		Values:      make([]*RecordValue, len(i.Values)),
	}

	for v := range i.Values {
		p.Values[v] = &RecordValue{
			Value: i.Values[v].Value,
			Name:  i.Values[v].Name,
		}
	}

	return p
}
