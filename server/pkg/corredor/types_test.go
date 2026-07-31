package corredor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScriptArgs_Decode(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]interface{}
		dec   map[string][]byte
		want  map[string]interface{}
	}{
		{
			// Corredor only returns the keys the script put on the returned
			// object, so this is the normal case, not an edge one
			name:  "script returns only some of the args",
			extra: map[string]interface{}{"validated": "false", "untouched": "keep"},
			dec:   map[string][]byte{"validated": []byte(`"true"`)},
			want:  map[string]interface{}{"validated": "true", "untouched": "keep"},
		},
		{
			name:  "script returns nothing",
			extra: map[string]interface{}{"validated": "false"},
			dec:   map[string][]byte{},
			want:  map[string]interface{}{"validated": "false"},
		},
		{
			name:  "script returns an empty value",
			extra: map[string]interface{}{"validated": "false"},
			dec:   map[string][]byte{"validated": {}},
			want:  map[string]interface{}{"validated": "false"},
		},
		{
			name:  "script returns another type than it was given",
			extra: map[string]interface{}{"count": "1"},
			dec:   map[string][]byte{"count": []byte(`42`)},
			want:  map[string]interface{}{"count": float64(42)},
		},
		{
			name:  "script returns a structured value",
			extra: map[string]interface{}{"cfg": "none"},
			dec:   map[string][]byte{"cfg": []byte(`{"a":true}`)},
			want:  map[string]interface{}{"cfg": map[string]interface{}{"a": true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				req = require.New(t)
				sa  = ExtendScriptArgs(mockEvent{}, tt.extra)
			)

			req.NoError(sa.Decode(tt.dec))
			req.Equal(tt.want, tt.extra)
		})
	}
}

func TestScriptArgs_DecodeMalformed(t *testing.T) {
	var (
		req   = require.New(t)
		extra = map[string]interface{}{"validated": "false"}
		sa    = ExtendScriptArgs(mockEvent{}, extra)
	)

	req.Error(sa.Decode(map[string][]byte{"validated": []byte(`{oops`)}))
}
