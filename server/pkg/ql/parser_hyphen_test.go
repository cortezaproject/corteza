package ql

import (
	"testing"
)

// Reducing an expression used to index with bestOpIx still at -1 and panic with
// "index out of range [-1]", which surfaced as an HTTP 500. Every one of these
// reaches the parser straight from a user-supplied query string.
func TestParseDoesNotPanic(t *testing.T) {
	tests := []string{
		// A module field literally named "bad-field" lexes as two idents around
		// '-'. This parses as arithmetic — eq(sub(bad, field), 'beta') — which
		// is the parser's job; rejecting the unknown symbol is the resolver's.
		"bad-field = 'beta'",
		"a-b-c = 1",
		"foo bar",
		"a = 1 b = 2",
		"1 2",
		"name ASC",
	}

	// These reduce as far as they can and come back as a "group" node rather
	// than a parse error; rejecting them outright needs a comma-aware adjacency
	// rule, since a list like "(1, 2)" reaches parserNodes with its commas
	// already dropped and would be falsely flagged.

	for _, q := range tests {
		t.Run(q, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Parse(%q) panicked: %v", q, r)
				}
			}()

			_, _ = NewParser().Parse(q)
		})
	}
}

// Guard against the fix over-reaching: these must all keep parsing.
func TestParseAcceptsValidExpressions(t *testing.T) {
	tests := []string{
		"a-b",
		"foo = 'bar'",
		"a = 1 AND b = 2",
		"foo IS NULL",
		"foo IS NOT NULL",
		"count(id) > 2",
		"a BETWEEN 1 AND 2",
		"foo LIKE 'bar%'",
		"foo IN (1, 2)",
		"(a = 1 OR b = 2) AND c = 3",
	}

	for _, q := range tests {
		t.Run(q, func(t *testing.T) {
			if _, err := NewParser().Parse(q); err != nil {
				t.Errorf("Parse(%q) errored: %v", q, err)
			}
		})
	}
}
