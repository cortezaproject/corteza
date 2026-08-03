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

// TestStaticAutomationToolsAreNotReadAsTAQIDs guards the fix for a live bug:
// Evaluate used to treat any automation_* name outside a hardcoded three-item
// exception list as a per-TAQ tool name, trimming the prefix and looking the
// remainder up as a TAQ ID. That denied the three static tools below with the
// reason `agent is not allowed to execute automation "taq_exec"`.
func TestStaticAutomationToolsAreNotReadAsTAQIDs(t *testing.T) {
	ctx := context.Background()

	// automation_taq_exec maps to a real resource via buildResource — that case
	// was unreachable behind the old trap — so it legitimately requires an allow
	// entry. The other two are scope-exempt and need none. What matters for the
	// regression is that none of them is mistaken for a TAQ ID.
	cases := map[string][]types.AgentAccessAllow{
		"automation_taq_exec":            {{NamespaceID: 1}},
		"automation_taq_executions":      nil,
		"automation_taq_execution_trace": nil,
	}

	for tool, allow := range cases {
		t.Run(tool, func(t *testing.T) {
			agent := &types.Agent{
				Access: types.AgentAccess{
					Tools: []types.AgentAccessTool{{Name: tool, Allow: allow}},
				},
			}
			d := Evaluate(ctx, agent, tool, MapValues{"taq": "1"}, nil)
			assert.True(t, d.Allowed, "reason: %s", d.Reason)
			assert.NotContains(t, d.Reason, "not allowed to execute automation")
		})
	}
}

// TestTaqExecStillRequiresAllowEntry pins the consequence of un-shadowing the
// buildResource case: with the trap gone, automation_taq_exec is resource-scoped
// and an agent without allow entries is denied on that basis — not on the old
// bogus TAQ-ID basis.
func TestTaqExecStillRequiresAllowEntry(t *testing.T) {
	ctx := context.Background()
	agent := &types.Agent{
		Access: types.AgentAccess{
			Tools: []types.AgentAccessTool{{Name: "automation_taq_exec"}},
		},
	}
	d := Evaluate(ctx, agent, "automation_taq_exec", MapValues{"taq": "1"}, nil)
	assert.False(t, d.Allowed)
	assert.Contains(t, d.Reason, "no allow entries")
}

func TestDynamicTAQRef(t *testing.T) {
	cases := map[string]struct {
		ref string
		ok  bool
	}{
		"automation_123":                 {"123", true},
		"automation_taq_exec":            {"", false},
		"automation_taq_execution_trace": {"", false},
		"automation_workflow_lookup":     {"", false},
		"automation_":                    {"", false},
		"automation_12a":                 {"", false},
		"compose_record_lookup":          {"", false},
	}
	for tool, want := range cases {
		ref, ok := dynamicTAQRef(tool)
		assert.Equal(t, want.ok, ok, "tool %q", tool)
		assert.Equal(t, want.ref, ref, "tool %q", tool)
	}
}

// TestIsClassified covers the CI-time half. An unmapped tool is deliberately
// still allowed at runtime — denying would break every newly added tool until
// policy.go caught up — but IsClassified reports it so a test can fail instead.
func TestIsClassified(t *testing.T) {
	classified := []string{
		"compose_record_lookup",          // mapped by buildResource
		"compose_module_update",          // mapped by buildResource
		"automation_taq_exec",            // mapped by buildResource
		"compose_chart_lookup",           // scope-exempt
		"automation_taq_execution_trace", // scope-exempt
		"discovery_search",               // scope-exempt
		"automation_907",                 // runtime-minted per-TAQ tool
	}
	for _, tool := range classified {
		assert.True(t, IsClassified(tool), "expected %q to be classified", tool)
	}

	// Whole system_* families are exempt by prefix — a user or a role has no
	// compose/automation dimension by construction, so per-name entries would
	// be forty lines of noise that bury the one-offs that carry meaning.
	byPrefix := []string{
		"system_user_lookup",
		"system_role_lookup",
		"system_user_group_member_add",
		"system_auth_client_delete",
		"system_application_flag",
	}
	for _, tool := range byPrefix {
		assert.True(t, IsClassified(tool), "expected %q to be classified by prefix", tool)
		assert.True(t, IsScopeExempt(tool), "expected %q to be scope-exempt", tool)
	}

	unclassified := []string{
		"system_settings_lookup", // no prefix entry: a real new family must be ruled on
		"totally_made_up",
	}
	for _, tool := range unclassified {
		assert.False(t, IsClassified(tool), "expected %q to be unclassified", tool)
	}
}

// TestUnmappedToolStillRuns pins the runtime behaviour: unclassified is not the
// same as denied.
func TestUnmappedToolStillRuns(t *testing.T) {
	ctx := context.Background()
	agent := &types.Agent{
		Access: types.AgentAccess{
			Tools: []types.AgentAccessTool{{Name: "system_user_lookup"}},
		},
	}
	d := Evaluate(ctx, agent, "system_user_lookup", MapValues{}, nil)
	assert.True(t, d.Allowed, "reason: %s", d.Reason)
}

func TestScopeExemptToolIsAllowed(t *testing.T) {
	ctx := context.Background()
	agent := &types.Agent{
		Access: types.AgentAccess{
			Tools: []types.AgentAccessTool{{Name: "compose_chart_lookup"}},
		},
	}
	d := Evaluate(ctx, agent, "compose_chart_lookup", MapValues{}, nil)
	assert.True(t, d.Allowed, "reason: %s", d.Reason)
}
