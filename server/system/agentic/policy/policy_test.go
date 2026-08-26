package policy

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// A named tool with no allow entries is not narrowed, and the agent reaches
	// whatever the invoking user reaches — it runs as that user. This used to
	// deny, which is what made an agent unusable until every tool was scoped by
	// hand.
	t.Run("a tool with no allow entries is not narrowed", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{
					{Name: "compose_record_lookup"},
				},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_lookup", MapValues{"namespaceID": "1", "recordID": "123"}, nil)
		assert.True(t, d.Allowed, d.Reason)
	})

	// The agent's own scope holds however its tools were granted — including
	// the ones it inherits, which carry no allow list of their own.
	t.Run("the agent scope narrows an unnarrowed tool", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Allow: []types.AgentAccessAllow{{NamespaceID: 100}},
				Tools: []types.AgentAccessTool{{Name: "compose_record_lookup"}},
			},
		}

		in := Evaluate(ctx, agent, "compose_record_lookup", MapValues{"namespaceID": "100"}, nil)
		assert.True(t, in.Allowed, in.Reason)

		out := Evaluate(ctx, agent, "compose_record_lookup", MapValues{"namespaceID": "200"}, nil)
		assert.False(t, out.Allowed)
	})

	t.Run("a tool set to deny is refused whatever its scope says", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{{
					Name:       "compose_record_delete",
					Permission: PermissionDeny,
					Allow:      []types.AgentAccessAllow{{NamespaceID: 100}},
				}},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_delete", MapValues{"namespaceID": "100"}, nil)
		assert.False(t, d.Allowed)
		assert.Contains(t, d.Reason, "deny")
	})

	t.Run("the permission mode travels with the decision", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{{Name: "compose_record_create", Permission: PermissionAsk}},
			},
		}
		d := Evaluate(ctx, agent, "compose_record_create", MapValues{"namespaceID": "100"}, nil)
		require.True(t, d.Allowed, d.Reason)
		assert.Equal(t, PermissionAsk, d.Permission)
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

	// automation_taq_exec is gated on access.taqs, like its workflow twin; the
	// other two are scope-exempt and need nothing. What matters for the
	// regression is that none of them is mistaken for a TAQ ID.
	cases := map[string]types.AgentAccess{
		"automation_taq_exec": {
			Tools: []types.AgentAccessTool{{Name: "automation_taq_exec"}},
			TAQs:  []types.AgentAccessTAQ{{ID: 1}},
		},
		"automation_taq_executions":      {Tools: []types.AgentAccessTool{{Name: "automation_taq_executions"}}},
		"automation_taq_execution_trace": {Tools: []types.AgentAccessTool{{Name: "automation_taq_execution_trace"}}},
	}

	for tool, access := range cases {
		t.Run(tool, func(t *testing.T) {
			d := Evaluate(ctx, &types.Agent{Access: access}, tool, MapValues{"taq": "1"}, nil)
			assert.True(t, d.Allowed, "reason: %s", d.Reason)
			assert.NotContains(t, d.Reason, "not allowed to execute automation")
		})
	}
}

// TestTaqExecNeedsTheTAQGranted pins what an allow entry can and cannot do for
// automation_taq_exec.
//
// A compose allow entry describes a namespace and its modules and says nothing
// about an automation, so it must not stand in for one: an agent granted the
// tool and a namespace could otherwise run every TAQ on the instance, which is
// what a "usage" group grant at write risk quietly bought.
func TestTaqExecNeedsTheTAQGranted(t *testing.T) {
	ctx := context.Background()

	t.Run("a namespace allow entry does not cover a TAQ", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{
				Tools: []types.AgentAccessTool{{
					Name:  "automation_taq_exec",
					Allow: []types.AgentAccessAllow{{NamespaceID: 100}},
				}},
			},
		}
		d := Evaluate(ctx, agent, "automation_taq_exec", MapValues{"taq": "999"}, nil)
		assert.False(t, d.Allowed)
		assert.Contains(t, d.Reason, "access.taqs")
	})

	t.Run("a granted TAQ runs", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{TAQs: []types.AgentAccessTAQ{{ID: 999}}},
		}
		d := Evaluate(ctx, agent, "automation_taq_exec", MapValues{"taq": "999"}, nil)
		assert.True(t, d.Allowed, d.Reason)
	})

	t.Run("another TAQ does not", func(t *testing.T) {
		agent := &types.Agent{
			Access: types.AgentAccess{TAQs: []types.AgentAccessTAQ{{ID: 999}}},
		}
		d := Evaluate(ctx, agent, "automation_taq_exec", MapValues{"taq": "1000"}, nil)
		assert.False(t, d.Allowed)
	})
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

// An automation does whatever its steps do and has no dry run. Granting one
// says the agent MAY run it, not that it may run it unannounced.
func TestAutomationRunsAskByDefault(t *testing.T) {
	ctx := context.Background()
	agent := &types.Agent{Access: types.AgentAccess{
		TAQs: []types.AgentAccessTAQ{
			{ID: 1},
			{ID: 2, Permission: PermissionAlways},
			{ID: 3, Permission: PermissionDeny},
		},
		Workflows: []types.AgentAccessWorkflow{{ID: 9}},
	}}

	t.Run("an unstated mode asks", func(t *testing.T) {
		d := Evaluate(ctx, agent, "automation_taq_exec", MapValues{"taq": "1"}, nil)
		require.True(t, d.Allowed, d.Reason)
		assert.Equal(t, PermissionAsk, d.Permission)
	})

	t.Run("the minted per-TAQ tool asks the same way", func(t *testing.T) {
		d := Evaluate(ctx, agent, "automation_1", MapValues{}, nil)
		require.True(t, d.Allowed, d.Reason)
		assert.Equal(t, PermissionAsk, d.Permission)
	})

	t.Run("a stated mode is kept", func(t *testing.T) {
		d := Evaluate(ctx, agent, "automation_taq_exec", MapValues{"taq": "2"}, nil)
		require.True(t, d.Allowed, d.Reason)
		assert.Equal(t, PermissionAlways, d.Permission)
	})

	t.Run("deny refuses the run", func(t *testing.T) {
		d := Evaluate(ctx, agent, "automation_taq_exec", MapValues{"taq": "3"}, nil)
		assert.False(t, d.Allowed)
	})

	t.Run("a workflow follows the same rule", func(t *testing.T) {
		d := Evaluate(ctx, agent, "automation_workflow_exec", MapValues{"workflow": "9"}, nil)
		require.True(t, d.Allowed, d.Reason)
		assert.Equal(t, PermissionAsk, d.Permission)
	})
}
