package revisions

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// Note: These tests use json.Number because when reading revisions from DB,
// SetValue uses decoder.UseNumber() which decodes JSON numbers as json.Number
// (not uint64 or float64). This preserves precision for large uint64 values.

func Test_normalizeChangeValues(t *testing.T) {
	cases := []struct {
		name     string
		input    []any
		expected []any
	}{
		{
			"empty",
			[]any{},
			nil,
		},
		{
			"nil",
			nil,
			nil,
		},
		{
			"json.Number to string",
			[]any{json.Number("123456789012345")},
			[]any{"123456789012345"},
		},
		{
			"strings unchanged",
			[]any{"test", "value"},
			[]any{"test", "value"},
		},
		{
			"booleans unchanged",
			[]any{true, false},
			[]any{true, false},
		},
		{
			"mixed types",
			[]any{json.Number("999"), "text", true, json.Number("12345678901234567890")},
			[]any{"999", "text", true, "12345678901234567890"},
		},
		{
			"regular numbers unchanged",
			[]any{int64(123), float64(45.67)},
			[]any{int64(123), float64(45.67)},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				result = normalizeChangeValues(c.input)
			)
			req.Equal(c.expected, result)
		})
	}
}

func Test_Change_MarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		change   *Change
		expected string
	}{
		{
			"strings",
			&Change{
				Key: "name",
				Old: []any{"old name"},
				New: []any{"new name"},
			},
			`{"key":"name","new":["new name"],"old":["old name"]}`,
		},
		{
			"json.Number to strings",
			&Change{
				Key: "ownedBy",
				Old: []any{json.Number("123456789012345")},
				New: []any{json.Number("987654321098765")},
			},
			`{"key":"ownedBy","new":["987654321098765"],"old":["123456789012345"]}`,
		},
		{
			"mixed types",
			&Change{
				Key: "status",
				Old: []any{json.Number("1"), "active", true},
				New: []any{json.Number("2"), "inactive", false},
			},
			`{"key":"status","new":["2","inactive",false],"old":["1","active",true]}`,
		},
		{
			"empty old",
			&Change{
				Key: "newField",
				Old: []any{},
				New: []any{"value"},
			},
			`{"key":"newField","new":["value"]}`,
		},
		{
			"empty new (deletion)",
			&Change{
				Key: "deletedField",
				Old: []any{"value"},
				New: []any{},
			},
			`{"key":"deletedField","old":["value"]}`,
		},
		{
			"large uint64 values",
			&Change{
				Key: "createdBy",
				Old: []any{json.Number("18446744073709551615")},
				New: []any{json.Number("18446744073709551614")},
			},
			`{"key":"createdBy","new":["18446744073709551614"],"old":["18446744073709551615"]}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				result, err = c.change.MarshalJSON()
			)
			req.NoError(err)
			req.JSONEq(c.expected, string(result))
		})
	}
}

func Test_Revision_SetValue_Delta_UseNumber(t *testing.T) {
	t.Run("large uint64 preserved as json.Number", func(t *testing.T) {
		var (
			req = require.New(t)
			r   = &Revision{}

			deltaJSON = []byte(`[
				{
					"key": "ownedBy",
					"old": [123456789012345],
					"new": [987654321098765]
				},
				{
					"key": "name",
					"old": ["old"],
					"new": ["new"]
				}
			]`)
		)

		err := r.SetValue("delta", 0, deltaJSON)
		req.NoError(err)
		req.Len(r.Changes, 2)

		req.Equal("ownedBy", r.Changes[0].Key)
		req.Len(r.Changes[0].Old, 1)
		req.Len(r.Changes[0].New, 1)

		oldVal, ok := r.Changes[0].Old[0].(json.Number)
		req.True(ok, "Old value should be json.Number")
		req.Equal("123456789012345", oldVal.String())

		newVal, ok := r.Changes[0].New[0].(json.Number)
		req.True(ok, "New value should be json.Number")
		req.Equal("987654321098765", newVal.String())

		req.Equal("name", r.Changes[1].Key)
		req.Equal("old", r.Changes[1].Old[0])
		req.Equal("new", r.Changes[1].New[0])
	})

	t.Run("round-trip preserves precision", func(t *testing.T) {
		var (
			req = require.New(t)
			original = &Change{
				Key: "userId",
				Old: []any{json.Number("18446744073709551615")},
				New: []any{json.Number("18446744073709551614")},
			}
		)

		marshaled, err := original.MarshalJSON()
		req.NoError(err)

		var intermediate map[string]any
		err = json.Unmarshal(marshaled, &intermediate)
		req.NoError(err)
		req.Equal("userId", intermediate["key"])

		oldSlice := intermediate["old"].([]any)
		req.Equal("18446744073709551615", oldSlice[0])

		newSlice := intermediate["new"].([]any)
		req.Equal("18446744073709551614", newSlice[0])
	})
}

func Test_Revision_MarshalJSON_Integration(t *testing.T) {
	t.Run("full revision with uint64 strings", func(t *testing.T) {
		var (
			req = require.New(t)
			r   = &Revision{
				ID:         1234567890123456789,
				ResourceID: 9876543210987654321,
				UserID:     1111111111111111111,
				Revision:   5,
				Operation:  "update",
				Comment:    "Test revision",
				Changes: []*Change{
					{
						Key: "ownedBy",
						Old: []any{json.Number("123456789012345")},
						New: []any{json.Number("987654321098765")},
					},
				},
			}
		)

		data, err := json.Marshal(r)
		req.NoError(err)

		var result map[string]any
		err = json.Unmarshal(data, &result)
		req.NoError(err)

		req.Equal("1234567890123456789", result["changeID"])
		req.Equal("9876543210987654321", result["resourceID"])
		req.Equal("1111111111111111111", result["userID"])

		changes := result["changes"].([]any)
		req.Len(changes, 1)

		change := changes[0].(map[string]any)
		req.Equal("ownedBy", change["key"])

		old := change["old"].([]any)
		req.Equal("123456789012345", old[0])

		new := change["new"].([]any)
		req.Equal("987654321098765", new[0])
	})
}
