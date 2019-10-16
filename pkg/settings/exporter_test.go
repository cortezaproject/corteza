package settings

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestExport(t *testing.T) {
	req := require.New(t)
	tests := []struct {
		name, rawValue string
		parsedValue    interface{}
	}{
		{
			name:        "v_string",
			rawValue:    `"string"`,
			parsedValue: yaml.MapItem{Key: "v_string", Value: "string"},
		},
		{
			name:        "v_int-as-float",
			rawValue:    "12.34",
			parsedValue: yaml.MapItem{Key: "v_int-as-float", Value: 12.34},
		},
		{
			name:        "v_float",
			rawValue:    "123",
			parsedValue: yaml.MapItem{Key: "v_float", Value: 123.0},
		},
		{
			name:        "v_bool_true",
			rawValue:    "true",
			parsedValue: yaml.MapItem{Key: "v_bool_true", Value: true},
		},
		{
			name:        "v_bool_false",
			rawValue:    "false",
			parsedValue: yaml.MapItem{Key: "v_bool_false", Value: false},
		},
		{
			name:        "v_slice",
			rawValue:    `[1, 1.23, "string", true, false]`,
			parsedValue: yaml.MapItem{Key: "v_slice", Value: []interface{}{1.0, 1.23, "string", true, false}},
		},
		{
			name:        "v_map",
			rawValue:    `{"k1": "string","k2": 1,"k3":1.23,"k4":true,"k5":false}`,
			parsedValue: yaml.MapItem{Key: "v_map", Value: map[string]interface{}{"k1": "string", "k2": 1.0, "k3": 1.23, "k4": true, "k5": false}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := Export(ValueSet{
				&Value{
					Name:  test.name,
					Value: types.JSONText(test.rawValue),
				},
			})
			req.Len(tt, 1)
			req.Equal(test.parsedValue, tt[0])
		})
	}
}
