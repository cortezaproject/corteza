package agentic

import (
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionIsUnset(t *testing.T) {
	for _, v := range []any{nil, "", "   ", "0", float64(0), []any{}, map[string]any{}} {
		assert.True(t, optionIsUnset(v), "%#v should read as unset", v)
	}
	for _, v := range []any{"123", float64(123), []any{"a"}, map[string]any{"k": 1}, true} {
		assert.False(t, optionIsUnset(v), "%#v should read as set", v)
	}
}

// The page builder refuses to save exactly this page; the API accepted it in
// silence, so a RecordList naming no module looked like one that worked.
func TestBlockReadinessNoteNamesEveryUnfinishedBlock(t *testing.T) {
	note := blockReadinessNote([]cmpTypes.PageBlock{
		{BlockID: 1, Kind: "AgentChat", Options: map[string]any{}},
		{BlockID: 2, Kind: "RecordList", Options: map[string]any{"moduleID": "0"}},
		{BlockID: 3, Kind: "Chart", Options: map[string]any{"chartID": "999"}},
	})

	require.NotEmpty(t, note)
	assert.Contains(t, note, "block 1 (AgentChat) needs allowedAgentIDs")
	assert.Contains(t, note, "block 2 (RecordList) needs moduleID")
	assert.NotContains(t, note, "Chart", "a configured block is not worth mentioning")
	assert.Contains(t, note, "2 of 3 blocks")
}

func TestBlockReadinessNoteIsSilentWhenThereIsNothingToSay(t *testing.T) {
	assert.Empty(t, blockReadinessNote(nil))
	assert.Empty(t, blockReadinessNote([]cmpTypes.PageBlock{
		{Kind: "RecordList", Options: map[string]any{"moduleID": "11"}},
		{Kind: "Record", Options: map[string]any{}},
	}), "a kind with no requirement is always ready")
}

// Before the server assigns blockIDs there is nothing to name a block by but
// its position, and "block 0" is not a thing the caller sent.
func TestBlockReadinessNoteRefersToPositionBeforeIDs(t *testing.T) {
	note := blockReadinessNote([]cmpTypes.PageBlock{
		{Kind: "Chart", Options: map[string]any{}},
	})
	assert.Contains(t, note, "block #1 (Chart)")
}

// A Metric with tiles passes the top-level check and still renders nothing
// when a tile names no module — the mistake people actually make.
func TestMissingBlockOptionsChecksMetricTiles(t *testing.T) {
	ok := missingBlockOptions(cmpTypes.PageBlock{Kind: "Metric", Options: map[string]any{
		"metrics": []any{map[string]any{"moduleID": "11"}},
	}})
	assert.Empty(t, ok)

	bad := missingBlockOptions(cmpTypes.PageBlock{Kind: "Metric", Options: map[string]any{
		"metrics": []any{
			map[string]any{"moduleID": "11"},
			map[string]any{"metricField": "count"},
		},
	}})
	assert.Equal(t, []string{"metrics.1.moduleID"}, bad)
}

// The table mirrors lib/js PageBlock.validate(); a kind added there and missed
// here is a block that goes back to failing silently.
func TestRequiredBlockOptionsCoversTheKindsThatHaveOne(t *testing.T) {
	for _, kind := range []string{
		"AgentChat", "Automation", "Calendar", "Chart", "ChatbotInbox", "Comment",
		"Content", "Geometry", "IFrame", "Metric", "Navigation", "Progress",
		"RecordList", "RecordOrganizer", "Tabs",
	} {
		assert.NotEmpty(t, requiredBlockOptions[kind], "%s must declare what it needs", kind)
	}

	for kind := range requiredBlockOptions {
		_, known := cmpTypes.PageBlockOptionSchemas[kind]
		assert.True(t, known, "%s is not a block kind this server has", kind)
	}
}
