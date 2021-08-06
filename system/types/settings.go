package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type (
	SettingValue struct {
		Name  string         `json:"name"`
		Value types.JSONText `json:"value"`

		// Setting owner, 0 for global settings
		OwnedBy uint64 `json:"-"`

		// Who updated & when
		UpdatedAt time.Time `json:"updatedAt"`
		UpdatedBy uint64    `json:"updatedBy"`
	}

	SettingsFilter struct {
		Prefix  string `json:"prefix"`
		OwnedBy uint64 `json:"ownedBy"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*SettingValue) (bool, error) `json:"-"`
	}

	SettingsKV map[string]types.JSONText
)

const (
	settingsFilterPerPageMax = 100
)

func MakeSettingValue(name string, value interface{}) *SettingValue {
	o := &SettingValue{Name: name}
	_ = o.SetValue(value)
	return o
}

func (v *SettingValue) SetRawValue(str string) error {
	var dummy interface{}
	// Test input to be sure we can save it...
	if err := json.Unmarshal([]byte(str), &dummy); err != nil {
		return err
	}

	v.Value = types.JSONText(str)
	return nil
}

func (v *SettingValue) SetValue(value interface{}) (err error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err = enc.Encode(value); err != nil {
		return err
	}

	v.Value = bytes.TrimSpace(buf.Bytes())
	return
}

func (v *SettingValue) String() (out string) {
	if v == nil {
		return ""
	}

	_ = v.Value.Unmarshal(&out)
	return
}

func (v *SettingValue) Bool() (out bool) {
	if v == nil {
		return false
	}

	_ = v.Value.Unmarshal(&out)
	return
}

func (v *SettingValue) NormalizeValue() {

}

func (v *SettingValue) Eq(c *SettingValue) bool {
	return v != nil &&
		c != nil &&
		v.Name == c.Name &&
		v.OwnedBy == c.OwnedBy &&
		fmt.Sprintf("%v", v.Value) == fmt.Sprintf("%v", c.Value)
}

func (set SettingValueSet) FilterByPrefix(prefix string) SettingValueSet {
	pf := SettingValueSet{}
	for _, v := range set {
		if strings.HasPrefix(v.Name, prefix) {
			pf = append(pf, v)
		}
	}

	return pf
}

func (set SettingValueSet) KV() SettingsKV {
	m := SettingsKV{}

	_ = set.Walk(func(v *SettingValue) error {
		m[v.Name] = v.Value
		return nil
	})

	return m
}

func (kv SettingsKV) Has(k string) (ok bool) {
	_, ok = kv[k]
	return
}

func (kv SettingsKV) Bool(k string) (out bool) {
	out = false
	if v, ok := kv[k]; ok {
		_ = v.Unmarshal(&out)
	}

	return
}

func (kv SettingsKV) String(k string) (out string) {
	out = ""
	if v, ok := kv[k]; ok {
		_ = v.Unmarshal(&out)
	}

	return
}

func (kv SettingsKV) Filter(prefix string) SettingsKV {
	var out = SettingsKV{}
	for k, v := range kv {
		if strings.Index(k, prefix) == 0 {
			out[k] = v
		}
	}

	return out
}

// CutPrefix returns values with matching prefix and removes the prefix from keys
func (kv SettingsKV) CutPrefix(prefix string) SettingsKV {
	var out = SettingsKV{}
	for k, v := range kv {
		if strings.Index(k, prefix) == 0 {
			out[k[len(prefix):]] = v
		}
	}

	return out
}

// Decode is a helper function on SettingsKV that calls DecodeKV() and passes on the dst
func (kv SettingsKV) Decode(dst interface{}) error {
	return DecodeKV(kv, dst)
}

// Replace finds and updates existing or appends new value
func (set *SettingValueSet) Replace(n *SettingValue) {
	for _, v := range *set {
		if v.Name == n.Name {
			v.Value = n.Value
			return
		}
	}

	*set = append(*set, n)
}

// Replace finds and updates existing or appends new value
func (set *SettingValueSet) Has(name string) bool {
	return set.First(name) != nil
}

// First finds and returns first value
func (set SettingValueSet) First(name string) *SettingValue {
	for _, v := range set {
		if v.Name == name {
			return v
		}
	}

	return nil
}

// Returns all valus that changed or do not exist in the original set
func (set SettingValueSet) Changed(in SettingValueSet) (out SettingValueSet) {
input:
	for _, i := range in {
		for _, s := range set {
			if s.Name != i.Name {
				// Different name, not interested
				continue
			}

			if s.Eq(i) {
				// SettingValue did not change, continue with next input set
				continue input
			}

			// SettingValue changed, break out the loop
			break
		}

		// Handle changed or missing value
		out = append(out, i)
	}

	return
}

// New returns all new values (that do not exist in the original set)
func (set SettingValueSet) New(in SettingValueSet) (out SettingValueSet) {
	org := set.KV()

	for _, v := range in {
		if !org.Has(v.Name) {
			out = append(out, v)
		}
	}

	return
}

// New returns all new values (that do not exist in the original set)
func (set SettingValueSet) Old(in SettingValueSet) (out SettingValueSet) {
	org := set.KV()

	for _, v := range in {
		if org.Has(v.Name) {
			out = append(out, v)
		}
	}

	return
}
