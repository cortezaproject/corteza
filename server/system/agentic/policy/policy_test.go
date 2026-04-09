package policy

import (
	"testing"

	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	t.Run("denied when tool not in allow-list", func(t *testing.T) {
		agent := &types.Agent{}
		d := Evaluate(agent, "compose_record_lookup", MapValues(nil))
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
		d := Evaluate(agent, "compose_record_lookup", MapValues{"recordID": "123"})
		assert.True(t, d.Allowed)
		assert.Equal(t, "123", d.SanitizedArgs["recordID"])
	})

	t.Run("global context defaults fill missing args", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Context: types.AgentAccessContext{
					Defaults: MapValues{"namespaceID": "ns1"},
				},
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup"},
				},
			},
		}
		d := Evaluate(agent, "compose_record_lookup", MapValues{})
		assert.True(t, d.Allowed)
		assert.Equal(t, "ns1", d.SanitizedArgs["namespaceID"])
	})

	t.Run("global context defaults do not overwrite existing args", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Context: types.AgentAccessContext{
					Defaults: MapValues{"namespaceID": "ns1"},
				},
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup"},
				},
			},
		}
		d := Evaluate(agent, "compose_record_lookup", MapValues{"namespaceID": "ns-custom"})
		assert.Equal(t, "ns-custom", d.SanitizedArgs["namespaceID"])
	})

	t.Run("tool-level defaults fill missing args", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{
						Name: "compose_record_lookup",
						Context: types.AgentAccessToolContext{
							Defaults: MapValues{"moduleID": "mod1"},
						},
					},
				},
			},
		}
		d := Evaluate(agent, "compose_record_lookup", MapValues{})
		assert.Equal(t, "mod1", d.SanitizedArgs["moduleID"])
	})

	t.Run("tool-level overrides always apply", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{
						Name: "compose_record_lookup",
						Context: types.AgentAccessToolContext{
							Overrides: MapValues{"namespaceID": "forced-ns"},
						},
					},
				},
			},
		}
		d := Evaluate(agent, "compose_record_lookup", MapValues{"namespaceID": "user-ns"})
		assert.Equal(t, "forced-ns", d.SanitizedArgs["namespaceID"])
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
		d := Evaluate(agent, "compose_record_create", MapValues{"namespaceID": "999", "moduleID": "200"})
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
		d := Evaluate(agent, "compose_module_lookup", MapValues{"namespaceID": "100", "moduleID": "999"})
		assert.True(t, d.Allowed)
	})

	t.Run("taq denied when id not in list", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				TAQs: []types.AgentAccessTAQ{{ID: 111}},
			},
		}
		d := Evaluate(agent, "automation_taq_exec", MapValues{"taq": "999"})
		assert.False(t, d.Allowed)
	})

	t.Run("taq allowed when id matches", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				TAQs: []types.AgentAccessTAQ{{ID: 111}},
			},
		}
		d := Evaluate(agent, "automation_taq_exec", MapValues{"taq": "111"})
		assert.True(t, d.Allowed)
	})
}

