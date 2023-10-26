package dal

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModelFindByRefs(t *testing.T) {
	tcc := []struct {
		name  string
		in    ModelSet
		refs  map[string]interface{}
		found bool
	}{
		{
			name: "one ref found",
			in: ModelSet{
				{Refs: map[string]any{"a": 1}},
			},
			refs:  map[string]any{"a": 1},
			found: true,
		},
		{
			name: "one ref not found",
			in: ModelSet{
				{Refs: map[string]any{"a": 1}},
			},
			refs:  map[string]any{"b": 1},
			found: false,
		},

		{
			name: "n refs found",
			in: ModelSet{
				{Refs: map[string]any{"a": 1, "b": 1, "c": 1}},
			},
			refs:  map[string]any{"a": 1, "b": 1},
			found: true,
		},
		{
			name: "n refs not found",
			in: ModelSet{
				{Refs: map[string]any{"a": 1, "b": 1}},
			},
			refs:  map[string]any{"a": 1, "b": 2},
			found: false,
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			m := tc.in.FindByRefs(tc.refs)
			if tc.found {
				require.NotNil(t, m)
			} else {
				require.Nil(t, m)
			}
		})
	}
}

func TestAttributeMarshling(t *testing.T) {
	a := &Attribute{
		Ident: "foo",
		Type:  &TypeBlob{},
		Store: &CodecPlain{},
	}

	bb, err := json.Marshal(a)
	require.NoError(t, err)

	b := &Attribute{}
	err = json.Unmarshal(bb, &b)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(a, b))
}
