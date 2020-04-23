package values

import (
	"reflect"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func Test_sanitizer_Run(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		options map[string]interface{}
		input   string
		output  string
		outref  uint64
	}{
		{
			name:   "numbers should be trimmed",
			kind:   "Number",
			input:  " 42 ",
			output: "42",
		},
		{
			name:   "object reference should be processed",
			kind:   "Record",
			input:  " 133569629112020995 ",
			output: "133569629112020995",
			outref: 133569629112020995,
		},
		{
			name:   "object reference should be numeric",
			kind:   "Record",
			input:  " foo ",
			output: "",
		},
		{
			name:   "user reference should be processed",
			kind:   "User",
			input:  " 133569629112020995 ",
			output: "133569629112020995",
			outref: 133569629112020995,
		},
		{
			name:   "user reference should be numeric",
			kind:   "User",
			input:  " foo ",
			output: "",
		},
		{
			name:   "strings should be kept intact",
			kind:   "String",
			input:  " The answer ",
			output: " The answer ",
		},
		{
			name:   "booleans should be converted (t)",
			kind:   "Bool",
			input:  "t",
			output: "1",
		},
		{
			name:   "booleans should be converted (false)",
			kind:   "Bool",
			input:  "false",
			output: "0",
		},
		{
			name:   "booleans should be converted (garbage)",
			kind:   "Bool",
			input:  "%%#)%)')$)'",
			output: "0",
		},
		{
			name:   "dates should be converted to ISO",
			kind:   "DateTime",
			input:  "Mon Jan 2 15:04:05 2006",
			output: "2006-01-02T15:04:05Z",
		},
		{
			name:   "dates should be converted to UTC",
			kind:   "DateTime",
			input:  "2020-03-02T20:20:20+05:00",
			output: "2020-03-02T15:20:20Z",
		},
		{
			name:   "micro/mili seconds should be cut off",
			kind:   "DateTime",
			input:  "2020-03-11T11:20:08.471Z",
			output: "2020-03-11T11:20:08Z",
		},
		{
			name:    "number precision",
			kind:    "Number",
			options: map[string]interface{}{"precision": 3},
			input:   "42.44455",
			output:  "42.445",
		},
		{
			name:    "number precision; bigger precision then provided fracture",
			kind:    "Number",
			options: map[string]interface{}{"precision": 3},
			input:   "42.4",
			output:  "42.4",
		},
		{
			name:    "number precision; default to 0",
			kind:    "Number",
			options: map[string]interface{}{},
			input:   "42.4",
			output:  "42",
		},
		{
			name:    "number precision; clamped between [0, 6]",
			kind:    "Number",
			options: map[string]interface{}{"precision": 10},
			input:   "42.5555555555",
			output:  "42.555556",
		},
		{
			name:    "number precision; clamped between [0, 6]",
			kind:    "Number",
			options: map[string]interface{}{"precision": -1},
			input:   "42.4",
			output:  "42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sanitizer{}
			m := &types.Module{Fields: types.ModuleFieldSet{&types.ModuleField{Name: "testField", Kind: tt.kind, Options: tt.options}}}
			v := types.RecordValueSet{&types.RecordValue{Name: "testField", Value: tt.input}}
			o := types.RecordValueSet{&types.RecordValue{Name: "testField", Value: tt.output, Ref: tt.outref}}

			// Need to mark values as updated to trigger sanitization.
			v.SetUpdatedFlag(true)
			o.SetUpdatedFlag(true)
			if sanitized := s.Run(m, v); !reflect.DeepEqual(sanitized, o) {
				t.Errorf("\ninput value:\n%v\n\nresult of sanitization:\n%v\n\nexpected:\n%v\n", tt.input, sanitized, o)
			}
		})
	}
}
