package proto

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/automation"
)

type (
	Runnable interface {
		IsAsync() bool
		GetName() string
		GetSource() string
		GetTimeout() uint32
	}
)

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

func FromModule(i *types.Module) *Module {
	if i == nil {
		return nil
	}

	var p = &Module{
		ModuleID:    i.ID,
		NamespaceID: i.NamespaceID,
		Name:        i.Name,
		CreatedAt:   fromTime(i.CreatedAt),
		UpdatedAt:   fromTime(i.UpdatedAt),
		DeletedAt:   fromTime(i.DeletedAt),
		Fields:      make([]*ModuleField, len(i.Fields)),
	}

	for f := range i.Fields {
		p.Fields[f] = &ModuleField{
			FieldID: i.Fields[f].ID,
			Name:    i.Fields[f].Name,
			Kind:    i.Fields[f].Kind,
		}
	}

	return p
}

func FromNamespace(i *types.Namespace) *Namespace {
	if i == nil {
		return nil
	}

	var p = &Namespace{
		NamespaceID: i.ID,
		Name:        i.Name,
		Slug:        i.Slug,
		Enabled:     i.Enabled,
		CreatedAt:   fromTime(i.CreatedAt),
		UpdatedAt:   fromTime(i.UpdatedAt),
		DeletedAt:   fromTime(i.DeletedAt),
	}

	return p
}

func FromAutomationScript(s *automation.Script) *Script {
	return &Script{
		Source:  s.Source,
		Name:    s.Name,
		Timeout: uint32(s.Timeout),
		Async:   s.Async,
	}
}

// Converts time.Time (ptr AND value) to *timestamp.Timestamp
//
// Intentionally ignoring
func fromTime(i interface{}) *timestamp.Timestamp {
	switch t := i.(type) {
	case *time.Time:
		if t == nil {
			return nil
		}
		return &timestamp.Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
	case time.Time:
		return &timestamp.Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
	default:
		return nil
	}
}
