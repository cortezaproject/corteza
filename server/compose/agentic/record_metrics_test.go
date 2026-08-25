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

func TestParseSortSpec(t *testing.T) {
	assert.Nil(t, parseSortSpec(""))
	assert.Nil(t, parseSortSpec("   "))

	s := parseSortSpec("total DESC")
	require.NotNil(t, s)
	assert.Equal(t, "total", s.key)
	assert.True(t, s.desc)

	s = parseSortSpec("count")
	require.NotNil(t, s)
	assert.Equal(t, "count", s.key)
	assert.False(t, s.desc)

	assert.True(t, parseSortSpec("total desc").desc)
	assert.False(t, parseSortSpec("total ASC").desc)
}

func TestSortAndLimitRows(t *testing.T) {
	rows := []map[string]any{
		{"dimension_0": "a", "total": 1.0},
		{"dimension_0": "b", "total": 3.0},
		{"dimension_0": "c", "total": 2.0},
	}

	t.Run("untouched when nothing is asked", func(t *testing.T) {
		assert.Equal(t, any(rows), sortAndLimitRows(rows, nil, 0))
	})

	t.Run("top N descending", func(t *testing.T) {
		got := sortAndLimitRows(rows, parseSortSpec("total DESC"), 2).([]map[string]any)
		require.Len(t, got, 2)
		assert.Equal(t, "b", got[0]["dimension_0"])
		assert.Equal(t, "c", got[1]["dimension_0"])
	})

	t.Run("ascending", func(t *testing.T) {
		got := sortAndLimitRows(rows, parseSortSpec("total"), 0).([]map[string]any)
		assert.Equal(t, "a", got[0]["dimension_0"])
	})

	t.Run("a limit alone keeps the aggregation's own order", func(t *testing.T) {
		got := sortAndLimitRows(rows, nil, 1).([]map[string]any)
		require.Len(t, got, 1)
		assert.Equal(t, "a", got[0]["dimension_0"])
	})

	// A group with no value must not outrank one that has a real figure.
	t.Run("a missing value never tops a ranking", func(t *testing.T) {
		with := []map[string]any{
			{"dimension_0": "a", "total": nil},
			{"dimension_0": "b", "total": 3.0},
		}
		got := sortAndLimitRows(with, parseSortSpec("total DESC"), 1).([]map[string]any)
		assert.Equal(t, "b", got[0]["dimension_0"])
	})
}

// A report with no limit returns every group; toolkit.Page's default page size
// would have silently kept the first fifty.
func TestReportLimit(t *testing.T) {
	assert.Equal(t, 0, reportLimit(map[string]any{}))
	assert.Equal(t, 0, reportLimit(map[string]any{"limit": ""}))
	assert.Equal(t, 0, reportLimit(map[string]any{"limit": "0"}))
	assert.Equal(t, 5, reportLimit(map[string]any{"limit": "5"}))
	assert.Equal(t, 5, reportLimit(map[string]any{"limit": float64(5)}))
}

// A field sent as an empty list contributes no values and would otherwise look
// like a field that was never named — which on a patching update means "leave
// it alone" rather than "clear it".
func TestParseValuesReportsNamedFields(t *testing.T) {
	_, named, err := parseValues(`{"notes":"hi","colors":[],"quantity":2}`)
	require.NoError(t, err)
	assert.Equal(t, []string{"colors", "notes", "quantity"}, named)
}
