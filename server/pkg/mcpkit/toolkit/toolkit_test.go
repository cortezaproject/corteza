package toolkit

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
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
		var coded Coded
		require.ErrorAs(t, err, &coded)
		assert.Equal(t, CodeTooLarge, coded.ToolErrorCode())
		assert.Contains(t, coded.ToolErrorNext(), "pageCursor")
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

func TestJSONResultWithAddsSiblingsWithoutNesting(t *testing.T) {
	type page struct {
		ID     string   `json:"pageID"`
		Blocks []string `json:"blocks"`
	}

	res, err := JSONResultWith(page{ID: "7", Blocks: []string{"a"}}, map[string]string{
		"url":     "http://localhost/compose/namespace/crm/pages/7",
		"editUrl": "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text, ok := mcp.AsTextContent(res.Content[0])
	if !ok {
		t.Fatal("result carries no text content")
	}

	var out map[string]any
	if err = json.Unmarshal([]byte(text.Text), &out); err != nil {
		t.Fatalf("result is not JSON: %v", err)
	}

	// Everything the value already had must still be readable where it was.
	if out["pageID"] != "7" {
		t.Errorf("pageID is %v, want 7 — the value was nested rather than extended", out["pageID"])
	}
	if _, ok := out["blocks"]; !ok {
		t.Error("blocks went missing")
	}
	if out["url"] != "http://localhost/compose/namespace/crm/pages/7" {
		t.Errorf("url is %v", out["url"])
	}

	// An empty link is no link. A blank key would read as "there is one, and it
	// is empty", which is a different and wrong answer.
	if _, ok := out["editUrl"]; ok {
		t.Error("an empty value was emitted as a key")
	}
}

func TestJSONResultWithRefusesAValueThatIsNotAnObject(t *testing.T) {
	if _, err := JSONResultWith([]string{"a"}, map[string]string{"url": "x"}); err == nil {
		t.Error("a slice has nowhere to put a sibling key; that must be an error, not a silent drop")
	}
}
