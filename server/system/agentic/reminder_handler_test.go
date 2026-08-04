package agentic

import (
	"testing"
	"time"

	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRegistrar struct {
	tools []mcp.Tool
}

func (s *stubRegistrar) RegisterTool(tool mcp.Tool, _ string, _ server.ToolHandlerFunc, _ ...hmcp.RegisterOption) {
	s.tools = append(s.tools, tool)
}

func (s *stubRegistrar) byName(name string) (mcp.Tool, bool) {
	for _, t := range s.tools {
		if t.Name == name {
			return t, true
		}
	}
	return mcp.Tool{}, false
}

func TestReminderHandlerRegistersTools(t *testing.T) {
	reg := &stubRegistrar{}
	ReminderHandler(reg)

	want := []string{
		"system_reminder_lookup",
		"system_reminder_create",
		"system_reminder_update",
		"system_reminder_delete",
		"system_reminder_dismiss",
		"system_reminder_undismiss",
		"system_reminder_snooze",
	}
	require.Len(t, reg.tools, len(want))
	for _, name := range want {
		_, ok := reg.byName(name)
		assert.True(t, ok, "expected %q to be registered", name)
	}

	// Reminder has no UndeleteByID on its service, so there is deliberately no
	// undelete tool. Asserted so nobody adds one by pattern-matching the
	// conventions without checking the service.
	_, ok := reg.byName("system_reminder_undelete")
	assert.False(t, ok)
}

func TestReminderToolsAreTagged(t *testing.T) {
	reg := &stubRegistrar{}
	ReminderHandler(reg)

	wantRisk := map[string]hmcp.Risk{
		"system_reminder_lookup":    hmcp.RiskRead,
		"system_reminder_create":    hmcp.RiskWrite,
		"system_reminder_update":    hmcp.RiskWrite,
		"system_reminder_delete":    hmcp.RiskDestructive,
		"system_reminder_dismiss":   hmcp.RiskWrite,
		"system_reminder_undismiss": hmcp.RiskWrite,
		"system_reminder_snooze":    hmcp.RiskWrite,
	}

	for _, tool := range reg.tools {
		assert.Equal(t, []hmcp.Group{hmcp.GroupUsage}, hmcp.GroupsOf(tool), tool.Name)
		assert.Equal(t, wantRisk[tool.Name], hmcp.RiskOf(tool), tool.Name)
	}
}

// TestReminderRiskDrivesAnnotations pins the behaviour that motivated making
// WithRisk the sole writer of the hints: mcp.NewTool defaults DestructiveHint
// and OpenWorldHint to true, so an untagged read tool would advertise itself as
// destructive.
func TestReminderRiskDrivesAnnotations(t *testing.T) {
	reg := &stubRegistrar{}
	ReminderHandler(reg)

	lookup, ok := reg.byName("system_reminder_lookup")
	require.True(t, ok)
	assert.True(t, *lookup.Annotations.ReadOnlyHint)
	assert.False(t, *lookup.Annotations.DestructiveHint)
	assert.False(t, *lookup.Annotations.OpenWorldHint)

	del, ok := reg.byName("system_reminder_delete")
	require.True(t, ok)
	assert.False(t, *del.Annotations.ReadOnlyHint)
	assert.True(t, *del.Annotations.DestructiveHint)
	assert.False(t, *del.Annotations.OpenWorldHint)
}

// TestReminderLookupContract covers the three lookup rules from CONVENTIONS.md
// §8.1, which every lookup tool in the fan-out must satisfy.
func TestReminderLookupContract(t *testing.T) {
	reg := &stubRegistrar{}
	ReminderHandler(reg)

	lookup, ok := reg.byName("system_reminder_lookup")
	require.True(t, ok)

	assert.NotContains(t, lookup.InputSchema.Required, "reminderID",
		"the reference param must not be required, or list mode is unreachable")

	for _, p := range []string{"limit", "pageCursor"} {
		assert.Contains(t, lookup.InputSchema.Properties, p)
	}

	// Every ID-shaped param crosses the boundary as a string.
	for _, p := range []string{"reminderID", "assignedTo"} {
		prop, ok := lookup.InputSchema.Properties[p].(map[string]any)
		require.True(t, ok, p)
		assert.Equal(t, "string", prop["type"], p)
	}
}

func TestOptTime(t *testing.T) {
	t.Run("absent and empty yield nil", func(t *testing.T) {
		v, err := optTime(map[string]any{}, "remindAt")
		require.NoError(t, err)
		assert.Nil(t, v)

		v, err = optTime(map[string]any{"remindAt": ""}, "remindAt")
		require.NoError(t, err)
		assert.Nil(t, v)
	})

	t.Run("parses RFC3339", func(t *testing.T) {
		v, err := optTime(map[string]any{"remindAt": "2026-08-03T09:00:00Z"}, "remindAt")
		require.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, 2026, v.Year())
		assert.Equal(t, time.August, v.Month())
	})

	t.Run("rejects other formats with a usable message", func(t *testing.T) {
		_, err := optTime(map[string]any{"remindAt": "tomorrow"}, "remindAt")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RFC3339")
	})
}
