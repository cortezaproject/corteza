package mcp

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScopeFromRequest(t *testing.T) {
	cases := map[string]struct {
		path string
		want Scope
	}{
		"bare mcp path is unnarrowed":  {"/api/mcp", Scope{}},
		"trailing slash is unnarrowed": {"/api/mcp/", Scope{}},
		"group segment":                {"/api/mcp/configuring", Scope{Group: GroupConfiguring}},
		"usage group":                  {"/api/mcp/usage", Scope{Group: GroupUsage}},
		"group with risk":              {"/api/mcp/configuring?maxRisk=read", Scope{Group: GroupConfiguring, MaxRisk: RiskRead}},
		"risk without group":           {"/api/mcp?maxRisk=write", Scope{MaxRisk: RiskWrite}},

		// A typo must not silently narrow the surface to nothing: an unknown
		// group is ignored, so the caller sees everything and notices, rather
		// than seeing an empty tool list and concluding the server is broken.
		"unknown group is ignored": {"/api/mcp/confguring", Scope{}},
		"unknown risk is ignored":  {"/api/mcp?maxRisk=yolo", Scope{}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest("POST", c.path, nil)
			assert.Equal(t, c.want, scopeFromRequest("/api/mcp", r))
		})
	}
}

func TestScopePermits(t *testing.T) {
	configuring := []Group{GroupConfiguring}
	both := []Group{GroupConfiguring, GroupUsage}

	t.Run("zero scope permits everything", func(t *testing.T) {
		assert.True(t, Scope{}.Permits(configuring, RiskDestructive))
	})

	t.Run("group narrows", func(t *testing.T) {
		s := Scope{Group: GroupUsage}
		assert.False(t, s.Permits(configuring, RiskRead))
		assert.True(t, s.Permits(both, RiskRead))
	})

	t.Run("risk ceiling caps", func(t *testing.T) {
		s := Scope{MaxRisk: RiskWrite}
		assert.True(t, s.Permits(configuring, RiskRead))
		assert.True(t, s.Permits(configuring, RiskWrite))
		assert.False(t, s.Permits(configuring, RiskDestructive))
	})

	t.Run("read-only ceiling admits only reads", func(t *testing.T) {
		s := Scope{MaxRisk: RiskRead}
		assert.True(t, s.Permits(configuring, RiskRead))
		assert.False(t, s.Permits(configuring, RiskWrite))
		assert.False(t, s.Permits(configuring, RiskDestructive))
	})

	// An untagged tool has no risk, and AtOrBelow treats an unknown level as
	// exceeding every ceiling — so a tool that skipped WithRisk cannot slip
	// through a capped session.
	t.Run("untagged risk never passes a ceiling", func(t *testing.T) {
		assert.False(t, Scope{MaxRisk: RiskDestructive}.Permits(configuring, ""))
		assert.True(t, Scope{}.Permits(configuring, ""))
	})
}

func TestRiskAtOrBelow(t *testing.T) {
	assert.True(t, RiskRead.AtOrBelow(RiskDestructive))
	assert.True(t, RiskWrite.AtOrBelow(RiskWrite))
	assert.False(t, RiskDestructive.AtOrBelow(RiskWrite))
	assert.False(t, Risk("nonsense").AtOrBelow(RiskDestructive))
	assert.False(t, RiskRead.AtOrBelow(Risk("nonsense")))
}
