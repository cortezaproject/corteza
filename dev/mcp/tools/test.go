package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// testReport is what dev_test_run returns.
//
// Failures only, deliberately. A green `go test ./automation/...` is four
// thousand lines to say "exit 0", and reading it into a model's context is the
// single largest avoidable cost in this repo's dev loop. What a caller needs is
// whether it passed and, if not, exactly what broke and where.
type testReport struct {
	Suite      string       `json:"suite"`
	Target     string       `json:"target"`
	Passed     bool         `json:"passed"`
	Packages   int          `json:"packages,omitempty"`
	Failures   []testFail   `json:"failures,omitempty"`
	Errors     []string     `json:"errors,omitempty"`
	Skipped    int          `json:"skipped,omitempty"`
	Aborted    bool         `json:"aborted,omitempty"`
	Unfinished []string     `json:"unfinished,omitempty"`
	Note       string       `json:"note,omitempty"`
	Timing     *testTimings `json:"timing,omitempty"`
}

type testFail struct {
	Package string `json:"package,omitempty"`
	Test    string `json:"test,omitempty"`
	File    string `json:"file,omitempty"`
	Output  string `json:"output"`

	// Tests that failed with this same error, named but not re-printed.
	AlsoFailing     []string `json:"alsoFailing,omitempty"`
	AlsoFailingMore int      `json:"alsoFailingMore,omitempty"`
}

type testTimings struct {
	Seconds float64 `json:"seconds"`
}

func registerTestRun(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_test_run",
			mcp.WithDescription(
				"Run a package's tests and get back only what failed. Go, vitest and mocha are all handled — "+
					"the runner is chosen from the target path, so 'server/automation/...' runs go test, "+
					"'client/web/unify' and 'lib/vue' run vitest, and 'lib/js' runs mocha because that is what "+
					"it is configured with. A passing run returns passed:true and nothing else, because the "+
					"output of a green suite is not information. A failing one returns each failure with its "+
					"package, test name and the assertion output, which is what you would have had to scroll "+
					"for. Run this rather than shelling out to go test, npx vitest or npx mocha: the raw "+
					"output of any of them is thousands of lines.",
			),
			mcp.WithString("target", mcp.Required(), mcp.Description(
				"Repo-relative path to test. Go: a package or pattern such as 'server/automation/...' or "+
					"'server/pkg/mcpkit'. JS: a workspace directory such as 'client/web/unify' or 'lib/vue', "+
					"optionally with a spec path after it.")),
			mcp.WithString("run", mcp.Description(
				"Only run tests whose name matches this. Go passes it to -run, vitest to -t. Use it to "+
					"re-run one failure without paying for the whole package.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskRead),
		),
		"Run tests",
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, err := toolkit.Args(req)
			if err != nil {
				return nil, err
			}

			target, err := toolkit.ReqStr(args, "target")
			if err != nil {
				return nil, err
			}

			var report testReport
			switch runnerFor(root, target) {
			case "go":
				module, _ := goModuleFor(target)
				report, err = runGoTests(ctx, root, module, target, toolkit.Str(args, "run"))
			case "mocha":
				report, err = runMocha(ctx, root, target, toolkit.Str(args, "run"))
			default:
				report, err = runVitest(ctx, root, target, toolkit.Str(args, "run"))
			}
			if err != nil {
				return nil, err
			}

			if folded := len(report.Failures); folded > 0 {
				report.Failures = collapseFailures(report.Failures)

				if folded -= len(report.Failures); folded > 0 {
					report.Note = strings.TrimSpace(report.Note + fmt.Sprintf(
						" %d further failure(s) carried an error already listed and were folded into"+
							" its alsoFailing rather than repeated. They are almost certainly one cause,"+
							" so fix the listed error and re-run before reading anything into the count.",
						folded))
				}
			}

			return toolkit.JSONResult(report)
		},
	)
}

// goModules maps a target prefix to the module directory its tests run in.
//
// Two Go modules, and the developer MCP is one of them — a tool that could not
// run its own package's tests would be an odd thing to ship.
var goModules = []struct{ prefix, dir string }{
	{"dev/mcp", "dev/mcp"},
	{"server", "server"},
}

// goModuleFor decides which runner to use and where to run it. An empty dir
// means this is not Go and a JS runner takes it.
func goModuleFor(target string) (dir string, ok bool) {
	clean := strings.TrimPrefix(target, "./")

	for _, m := range goModules {
		if clean == m.prefix || strings.HasPrefix(clean, m.prefix+"/") {
			return m.dir, true
		}
	}

	return "", false
}

// runnerFor names the suite a target belongs to: "go", "mocha" or "vitest".
//
// Which JS runner a workspace uses is read off disk rather than listed here,
// because it is a fact about the workspace and drifts with it.
func runnerFor(root, target string) string {
	if _, ok := goModuleFor(target); ok {
		return "go"
	}

	if dir, _ := splitWorkspace(target); mochaConfig(root, dir) != "" {
		return "mocha"
	}

	return "vitest"
}

// mochaConfig returns the workspace's mocha config, or "" if it has none.
//
// Which JS runner a workspace uses is a fact about the workspace, so it is read
// off disk rather than hard-coded: lib/js runs mocha while lib/vue and the
// webapp run vitest, and sending a mocha workspace to vitest produces "no tests
// ran: nothing matched <target>" — a message that reads as a wrong path and
// sends the caller hunting for a target that was right all along.
func mochaConfig(root, dir string) string {
	if dir == "" {
		return ""
	}

	for _, name := range []string{".mocharc.js", ".mocharc.cjs", ".mocharc.json", ".mocharc.yml"} {
		path := filepath.Join(root, dir, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// runGoTests runs go test and keeps only the failures.
//
// -json is not a convenience here, it is the whole mechanism: parsing human
// output means guessing at line prefixes, while the event stream says exactly
// which test failed, in which package, with which output attached. A non-zero
// exit is expected on failure and is not itself an error.
func runGoTests(ctx context.Context, root, module, target, run string) (testReport, error) {
	out := testReport{Suite: "go", Target: target}

	clean := strings.TrimPrefix(target, "./")

	pkg := "./" + strings.TrimPrefix(strings.TrimPrefix(clean, module), "/")
	if pkg == "./" {
		pkg = "./..."
	}

	argv := []string{"test", "-json", pkg}
	if run != "" {
		argv = append(argv, "-run", run)
	}

	cmd := exec.CommandContext(ctx, "go", argv...)
	cmd.Dir = root + "/" + module

	var stderr strings.Builder
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()

	// A build failure produces no usable events, so it surfaces as an error
	// rather than an empty pass — the failure mode that would otherwise report
	// "passed" for code that does not compile.
	if len(stdout) == 0 {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return out, fmt.Errorf("go test %s: %s", pkg, firstLines(msg, 20))
		}
		if err != nil {
			return out, fmt.Errorf("go test %s: %w", pkg, err)
		}
	}

	type event struct {
		Action string `json:"Action"`

		Package string  `json:"Package"`
		Test    string  `json:"Test"`
		Output  string  `json:"Output"`
		Elapsed float64 `json:"Elapsed"`

		// A compile error is not a test failure and does not arrive as one: go
		// test emits build-output/build-fail events keyed by ImportPath, and
		// the package's own event says only "[build failed]". Without these the
		// tool would report broken code with no idea where.
		ImportPath  string `json:"ImportPath"`
		FailedBuild string `json:"FailedBuild"`
	}

	// Output arrives in fragments, so it is accumulated per test and only kept
	// for the ones that end up failing.
	buf := map[string]*strings.Builder{}
	seen := map[string]bool{}
	builds := map[string]*strings.Builder{}

	// A panic takes the whole test binary down, so every test after it never
	// runs and go test never reports them. Tracking which tests started but
	// never resolved is the only way to notice: without it the tool says
	// "1 failure" for a run where five more tests silently did not execute,
	// which reads as "the rest were fine".
	started := map[string]bool{}
	aborted := false

	for _, line := range strings.Split(string(stdout), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var e event
		if json.Unmarshal([]byte(line), &e) != nil {
			// A line that is not an event is build noise; go test interleaves
			// it, and it matters only when nothing else explains a failure.
			out.Errors = appendCapped(out.Errors, strings.TrimSpace(line))
			continue
		}

		if e.Action == "build-output" {
			if _, ok := builds[e.ImportPath]; !ok {
				builds[e.ImportPath] = &strings.Builder{}
			}
			builds[e.ImportPath].WriteString(e.Output)
			continue
		}

		key := e.Package + "\x00" + e.Test

		switch e.Action {
		case "run":
			if e.Test != "" {
				started[key] = true
			}

		case "output":
			if strings.HasPrefix(e.Output, "panic:") || strings.Contains(e.Output, "[signal SIGSEGV") {
				aborted = true
			}
			if _, ok := buf[key]; !ok {
				buf[key] = &strings.Builder{}
			}
			buf[key].WriteString(e.Output)

		case "pass":
			if e.Test == "" {
				seen[e.Package] = true
				out.Timing = addSeconds(out.Timing, e.Elapsed)
			}
			delete(started, key)
			delete(buf, key)

		case "skip":
			if e.Test != "" {
				out.Skipped++
			}
			delete(started, key)
			delete(buf, key)

		case "fail":
			if e.Test == "" {
				seen[e.Package] = true
				// A package-level failure with no failing test inside it is a
				// build or a panic; keep its output, since nothing else has it.
				if !hasFailureIn(out.Failures, e.Package) {
					fail := testFail{Package: e.Package, Output: trimOutput(buf[key])}

					// FailedBuild names the build whose output holds the actual
					// compiler error, which is the only useful thing here.
					if b, ok := builds[e.FailedBuild]; ok {
						fail.Output = trimString(b.String())
						fail.File = firstFileRef(b)
					}

					out.Failures = append(out.Failures, fail)
				}
				delete(buf, key)
				continue
			}

			out.Failures = append(out.Failures, testFail{
				Package: e.Package,
				Test:    e.Test,
				File:    firstFileRef(buf[key]),
				Output:  trimOutput(buf[key]),
			})
			delete(started, key)
			delete(buf, key)
		}
	}

	for key := range started {
		if _, test, ok := strings.Cut(key, "\x00"); ok && test != "" {
			out.Unfinished = append(out.Unfinished, test)
		}
	}
	sort.Strings(out.Unfinished)

	out.Aborted = aborted || len(out.Unfinished) > 0

	out.Packages = len(seen)
	out.Passed = len(out.Failures) == 0 && !out.Aborted

	// A setup failure — a target that matches nothing, a bad flag — has no
	// build output to borrow, and go test explains it on stderr.
	if msg := strings.TrimSpace(stderr.String()); msg != "" {
		for i, f := range out.Failures {
			if strings.Contains(f.Output, "setup failed") {
				out.Failures[i].Output = strings.TrimSpace(f.Output + "\n" + firstLines(msg, maxOutputLines))
			}
		}
	}

	switch {
	case out.Aborted:
		// The distinction matters: a test left mid-flight can at least be
		// named, while a test that never started emits no event at all and is
		// invisible. Both mean the same thing for the caller — this result is
		// not the whole package — but only one of them can be listed.
		out.Note = "The test binary aborted, most likely a panic. "
		if len(out.Unfinished) > 0 {
			out.Note += fmt.Sprintf("%d test(s) were left mid-flight (see unfinished), and ", len(out.Unfinished))
		}
		out.Note += "Any test scheduled after the abort never started — go test emits nothing for those, so " +
			"they cannot be listed and this result is NOT a full picture of the package. Re-run with the " +
			"'run' argument excluding the failing test to see the rest."
	case out.Passed && out.Packages == 0:
		out.Note = "no test packages matched " + pkg + " — check the target path"
	}

	return out, nil
}

// runMocha runs a mocha workspace's suite and keeps only the failures.
//
// mocha's json reporter reports `stats` alongside a `failures` array, which is
// everything needed; as with vitest the object is located in the stream rather
// than assumed to be all of it, since a test's own console output shares stdout.
func runMocha(ctx context.Context, root, target, run string) (testReport, error) {
	out := testReport{Suite: "mocha", Target: target}

	dir, spec := splitWorkspace(target)

	argv := []string{"mocha", "--reporter", "json"}
	if spec != "" {
		argv = append(argv, spec)
	}
	if run != "" {
		argv = append(argv, "--grep", run)
	}

	cmd := exec.CommandContext(ctx, "npx", argv...)
	cmd.Dir = filepath.Join(root, dir)

	var stderr strings.Builder
	cmd.Stderr = &stderr
	stdout, _ := cmd.Output()

	raw := extractJSON(string(stdout))
	if raw == "" {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(string(stdout))
		}
		return out, fmt.Errorf("mocha in %s produced no report: %s", dir, firstLines(msg, 20))
	}

	var report struct {
		Stats struct {
			Suites   int `json:"suites"`
			Tests    int `json:"tests"`
			Pending  int `json:"pending"`
			Failures int `json:"failures"`
		} `json:"stats"`
		Failures []struct {
			FullTitle string `json:"fullTitle"`
			File      string `json:"file"`
			Err       struct {
				Message string `json:"message"`
				Stack   string `json:"stack"`
			} `json:"err"`
		} `json:"failures"`
	}

	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		return out, fmt.Errorf("cannot read mocha report: %w", err)
	}

	for _, f := range report.Failures {
		output := f.Err.Stack
		if output == "" {
			output = f.Err.Message
		}

		out.Failures = append(out.Failures, testFail{
			Test:   f.FullTitle,
			File:   relativeTo(root, f.File),
			Output: trimString(output),
		})
	}

	out.Packages = report.Stats.Suites
	out.Skipped = report.Stats.Pending
	out.Passed = len(out.Failures) == 0 && report.Stats.Failures == 0

	// Zero tests is not a pass, for the same reason it is not one under vitest:
	// a spec path that matches nothing exits clean with an empty report.
	if report.Stats.Tests == 0 {
		out.Passed = false
		out.Note = "no tests ran: nothing matched " + target +
			". Check the path is inside the config's spec globs and that the spec name is right."
		return out, nil
	}

	// bail stops the run at the first failure, so everything after it never
	// executed and is absent from the report rather than passing. Reporting one
	// failure without saying so reads as "the rest were fine".
	if len(out.Failures) > 0 && mochaBails(mochaConfig(root, dir)) {
		out.Note = "this workspace's mocha config sets bail, so the run STOPPED at the first failure — " +
			"tests after it never ran and this result is NOT a full picture of the suite. Fix this one and " +
			"re-run to see the rest."
	}

	return out, nil
}

func mochaBails(configPath string) bool {
	if configPath == "" {
		return false
	}

	raw, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}

	// Good enough for the note it drives: a false negative only costs the
	// caveat, and the config is a literal in every form mocha accepts.
	body := strings.Join(strings.Fields(string(raw)), "")

	return strings.Contains(body, "bail:true") || strings.Contains(body, `"bail":true`)
}

// runVitest runs the JS suite for a workspace.
//
// vitest's json reporter writes to stdout alongside its own progress output, so
// the report is found by locating the JSON object rather than assuming the
// whole of stdout is JSON.
func runVitest(ctx context.Context, root, target, run string) (testReport, error) {
	out := testReport{Suite: "vitest", Target: target}

	dir, spec := splitWorkspace(target)

	argv := []string{"vitest", "run", "--reporter=json"}
	if spec != "" {
		argv = append(argv, spec)
	}
	if run != "" {
		argv = append(argv, "-t", run)
	}

	cmd := exec.CommandContext(ctx, "npx", argv...)
	cmd.Dir = root + "/" + dir

	var stderr strings.Builder
	cmd.Stderr = &stderr
	stdout, _ := cmd.Output()

	raw := extractJSON(string(stdout))
	if raw == "" {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(string(stdout))
		}
		return out, fmt.Errorf("vitest in %s produced no report: %s", dir, firstLines(msg, 20))
	}

	var report struct {
		NumTotalTestSuites int  `json:"numTotalTestSuites"`
		NumTotalTests      int  `json:"numTotalTests"`
		NumPendingTests    int  `json:"numPendingTests"`
		Success            bool `json:"success"`
		TestResults        []struct {
			Name             string `json:"name"`
			AssertionResults []struct {
				FullName        string   `json:"fullName"`
				Status          string   `json:"status"`
				FailureMessages []string `json:"failureMessages"`
			} `json:"assertionResults"`
		} `json:"testResults"`
	}

	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		return out, fmt.Errorf("cannot read vitest report: %w", err)
	}

	for _, file := range report.TestResults {
		for _, a := range file.AssertionResults {
			if a.Status != "failed" {
				continue
			}
			out.Failures = append(out.Failures, testFail{
				Test:   a.FullName,
				File:   relativeTo(root, file.Name),
				Output: trimString(strings.Join(a.FailureMessages, "\n")),
			})
		}
	}

	out.Packages = report.NumTotalTestSuites
	out.Skipped = report.NumPendingTests
	out.Passed = len(out.Failures) == 0 && report.Success

	// Zero tests is not a pass. A spec path that matches nothing — outside the
	// config's include globs, or simply misspelt — makes vitest exit 0 with an
	// empty report, and reporting that as green is the exact false negative
	// this tool exists to prevent.
	if report.NumTotalTests == 0 {
		out.Passed = false
		out.Note = "no tests ran: nothing matched " + target +
			". Check the path is inside the workspace's include globs and that the spec name is right."
		return out, nil
	}

	// success:false with nothing parsed means the run broke in a way the
	// assertion list does not describe — an import error, a config failure.
	// Saying so beats returning an empty failure list next to passed:false.
	if !report.Success && len(out.Failures) == 0 {
		out.Note = "vitest reported failure with no failing assertion — likely an import or config error; " +
			"re-run it directly in " + dir + " to see the output"
	}

	return out, nil
}
