package policy

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	t.Run("denied when tool not in allow-list", func(t *testing.T) {
		agent := &types.Agent{}
		d := Evaluate(agent, "compose_record_lookup", nil)
		assert.False(t, d.Allowed)
		assert.Contains(t, d.Reason, "not in the agent's allow-list")
	})

	t.Run("allowed when tool is in allow-list", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup"},
				},
			},
		}
		d := Evaluate(agent, "compose_record_lookup", map[string]any{"recordID": "123"})
		assert.True(t, d.Allowed)
		assert.Equal(t, "123", d.SanitizedArgs["recordID"])
	})

	t.Run("global context defaults fill missing args", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Context: types.AgentAccessContext{
					Defaults: map[string]any{"namespaceID": "ns1"},
				},
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup"},
				},
			},
		}
		d := Evaluate(agent, "compose_record_lookup", map[string]any{})
		assert.True(t, d.Allowed)
		assert.Equal(t, "ns1", d.SanitizedArgs["namespaceID"])
	})

	t.Run("global context defaults do not overwrite existing args", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Context: types.AgentAccessContext{
					Defaults: map[string]any{"namespaceID": "ns1"},
				},
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup"},
				},
			},
		}
		d := Evaluate(agent, "compose_record_lookup", map[string]any{"namespaceID": "ns-custom"})
		assert.Equal(t, "ns-custom", d.SanitizedArgs["namespaceID"])
	})

	t.Run("tool-level defaults fill missing args", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{
						Name: "compose_record_lookup",
						Context: types.AgentAccessToolContext{
							Defaults: map[string]any{"moduleID": "mod1"},
						},
					},
				},
			},
		}
		d := Evaluate(agent, "compose_record_lookup", map[string]any{})
		assert.Equal(t, "mod1", d.SanitizedArgs["moduleID"])
	})

	t.Run("tool-level overrides always apply", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{
						Name: "compose_record_lookup",
						Context: types.AgentAccessToolContext{
							Overrides: map[string]any{"namespaceID": "forced-ns"},
						},
					},
				},
			},
		}
		d := Evaluate(agent, "compose_record_lookup", map[string]any{"namespaceID": "user-ns"})
		assert.Equal(t, "forced-ns", d.SanitizedArgs["namespaceID"])
	})
}

func TestFilterResponse(t *testing.T) {
	ctx := context.Background()

	t.Run("no matching allow entry returns empty map", func(t *testing.T) {
		agent := &types.Agent{}
		result := FilterResponse(ctx, agent, "compose_record_lookup", map[string]any{"recordID": "123", "values": "data"})
		assert.Empty(t, result)
	})

	t.Run("matching entry with no properties returns all fields", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Allow: []types.AgentAccessAllow{
					{Resource: "compose_record_lookup"},
				},
			},
		}
		data := map[string]any{"recordID": "123", "values": "data"}
		result := FilterResponse(ctx, agent, "compose_record_lookup", data)
		assert.Equal(t, data, result)
	})

	t.Run("only allowed properties are returned", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Allow: []types.AgentAccessAllow{
					{
						Resource: "compose_record_lookup",
						Properties: []types.AgentAccessAllowProperty{
							{Name: "recordID", Access: "allow"},
						},
					},
				},
			},
		}
		result := FilterResponse(ctx, agent, "compose_record_lookup", map[string]any{"recordID": "123", "secret": "hidden"})
		assert.Equal(t, "123", result["recordID"])
		assert.NotContains(t, result, "secret")
	})

	t.Run("row-level filter match returns data", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Allow: []types.AgentAccessAllow{
					{
						Resource: "compose_record_lookup",
						Filter:   "status == \"active\"",
					},
				},
			},
		}
		result := FilterResponse(ctx, agent, "compose_record_lookup", map[string]any{"status": "active", "recordID": "123"})
		assert.Equal(t, "123", result["recordID"])
	})

	t.Run("row-level filter mismatch returns empty map", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Allow: []types.AgentAccessAllow{
					{
						Resource: "compose_record_lookup",
						Filter:   "status == \"active\"",
					},
				},
			},
		}
		result := FilterResponse(ctx, agent, "compose_record_lookup", map[string]any{"status": "inactive", "recordID": "123"})
		assert.Empty(t, result)
	})
}
