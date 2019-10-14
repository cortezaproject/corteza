package settings

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
)

func TestExport(t *testing.T) {
	ss := ValueSet{}

	ss = append(ss, ([]*Value{
		&Value{
			Name:  "v_string",
			Value: types.JSONText("\"string\""),
		},
		&Value{
			Name:  "v_float",
			Value: types.JSONText("123"),
		},
		&Value{
			Name:  "v_float",
			Value: types.JSONText("12.34"),
		},
		&Value{
			Name:  "v_bool",
			Value: types.JSONText("true"),
		},
		&Value{
			Name:  "v_slice",
			Value: types.JSONText("[1, \"string\", true]"),
		},
		&Value{
			Name:  "v_map",
			Value: types.JSONText("{\"k1\": \"v1\",\"k2\": 2}"),
		},
	})...)

	tt := Export(ss)

	if _, ok := tt[0].Value.(string); !ok {
		t.Errorf("Expecting %v to be string", tt[0].Value)
	}

	if _, ok := tt[1].Value.(float64); !ok {
		t.Errorf("Expecting %v to be float64", tt[1].Value)
	}

	if _, ok := tt[2].Value.(float64); !ok {
		t.Errorf("Expecting %v to be float64", tt[2].Value)
	}

	if _, ok := tt[3].Value.(bool); !ok {
		t.Errorf("Expecting %v to be bool", tt[3].Value)
	}

	if _, ok := tt[4].Value.([]interface{}); !ok {
		t.Errorf("Expecting %v to be []interface{}", tt[4].Value)
	}

	if _, ok := tt[5].Value.(map[string]interface{}); !ok {
		t.Errorf("Expecting %v to be map[string]interface {}", tt[5].Value)
	}
}
