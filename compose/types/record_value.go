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
	}
)

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
		return append(append(set[:i], v), set[i+1:]...)
	}

	// Append new value
	return append(set, v)
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
