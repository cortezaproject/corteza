// Package toolkit holds the shared plumbing every MCP tool handler repeats:
// argument unwrapping, ID parsing, paging, result marshalling and error
// wrapping.
//
// Beyond removing duplication, the point is to create single choke points. Most
// importantly JSONResult: it is the only path from a Go value to a tool result,
// so a result-size ceiling lives in exactly one place, and if tool output ever
// needs to mark third-party or user-authored content as data rather than
// instruction, that is one function to change instead of ~200 call sites.
package toolkit

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	// DefaultLimit is applied when a lookup declares paging but the caller
	// omits it. Chosen well below MaxLimit so the common case is cheap.
	DefaultLimit uint = 50

	// MaxLimit caps what a caller may ask for. The REST record search caps at
	// 1000; tool results are read into a model's context, so the ceiling here
	// is deliberately lower.
	MaxLimit uint = 200

	// MaxResultBytes is the backstop on marshalled tool output. Schema-level
	// paging is what a well-behaved tool uses; this catches the ones that do
	// not, including any lookup that reaches a service with unbounded paging.
	MaxResultBytes = 256 * 1024
)

// Args unwraps a tool request's arguments.
func Args(req mcp.CallToolRequest) (map[string]any, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid request: arguments must be an object")
	}
	return args, nil
}

// Str reads an optional string argument.
func Str(args map[string]any, key string) string {
	s, _ := args[key].(string)
	return s
}

// ReqStr reads a required string argument.
func ReqStr(args map[string]any, key string) (string, error) {
	s, _ := args[key].(string)
	if s == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return s, nil
}

// Bool reads an optional boolean argument.
//
// Tool schemas declare flags with mcp.WithBoolean — the schema should say what
// a param actually is. A string is accepted too, because a model that sends
// "true" against a boolean schema has expressed the intent unambiguously and
// failing it would buy nothing. Anything else is false.
//
// This is the one place the two wire forms are reconciled; handlers must not
// compare against "true" themselves.
func Bool(args map[string]any, key string) bool {
	switch v := args[key].(type) {
	case bool:
		return v
	case string:
		b, err := strconv.ParseBool(v)
		return err == nil && b
	default:
		return false
	}
}

// Ref reads an argument that may be either an ID or a handle.
//
// Neither Str nor ID fits: ID cannot parse a handle, and Str silently returns
// "" for a JSON number, which would quietly drop a lookup into list mode
// instead of erroring. So a numeric argument is rejected here for the same
// reason ID rejects one — a JSON number has already lost precision, and the
// caller meant a specific record.
func Ref(args map[string]any, key string) (string, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return "", nil
	}

	s, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string to avoid precision loss, got %T", key, raw)
	}
	return s, nil
}

// ReqRef is Ref for a required reference.
func ReqRef(args map[string]any, key string) (string, error) {
	s, err := Ref(args, key)
	if err != nil {
		return "", err
	}
	if s == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return s, nil
}

// ID parses a string-encoded uint64 identifier.
//
// Human IDs are uint64 and exceed JavaScript's safe integer range, so every ID
// crosses the tool boundary as a string. Accepting a JSON number here would
// silently truncate, which is close to impossible to diagnose downstream — so
// a numeric argument is rejected rather than coerced.
func ID(args map[string]any, key string) (uint64, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return 0, nil
	}

	s, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("%s must be a string to avoid precision loss, got %T", key, raw)
	}
	if s == "" {
		return 0, nil
	}

	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return id, nil
}

// ReqID is ID for a required identifier.
func ReqID(args map[string]any, key string) (uint64, error) {
	id, err := ID(args, key)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, fmt.Errorf("%s is required", key)
	}
	return id, nil
}

// Paging carries the normalised paging arguments for a lookup tool.
type Paging struct {
	Limit  uint
	Cursor string
}

// Page reads limit and pageCursor, applying the default and the cap.
//
// A caller asking for more than MaxLimit gets MaxLimit rather than an error:
// the request is reasonable, the number is not, and failing the call would just
// cost a round trip.
func Page(args map[string]any) Paging {
	p := Paging{Limit: DefaultLimit, Cursor: Str(args, "pageCursor")}

	switch v := args["limit"].(type) {
	case string:
		if n, err := strconv.ParseUint(v, 10, 32); err == nil && n > 0 {
			p.Limit = uint(n)
		}
	case float64: // JSON numbers decode as float64
		if v > 0 {
			p.Limit = uint(v)
		}
	}

	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
	return p
}

// JSONResult marshals v as the tool's result.
//
// This is the only sanctioned way to turn a Go value into a tool result. See
// the package comment for why that matters.
func JSONResult(v any) (*mcp.CallToolResult, error) {
	out, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	if len(out) > MaxResultBytes {
		return nil, fmt.Errorf(
			"result too large: %d bytes exceeds the %d byte limit — narrow the query, or use limit and pageCursor",
			len(out), MaxResultBytes,
		)
	}

	return mcp.NewToolResultText(string(out)), nil
}

// TextResult returns a plain-text tool result. Used where JSON would be noise —
// delete handlers acknowledging an action, for instance.
func TextResult(format string, a ...any) *mcp.CallToolResult {
	return mcp.NewToolResultText(fmt.Sprintf(format, a...))
}

// Errf wraps a service error in the house style: lowercase subject and verb,
// wrapped cause, no trailing punctuation.
//
// Tool errors are read by a model deciding what to do next, so they say what
// failed rather than where in the code it failed.
func Errf(subject string, err error) error {
	return fmt.Errorf("%s failed: %w", subject, err)
}
