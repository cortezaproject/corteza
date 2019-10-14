package settings

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestSettingImport_CastSet(t *testing.T) {
	type tv struct {
		name  string
		value string
	}

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
	req.Len(imp.settings, 6)

	tests := []tv{
		tv{
			name:  "array",
			value: "[\"v1\",\"v2\"]",
		},
		tv{
			name:  "object",
			value: "{\"k1\":\"v1\",\"k2\":\"v2\"}",
		},
		tv{
			name:  "plain_b",
			value: "true",
		},
		tv{
			name:  "plain_double",
			value: "12.34",
		},
		tv{
			name:  "plain_double_2",
			value: "12",
		},
		tv{
			name:  "plain_s",
			value: "\"string\"",
		},
	}

	for _, t := range tests {
		v := imp.settings.First(t.name)
		vv, _ := types.JSONText(t.value).MarshalJSON()
		ee, _ := v.Value.MarshalJSON()
		req.Equal(vv, ee)
	}
}
