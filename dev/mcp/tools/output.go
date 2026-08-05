package tools

import (
	"os"
	"regexp"
	"sort"
	"strings"
)

// The point of this package is to hand back less than a shell would, so the
// trimming rules live in one place rather than being re-invented per tool.
const (
	// maxOutputLines is what is kept of a single failure. Go's failure output
	// is the assertion plus a stack; the assertion is at the top and is what
	// identifies the problem, so keeping the head is right.
	maxOutputLines = 40

	// maxErrors bounds unparseable noise. Past a handful it is never the
	// signal, and an unbounded list would defeat the purpose of the tool.
	maxErrors = 10
)

// fileRef matches the file:line prefix Go puts on a failure, so a caller gets
// somewhere to go without reading the whole output.
var fileRef = regexp.MustCompile(`(?m)^\s*([\w./-]+\.go):(\d+):`)

func trimOutput(b *strings.Builder) string {
	if b == nil {
		return ""
	}
	return trimString(b.String())
}

func trimString(s string) string {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return ""
	}

	lines := strings.Split(s, "\n")
	if len(lines) <= maxOutputLines {
		return s
	}

	kept := append(lines[:maxOutputLines:maxOutputLines],
		"... truncated, re-run this one test with the 'run' argument for the rest")

	return strings.Join(kept, "\n")
}

func firstFileRef(b *strings.Builder) string {
	if b == nil {
		return ""
	}
	if m := fileRef.FindStringSubmatch(b.String()); m != nil {
		return m[1] + ":" + m[2]
	}
	return ""
}

func firstLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) <= n {
		return s
	}
	return strings.Join(lines[:n], "\n") + "\n... truncated"
}

func appendCapped(dst []string, v string) []string {
	if v == "" || len(dst) >= maxErrors {
		return dst
	}
	return append(dst, v)
}

func hasFailureIn(ff []testFail, pkg string) bool {
	for _, f := range ff {
		if f.Package == pkg {
			return true
		}
	}
	return false
}

func addSeconds(t *testTimings, s float64) *testTimings {
	if t == nil {
		t = &testTimings{}
	}
	t.Seconds += s
	return t
}

// splitWorkspace separates the package directory from an optional spec path.
//
// The workspace is where npx has to run, since that is where vitest and its
// config live; anything after it is passed through as a filter.
func splitWorkspace(target string) (dir, spec string) {
	target = strings.Trim(target, "/")

	for _, ws := range []string{"client/web/unify", "lib/vue", "lib/js"} {
		if target == ws {
			return ws, ""
		}
		if strings.HasPrefix(target, ws+"/") {
			return ws, strings.TrimPrefix(target, ws+"/")
		}
	}

	return target, ""
}

// extractJSON finds the report object in output that also carries progress
// chatter. vitest prints both to stdout, so taking the whole stream would fail
// to parse.
func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start < 0 || end <= start {
		return ""
	}
	return s[start : end+1]
}

func relativeTo(root, path string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, root), "/")
}

func sortedStrings(in []string) []string {
	sort.Strings(in)
	return in
}

// readFirstMatch pulls the text between two markers out of a file, for reading
// a default out of a config rather than duplicating it here.
func readFirstMatch(path, open, close string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	_, rest, ok := strings.Cut(string(raw), open)
	if !ok {
		return "", nil
	}

	value, _, ok := strings.Cut(rest, close)
	if !ok {
		return "", nil
	}

	return value, nil
}
