package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

type formatReport struct {
	Formatted []string `json:"formatted,omitempty"`
	Skipped   []string `json:"skipped,omitempty"`
	Note      string   `json:"note"`
}

type intentReport struct {
	Clean   bool          `json:"clean"`
	Drifted []intentDrift `json:"drifted,omitempty"`
	Total   int           `json:"total"`
	Note    string        `json:"note,omitempty"`
	Yours   []intentDrift `json:"yours,omitempty"`
}

type intentDrift struct {
	File   string `json:"file"`
	Reason string `json:"reason"`
	Doc    string `json:"doc"`
}

func registerFormat(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_format_run",
			mcp.WithDescription(
				"Format the files you changed, and only those. gofmt for Go, prettier for everything else, "+
					"chosen per file. Running a formatter over a whole directory in this repo reformats files "+
					"that were already unformatted before you arrived, which turns a small diff into an "+
					"unreviewable one — that has happened twice and been reverted twice. With no arguments this "+
					"formats what git reports as changed, which is almost always what you want after editing.",
			),
			mcp.WithString("files", mcp.Description(
				"Space-separated repo-relative paths. Omit to format every file git reports as modified, "+
					"staged or untracked.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskWrite),
		),
		"Format changed files",
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, err := toolkit.Args(req)
			if err != nil {
				return nil, err
			}

			files := strings.Fields(toolkit.Str(args, "files"))
			if len(files) == 0 {
				if files, err = changedFiles(ctx, root); err != nil {
					return nil, err
				}
			}

			return toolkit.JSONResult(formatFiles(ctx, root, files))
		},
	)
}

// formatFiles runs the right formatter per file.
//
// A deleted path is skipped rather than failing the whole call: git reports
// deletions as changes, and a caller asking to format its working tree does not
// want the request refused because one file went away.
func formatFiles(ctx context.Context, root string, files []string) formatReport {
	out := formatReport{}

	var goFiles, jsFiles []string
	for _, f := range files {
		switch {
		case strings.HasSuffix(f, ".go"):
			if strings.Contains(f, "/vendor/") || strings.HasSuffix(f, ".gen.go") {
				out.Skipped = append(out.Skipped, f)
				continue
			}
			goFiles = append(goFiles, f)
		case hasAnySuffix(f, ".ts", ".js", ".vue", ".json", ".yaml", ".yml", ".md", ".css", ".scss"):
			if strings.Contains(f, "/node_modules/") {
				out.Skipped = append(out.Skipped, f)
				continue
			}
			jsFiles = append(jsFiles, f)
		default:
			out.Skipped = append(out.Skipped, f)
		}
	}

	if len(goFiles) > 0 {
		if _, err := run(ctx, root, "gofmt", append([]string{"-w"}, goFiles...)...); err == nil {
			out.Formatted = append(out.Formatted, goFiles...)
		} else {
			out.Skipped = append(out.Skipped, goFiles...)
		}
	}

	if len(jsFiles) > 0 {
		argv := append([]string{"prettier", "--write", "--ignore-unknown"}, jsFiles...)
		if _, err := run(ctx, root, "npx", argv...); err == nil {
			out.Formatted = append(out.Formatted, jsFiles...)
		} else {
			out.Skipped = append(out.Skipped, jsFiles...)
		}
	}

	switch {
	case len(out.Formatted) == 0:
		out.Note = "nothing to format"
	default:
		out.Note = fmt.Sprintf("formatted %d file(s); anything under skipped was generated, vendored, "+
			"deleted, or of a kind no formatter here owns", len(out.Formatted))
	}

	return out
}

func registerIntentCheck(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_intent_check",
			mcp.WithDescription(
				"Check which intent docs have drifted from the code they govern. This repo carries a lot of "+
					"pre-existing drift that is nobody's fault today, so the result separates what your own "+
					"changed files are responsible for from the rest — reconcile yours, and do not try to fix "+
					"the baseline. The intent system is opt-in: this reports, it never edits a doc.",
			),
			mcp.WithBoolean("all", mcp.Description(
				"Return every drifted file rather than just the ones your working tree touches. Useful for an "+
					"audit, noisy for ordinary work.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskRead),
		),
		"Check intent drift",
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, err := toolkit.Args(req)
			if err != nil {
				return nil, err
			}

			return intentCheck(ctx, root, toolkit.Bool(args, "all"))
		},
	)
}

// intentCheck runs the repo's own intent CLI and splits the result by blame.
//
// A non-zero exit means drift, which is the normal case here and not an error.
func intentCheck(ctx context.Context, root string, all bool) (*mcp.CallToolResult, error) {
	drifted := intentDrifted(ctx, root)

	mine := map[string]bool{}
	if changed, err := changedFiles(ctx, root); err == nil {
		for _, f := range changed {
			mine[f] = true
		}
	}

	out := intentReport{}

	for _, d := range drifted {
		out.Total++

		if mine[d.File] {
			out.Yours = append(out.Yours, d)
		}
		if all {
			out.Drifted = append(out.Drifted, d)
		}
	}

	out.Clean = out.Total == 0

	switch {
	case out.Clean:
		out.Note = "no drift"
	case len(out.Yours) == 0:
		out.Note = fmt.Sprintf("%d drifted file(s), none of them yours — this is the repo's pre-existing "+
			"baseline, leave it alone", out.Total)
	default:
		out.Note = fmt.Sprintf("%d of %d drifted file(s) are in your working tree: reconcile the docs listed "+
			"under yours, then run 'node .intent/intent.mjs sync <files>'", len(out.Yours), out.Total)
	}

	return toolkit.JSONResult(out)
}

// intentDrifted runs the repo's own intent CLI and returns every file whose doc
// has drifted from it. Shared with the commit tool, which uses it to notice an
// intent change left out of the commit that made it necessary.
func intentDrifted(ctx context.Context, root string) []intentDrift {
	// Non-zero exit is the signal, not an error: intent check fails precisely
	// when it found drift to report.
	stdout, _ := runAllowFail(ctx, root, "node", ".intent/intent.mjs", "check")

	var out []intentDrift

	for _, line := range strings.Split(stdout, "\n") {
		// Each drift line is "<file>  (<reason>)  → update & sync: <doc>". The
		// arrow is what distinguishes it from headings and the summary.
		arrow := strings.Index(line, "→")
		if arrow < 0 {
			continue
		}

		head := strings.TrimSpace(line[:arrow])
		doc := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line[arrow+len("→"):]), "update & sync:"))

		file, reason := head, ""
		if i := strings.Index(head, "("); i >= 0 {
			file = strings.TrimSpace(head[:i])
			reason = strings.Trim(strings.TrimSpace(head[i:]), "()")
		}

		out = append(out, intentDrift{File: file, Reason: reason, Doc: doc})
	}

	return out
}

// changedFiles is every path git considers changed: staged, modified or
// untracked, with deletions left out because there is nothing to act on.
func changedFiles(ctx context.Context, root string) ([]string, error) {
	git := gitRunner(root)

	status, err := git(ctx, "status", "--porcelain=v1")
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	out := make([]string, 0, 8)

	for _, line := range strings.Split(status, "\n") {
		if len(line) < 4 {
			continue
		}

		x, y, path := line[0], line[1], strings.TrimSpace(line[3:])
		if i := strings.Index(path, " -> "); i >= 0 {
			path = path[i+len(" -> "):]
		}

		if x == 'D' || y == 'D' || seen[path] {
			continue
		}

		// An untracked directory is reported as one entry with a trailing
		// slash; expanding it is git's job, not ours.
		if strings.HasSuffix(path, "/") {
			if listed, err := git(ctx, "ls-files", "--others", "--exclude-standard", path); err == nil {
				for _, f := range strings.Split(listed, "\n") {
					if f = strings.TrimSpace(f); f != "" && !seen[f] {
						seen[f] = true
						out = append(out, f)
					}
				}
			}
			continue
		}

		seen[path] = true
		out = append(out, path)
	}

	return out, nil
}

func hasAnySuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}
