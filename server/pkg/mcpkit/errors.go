package mcpkit

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ToolError is what a failed call returns: a tool result with isError set,
// whose text and structuredContent are {"error": ToolError}.
//
// A handler error used to leave as a JSON-RPC INTERNAL_ERROR, which some
// clients show as a transport failure; the model then retries the same call
// blind. A result it can read carries the kind of failure (code) and the move
// that follows (next).
type ToolError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Next    string `json:"next,omitempty"`
}

// Classifier maps a domain error to a code and next step. It returns "" for a
// code it cannot name; the domain registers one at boot because mcpkit may not
// know the domain's error types.
type Classifier func(err error) (code, next string)

// SetErrorClassifier installs the domain's classifier.
func (r *Registry) SetErrorClassifier(fn Classifier) { r.classify = fn }

// defaultNext is the move for a code when nobody supplied a more specific one.
func defaultNext(code, tool string) string {
	switch code {
	case toolkit.CodeInvalidArgument, toolkit.CodeUnknownArgument:
		return "call " + toolLoadName + " with " + tool + " for the parameter documentation, then retry"
	case toolkit.CodeNotFound:
		return "resolve the resource with its *_lookup tool first; a deleted one is unreachable by name and needs its ID with the *_undelete tool"
	case toolkit.CodeForbidden:
		return "the caller's roles do not allow this; no change to the arguments will, ask for the permission or use a different account"
	case toolkit.CodeConflict:
		return "re-read the resource and retry with its current state, or pick a handle that is not taken"
	case toolkit.CodeUnauthenticated:
		return "the session has no usable identity; reconnect with a valid token"
	}
	return ""
}

// classifyError turns err into the ToolError the wire carries: a Coded error
// names its own code, else the domain classifier, else "failed".
func (r *Registry) classifyError(err error, tool string) ToolError {
	te := ToolError{Code: toolkit.CodeFailed, Message: err.Error()}

	var coded toolkit.Coded
	if errors.As(err, &coded) {
		te.Code, te.Next = coded.ToolErrorCode(), coded.ToolErrorNext()
	} else if r.classify != nil {
		if code, next := r.classify(err); code != "" {
			te.Code, te.Next = code, next
		}
	}

	if te.Next == "" {
		te.Next = defaultNext(te.Code, tool)
	}
	return te
}

// ErrorResult renders a ToolError as the isError tool result.
func ErrorResult(te ToolError) *mcp.CallToolResult {
	body, _ := json.Marshal(map[string]ToolError{"error": te})
	res := mcp.NewToolResultStructured(json.RawMessage(body), string(body))
	res.IsError = true
	return res
}

// errorsAsResults is the outermost tool middleware: whatever fails inside it,
// including the risk ceiling and the unknown-argument check, reaches the client
// as a result rather than a protocol error. The in-process runtime does not
// pass through here and keeps Go errors (Registry.ExecuteTool).
func (m *MCPServer) errorsAsResults() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := next(ctx, req)
			if err == nil {
				return res, nil
			}
			return ErrorResult(m.reg.classifyError(err, req.Params.Name)), nil
		}
	}
}
