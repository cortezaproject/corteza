package tools

import "testing"

// One broken import takes every test in a file down with it, each carrying a
// full copy of the same stack. The caller needs the error once and the names of
// what it took with it — not the stack fifteen times.

func TestCollapseFailuresFoldsOneErrorAcrossTests(t *testing.T) {
	const boom = `Error: [vitest] No "useHistoryBack" export is defined on the "@planetcrust/human-vue" mock.
    at VitestMocker.createError (…/execute.js:284:17)
    at Object.get (…/execute.js:330:16)`

	in := []testFail{
		{Test: "layout evaluates conditions", File: "RecordView.layout.test.js", Output: boom},
		{Test: "layout re-evaluates on swap", File: "RecordView.layout.test.js", Output: boom},
		{Test: "title uses the page title", File: "RecordView.title.test.js", Output: boom},
	}

	got := collapseFailures(in)

	if len(got) != 1 {
		t.Fatalf("expected the three to fold into one, got %d entries", len(got))
	}
	if got[0].Output != boom {
		t.Errorf("the surviving entry lost its output")
	}
	if len(got[0].AlsoFailing) != 2 {
		t.Fatalf("expected 2 names in alsoFailing, got %v", got[0].AlsoFailing)
	}
	if got[0].AlsoFailing[0] != "layout re-evaluates on swap" {
		t.Errorf("unexpected roll-call: %v", got[0].AlsoFailing)
	}
}

func TestCollapseFailuresKeepsDistinctErrorsApart(t *testing.T) {
	in := []testFail{
		{Test: "a", Output: "Error: expected 1 to be 2"},
		{Test: "b", Output: "Error: cannot read properties of undefined"},
		{Test: "c", Output: "Error: expected 1 to be 2"},
	}

	got := collapseFailures(in)

	if len(got) != 2 {
		t.Fatalf("expected 2 distinct errors, got %d: %+v", len(got), got)
	}
	if len(got[0].AlsoFailing) != 1 || got[0].AlsoFailing[0] != "c" {
		t.Errorf("the matching pair did not fold: %+v", got[0])
	}
	if len(got[1].AlsoFailing) != 0 {
		t.Errorf("the unrelated failure picked up a passenger: %+v", got[1])
	}
}

// Go stamps a per-test file:line on an otherwise identical assertion, so the
// prefix comes off before two failures are compared.
func TestCollapseFailuresIgnoresPerTestFileLine(t *testing.T) {
	in := []testFail{
		{Package: "p", Test: "A", Output: "    store_test.go:41: connect: no such host"},
		{Package: "p", Test: "B", Output: "    store_test.go:88: connect: no such host"},
	}

	if got := collapseFailures(in); len(got) != 1 {
		t.Fatalf("expected the shared cause to fold, got %d entries", len(got))
	}
}

func TestCollapseFailuresKeepsPackagesApart(t *testing.T) {
	in := []testFail{
		{Package: "one", Test: "A", Output: "Error: boom"},
		{Package: "two", Test: "B", Output: "Error: boom"},
	}

	if got := collapseFailures(in); len(got) != 2 {
		t.Fatalf("failures in different packages must not fold, got %d", len(got))
	}
}

func TestCollapseFailuresCountsPastTheCap(t *testing.T) {
	in := []testFail{{Test: "first", Output: "Error: boom"}}
	for i := 0; i < maxAlsoFailing+7; i++ {
		in = append(in, testFail{Test: "t", Output: "Error: boom"})
	}

	got := collapseFailures(in)

	if len(got) != 1 {
		t.Fatalf("expected one entry, got %d", len(got))
	}
	if len(got[0].AlsoFailing) != maxAlsoFailing {
		t.Errorf("expected the list capped at %d, got %d", maxAlsoFailing, len(got[0].AlsoFailing))
	}
	if got[0].AlsoFailingMore != 7 {
		t.Errorf("expected 7 counted past the cap, got %d", got[0].AlsoFailingMore)
	}
}

// A failure with no output cannot be compared, and folding it into an unrelated
// one would hide a real second cause.
func TestCollapseFailuresKeepsOutputlessFailures(t *testing.T) {
	in := []testFail{
		{Test: "a", Output: ""},
		{Test: "b", Output: ""},
	}

	if got := collapseFailures(in); len(got) != 2 {
		t.Fatalf("outputless failures must stay separate, got %d", len(got))
	}
}
