package ql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// A filter the parser cannot read says what to write instead.
// Example: an agent wrote status IN ("open","waiting") and was told only
// found an illegal token ({code:0 literal:" line:0 char:15}).
func TestParserIllegalTokenErrors(t *testing.T) {
	for q, want := range map[string]string{
		`status IN ("open", "waiting")`: "text is quoted with single quotes",
		`status = "open"`:               "text is quoted with single quotes",
		`status = 'open`:                "quoted text that is not closed",
		`status = #`:                    `found an illegal token "#"`,
	} {
		_, err := NewParser().Parse(q)
		require.ErrorContains(t, err, want, q)
	}
}
