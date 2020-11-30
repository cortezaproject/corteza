package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type (
	// RecordValue is a stored row in the `record_value` table
	RecordValue struct {
		RecordID  uint64     `json:"-"`
		Name      string     `json:"name"`
		Value     string     `json:"value,omitempty"`
		Ref       uint64     `json:"-"`
		Place     uint       `json:"-"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`

		Updated  bool   `json:"-"`
		OldValue string `json:"-"`
	}

	RecordValueFilter struct {
		RecordID []uint64
		Deleted  filter.State `json:"deleted"`
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

func (set RecordValueSet) Clone() (vv RecordValueSet) {
	vv = make(RecordValueSet, len(vv))
	for i := range set {
		vv[i] = set[i].Clone()
	}

	return
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

// Replaces existing values, remove extra
func (set RecordValueSet) Replace(name string, values ...string) (vv RecordValueSet) {
	for i := range set {
		if set[i].Name != name {
			// copy values from other fields
			vv = append(vv, set[i])
		}
	}

	for p, v := range values {
		vv = append(vv, &RecordValue{
			Name:  name,
			Value: v,
			Place: uint(p),
		})
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

func (set RecordValueSet) SetRecordID(recordID uint64) {
	for i := range set {
		set[i].RecordID = recordID
	}
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
		// Mark all old as updated
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
				// Value is the same
				ex.Updated = false
			} else if !ex.Updated {
				// Value changed and old one was not marked as updated before
				// See if values really changed and update old value on existing value
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

// Returns structured representation of values casted to the appropriate types
func (set RecordValueSet) Dict(fields ModuleFieldSet) map[string]interface{} {
	var (
		rval = make(map[string]interface{})

		format = func(f *ModuleField, v string) interface{} {
			switch strings.ToLower(f.Kind) {
			case "bool":
				return payload.ParseBool(v)
			case "number":
				if f.Options.Precision() > 0 {
					num, _ := strconv.ParseFloat(v, 64)
					return num
				}

				num, _ := strconv.ParseInt(v, 10, 64)
				return num
			}

			return v
		}
	)

	if len(fields) == 0 {
		return rval
	}

	_ = fields.Walk(func(f *ModuleField) error {
		if f.Multi {
			var (
				rv = set.FilterByName(f.Name)
				vv = make([]interface{}, len(rv))
			)
			for i, val := range rv {
				vv[i] = format(f, val.Value)
			}
			rval[f.Name] = vv
		} else if v := set.Get(f.Name, 0); v != nil {
			rval[f.Name] = format(f, v.Value)
		}

		return nil
	})

	return rval
}

func (set RecordValueSet) Len() int           { return len(set) }
func (set RecordValueSet) Swap(i, j int)      { set[i], set[j] = set[j], set[i] }
func (set RecordValueSet) Less(i, j int) bool { return set[i].Place < set[j].Place }
