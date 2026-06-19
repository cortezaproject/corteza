package service

import (
	"context"
	"testing"

	automationTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

// noExternal is the loadExternal stub for tests with no out-of-scope targets.
func noExternal(context.Context, string, uint64) (*GraphSource, error) { return nil, nil }

func TestAssembleProjectGraphIDRefs(t *testing.T) {
	sources := []*GraphSource{
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 1, Name: "Lead", Sensitivity: "internal"},
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 2, Name: "Opportunity", Refs: []resourceref.Ref{
			resourceref.Make(resourceref.KindComposeModule, 1, resourceref.ReasonModuleFieldRef),
		}},
	}

	g, err := assembleProjectGraph(context.Background(), sources, noExternal)
	require.NoError(t, err)

	require.Equal(t, []*types.ProjectGraphNode{
		{ID: 1, Kind: "module", Name: "Lead", Sensitivity: "internal"},
		{ID: 2, Kind: "module", Name: "Opportunity"},
	}, g.Nodes)

	require.Equal(t, []*types.ProjectGraphEdge{
		{SourceID: 2, TargetID: 1, Reason: resourceref.ReasonModuleFieldRef},
	}, g.Edges)

	require.Empty(t, g.Missing)
}

func TestAssembleProjectGraphHandleResolution(t *testing.T) {
	sources := []*GraphSource{
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 1, Name: "Lead", Handle: "lead"},
		{ResourceType: resourceref.KindNgAutomation, Kind: "automation", ID: 5, Name: "Lead Scoring", Refs: []resourceref.Ref{
			resourceref.MakeIdent(resourceref.KindComposeModule, "lead", resourceref.ReasonTriggerModule),
		}},
	}

	g, err := assembleProjectGraph(context.Background(), sources, noExternal)
	require.NoError(t, err)

	require.Equal(t, []*types.ProjectGraphEdge{
		{SourceID: 5, TargetID: 1, Reason: resourceref.ReasonTriggerModule},
	}, g.Edges)

	require.Empty(t, g.Missing)
}

func TestAssembleProjectGraphExternalTarget(t *testing.T) {
	sources := []*GraphSource{
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 1, Name: "Lead", Refs: []resourceref.Ref{
			resourceref.Make(resourceref.KindDalConnection, 40, resourceref.ReasonModuleConnection),
		}},
		// second ref to same connection must reuse the node and dedup the edge
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 2, Name: "Quote", Refs: []resourceref.Ref{
			resourceref.Make(resourceref.KindDalConnection, 40, resourceref.ReasonModuleConnection),
		}},
	}

	loadExternal := func(_ context.Context, rt string, id uint64) (*GraphSource, error) {
		if rt != resourceref.KindDalConnection || id != 40 {
			return nil, nil
		}
		return &GraphSource{ResourceType: rt, Kind: "connection", ID: 40, Name: "Google Sheets", Sensitivity: "confidential"}, nil
	}

	g, err := assembleProjectGraph(context.Background(), sources, loadExternal)
	require.NoError(t, err)

	require.Equal(t, []*types.ProjectGraphNode{
		{ID: 1, Kind: "module", Name: "Lead"},
		{ID: 2, Kind: "module", Name: "Quote"},
		{ID: 40, Kind: "connection", Name: "Google Sheets", Sensitivity: "confidential", External: true},
	}, g.Nodes)

	require.Equal(t, []*types.ProjectGraphEdge{
		{SourceID: 1, TargetID: 40, Reason: resourceref.ReasonModuleConnection},
		{SourceID: 2, TargetID: 40, Reason: resourceref.ReasonModuleConnection},
	}, g.Edges)
}

func TestAssembleProjectGraphMissingRefs(t *testing.T) {
	sources := []*GraphSource{
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 1, Name: "Lead", Refs: []resourceref.Ref{
			// ID target absent, no external loader hit
			resourceref.Make(resourceref.KindComposeModule, 99, resourceref.ReasonModuleFieldRef),
			// handle resolves to nothing
			resourceref.MakeIdent(resourceref.KindComposeModule, "ghost", resourceref.ReasonModuleFieldRef),
			// target kind the external loader does not know
			resourceref.Make(resourceref.KindLlmProvider, 7, resourceref.ReasonAgentLlmProvider),
			// external loader returns nil (deleted resource)
			resourceref.Make(resourceref.KindDalConnection, 41, resourceref.ReasonModuleConnection),
		}},
	}

	// every external lookup misses
	g, err := assembleProjectGraph(context.Background(), sources, noExternal)
	require.NoError(t, err)
	require.Empty(t, g.Edges)

	require.Equal(t, []*types.ProjectGraphMissingRef{
		{SourceID: 1, Kind: resourceref.KindComposeModule, TargetID: 99, Reason: resourceref.ReasonModuleFieldRef},
		{SourceID: 1, Kind: resourceref.KindComposeModule, TargetIdent: "ghost", Reason: resourceref.ReasonModuleFieldRef},
		{SourceID: 1, Kind: resourceref.KindLlmProvider, TargetID: 7, Reason: resourceref.ReasonAgentLlmProvider},
		{SourceID: 1, Kind: resourceref.KindDalConnection, TargetID: 41, Reason: resourceref.ReasonModuleConnection},
	}, g.Missing)
}

func TestAssembleProjectGraphEdgeDedup(t *testing.T) {
	sources := []*GraphSource{
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 1, Name: "Lead"},
		// two fields referencing the same module collapse into one edge
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 2, Name: "Opportunity", Refs: []resourceref.Ref{
			resourceref.Make(resourceref.KindComposeModule, 1, resourceref.ReasonModuleFieldRef),
			resourceref.Make(resourceref.KindComposeModule, 1, resourceref.ReasonModuleFieldRef),
		}},
	}

	g, err := assembleProjectGraph(context.Background(), sources, noExternal)
	require.NoError(t, err)
	require.Len(t, g.Edges, 1)
}

func TestAssembleProjectGraphDynamicWarnings(t *testing.T) {
	sources := []*GraphSource{
		{ResourceType: resourceref.KindNgAutomation, Kind: "automation", ID: 5, Name: "Scoring", Refs: []resourceref.Ref{
			resourceref.MakeDynamic(resourceref.KindComposeModule, resourceref.ReasonStepArgument),
		}},
	}

	g, err := assembleProjectGraph(context.Background(), sources, noExternal)
	require.NoError(t, err)
	require.Empty(t, g.Edges)
	require.Empty(t, g.Missing)

	require.Equal(t, []*types.ProjectGraphWarning{
		{SourceID: 5, Kind: resourceref.KindComposeModule, Reason: resourceref.ReasonStepArgument},
	}, g.Warnings)
}

// TestProjectGraphNgAutomationResolution drives a realistic NgAutomation config
// through the real extractor and assembler: trigger constraint by handle,
// constant step argument, connection function, computed argument and a dangling
// target all resolve to their respective graph outputs.
func TestProjectGraphNgAutomationResolution(t *testing.T) {
	au := automationTypes.NgAutomation{
		ID: 500,
		Triggers: automationTypes.NgAutomationTriggerSet{
			{
				ResourceType: "compose:record",
				EventType:    "afterCreate",
				Constraints: []automationTypes.NgTriggerConstraint{
					{Name: "module", Values: []automationTypes.NgTriggerConstraintValue{{Value: "lead"}}},
				},
			},
		},
		Steps: automationTypes.NgAutomationStepSet{
			{Kind: "function", Ref: "composeRecordsCreate", Arguments: []*automationTypes.Expr{
				{Target: "module", Value: "1002"},
			}},
			{Kind: "function", Ref: "conn_40_sheetAppend", Arguments: []*automationTypes.Expr{
				{Target: "range", Value: "A1"},
			}},
			{Kind: "function", Ref: "composeRecordsSearch", Arguments: []*automationTypes.Expr{
				{Target: "module", Expr: "scope.mod"},
			}},
			// references a module that no longer exists
			{Kind: "function", Ref: "composeRecordsUpdate", Arguments: []*automationTypes.Expr{
				{Target: "module", Value: "9999"},
			}},
		},
	}

	sources := []*GraphSource{
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 1001, Name: "Lead", Handle: "lead"},
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 1002, Name: "Opportunity", Handle: "opportunity"},
		// what a real store-backed loader will produce
		{ResourceType: resourceref.KindNgAutomation, Kind: "automation", ID: au.ID, Name: "Lead Scoring", Refs: au.ResourceRefs()},
	}

	loadExternal := func(_ context.Context, rt string, id uint64) (*GraphSource, error) {
		if rt != resourceref.KindConfiguredConnection || id != 40 {
			return nil, nil
		}
		return &GraphSource{ResourceType: rt, Kind: "connection", ID: 40, Name: "Google Sheets", Sensitivity: "confidential"}, nil
	}

	g, err := assembleProjectGraph(context.Background(), sources, loadExternal)
	require.NoError(t, err)

	require.Equal(t, []*types.ProjectGraphEdge{
		// trigger constraint resolved through the module handle index
		{SourceID: 500, TargetID: 1001, Reason: resourceref.ReasonTriggerModule},
		// constant step argument resolved by ID
		{SourceID: 500, TargetID: 1002, Reason: resourceref.ReasonStepArgument},
		// connection function ref resolved to an external connection node
		{SourceID: 500, TargetID: 40, Reason: resourceref.ReasonStepConnection},
	}, g.Edges)

	require.Equal(t, []*types.ProjectGraphWarning{
		{SourceID: 500, Kind: resourceref.KindComposeModule, Reason: resourceref.ReasonStepArgument},
	}, g.Warnings)

	require.Equal(t, []*types.ProjectGraphMissingRef{
		{SourceID: 500, Kind: resourceref.KindComposeModule, TargetID: 9999, Reason: resourceref.ReasonStepArgument},
	}, g.Missing)

	byID := make(map[uint64]*types.ProjectGraphNode)
	for _, n := range g.Nodes {
		byID[n.ID] = n
	}
	require.Len(t, g.Nodes, 4)
	require.True(t, byID[40].External)
	require.Equal(t, "confidential", byID[40].Sensitivity)
}

// TestRoleRuleRefs drives RBAC rules through the role->resource extractor:
// specific resources resolve by ID, whole-kind grants (trailing wildcard)
// become wildcard refs, per-operation duplicates collapse, and rules on
// non-graph resources (users, records) are dropped.
func TestRoleRuleRefs(t *testing.T) {
	rules := rbac.RuleSet{
		{RoleID: 1, Resource: "corteza::compose:module/12/34", Operation: "read"},
		{RoleID: 1, Resource: "corteza::compose:module/12/34", Operation: "update"}, // dup -> collapse
		{RoleID: 1, Resource: "corteza::compose:module/*/*", Operation: "read"},      // whole kind -> wildcard
		{RoleID: 1, Resource: "corteza::compose:page/*/55", Operation: "delete"},     // specific page, any namespace
		{RoleID: 1, Resource: "corteza::system:user/*", Operation: "read"},           // non-graph -> skip
		{RoleID: 1, Resource: "corteza::compose:record/12/34/7", Operation: "read"},  // non-graph -> skip
	}

	require.Equal(t, []resourceref.Ref{
		resourceref.Make(resourceref.KindComposeModule, 34, resourceref.ReasonRoleRbac),
		resourceref.MakeWildcard(resourceref.KindComposeModule, resourceref.ReasonRoleRbac),
		resourceref.Make(resourceref.KindComposePage, 55, resourceref.ReasonRoleRbac),
	}, roleRuleRefs(rules))
}

// TestAssembleProjectGraphWildcardFanout verifies a wildcard role ref expands
// to one edge per in-scope node of the kind and skips the role's self-edge.
func TestAssembleProjectGraphWildcardFanout(t *testing.T) {
	sources := []*GraphSource{
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 1, Name: "Lead"},
		{ResourceType: resourceref.KindComposeModule, Kind: "module", ID: 2, Name: "Quote"},
		{ResourceType: resourceref.KindRole, Kind: "role", ID: 10, Name: "Admin", Refs: []resourceref.Ref{
			resourceref.MakeWildcard(resourceref.KindComposeModule, resourceref.ReasonRoleRbac),
			// wildcard over roles must not produce a 10 -> 10 self-edge
			resourceref.MakeWildcard(resourceref.KindRole, resourceref.ReasonRoleRbac),
		}},
		{ResourceType: resourceref.KindRole, Kind: "role", ID: 11, Name: "Editor", Refs: []resourceref.Ref{
			resourceref.Make(resourceref.KindComposeModule, 1, resourceref.ReasonRoleRbac),
		}},
	}

	g, err := assembleProjectGraph(context.Background(), sources, noExternal)
	require.NoError(t, err)

	require.Equal(t, []*types.ProjectGraphEdge{
		{SourceID: 10, TargetID: 1, Reason: resourceref.ReasonRoleRbac},
		{SourceID: 10, TargetID: 2, Reason: resourceref.ReasonRoleRbac},
		{SourceID: 10, TargetID: 11, Reason: resourceref.ReasonRoleRbac},
		{SourceID: 11, TargetID: 1, Reason: resourceref.ReasonRoleRbac},
	}, g.Edges)

	require.Empty(t, g.Missing)
}
