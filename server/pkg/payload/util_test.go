package payload

import (
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseFilterState(t *testing.T) {
	var (
		req = require.New(t)
	)
	tests := []struct {
		name     string
		input    string
		expected filter.State
	}{
		{
			"invalid string should be Excluded",
			"zero",
			filter.StateExcluded,
		},
		{
			"empty string should be Excluded",
			"0",
			filter.StateExcluded,
		},
		{
			"Excluded",
			"0",
			filter.StateExcluded,
		},
		{
			"Excluded",
			"1",
			filter.StateInclusive,
		},
		{
			"Excluded",
			"2",
			filter.StateExclusive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req.Equal(tt.expected, ParseFilterState(tt.input))
		})
	}
}

func TestParseMeta(t *testing.T) {
	var (
		req = require.New(t)
	)
	tests := []struct {
		name     string
		input    []string
		expected map[string]any
		hasError bool
	}{
		{
			"key-value pair",
			[]string{"foo=bar"},
			map[string]any{"foo": "bar"},
			false,
		},
		{
			"json object",
			[]string{`{"foo":"bar"}`},
			map[string]any{"foo": "bar"},
			false,
		},
		{
			"json object with a key the handle validation would reject",
			[]string{`{"a":"bar"}`},
			map[string]any{"a": "bar"},
			false,
		},
		{
			"json object with a quoted key",
			[]string{`{"fo'o":"bar"}`},
			nil,
			true,
		},
		{
			"json object with a backslash in the key",
			[]string{`{"fo\\o":"bar"}`},
			map[string]any{`fo\o`: "bar"},
			false,
		},
		{
			"json object with a double quote in the key",
			[]string{`{"fo\"o":"bar"}`},
			map[string]any{`fo"o`: "bar"},
			false,
		},
		{
			"json object with a control character in the key",
			[]string{"{\"fo\\u0000o\":\"bar\"}"},
			map[string]any{"fo\x00o": "bar"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := ParseMeta(tt.input)

			if tt.hasError {
				req.Error(err)
				req.Nil(m)
				return
			}

			req.NoError(err)
			req.Equal(tt.expected, m)
		})
	}
}
