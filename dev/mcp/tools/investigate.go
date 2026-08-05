package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

type governingReport struct {
	Files []governedFile `json:"files"`
	Note  string         `json:"note"`
}

type governedFile struct {
	File string         `json:"file"`
	Docs []governingDoc `json:"docs,omitempty"`
	Note string         `json:"note,omitempty"`
}

// governingDoc carries the parts of a contract that must not be paraphrased.
//
// Locked and WIP are returned verbatim rather than summarised on purpose: a
// locked contract restated in someone's own words is how it stops being the
// contract, and a WIP marker is a warning that only works if the caller reads
// what it actually says.
type governingDoc struct {
	Doc    string `json:"doc"`
	Kind   string `json:"kind"`
	Scope  string `json:"scope"`
	Covers string `json:"covers"`

	Locked  []string `json:"lockedContracts,omitempty"`
	WIP     []string `json:"wip,omitempty"`
	Sibling bool     `json:"sibling,omitempty"`
}

type affectedReport struct {
	Specs     []string `json:"specs"`
	UnitTests []string `json:"unitTests,omitempty"`
	Note      string   `json:"note"`
}

func registerIntentGoverning(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_intent_governing",
			mcp.WithDescription(
				"Find the intent docs that govern the files you are about to change, and get their locked "+
					"contracts and WIP markers back word for word. Read this before editing, not after: a "+
					"locked contract is a change the human has to rule on rather than a change you make, and "+
					"a WIP zone is ground that has not been decided yet. Docs are returned nearest-first — a "+
					"sibling doc for the file itself, then each folder doc up to the repo root, which is "+
					"where an area's constitution lives.",
			),
			mcp.WithString("files", mcp.Required(), mcp.Description(
				"Space-separated repo-relative paths you intend to change.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskRead),
		),
		"Find governing intent docs",
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, err := toolkit.Args(req)
			if err != nil {
				return nil, err
			}

			raw, err := toolkit.ReqStr(args, "files")
			if err != nil {
				return nil, err
			}

			return toolkit.JSONResult(governingDocs(root, strings.Fields(raw)))
		},
	)
}

func governingDocs(root string, files []string) governingReport {
	out := governingReport{Files: make([]governedFile, 0, len(files))}

	governed := 0
	for _, f := range files {
		g := governedFile{File: f, Docs: docsFor(root, f)}
		if len(g.Docs) == 0 {
			g.Note = "no intent doc governs this path — it is outside the enforced roots, so there is no " +
				"recorded contract to check against"
		} else {
			governed++
		}
		out.Files = append(out.Files, g)
	}

	locked := 0
	for _, f := range out.Files {
		for _, d := range f.Docs {
			locked += len(d.Locked)
		}
	}

	switch {
	case governed == 0:
		out.Note = "none of these paths are covered by an intent doc"
	case locked > 0:
		out.Note = fmt.Sprintf("%d locked contract section(s) apply. A change that touches one is an intent "+
			"change, not a code change: the human rules on it first. Docs marked covers:governs bind this "+
			"file directly; covers:area-context are the area's constitution, which does not formally cover "+
			"the path but whose contracts still hold", locked)
	default:
		out.Note = "covered, with no locked contracts in scope"
	}

	return out
}

// docsFor walks from the file's directory to the repo root, collecting the
// intent docs that apply, nearest first.
func docsFor(root, file string) []governingDoc {
	var out []governingDoc

	clean := filepath.Clean(strings.TrimPrefix(file, "./"))
	dir := filepath.Dir(clean)

	for {
		abs := filepath.Join(root, dir)

		entries, err := os.ReadDir(abs)
		if err == nil {
			names := make([]string, 0, 2)
			for _, e := range entries {
				if !e.IsDir() && strings.HasSuffix(e.Name(), ".intent.md") {
					names = append(names, e.Name())
				}
			}
			sort.Strings(names)

			for _, name := range names {
				doc, err := readGoverningDoc(filepath.Join(abs, name), filepath.Join(dir, name))
				if err != nil {
					continue
				}

				// The doc's own frontmatter says what it governs, and guessing
				// from the filename instead was wrong twice over: a kind:file
				// doc names its target in covers (which may differ in extension
				// from the doc's own name), and treating every other doc in the
				// directory as a folder doc returned all seventeen store
				// contracts for one store.
				if doc.Kind == "file" {
					if filepath.Base(clean) != doc.Scope {
						continue
					}
					doc.Sibling = true
				}

				// Formal coverage and useful context are not the same thing, so
				// both are returned and labelled rather than merged.
				//
				// A folder doc governs a nested file only when it covers
				// recursively. An ancestor that covers just its own directory —
				// which is how an area constitution like project.intent.md is
				// declared — does not govern the file by the intent system's
				// own rules, but its locked contracts still bind the area and
				// the change procedure says to read it every time. Dropping it
				// silently would hide exactly the contracts a caller must not
				// break.
				switch {
				case doc.Sibling || (doc.Kind == "folder" && dir == filepath.Dir(clean)):
					doc.Covers = "governs"
				case doc.Kind == "folder" && doc.Scope == "recursive":
					doc.Covers = "governs"
				case doc.Kind == "folder" && (len(doc.Locked) > 0 || len(doc.WIP) > 0):
					doc.Covers = "area-context"
				default:
					continue
				}

				out = append(out, doc)
			}
		}

		if dir == "." || dir == string(filepath.Separator) {
			break
		}
		dir = filepath.Dir(dir)
	}

	return out
}

func readGoverningDoc(abs, rel string) (governingDoc, error) {
	raw, err := os.ReadFile(abs)
	if err != nil {
		return governingDoc{}, err
	}

	body := string(raw)
	doc := governingDoc{
		Doc:   rel,
		Kind:  frontmatterValue(body, "kind"),
		Scope: strings.Trim(frontmatterValue(body, "covers"), `'"`),
	}

	doc.Locked = sectionsMatching(body, "locked")
	doc.WIP = wipMarkers(body)

	return doc, nil
}

// frontmatterValue reads one key out of the leading --- block.
func frontmatterValue(body, key string) string {
	if !strings.HasPrefix(body, "---") {
		return ""
	}

	end := strings.Index(body[3:], "\n---")
	if end < 0 {
		return ""
	}

	for _, line := range strings.Split(body[3:3+end], "\n") {
		if k, v, ok := strings.Cut(line, ":"); ok && strings.TrimSpace(k) == key {
			return strings.TrimSpace(v)
		}
	}

	return ""
}

// sectionsMatching returns whole markdown sections whose heading contains the
// word, verbatim including the heading.
func sectionsMatching(body, word string) []string {
	var out []string

	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, "#") || !strings.Contains(strings.ToLower(line), word) {
			continue
		}

		depth := len(line) - len(strings.TrimLeft(line, "#"))

		section := []string{line}
		for _, next := range lines[i+1:] {
			if strings.HasPrefix(next, "#") {
				if d := len(next) - len(strings.TrimLeft(next, "#")); d <= depth {
					break
				}
			}
			section = append(section, next)
		}

		out = append(out, strings.TrimSpace(strings.Join(section, "\n")))
	}

	return out
}

// wipMarkers returns the lines that flag undecided ground, verbatim.
func wipMarkers(body string) []string {
	var out []string

	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.Contains(trimmed, "WIP") || strings.Contains(trimmed, "DRIFT") {
			out = append(out, trimmed)
		}
	}

	return out
}

func registerIntentAffected(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_intent_affected",
			mcp.WithDescription(
				"List the e2e specs a set of changed files affects. e2e is slow and human-run in this repo, "+
					"so the point is to tell the human exactly which specs to run rather than to run them: "+
					"report the list, do not launch playwright yourself.",
			),
			mcp.WithString("files", mcp.Description(
				"Space-separated repo-relative paths. Omit to use whatever git reports as changed.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskRead),
		),
		"Find affected e2e specs",
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

			if len(files) == 0 {
				return toolkit.JSONResult(affectedReport{Specs: []string{}, Note: "nothing changed"})
			}

			stdout, _ := runAllowFail(ctx, root, "node", append([]string{".intent/intent.mjs", "affected"}, files...)...)

			// The intent CLI answers with every test file it associates with the
			// input, which is not all e2e: a server change comes back with
			// lib/js unit tests. Telling the human to run playwright on a
			// vitest file wastes their time and makes the tool look careless,
			// so the two are separated by path and only e2e is handed over.
			out := affectedReport{Specs: []string{}, UnitTests: []string{}}
			for _, line := range strings.Split(stdout, "\n") {
				line = strings.TrimSpace(line)
				if !strings.HasSuffix(line, ".ts") && !strings.HasSuffix(line, ".js") {
					continue
				}

				if strings.Contains(line, "/e2e/") {
					out.Specs = append(out.Specs, line)
				} else {
					out.UnitTests = append(out.UnitTests, line)
				}
			}

			switch {
			case len(out.Specs) > 0:
				out.Note = fmt.Sprintf("%d e2e spec(s) affected. Hand these to the human to run with "+
					"'npx playwright test <specs>' — do not run them yourself, they are slow", len(out.Specs))
			case len(out.UnitTests) > 0:
				out.Note = fmt.Sprintf("no e2e specs, but %d unit test file(s) are associated. Run those "+
					"with dev_test_run", len(out.UnitTests))
			default:
				out.Note = "no tests are associated with these files"
			}

			return toolkit.JSONResult(out)
		},
	)
}
