package agentic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateMetrics(t *testing.T) {
	ok := []string{
		"",
		"SUM(line_value)",
		"SUM(line_value) AS total",
		"sum(a) as x, AVG(b) AS y, MIN(c), MAX(d), COUNT(e)",
		"SUM(quantity * unit_value) AS total",
	}
	for _, m := range ok {
		assert.NoError(t, validateMetrics(m), "should accept %q", m)
	}

	// The case this exists for: a bare field name reaches the database
	// ungrouped and comes back as a driver complaint about a column the caller
	// never named.
	t.Run("a bare field is refused by name", func(t *testing.T) {
		err := validateMetrics("missing_count, SUM(deck_value) AS v")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing_count")
		assert.Contains(t, err.Error(), "SUM(missing_count)")
	})

	t.Run("a comma inside a call is not a separator", func(t *testing.T) {
		assert.NoError(t, validateMetrics("SUM(COALESCE(a, 0)) AS total"))
	})
}

func TestSplitMetrics(t *testing.T) {
	assert.Equal(t, []string{"SUM(a)", "AVG(b)"}, splitMetrics("SUM(a), AVG(b)"))
	assert.Equal(t, []string{"SUM(f(a, b))"}, splitMetrics("SUM(f(a, b))"))
	assert.Nil(t, splitMetrics(""))
}
