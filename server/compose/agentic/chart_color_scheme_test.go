package agentic

import (
	"strings"
	"testing"
)

func TestValidateChartColorScheme(t *testing.T) {
	cases := []struct {
		name     string
		scheme   string
		ok       bool
		mentions []string
	}{
		{name: "empty is the default palette", scheme: "", ok: true},
		{name: "known tableau key", scheme: "tableau.Tableau10", ok: true},
		{name: "known brewer key", scheme: "brewer.Blues3", ok: true},
		{name: "known office key", scheme: "office.Apex6", ok: true},

		// The scheme that produced a legend with no arcs on a real dashboard.
		// Its swatch count is part of the name, and 13 is the one that exists.
		{
			name:   "wrong swatch count is refused and corrected",
			scheme: "tableau.ClassicOrangeBlue7", ok: false,
			mentions: []string{"tableau.ClassicOrangeBlue13", "no visible series"},
		},
		{
			name:   "unqualified name is refused and qualified",
			scheme: "Tableau10", ok: false,
			mentions: []string{"tableau.Tableau10"},
		},
		{
			name:   "wrong case is refused and corrected",
			scheme: "tableau.tableau10", ok: false,
			mentions: []string{"tableau.Tableau10"},
		},
		{
			name:   "unrecognisable name still names the families",
			scheme: "tableau.NotAScheme", ok: false,
			mentions: []string{"brewer/office/tableau"},
		},

		// Custom schemes live on the namespace, not in these tables; the webapp
		// selects them by the same substring, so they must pass untouched.
		{name: "custom scheme", scheme: "custom.myPalette", ok: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateChartColorScheme(c.scheme)

			if c.ok {
				if err != nil {
					t.Fatalf("expected %q to be accepted, got: %v", c.scheme, err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected %q to be refused, it was accepted", c.scheme)
			}
			for _, want := range c.mentions {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error should mention %q, got: %v", want, err)
				}
			}
		})
	}
}

// The check has to sit on the path compose_chart_create actually takes, not
// only in its own function: an otherwise valid config is what carries a bad
// scheme in practice, and that is the config that used to be stored.
func TestParseChartConfigChecksColorScheme(t *testing.T) {
	cfg := func(scheme string) string {
		return `{"colorScheme":"` + scheme + `","reports":[{"moduleID":"1","filter":"",` +
			`"dimensions":[{"field":"stage","modifier":"(no grouping / buckets)","conditions":{}}],` +
			`"metrics":[{"field":"count","type":"doughnut"}]}]}`
	}

	if _, err := parseChartConfig(cfg("tableau.Tableau10")); err != nil {
		t.Fatalf("a valid scheme must still pass: %v", err)
	}

	_, err := parseChartConfig(cfg("tableau.ClassicOrangeBlue7"))
	if err == nil {
		t.Fatal("a config with an unresolvable colorScheme was accepted")
	}
	if !strings.Contains(err.Error(), "tableau.ClassicOrangeBlue13") {
		t.Errorf("refusal should point at the real scheme, got: %v", err)
	}

	// No scheme at all is the default palette and must stay valid.
	if _, err := parseChartConfig(`{"reports":[{"moduleID":"1","filter":"","dimensions":[{"field":"stage","modifier":"(no grouping / buckets)","conditions":{}}],"metrics":[{"field":"count","type":"doughnut"}]}]}`); err != nil {
		t.Fatalf("a config without a colorScheme must pass: %v", err)
	}
}

// The embedded list is generated from the webapp's tables; if it were empty or
// truncated every scheme would be refused, which is worse than the bug it fixes.
func TestChartColorSchemesEmbedded(t *testing.T) {
	if len(chartColorSchemes) < 400 {
		t.Fatalf("embedded scheme list looks truncated: %d names", len(chartColorSchemes))
	}

	fams := map[string]int{}
	for _, n := range chartColorSchemes {
		i := strings.IndexByte(n, '.')
		if i <= 0 {
			t.Fatalf("scheme %q is not family-qualified", n)
		}
		fams[n[:i]]++
	}

	for _, want := range []string{"brewer", "office", "tableau"} {
		if fams[want] == 0 {
			t.Errorf("no schemes from the %q family", want)
		}
	}
}
