package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type (
	// RecordValue is a stored row in the `record_value` table
	RecordValue struct {
		RecordID  uint64     `db:"record_id"  json:"-"`
		Name      string     `db:"name"       json:"name"`
		Value     string     `db:"value"      json:"value,omitempty"`
		Ref       uint64     `db:"ref"        json:"-"`
		Place     uint       `db:"place"      json:"-"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`

		Updated  bool   `db:"-"      json:"-"`
		OldValue string `db:"-"      json:"-"`
	}
)

func (v RecordValue) IsUpdated() bool {
	return v.Updated
}

func (v RecordValue) IsDeleted() bool {
	return v.DeletedAt != nil
}

func (v RecordValue) Clone() *RecordValue {
	return &RecordValue{
		RecordID:  v.RecordID,
		Name:      v.Name,
		Value:     v.Value,
		Ref:       v.Ref,
		Place:     v.Place,
		DeletedAt: v.DeletedAt,
		Updated:   v.Updated,
		OldValue:  v.OldValue,
	}
}

func (set RecordValueSet) FilterByName(name string) (vv RecordValueSet) {
	for i := range set {
		if set[i].Name == name {
			vv = append(vv, set[i])
		}
	}

	return
}

func (set RecordValueSet) FilterByRecordID(recordID uint64) (vv RecordValueSet) {
	// Make sure we never return nil
	vv = RecordValueSet{}

	for i := range set {
		if set[i].RecordID == recordID {
			vv = append(vv, set[i])
		}
	}

	return
}

// Set updates existing value or creates a new one
func (set RecordValueSet) Set(v *RecordValue) RecordValueSet {
	for i := range set {
		if set[i].Name != v.Name {
			continue
		}
		if set[i].Place != v.Place {
			continue
		}

		//  Update existing entry
		set[i] = v
		return set
	}

	// Append new value
	return append(set, v)
}

// Has value set?
func (set RecordValueSet) Get(name string, place uint) *RecordValue {
	for i := range set {
		if set[i].Name != name {
			continue
		}
		if set[i].Place != place {
			continue
		}

		return set[i]
	}

	return nil
}

// Has value set?
func (set RecordValueSet) Has(name string, place uint) bool {
	for i := range set {
		if set[i].Name != name {
			continue
		}
		if set[i].Place != place {
			continue
		}

		return true
	}

	return false
}

func (set RecordValueSet) SetUpdatedFlag(updated bool) {
	for i := range set {
		set[i].Updated = updated
	}
}

func (set RecordValueSet) GetUpdated() (out RecordValueSet) {
	out = make([]*RecordValue, 0, len(set))
	for i := range set {
		if !set[i].Updated {
			continue
		}

		out = append(out, set[i])
	}

	// Append new value
	return out
}

func (set RecordValueSet) GetClean() (out RecordValueSet) {
	out = make([]*RecordValue, 0, len(set))
	for s := range set {
		if set[s].DeletedAt != nil {
			continue
		}

		out = append(out, &RecordValue{
			Name:  set[s].Name,
			Value: set[s].Value,
			Ref:   set[s].Ref,
			Place: set[s].Place,
		})
	}

	return out
}

// Merge merges old value set with new one and expects unchanged values to be in the new set
//
// This satisfies current requirements where record values are always
// manipulated as a whole (not partial)
//
func (set RecordValueSet) Merge(new RecordValueSet) (out RecordValueSet) {
	if len(set) == 0 {
		// Empty set, copy all new values and return them
		for i := range new {
			new[i].Updated = true
		}

		return new
	}

	out = make([]*RecordValue, 0)
	for s := range set {
		// Mark all old as deleted
		out = append(out, &RecordValue{
			Name:      set[s].Name,
			Value:     set[s].Value,
			Ref:       set[s].Ref,
			Place:     set[s].Place,
			DeletedAt: &time.Time{},
			Updated:   true,
			OldValue:  set[s].Value,
		})
	}

	for n := range new {
		if ex := out.Get(new[n].Name, new[n].Place); ex != nil {
			// Reset deleted flag
			ex.DeletedAt = new[n].DeletedAt

			if ex.OldValue == new[n].Value {
				ex.Updated = false
			} else if !ex.Updated {
				// Did value change?
				ex.Updated = ex.Value != new[n].Value
				ex.OldValue = ex.Value
			}

			ex.Value = new[n].Value
			ex.Ref = new[n].Ref
		} else {
			// Value not previously set, make new
			out = append(out, &RecordValue{
				Name:    new[n].Name,
				Value:   new[n].Value,
				Ref:     new[n].Ref,
				Place:   new[n].Place,
				Updated: true,

				// verbose & explicit for clarity
				OldValue:  "",
				DeletedAt: nil,
			})
		}
	}

	return out
}

func (set *RecordValueSet) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*set = RecordValueSet{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), set); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into RecordValueSet", value)
		}
	}

	return nil
}

func (set RecordValueSet) Value() (driver.Value, error) {
	return json.Marshal(set)
}

// Simple RVS as string output utility fn that
// can help with debugging
func (set RecordValueSet) String() (o string) {
	if set == nil {
		return "<RecordValueSet = nil>"
	}

	is := func(in interface{}) string {
		switch in := in.(type) {
		case bool:
			if in {
				return "✔"
			}
		case *time.Time:
			if in != nil {
				return "✔"
			}
		}

		return "x"
	}

	o += "━━━━━━━━━━━┳━━━━┳━━━┳━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	o += "name       ┃ ## ┃ u ┃ d ┃ value                     ┃ ref                  ┃ old value  \n"
	o += "━━━━━━━━━━━╋━━━━╋━━━╋━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	const tpl = "%-10s ┃ %2d ┃ %s ┃ %s ┃ %-25s ┃ %-20d ┃ %-25s\n"
	for _, v := range set {
		if v == nil {
			o += "<--------> ┃ -- ┃ - ┃ - ┃ <------------------------> ┃ <------------------> ┃ <------------------> \n"
			continue
		}

		o += fmt.Sprintf(
			tpl,
			v.Name,
			v.Place,
			is(v.Updated),
			is(v.DeletedAt),
			v.Value,
			v.Ref,
			v.OldValue,
		)
	}
	o += "━━━━━━━━━━━┻━━━━┻━━━┻━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"

	return o
}

func (set RecordValueSet) Len() int           { return len(set) }
func (set RecordValueSet) Swap(i, j int)      { set[i], set[j] = set[j], set[i] }
func (set RecordValueSet) Less(i, j int) bool { return set[i].Place < set[j].Place }
