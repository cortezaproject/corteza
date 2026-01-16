package types

import (
	"database/sql/driver"

	"github.com/cortezaproject/corteza/server/pkg/sql"

	"encoding/json"
)

type (
	LabelValue struct {
		Val    string   `json:"value,omitempty"`
		Values []string `json:"values,omitempty"`
	}
	Label struct {
		// Kind of the labeled resource
		Kind string

		// ID of the labeled resource
		ResourceID uint64

		Name  string
		Value LabelValue
	}

	LabelFilter struct {
		Kind       string
		ResourceID []uint64
		Filter     map[string][]string
		Limit      uint
		Name       string
		Value      []string
	}
)

const (
	LabelResourceType = "corteza::generic:label"
)

func (set LabelSet) ResourceIDs() (rr []uint64) {
	rr = make([]uint64, len(set))
	for r := range set {
		rr[r] = set[r].ResourceID
	}

	return
}

func (set LabelSet) FilterByResource(kind string, ID uint64) map[string]LabelValue {
	var kv = make(map[string]LabelValue)
	for _, label := range set {
		if kind == label.Kind && ID == label.ResourceID {
			kv[label.Name] = label.Value
		}
	}

	return kv
}
func (lv *LabelValue) Scan(src any) error { return sql.ParseJSON(src, lv) }

// Value implements driver.Value for database storage
// Uses explicit object format {"value":...} or {"values":...} for DB storage
// (different from MarshalJSON which outputs simplified format for API)
func (lv LabelValue) Value() (driver.Value, error) {
	if len(lv.Values) > 0 {
		return json.Marshal(struct {
			Values []string `json:"values"`
		}{Values: lv.Values})
	}
	return json.Marshal(struct {
		Val string `json:"value"`
	}{Val: lv.Val})
}

func (lv LabelValue) MarshalJSON() ([]byte, error) {
	if len(lv.Values) > 0 {
		return json.Marshal(lv.Values)
	}
	return json.Marshal(lv.Val)
}

func (lv *LabelValue) UnmarshalJSON(data []byte) error {
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		lv.Val = strVal
		lv.Values = nil
		return nil
	}
	var arrVal []string
	if err := json.Unmarshal(data, &arrVal); err == nil {
		lv.Values = arrVal
		lv.Val = ""
		return nil
	}
	var obj struct {
		Val    string   `json:"value,omitempty"`
		Values []string `json:"values,omitempty"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		lv.Val = obj.Val
		lv.Values = obj.Values
		return nil
	}

	return nil
}

func (lv LabelValue) Equal(other LabelValue) bool {
	if lv.Val != "" && other.Val != "" {
		return lv.Val == other.Val
	}

	if len(lv.Values) > 0 && len(other.Values) > 0 {
		if len(lv.Values) != len(other.Values) {
			return false
		}
		for i := range lv.Values {
			if lv.Values[i] != other.Values[i] {
				return false
			}
		}
		return true
	}

	return lv.Val == "" && len(lv.Values) == 0 && other.Val == "" && len(other.Values) == 0
}
