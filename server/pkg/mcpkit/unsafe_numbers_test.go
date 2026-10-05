package mcpkit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestRejectUnsafeNumbers(t *testing.T) {
	// call decodes args the way the transport does, so a large number arrives
	// as an already rounded float64
	call := func(t *testing.T, args string) (bool, error) {
		var decoded map[string]any
		require.NoError(t, json.Unmarshal([]byte(args), &decoded))

		reached := false
		next := func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			reached = true
			return &mcp.CallToolResult{}, nil
		}

		req := mcp.CallToolRequest{}
		req.Params.Name = "t"
		req.Params.Arguments = decoded

		_, err := (&MCPServer{}).rejectUnsafeNumbers()(next)(context.Background(), req)
		return reached, err
	}

	t.Run("a nested ID sent as a number is refused, naming its path", func(t *testing.T) {
		reached, err := call(t, `{"fields":[{"name":"ref","options":{"moduleID":516687706984677377}}]}`)
		require.False(t, reached)
		require.ErrorContains(t, err, "fields[0].options.moduleID is a number too large")
	})

	t.Run("an ID inside a JSON string argument is refused", func(t *testing.T) {
		reached, err := call(t, `{"input":"{\"recordID\": 516687706984677377}"}`)
		require.False(t, reached)
		require.ErrorContains(t, err, "input.recordID")
	})

	t.Run("a top-level ID sent as a number is refused", func(t *testing.T) {
		_, err := call(t, `{"recordID":516687706984677377}`)
		require.ErrorContains(t, err, "recordID is a number too large")
	})

	t.Run("IDs as strings and small numbers pass", func(t *testing.T) {
		reached, err := call(t, `{"recordID":"516687706984677377","weight":3,"values":{"amount":12.5},"text":"[not json"}`)
		require.NoError(t, err)
		require.True(t, reached)
	})
}
