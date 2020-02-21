package types

import (
	"database/sql/driver"
	"encoding/json"
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

		updated bool
	}
)

func (v RecordValue) IsUpdated() bool {
	return v.updated
}

func (v RecordValue) IsDeleted() bool {
	return v.DeletedAt != nil
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
}

func (set RecordValueSet) GetUpdated() (out RecordValueSet) {
	out = make([]*RecordValue, 0, len(set))
	for i := range set {
		if !set[i].updated {
			continue
		}

		out = append(out, set[i])
	}

	// Append new value
	return out
}

// Replace old value set with new one
//
// Will remove all values that are not set in new set
//
// This satisfies current requirements where record is always
// manipulated as a whole
func (set RecordValueSet) Replace(new RecordValueSet) (out RecordValueSet) {
	if len(new) == 0 {
		return nil
	}

	if len(set) == 0 {
		for i := range new {
			new[i].updated = true
		}

		return new
	}

	out = make([]*RecordValue, 0, len(set)+len(new))
	for s := range set {
		// Mark all old as deleted
		out = append(out, &RecordValue{
			Name:      set[s].Name,
			Value:     set[s].Value,
			Place:     set[s].Place,
			DeletedAt: &time.Time{},
			updated:   true,
		})
	}

	for n := range new {
		if ex := out.Get(new[n].Name, new[n].Place); ex != nil {
			// Reset deleted flag
			ex.DeletedAt = nil

			// Did value change?
			ex.updated = ex.updated || ex.Value != new[n].Value

			ex.Value = new[n].Value
		} else {
			out = append(out, &RecordValue{
				Name:    new[n].Name,
				Value:   new[n].Value,
				Place:   new[n].Place,
				updated: true,
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

func (set RecordValueSet) Len() int           { return len(set) }
func (set RecordValueSet) Swap(i, j int)      { set[i], set[j] = set[j], set[i] }
func (set RecordValueSet) Less(i, j int) bool { return set[i].Place < set[j].Place }
