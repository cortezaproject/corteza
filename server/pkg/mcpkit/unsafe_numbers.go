package mcpkit

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// rejectUnsafeNumbers refuses a call whose arguments hold a number at or
// beyond 2^53, whether sent as JSON or inside an argument that is a JSON
// object or array written as a string. Either way the number is decoded as a
// float64, so such a number — in practice an ID sent unquoted — is already a
// different value.
func (m *MCPServer) rejectUnsafeNumbers() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, ok := req.Params.Arguments.(map[string]any)
			if !ok {
				return next(ctx, req)
			}

			for k, v := range args {
				if s, is := v.(string); is {
					v = jsonInString(s)
				}

				if err := toolkit.RejectUnsafeNumbers(k, v); err != nil {
					return nil, err
				}
			}

			return next(ctx, req)
		}
	}
}

// jsonInString decodes s when it holds a JSON object or array, and returns nil
// otherwise.
func jsonInString(s string) any {
	s = strings.TrimSpace(s)
	if s == "" || (s[0] != '{' && s[0] != '[') {
		return nil
	}

	var out any
	if json.Unmarshal([]byte(s), &out) != nil {
		return nil
	}

	return out
}
