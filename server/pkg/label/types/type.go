package types

import (
	"github.com/cortezaproject/corteza/server/pkg/sql"
	"database/sql/driver"

	"encoding/json"
	
)
type (
	LabelValue struct{
			Val string	`json:"value,omitempty"`
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
		Name      string
		Value 	  []string //use a slice of string
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
func (lv LabelValue) Value() (driver.Value, error) { return json.Marshal(lv) }

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