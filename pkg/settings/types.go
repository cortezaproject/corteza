package settings

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type (
	Value struct {
		Name  string         `json:"name"  db:"name"`
		Value types.JSONText `json:"value" db:"value"`

		// Setting owner, 0 for global settings
		OwnedBy uint64 `json:"-" db:"rel_owner"`

		// Who updated & when
		UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
		UpdatedBy uint64    `json:"updatedBy" db:"updated_by"`
	}

	Filter struct {
		Prefix  string `json:"prefix"`
		OwnedBy uint64 `json:"ownedBy"`
		Page    uint64 `json:"page"`
		PerPage uint64 `json:"perPage"`
	}

	KV map[string]types.JSONText
)

const (
	settingsFilterPerPageMax = 100
)

func (f *Filter) Normalize() {
	f.Prefix = strings.TrimSpace(f.Prefix)
	if f.PerPage > settingsFilterPerPageMax {
		f.PerPage = settingsFilterPerPageMax
	}
}

func (v *Value) SetRawValue(str string) error {
	var dummy interface{}
	// Test input to be sure we can save it...
	if err := json.Unmarshal([]byte(str), &dummy); err != nil {
		return err
	}

	v.Value = types.JSONText(str)
	return nil
}

func (v *Value) SetValue(value interface{}) (err error) {
	v.Value, err = json.Marshal(value)
	return
}

func (v *Value) String() (out string) {
	if v == nil {
		return ""
	}

	_ = v.Value.Unmarshal(&out)
	return
}

func (v *Value) Bool() (out bool) {
	if v == nil {
		return false
	}

	_ = v.Value.Unmarshal(&out)
	return
}

func (set ValueSet) KV() KV {
	m := KV{}

	_ = set.Walk(func(v *Value) error {
		m[v.Name] = v.Value
		return nil
	})

	return m
}

func (kv KV) Has(k string) (ok bool) {
	_, ok = kv[k]
	return
}

func (kv KV) Bool(k string) (out bool) {
	out = false
	if v, ok := kv[k]; ok {
		_ = v.Unmarshal(&out)
	}

	return
}

func (kv KV) String(k string) (out string) {
	out = ""
	if v, ok := kv[k]; ok {
		_ = v.Unmarshal(&out)
	}

	return
}

func (kv KV) Filter(prefix string) KV {
	var out = KV{}
	for k, v := range kv {
		if strings.Index(k, prefix) == 0 {
			out[k] = v
		}
	}

	return out
}

// Replace finds and updates existing or appends new value
func (set *ValueSet) Replace(n *Value) {
	for _, v := range *set {
		if v.Name == n.Name {
			v.Value = n.Value
			return
		}
	}

	*set = append(*set, n)
}

// Replace finds and updates existing or appends new value
func (set *ValueSet) Has(name string) bool {
	return set.First(name) != nil
}

// First finds and returns first value
func (set ValueSet) First(name string) *Value {
	for _, v := range set {
		if v.Name == name {
			return v
		}
	}

	return nil
}

// Returns all valus that changed or do not exist in the original set
func (set ValueSet) Changed(in ValueSet) (out ValueSet) {
input:
	for _, i := range in {
		for _, s := range set {
			if s.Name != i.Name {
				// Different name, not interested
				continue
			}

			if s.String() == i.String() {
				// Value did not change, continue with next input set
				continue input
			}

			// Value changed, break out the loop
			break
		}

		// Hande changed or missing value
		out = append(out, i)
	}

	return
}
