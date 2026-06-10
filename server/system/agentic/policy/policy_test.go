package policy

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	ctx := context.Background()

	t.Run("denied when tool not in allow-list", func(t *testing.T) {
		agent := &types.Agent{}
		d := Evaluate(ctx, agent, "compose_record_lookup", MapValues(nil), nil)
		assert.False(t, d.Allowed)
		assert.Contains(t, d.Reason, "not in the agent's allow-list")
	})

	t.Run("allowed when tool is in allow-list with allow entry", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup", Allow: []types.AgentAccessAllow{{NamespaceID: 1}}},
				},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_lookup", MapValues{"namespaceID": "1", "recordID": "123"}, nil)
		assert.True(t, d.Allowed)
		assert.Equal(t, "123", d.SanitizedArgs["recordID"])
	})

	t.Run("denied when tool has no allow entries", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup"},
				},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_lookup", MapValues{"namespaceID": "1", "recordID": "123"}, nil)
		assert.False(t, d.Allowed)
	})

	t.Run("global context defaults fill missing args", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Context: types.AgentAccessContext{
					Defaults: MapValues{"namespaceID": "1"},
				},
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup", Allow: []types.AgentAccessAllow{{NamespaceID: 1}}},
				},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_lookup", MapValues{"namespaceID": "1"}, nil)
		assert.True(t, d.Allowed)
		assert.Equal(t, "1", d.SanitizedArgs["namespaceID"])
	})

	t.Run("global context defaults do not overwrite existing args", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Context: types.AgentAccessContext{
					Defaults: MapValues{"namespaceID": "1"},
				},
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup", Allow: []types.AgentAccessAllow{{NamespaceID: 1}}},
				},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_lookup", MapValues{"namespaceID": "1"}, nil)
		assert.Equal(t, "1", d.SanitizedArgs["namespaceID"])
	})

	t.Run("tool-level defaults fill missing args", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{
						Name:  "compose_record_lookup",
						Allow: []types.AgentAccessAllow{{NamespaceID: 1}},
						Context: types.AgentAccessToolContext{
							Defaults: MapValues{"moduleID": "mod1"},
						},
					},
				},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_lookup", MapValues{"namespaceID": "1"}, nil)
		assert.Equal(t, "mod1", d.SanitizedArgs["moduleID"])
	})

	t.Run("tool-level overrides always apply", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{
						Name:  "compose_record_lookup",
						Allow: []types.AgentAccessAllow{{NamespaceID: 1}},
						Context: types.AgentAccessToolContext{
							Overrides: MapValues{"someField": "forced-value"},
						},
					},
				},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_lookup", MapValues{"namespaceID": "1"}, nil)
		assert.True(t, d.Allowed)
		assert.Equal(t, "forced-value", d.SanitizedArgs["someField"])
	})

	t.Run("denied when namespaceID does not match allow entry", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{
						Name:  "compose_record_create",
						Allow: []types.AgentAccessAllow{{NamespaceID: 100, ModuleIDs: types.AgentAccessIDList{200}}},
					},
				},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_create", MapValues{"namespaceID": "999", "moduleID": "200"}, nil)
		assert.False(t, d.Allowed)
	})

	t.Run("allowed when namespace-level allow covers module tool", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{
						Name:  "compose_module_lookup",
						Allow: []types.AgentAccessAllow{{NamespaceID: 100}},
					},
				},
			},
		}
		d := Evaluate(ctx, agent, "compose_module_lookup", MapValues{"namespaceID": "100", "moduleID": "999"}, nil)
		assert.True(t, d.Allowed)
	})

	t.Run("denied when namespace is wildcard with allow entries present", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{
						Name:  "compose_namespace_lookup",
						Allow: []types.AgentAccessAllow{{NamespaceID: 100}},
					},
				},
			},
		}
		// No namespaceID resolved — wildcard is denied when allow entries are present
		d := Evaluate(ctx, agent, "compose_namespace_lookup", MapValues{}, nil)
		assert.False(t, d.Allowed)
	})

	t.Run("taq denied when id not in list", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				TAQs: []types.AgentAccessTAQ{{ID: 111}},
			},
		}
		d := Evaluate(ctx, agent, "automation_999", MapValues{}, nil)
		assert.False(t, d.Allowed)
	})

	t.Run("taq allowed when id matches", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				TAQs: []types.AgentAccessTAQ{{ID: 111}},
			},
		}
		d := Evaluate(ctx, agent, "automation_111", MapValues{}, nil)
		assert.True(t, d.Allowed)
	})

	t.Run("ownsTarget fallback overrides allow-list miss", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{Name: "compose_module_lookup", Allow: []types.AgentAccessAllow{{NamespaceID: 100}}},
				},
			},
		}
		owns := func(_ context.Context, _ ValueGetter) bool { return true }
		d := Evaluate(ctx, agent, "compose_module_lookup", MapValues{"namespaceID": "999", "moduleID": "1"}, owns)
		assert.True(t, d.Allowed)
	})
}
