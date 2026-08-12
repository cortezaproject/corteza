package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// uiResult is what the browser saw.
//
// Console errors and failed requests are the payload, not the decoration: a
// click that "does nothing" almost always leaves one or the other behind, and
// neither is visible from the final URL.
type uiResult struct {
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`

	// Where the caller asked to go and where the app was once it had settled.
	// They differ whenever a route redirects, and — the case worth catching —
	// whenever no route matched at all, since the webapp's catch-all sends an
	// unknown path to the home section instead of showing anything broken.
	RequestedPath string `json:"requestedPath,omitempty"`
	LandedPath    string `json:"landedPath,omitempty"`

	ConsoleErrors  []string `json:"consoleErrors,omitempty"`
	FailedRequests []string `json:"failedRequests,omitempty"`

	// Errors belonging to a page load that was replaced before it finished,
	// kept apart so they cannot be mistaken for the rendered page's own.
	SupersededErrors []string `json:"supersededErrors,omitempty"`
	SupersededNote   string   `json:"supersededNote,omitempty"`

	Steps      []uiStep `json:"steps,omitempty"`
	Screenshot string   `json:"screenshot,omitempty"`
	Error      string   `json:"error,omitempty"`
	Hint       string   `json:"hint,omitempty"`
	Note       string   `json:"note,omitempty"`
}

type uiStep struct {
	Action   string `json:"action,omitempty"`
	Selector string `json:"selector,omitempty"`
	Value    string `json:"value,omitempty"`
	OK       bool   `json:"ok"`
	URLAfter string `json:"urlAfter,omitempty"`
	Error    string `json:"error,omitempty"`
}

func registerUIVerify(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_ui_verify",
			mcp.WithDescription(
				"Drive the webapp in a real browser and report what happened: final URL, console errors, "+
					"failed requests, and a screenshot. This is the only tool that exercises frontend code — "+
					"the Human MCP tools reach the REST API and never load the app, so an MCP-green result "+
					"says nothing about a defect a user clicks into. Use it to reproduce a UI bug before "+
					"diagnosing it, and to verify the fix afterwards. It logs in as the agent-dev user; run "+
					"dev/agent/bootstrap.sh once if it reports no password.",
			),
			mcp.WithString("path", mcp.Description(
				"App path to open, e.g. '/compose/namespace/my-namespace/pages/123' — every compose route "+
					"is under /compose, and a namespace is addressed by slug. Defaults to '/'. A path no route "+
					"matches is not an error: the app's catch-all redirects it to the home section, which the "+
					"result reports as landedPath differing from requestedPath.")),
			mcp.WithString("steps", mcp.Description(
				"JSON array of interactions, run in order, e.g. "+
					`[{"action":"click","selector":"button:has-text('Add')"}]. `+
					"action is 'click' or 'fill'; fill takes a value. Selectors are playwright selectors, so "+
					"text matching works. The run stops at the first step that fails and reports why.")),
			mcp.WithBoolean("screenshot", mcp.Description(
				"Capture a screenshot. On by default; the path is returned and can be read with the Read tool.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskWrite),
		),
		"Drive the webapp",
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, err := toolkit.Args(req)
			if err != nil {
				return nil, err
			}

			payload := map[string]any{
				"root":    root,
				"path":    toolkit.Str(args, "path"),
				"baseURL": webappURL(root),
			}

			if v, present := toolkit.OptBool(args, "screenshot"); present {
				payload["screenshot"] = v
			}

			var steps []map[string]any
			if _, err := toolkit.JSONArg(args, "steps", "steps", &steps); err != nil {
				return nil, err
			}
			if steps != nil {
				payload["steps"] = steps
			}

			encoded, err := json.Marshal(payload)
			if err != nil {
				return nil, fmt.Errorf("cannot encode the request: %w", err)
			}

			// npx from the webapp workspace, which is where playwright and its
			// browsers are installed.
			raw, runErr := runAllowFail(ctx, filepath.Join(root, "client", "web", "unify"),
				"node", filepath.Join(root, "dev", "mcp", "uiverify.mjs"), string(encoded))

			out, parseErr := parseUIResult(raw)
			if parseErr != nil {
				return nil, fmt.Errorf("the browser driver produced no result (%v): %s",
					runErr, firstLines(strings.TrimSpace(raw), 15))
			}

			out.Note = uiNote(out)

			return toolkit.JSONResult(out)
		},
	)
}

// parseUIResult reads the driver's JSON, which is the last line of its output —
// playwright and node write their own noise to stdout on occasion.
func parseUIResult(raw string) (uiResult, error) {
	var out uiResult

	lines := nonEmptyLines(raw)
	for i := len(lines) - 1; i >= 0; i-- {
		if err := json.Unmarshal([]byte(lines[i]), &out); err == nil {
			return out, nil
		}
	}

	return out, fmt.Errorf("no JSON in output")
}

func uiNote(out uiResult) string {
	if out.Error != "" {
		return "the run did not complete; nothing here describes working behaviour"
	}

	var failed *uiStep
	for i := range out.Steps {
		if !out.Steps[i].OK {
			failed = &out.Steps[i]
			break
		}
	}

	switch {
	case failed != nil:
		return fmt.Sprintf("step %q failed, so any step after it never ran", failed.Selector)
	case leftRequestedPath(out):
		// Ranked above console errors on purpose: everything else in the result
		// describes a page that is not the one that was asked for, so reading it
		// as a verdict on that page is the wrong conclusion to draw first.
		return fmt.Sprintf(
			"the app did not stay on %s — it settled on %s, so the screenshot and every error below "+
				"describe that page instead. Either the route redirects deliberately or nothing matched the "+
				"path and the catch-all sent the app home; check the path against the section's routes "+
				"(compose pages live under /compose/namespace/<slug>/pages/<pageID>)",
			out.RequestedPath, out.LandedPath,
		)
	case len(out.ConsoleErrors) > 0:
		return "the steps ran, but the page logged console errors — read them before calling this a pass, " +
			"since a click that appears to do nothing usually leaves one here"
	case len(out.Steps) == 0:
		return "page loaded; no interactions were requested"
	default:
		return "all steps ran with no console errors. Check url and the screenshot against what you expected — " +
			"a step that clicks successfully has not necessarily done the right thing"
	}
}

// leftRequestedPath reports whether the app ended up somewhere other than the
// path the caller asked for. A trailing slash is not a difference, and neither
// is a query string — the app appends its own list parameters after mount.
func leftRequestedPath(out uiResult) bool {
	if out.RequestedPath == "" || out.LandedPath == "" {
		return false
	}

	trim := func(p string) string {
		if i := strings.IndexAny(p, "?#"); i >= 0 {
			p = p[:i]
		}
		return strings.TrimSuffix(p, "/")
	}

	return trim(out.RequestedPath) != trim(out.LandedPath)
}

// webappURL is where vite serves the app.
func webappURL(root string) string {
	if raw, err := readFirstMatch(filepath.Join(root, "client", "web", "unify", "playwright.config.ts"),
		"E2E_BASE_URL || '", "'"); err == nil && raw != "" {
		return raw
	}

	return "http://localhost:5173"
}
