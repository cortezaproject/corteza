package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
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
	Suite    string       `json:"suite"`
	Target   string       `json:"target"`
	Passed   bool         `json:"passed"`
	Packages int          `json:"packages,omitempty"`
	Failures []testFail   `json:"failures,omitempty"`
	Errors   []string     `json:"errors,omitempty"`
	Skipped  int          `json:"skipped,omitempty"`
	Note     string       `json:"note,omitempty"`
	Timing   *testTimings `json:"timing,omitempty"`
}

type testFail struct {
	Package string `json:"package,omitempty"`
	Test    string `json:"test,omitempty"`
	File    string `json:"file,omitempty"`
	Output  string `json:"output"`
}

type testTimings struct {
	Seconds float64 `json:"seconds"`
}

func registerTestRun(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_test_run",
			mcp.WithDescription(
				"Run a package's tests and get back only what failed. Go and vitest are both handled — the "+
					"suite is chosen from the target path, so 'server/automation/...' runs go test and "+
					"'client/web/unify' runs vitest. A passing run returns passed:true and nothing else, "+
					"because the output of a green suite is not information. A failing one returns each "+
					"failure with its package, test name and the assertion output, which is what you would "+
					"have had to scroll for. Run this rather than shelling out to go test or npx vitest: the "+
					"raw output of either is thousands of lines.",
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
			if module, ok := goModuleFor(target); ok {
				report, err = runGoTests(ctx, root, module, target, toolkit.Str(args, "run"))
			} else {
				report, err = runVitest(ctx, root, target, toolkit.Str(args, "run"))
			}
			if err != nil {
				return nil, err
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
// means this is not Go and the JS runner takes it.
func goModuleFor(target string) (dir string, ok bool) {
	clean := strings.TrimPrefix(target, "./")

	for _, m := range goModules {
		if clean == m.prefix || strings.HasPrefix(clean, m.prefix+"/") {
			return m.dir, true
		}
	}

	return "", false
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
		case "output":
			if _, ok := buf[key]; !ok {
				buf[key] = &strings.Builder{}
			}
			buf[key].WriteString(e.Output)

		case "pass":
			if e.Test == "" {
				seen[e.Package] = true
				out.Timing = addSeconds(out.Timing, e.Elapsed)
			}
			delete(buf, key)

		case "skip":
			if e.Test != "" {
				out.Skipped++
			}
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
			delete(buf, key)
		}
	}

	out.Packages = len(seen)
	out.Passed = len(out.Failures) == 0

	// A setup failure — a target that matches nothing, a bad flag — has no
	// build output to borrow, and go test explains it on stderr.
	if msg := strings.TrimSpace(stderr.String()); msg != "" {
		for i, f := range out.Failures {
			if strings.Contains(f.Output, "setup failed") {
				out.Failures[i].Output = strings.TrimSpace(f.Output + "\n" + firstLines(msg, maxOutputLines))
			}
		}
	}

	if out.Passed && out.Packages == 0 {
		out.Note = "no test packages matched " + pkg + " — check the target path"
	}

	return out, nil
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
