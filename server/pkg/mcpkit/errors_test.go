package mcpkit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func failingTool(err error) server.ToolHandlerFunc {
	return func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) { return nil, err }
}

func callWith(t *testing.T, m *MCPServer, tool string, err error) ToolError {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Name = tool
	res, callErr := m.errorsAsResults()(failingTool(err))(context.Background(), req)
	require.NoError(t, callErr, "a handler error must become a result, never a protocol error")
	require.NotNil(t, res)
	assert.True(t, res.IsError)

	text, ok := res.Content[0].(mcp.TextContent)
	require.True(t, ok)
	var body struct{ Error ToolError }
	require.NoError(t, json.Unmarshal([]byte(text.Text), &body))

	structured, err2 := json.Marshal(res.StructuredContent)
	require.NoError(t, err2)
	assert.JSONEq(t, text.Text, string(structured), "text and structuredContent carry the same error")
	return body.Error
}

// TestErrorsAsResults pins the wire shape of a failure and where each code
// comes from: a Coded error names its own, the domain classifier names the
// rest, and anything left is "failed".
func TestErrorsAsResults(t *testing.T) {
	reg := NewRegistry()
	m := &MCPServer{reg: reg}

	t.Run("coded error keeps its code and next", func(t *testing.T) {
		te := callWith(t, m, "x_y_create", toolkit.WithCode(errors.New("too big"), toolkit.CodeTooLarge, "narrow it"))
		assert.Equal(t, ToolError{Code: toolkit.CodeTooLarge, Message: "too big", Next: "narrow it"}, te)
	})

	t.Run("an invalid argument points at the tool's documentation", func(t *testing.T) {
		_, err := toolkit.ReqStr(map[string]any{}, "namespace")
		te := callWith(t, m, "compose_module_create", err)
		assert.Equal(t, toolkit.CodeInvalidArgument, te.Code)
		assert.Equal(t, "namespace is required", te.Message)
		assert.Contains(t, te.Next, toolLoadName)
		assert.Contains(t, te.Next, "compose_module_create")
	})

	t.Run("an unclassified error is failed with no next", func(t *testing.T) {
		te := callWith(t, m, "x", fmt.Errorf("module create failed: %w", errors.New("disk on fire")))
		assert.Equal(t, ToolError{Code: toolkit.CodeFailed, Message: "module create failed: disk on fire"}, te)
	})

	t.Run("the domain classifier names the code and gets the default next", func(t *testing.T) {
		sentinel := errors.New("namespace does not exist")
		reg.SetErrorClassifier(func(err error) (string, string) {
			if errors.Is(err, sentinel) {
				return toolkit.CodeNotFound, ""
			}
			return "", ""
		})
		te := callWith(t, m, "x", fmt.Errorf("namespace lookup failed: %w", sentinel))
		assert.Equal(t, toolkit.CodeNotFound, te.Code)
		assert.Contains(t, te.Next, "_lookup")
	})

	t.Run("a passing call is untouched", func(t *testing.T) {
		ok := func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultText("fine"), nil
		}
		res, err := m.errorsAsResults()(ok)(context.Background(), mcp.CallToolRequest{})
		require.NoError(t, err)
		assert.False(t, res.IsError)
	})
}

// TestRiskCeilingAndUnknownArgumentsAreCoded: the two refusals mcpkit raises
// itself carry their codes, so the outer middleware does not have to recognise
// them by message.
func TestRiskCeilingAndUnknownArgumentsAreCoded(t *testing.T) {
	var coded toolkit.Coded

	err := unknownArgumentError("compose_record_lookup", []string{"sort"}, map[string]any{"query": 1, "limit": 1})
	require.ErrorAs(t, err, &coded)
	assert.Equal(t, toolkit.CodeUnknownArgument, coded.ToolErrorCode())

	reg := NewRegistry()
	reg.RegisterTool(mcp.NewTool("x_thing_delete", WithRisk(RiskDestructive)), "Delete", failingTool(nil))
	m := &MCPServer{reg: reg}
	req := mcp.CallToolRequest{}
	req.Params.Name = "x_thing_delete"
	ctx := WithScope(context.Background(), Scope{MaxRisk: RiskWrite})
	_, err = m.riskCeiling()(failingTool(nil))(ctx, req)
	require.ErrorAs(t, err, &coded)
	assert.Equal(t, toolkit.CodeRiskCapped, coded.ToolErrorCode())
	assert.Contains(t, coded.ToolErrorNext(), "maxRisk")
}

// TestExecuteToolKeepsGoErrors: the in-process runtime is not behind the HTTP
// middleware and still receives the error itself.
func TestExecuteToolKeepsGoErrors(t *testing.T) {
	reg := NewRegistry()
	boom := errors.New("boom")
	reg.RegisterTool(mcp.NewTool("x_thing_lookup", WithRisk(RiskRead)), "Lookup", failingTool(boom))
	_, err := reg.ExecuteTool(context.Background(), "x_thing_lookup", nil)
	assert.ErrorIs(t, err, boom)
}
