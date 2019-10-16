package settings

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestSettingImport_CastSet(t *testing.T) {
	var aux interface{}
	req := require.New(t)
	f, err := os.Open("testdata/settings.yaml")
	defer f.Close()
	req.NoError(err)
	req.NoError(yaml.NewDecoder(f).Decode(&aux))
	req.NotNil(aux)

	imp := NewImporter()
	err = imp.CastSet(aux)
	req.Nil(err)

	tests := map[string]string{
		"v_string":       "\"string\"",
		"v_float":        "12.34",
		"v_int-as-float": "123",
		"v_bool_true":    "true",
		"v_bool_false":   "false",
		"v_slice":        "[1,1.23,\"string\",true,false]",
		"v_map":          "{\"k1\":\"string\",\"k2\":1,\"k3\":1.23,\"k4\":true,\"k5\":false}",
	}

	for k, test := range tests {
		t.Run(k, func(t *testing.T) {
			v := imp.settings.First(k)
			vv, _ := types.JSONText(test).MarshalJSON()
			ee, _ := v.Value.MarshalJSON()
			req.Equal(ee, vv)
		})
	}
}
