package cast2

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMap(t *testing.T) {
	var (
		req = require.New(t)
	)

	{
		target := make(map[string]any)
		req.NoError(Meta([]byte(`{"a":"b"}`), &target))
		req.Equal(map[string]any{"a": "b"}, target)
	}

	{
		target := make(map[string]any)
		req.NoError(Meta(`{"a":"b"}`, &target))
		req.Equal(map[string]any{"a": "b"}, target)
	}

	{
		target := make(map[string]any)
		req.NoError(Meta(map[string]any{"a": "b"}, &target))
		req.Equal(map[string]any{"a": "b"}, target)
	}
}
