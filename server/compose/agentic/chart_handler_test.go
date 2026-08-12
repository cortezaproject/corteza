package agentic

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// TestChartConfigCasesContract is one half of a cross-language contract; the
// other is lib/js/src/compose/types/chart/config-contract.test.ts, which asserts
// the webapp's chart classes reach the same verdicts on the same fixture. The
// renderer validates a chart only when it draws it, so anything these two
// disagree about is stored happily here and fails in front of a person there.
func TestChartConfigCasesContract(t *testing.T) {
	var cases struct {
		Rejected []struct {
			Name   string          `json:"name"`
			Error  string          `json:"error"`
			Config json.RawMessage `json:"config"`
		} `json:"rejected"`
		Accepted []struct {
			Name   string          `json:"name"`
			Config json.RawMessage `json:"config"`
		} `json:"accepted"`
	}

	raw, err := os.ReadFile("testdata/chart_config_cases.json")
	if err != nil {
		t.Fatalf("cannot read the shared fixture: %v", err)
	}
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("cannot decode the shared fixture: %v", err)
	}
	if len(cases.Rejected) == 0 || len(cases.Accepted) == 0 {
		t.Fatal("the shared fixture must carry both rejected and accepted cases")
	}

	for _, c := range cases.Rejected {
		t.Run("rejected/"+c.Name, func(t *testing.T) {
			if _, err := parseChartConfig(string(c.Config)); err == nil {
				t.Errorf("accepted a config the webapp rejects with %s", c.Error)
			}
		})
	}

	for _, c := range cases.Accepted {
		t.Run("accepted/"+c.Name, func(t *testing.T) {
			if _, err := parseChartConfig(string(c.Config)); err != nil {
				t.Errorf("rejected a config the webapp renders: %v", err)
			}
		})
	}
}

func TestParseChartConfigRejectsUnrenderableReports(t *testing.T) {
	tests := []struct {
		name   string
		config string
		errHas string
	}{
		{
			// What agents were writing, and what renders as "Metrics aggregate
			// not defined" in front of a person.
			name:   "metric on a module field without an aggregate",
			config: `{"reports":[{"moduleID":"1","dimensions":[{"field":"status","modifier":"(no grouping / buckets)"}],"metrics":[{"field":"amount","type":"bar"}]}]}`,
			errHas: `"aggregate" is required`,
		},
		{
			name:   "unknown aggregate",
			config: `{"reports":[{"moduleID":"1","metrics":[{"field":"amount","aggregate":"MEDIAN","type":"bar"}]}]}`,
			errHas: `aggregate "MEDIAN" is not one of`,
		},
		{
			name:   "metric without a field",
			config: `{"reports":[{"moduleID":"1","metrics":[{"type":"bar"}]}]}`,
			errHas: `"field" is required`,
		},
		{
			name:   "metric without a type",
			config: `{"reports":[{"moduleID":"1","metrics":[{"field":"count"}]}]}`,
			errHas: `"type" is required`,
		},
		{
			name:   "unknown chart type",
			config: `{"reports":[{"moduleID":"1","metrics":[{"field":"count","type":"sunburst"}]}]}`,
			errHas: `type "sunburst" is not one of`,
		},
		{
			// A handle cannot even decode into the uint64 moduleID, so this is
			// caught earlier than the rest — the message still has to name the fix.
			name:   "module handle where the ID belongs",
			config: `{"reports":[{"moduleID":"task","metrics":[{"field":"count","type":"pie"}]}]}`,
			errHas: "not a module handle",
		},
		{
			name:   "report without a module",
			config: `{"reports":[{"metrics":[{"field":"count","type":"pie"}]}]}`,
			errHas: "moduleID is required",
		},
		{
			name:   "report with no metrics",
			config: `{"reports":[{"moduleID":"1","dimensions":[{"field":"status"}]}]}`,
			errHas: "at least one metric",
		},
		{
			name:   "dimension without a field",
			config: `{"reports":[{"moduleID":"1","dimensions":[{"modifier":"MONTH"}],"metrics":[{"field":"count","type":"line"}]}]}`,
			errHas: `"field" is required`,
		},
		{
			name:   "unknown dimension modifier",
			config: `{"reports":[{"moduleID":"1","dimensions":[{"field":"createdAt","modifier":"DECADE"}],"metrics":[{"field":"count","type":"line"}]}]}`,
			errHas: `modifier "DECADE" is not one of`,
		},
		{
			name:   "gauge without steps",
			config: `{"reports":[{"moduleID":"1","dimensions":[{"field":"status"}],"metrics":[{"field":"count","type":"gauge"}]}]}`,
			errHas: "meta.steps",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseChartConfig(tc.config)
			if err == nil {
				t.Fatalf("expected an error mentioning %q, got none", tc.errHas)
			}
			if !strings.Contains(err.Error(), tc.errHas) {
				t.Errorf("error %q does not mention %q", err, tc.errHas)
			}
		})
	}
}

func TestParseChartConfigAcceptsRenderableReports(t *testing.T) {
	for _, cfg := range []string{
		// Counting records, grouped by a Select field.
		`{"reports":[{"moduleID":"1","dimensions":[{"field":"status","modifier":"(no grouping / buckets)"}],"metrics":[{"field":"count","type":"doughnut"}]}]}`,
		// Summing a numeric field over months.
		`{"reports":[{"moduleID":"1","dimensions":[{"field":"createdAt","modifier":"MONTH"}],"metrics":[{"field":"amount","aggregate":"SUM","type":"bar"}]}]}`,
		// A gauge with its bands.
		`{"reports":[{"moduleID":"1","dimensions":[{"field":"status","meta":{"steps":[{"value":0},{"value":100}]}}],"metrics":[{"field":"count","type":"gauge"}]}]}`,
	} {
		if _, err := parseChartConfig(cfg); err != nil {
			t.Errorf("config %s: unexpected error %v", cfg, err)
		}
	}
}

func TestParseChartConfigNormalizesWhatItCan(t *testing.T) {
	cfg, err := parseChartConfig(
		`{"reports":[{"moduleID":"1","dimensions":[{"field":"createdAt","modifier":"month"},{"field":"status"}],"metrics":[{"field":"amount","aggregate":"sum","type":"Bar"}]}]}`,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := cfg.Reports[0]
	if got := r.Metrics[0]["aggregate"]; got != "SUM" {
		t.Errorf("aggregate: got %v, want SUM", got)
	}
	if got := r.Metrics[0]["type"]; got != "bar" {
		t.Errorf("type: got %v, want bar", got)
	}
	if got := r.Dimensions[0]["modifier"]; got != "MONTH" {
		t.Errorf("modifier: got %v, want MONTH", got)
	}
	// An absent modifier is the plainest grouping, not an error.
	if got := r.Dimensions[1]["modifier"]; got != chartNoGrouping {
		t.Errorf("defaulted modifier: got %v, want %q", got, chartNoGrouping)
	}
}
