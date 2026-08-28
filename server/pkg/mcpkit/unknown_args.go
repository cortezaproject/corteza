package mcpkit

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// rejectUnknownArguments refuses a call carrying a parameter the tool does not
// declare.
//
// An undeclared key used to be dropped without a word, so a caller that misread
// a tool got the call it did not ask for and no way to tell. The case that made
// this worth doing is 'sort': compose_record_report declares it and orders by
// it, compose_record_lookup did not, and passing it there returned records in
// arbitrary order while reporting success. Every typo in a parameter name had
// the same shape — the tool ran, something plausible came back, and the
// instruction went nowhere.
//
// Refusing rather than warning is deliberate. A warning alongside a result is
// still a result, and a caller that did not read the parameter list is not the
// caller who reads the warning.
func (m *MCPServer) rejectUnknownArguments() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, ok := req.Params.Arguments.(map[string]any)
			if !ok || len(args) == 0 {
				return next(ctx, req)
			}

			t, ok := m.reg.tools[ResolveToolAlias(req.Params.Name)]
			if !ok {
				return next(ctx, req)
			}

			declared := t.Tool.InputSchema.Properties
			if len(declared) == 0 {
				// A tool that declares no properties takes whatever it is
				// given; there is nothing to check it against.
				return next(ctx, req)
			}

			var unknown []string
			for k := range args {
				if _, ok := declared[k]; !ok {
					unknown = append(unknown, k)
				}
			}
			if len(unknown) == 0 {
				return next(ctx, req)
			}
			sort.Strings(unknown)

			return nil, unknownArgumentError(req.Params.Name, unknown, declared)
		}
	}
}

func unknownArgumentError(tool string, unknown []string, declared map[string]any) error {
	names := make([]string, 0, len(declared))
	for k := range declared {
		names = append(names, k)
	}
	sort.Strings(names)

	var b strings.Builder
	quoted := make([]string, 0, len(unknown))
	for _, u := range unknown {
		quoted = append(quoted, fmt.Sprintf("%q", u))
	}

	if len(unknown) == 1 {
		fmt.Fprintf(&b, "%s does not take a parameter named %s", tool, quoted[0])
	} else {
		fmt.Fprintf(&b, "%s does not take the parameters %s", tool, strings.Join(quoted, ", "))
	}

	if near := nearestArgument(unknown[0], names); near != "" {
		fmt.Fprintf(&b, "; did you mean %q", near)
	}

	fmt.Fprintf(&b, ". It takes: %s", strings.Join(names, ", "))
	return fmt.Errorf("%s", b.String())
}

// nearestArgument picks the declared name a misspelling most likely meant.
// Case first, then a prefix, then one or two wrong characters — beyond that a
// suggestion is a guess, and a wrong guess costs more than none.
func nearestArgument(given string, names []string) string {
	lower := strings.ToLower(given)

	for _, n := range names {
		if strings.EqualFold(n, given) {
			return n
		}
	}
	for _, n := range names {
		ln := strings.ToLower(n)
		if strings.HasPrefix(ln, lower) || strings.HasPrefix(lower, ln) {
			return n
		}
	}

	best, bestDist := "", 3
	for _, n := range names {
		if d := editDistance(lower, strings.ToLower(n)); d < bestDist {
			best, bestDist = n, d
		}
	}
	return best
}

func editDistance(a, b string) int {
	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(min(curr[j-1]+1, prev[j]+1), prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}
