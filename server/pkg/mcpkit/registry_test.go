package mcpkit

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A caller choosing tools to GRANT needs to know what each one is and what it
// costs; a caller choosing tools to RUN does not. Select serves both, so the
// classification has to survive the projection.
func TestSelectCarriesGroupsAndRisk(t *testing.T) {
	r := NewRegistry()
	r.RegisterTool(
		mcp.NewTool("thing_delete",
			mcp.WithDescription("deletes a thing"),
			InGroup(GroupConfiguring),
			WithRisk(RiskDestructive),
		),
		"Delete thing",
		func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) { return nil, nil },
	)

	defs, err := r.Select(nil)
	require.NoError(t, err)
	require.Len(t, defs, 1)
	assert.Equal(t, []Group{GroupConfiguring}, defs[0].Groups)
	assert.Equal(t, RiskDestructive, defs[0].Risk)
}

// What a registrant adds reaches the initialize text after the server's own
// paragraph, in order; blank additions leave no empty paragraph behind.
func TestInstructionsCarryRegistrantParagraphs(t *testing.T) {
	reg := NewRegistry()
	reg.AddInstructions("  first  ")
	reg.AddInstructions("")
	reg.AddInstructions("second")

	assert.Equal(t, []string{"first", "second"}, reg.Instructions())
	assert.Equal(t, serverInstructions+"\n\nfirst\n\nsecond", instructions(reg))
	assert.Equal(t, serverInstructions, instructions(NewRegistry()))
}
