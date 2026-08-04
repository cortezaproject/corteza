package toolkit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID(t *testing.T) {
	t.Run("parses a string id", func(t *testing.T) {
		id, err := ID(map[string]any{"recordID": "18446744073709551615"}, "recordID")
		require.NoError(t, err)
		assert.Equal(t, uint64(18446744073709551615), id)
	})

	t.Run("absent and empty are zero, not errors", func(t *testing.T) {
		id, err := ID(map[string]any{}, "recordID")
		require.NoError(t, err)
		assert.Zero(t, id)

		id, err = ID(map[string]any{"recordID": ""}, "recordID")
		require.NoError(t, err)
		assert.Zero(t, id)
	})

	// The whole reason IDs cross the boundary as strings: a JSON number has
	// already lost precision by the time it reaches us, so coercing it would
	// launder a corrupted value into the service layer.
	t.Run("rejects a numeric id rather than coercing", func(t *testing.T) {
		_, err := ID(map[string]any{"recordID": float64(1.8446744073709552e19)}, "recordID")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be a string")
	})

	t.Run("rejects a non-numeric string", func(t *testing.T) {
		_, err := ID(map[string]any{"recordID": "abc"}, "recordID")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid recordID")
	})
}

func TestReqID(t *testing.T) {
	_, err := ReqID(map[string]any{}, "recordID")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestPage(t *testing.T) {
	cases := map[string]struct {
		args map[string]any
		want uint
	}{
		"default when absent": {map[string]any{}, DefaultLimit},
		"string limit":        {map[string]any{"limit": "10"}, 10},
		"numeric limit":       {map[string]any{"limit": float64(10)}, 10},
		"capped at max":       {map[string]any{"limit": "100000"}, MaxLimit},
		"zero falls back":     {map[string]any{"limit": "0"}, DefaultLimit},
		"garbage falls back":  {map[string]any{"limit": "abc"}, DefaultLimit},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.want, Page(c.args).Limit)
		})
	}

	assert.Equal(t, "abc", Page(map[string]any{"pageCursor": "abc"}).Cursor)
}

func TestJSONResult(t *testing.T) {
	t.Run("marshals", func(t *testing.T) {
		res, err := JSONResult(map[string]any{"a": 1})
		require.NoError(t, err)
		require.Len(t, res.Content, 1)
	})

	t.Run("refuses oversized output", func(t *testing.T) {
		_, err := JSONResult(strings.Repeat("x", MaxResultBytes+1))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "result too large")
		assert.Contains(t, err.Error(), "pageCursor")
	})
}

func TestErrf(t *testing.T) {
	err := Errf("record lookup", assert.AnError)
	assert.Contains(t, err.Error(), "record lookup failed:")
	assert.ErrorIs(t, err, assert.AnError)
}

func TestBool(t *testing.T) {
	cases := map[string]struct {
		args map[string]any
		want bool
	}{
		"absent":       {map[string]any{}, false},
		"string true":  {map[string]any{"f": "true"}, true},
		"string false": {map[string]any{"f": "false"}, false},
		"native true":  {map[string]any{"f": true}, true},
		"garbage":      {map[string]any{"f": "yes please"}, false},
		"number":       {map[string]any{"f": float64(1)}, false},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.want, Bool(c.args, "f"))
		})
	}
}

func TestRef(t *testing.T) {
	t.Run("accepts a handle or an id as string", func(t *testing.T) {
		s, err := Ref(map[string]any{"r": "my-handle"}, "r")
		require.NoError(t, err)
		assert.Equal(t, "my-handle", s)

		s, err = Ref(map[string]any{"r": "123"}, "r")
		require.NoError(t, err)
		assert.Equal(t, "123", s)
	})

	t.Run("absent is empty, not an error", func(t *testing.T) {
		s, err := Ref(map[string]any{}, "r")
		require.NoError(t, err)
		assert.Empty(t, s)
	})

	// Without this a numeric ref would read as "" and silently drop a
	// single-item lookup into list mode.
	t.Run("rejects a numeric ref", func(t *testing.T) {
		_, err := Ref(map[string]any{"r": float64(123)}, "r")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be a string")
	})

	t.Run("ReqRef demands a value", func(t *testing.T) {
		_, err := ReqRef(map[string]any{}, "r")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "required")
	})
}
