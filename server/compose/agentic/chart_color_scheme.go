package agentic

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Valid colour scheme names, generated from the webapp's own scheme tables
// (lib/js/src/shared/types/chart/colorschemes/{brewer,office,tableau}.ts).
// lib/js keeps this honest: chart-color-schemes.test.ts fails if the two drift.
//
//go:embed chart_color_schemes.json
var chartColorSchemesJSON []byte

var chartColorSchemes = func() []string {
	var names []string
	if err := json.Unmarshal(chartColorSchemesJSON, &names); err != nil {
		panic("compose/agentic: cannot decode chart_color_schemes.json: " + err.Error())
	}
	return names
}()

// A scheme the webapp cannot resolve is not an error anywhere in the stack: it
// reaches echarts as an undefined palette, so the chart draws its legend and
// none of its series — a blank panel that reads as "no data". Nothing in a
// render check catches it either, because the block does render. Refusing the
// write is the only place the mistake is still cheap.
func validateChartColorScheme(name string) error {
	if name == "" {
		return nil
	}

	// Custom schemes resolve against the namespace's own customColorSchemes
	// rather than these tables; the webapp selects them by the same substring.
	if strings.Contains(strings.ToLower(name), "custom") {
		return nil
	}

	for _, known := range chartColorSchemes {
		if known == name {
			return nil
		}
	}

	msg := fmt.Sprintf("config: unknown colorScheme %q — the webapp resolves it to no palette, so the chart renders its legend with no visible series", name)
	if near := nearestColorSchemes(name); len(near) > 0 {
		return fmt.Errorf("%s. Did you mean %s?", msg, strings.Join(quoteAll(near), " or "))
	}
	return fmt.Errorf("%s. Valid names are %q-prefixed keys from lib/js colorschemes, e.g. \"tableau.Tableau10\"", msg, familiesOf())
}

// The suffix is the usual casualty — a scheme is remembered by name and given
// the wrong swatch count ("tableau.ClassicOrangeBlue7" for ...13). Matching on
// the name with its digits removed turns that into a one-line correction.
func nearestColorSchemes(name string) []string {
	var (
		want = foldScheme(name)
		bare = want
		out  []string
	)

	// An unqualified name ("Tableau10") should still find its family.
	if i := strings.IndexByte(bare, '.'); i >= 0 {
		bare = bare[i+1:]
	}

	for _, known := range chartColorSchemes {
		k := foldScheme(known)
		if k == want || k[strings.IndexByte(k, '.')+1:] == bare {
			out = append(out, known)
		}
	}

	if len(out) == 0 {
		return nil
	}
	sort.Strings(out)
	if len(out) > 3 {
		out = out[:3]
	}
	return out
}

// Lower-cased and stripped of digits, so a name differing only in casing or in
// its swatch count folds onto the real one.
func foldScheme(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if r < '0' || r > '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func familiesOf() string {
	seen := map[string]bool{}
	var fams []string
	for _, n := range chartColorSchemes {
		if i := strings.IndexByte(n, '.'); i > 0 && !seen[n[:i]] {
			seen[n[:i]] = true
			fams = append(fams, n[:i])
		}
	}
	sort.Strings(fams)
	return strings.Join(fams, "/")
}

func quoteAll(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = fmt.Sprintf("%q", s)
	}
	return out
}
