package runtime

import (
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/assert"
)

// A grant may pin arguments so the agent can run an automation only one way.
// Filling in rather than overwriting would leave the pin something the model
// can talk its way past, which is no pin at all.
func TestPinTAQParams(t *testing.T) {
	agent := &types.Agent{
		Access: types.AgentAccess{
			TAQs: []types.AgentAccessTAQ{
				{ID: 10, Params: map[string]string{"queue": "urgent"}},
				{ID: 20},
			},
		},
	}

	t.Run("a pinned value replaces what was sent", func(t *testing.T) {
		out := pinTAQParams(agent, "10", map[string]any{"queue": "bulk", "note": "hi"})
		assert.Equal(t, "urgent", out["queue"])
		assert.Equal(t, "hi", out["note"])
	})

	t.Run("a pinned value is added when nothing was sent", func(t *testing.T) {
		out := pinTAQParams(agent, "10", map[string]any{})
		assert.Equal(t, "urgent", out["queue"])
	})

	t.Run("the caller's map is not mutated", func(t *testing.T) {
		in := map[string]any{"queue": "bulk"}
		pinTAQParams(agent, "10", in)
		assert.Equal(t, "bulk", in["queue"])
	})

	t.Run("a grant with no params changes nothing", func(t *testing.T) {
		in := map[string]any{"queue": "bulk"}
		assert.Equal(t, in, pinTAQParams(agent, "20", in))
	})

	t.Run("an ungranted TAQ changes nothing", func(t *testing.T) {
		in := map[string]any{"queue": "bulk"}
		assert.Equal(t, in, pinTAQParams(agent, "99", in))
	})
}
