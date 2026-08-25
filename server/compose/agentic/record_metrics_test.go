package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"

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
		assert.NoError(t, validateMetrics(m, nil), "should accept %q", m)
	}

	// The case this exists for: a bare field name reaches the database
	// ungrouped and comes back as a driver complaint about a column the caller
	// never named.
	t.Run("a bare field is refused by name", func(t *testing.T) {
		err := validateMetrics("missing_count, SUM(deck_value) AS v", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing_count")
		assert.Contains(t, err.Error(), "SUM(missing_count)")
	})

	t.Run("a comma inside a call is not a separator", func(t *testing.T) {
		assert.NoError(t, validateMetrics("SUM(COALESCE(a, 0)) AS total", nil))
	})
}

func TestSplitMetrics(t *testing.T) {
	assert.Equal(t, []string{"SUM(a)", "AVG(b)"}, splitMetrics("SUM(a), AVG(b)"))
	assert.Equal(t, []string{"SUM(f(a, b))"}, splitMetrics("SUM(f(a, b))"))
	assert.Nil(t, splitMetrics(""))
}

// The model reaches for a batch by ID in whatever shape it guessed, and the
// query language accepts none of them. Every shape it tried has to land here.
func TestParseIDList(t *testing.T) {
	want := map[uint64]struct{}{101: {}, 102: {}, 103: {}}

	for _, raw := range []string{
		"101,102,103",
		" 101, 102 ,103 ",
		"[101, 102, 103]",
		"(101,102,103)",
		`['101', '102', '103']`,
		`"101","102","103"`,
	} {
		assert.Equal(t, want, parseIDList(raw), "should parse %q", raw)
	}

	for _, raw := range []string{"", "   ", "[]", "abc", "0"} {
		assert.Nil(t, parseIDList(raw), "should reject %q", raw)
	}
}

// A bare field listed as a metric is nearly always one meant to label the
// groups. When the module says the field cannot be aggregated at all, say
// which parameter it belongs in rather than offering to sum a card reference.
func TestValidateMetricsPointsAtDimension(t *testing.T) {
	mod := &cmpTypes.Module{Fields: cmpTypes.ModuleFieldSet{
		{Name: "card", Kind: "Record"},
		{Name: "price", Kind: "Number"},
	}}

	err := validateMetrics("MAX(change_abs) AS max_change, card", mod)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dimension")
	assert.Contains(t, err.Error(), "Record")
	assert.NotContains(t, err.Error(), "SUM(card)")

	// A numeric field still gets the wrap-it suggestion — that one is right.
	err = validateMetrics("price", mod)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SUM(price)")
}
