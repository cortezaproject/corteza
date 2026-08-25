package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func refMod() *cmpTypes.Module {
	return &cmpTypes.Module{
		NamespaceID: 1,
		Fields: cmpTypes.ModuleFieldSet{
			{Name: "card", Kind: "Record", Options: cmpTypes.ModuleFieldOptions{"moduleID": "22"}},
			{Name: "storage", Kind: "Select"},
			{Name: "quantity", Kind: "Number"},
		},
	}
}

// The service is unreachable from a unit test, so what is pinned here is
// everything decided before it: which terms are recognised, which are left
// alone, and which are refused outright.
func TestResolveRefPathsLeavesUnrelatedExpressions(t *testing.T) {
	for _, expr := range []string{
		"",
		"storage = 'Binder'",
		"quantity > 2 AND storage = 'Binder'",
		// Not a Record field, so not this rewrite's business.
		"storage.name = 'Binder'",
		// A dot inside a literal must not be read as a path.
		"storage = 'Deck.Box'",
	} {
		got, err := resolveRefPaths(nil, refMod(), expr)
		require.NoError(t, err, expr)
		assert.Equal(t, expr, got, "should leave %q alone", expr)
	}
}

// An operator the rewrite cannot express must be refused, not silently passed
// through to fail as "unknown attribute".
func TestResolveRefPathsRefusesUnsupportedOperators(t *testing.T) {
	for _, expr := range []string{
		"card.name != 'Bolt'",
		"card.mana_value > 3",
		"card.name NOT LIKE 'B%'",
	} {
		_, err := resolveRefPaths(nil, refMod(), expr)
		require.Error(t, err, expr)
		assert.Contains(t, err.Error(), "card")
	}
}

// With no matches the term must still select nothing. An empty replacement
// would leave a dangling AND, or worse, widen the filter to every record.
func TestIDDisjunction(t *testing.T) {
	assert.Equal(t, "card = '0'", idDisjunction("card", nil))
	assert.Equal(t, "(card = '7')", idDisjunction("card", []uint64{7}))
	assert.Equal(t, "(card = '7' OR card = '8')", idDisjunction("card", []uint64{7, 8}))
}

func TestRefPathTermMatches(t *testing.T) {
	cases := map[string]bool{
		"card.name = 'Bolt'":     true,
		"card.name='Bolt'":       true,
		"card.name LIKE 'Bo%'":   true,
		"card.name like 'Bo%'":   true,
		"card.name = 'O''Brien'": true,
		"card = '123'":           false,
		"name = 'Bolt'":          false,
		"card.name = 3":          false,
	}
	for expr, want := range cases {
		assert.Equal(t, want, refPathTerm.MatchString(expr), "for %q", expr)
	}
}

func TestLiteralRanges(t *testing.T) {
	assert.Nil(t, literalRanges("a = b"))
	assert.Equal(t, [][2]int{{4, 9}}, literalRanges("a = 'bolt'"))

	// A backslash escapes the quote; the literal runs to the closing one.
	assert.Equal(t, [][2]int{{4, 11}}, literalRanges(`a = 'O\'Bri'`))

	// Doubling does NOT escape here, whatever SQL habit suggests: it closes one
	// literal and opens the next, which is why `'Urza''s Saga'` matches nothing
	// instead of erroring.
	assert.Len(t, literalRanges("a = 'O''Bri'"), 2)

	// An unterminated literal runs to the end rather than being ignored.
	assert.Equal(t, [][2]int{{4, 7}}, literalRanges("a = 'bo"))
}

func TestInRanges(t *testing.T) {
	rr := [][2]int{{4, 9}}
	assert.False(t, inRanges(rr, 0))
	assert.False(t, inRanges(rr, 4), "the opening quote itself is not inside")
	assert.True(t, inRanges(rr, 5))
	assert.False(t, inRanges(rr, 9))
	assert.False(t, inRanges(rr, 20))
}

// A path written inside a string literal is data being compared against, not an
// expression; rewriting it would corrupt the value.
func TestResolveRefPathsIgnoresPathsInsideLiterals(t *testing.T) {
	expr := "notes = 'see card.name = ' AND storage = 'Binder'"
	got, err := resolveRefPaths(nil, refMod(), expr)
	require.NoError(t, err)
	assert.Equal(t, expr, got)
}
