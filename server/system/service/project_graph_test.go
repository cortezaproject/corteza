package service

import (
	"context"
	"fmt"
	"testing"

	automationTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestAssembleProjectGraphIDRefs(t *testing.T) {
	registry := []*GraphKind{
		{
			ResourceType: resourceref.KindComposeModule,
			Kind:         "module",
			Load: staticGraphSources(
				&GraphSource{ID: 1, Name: "Lead", Sensitivity: "internal"},
				&GraphSource{ID: 2, Name: "Opportunity", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindComposeModule, 1, resourceref.ReasonModuleFieldRef),
				}},
			),
		},
	}

	g, err := assembleProjectGraph(context.Background(), 0, registry)
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
	registry := []*GraphKind{
		{
			ResourceType: resourceref.KindComposeModule,
			Kind:         "module",
			Load: staticGraphSources(
				&GraphSource{ID: 1, Name: "Lead", Handle: "lead"},
			),
		},
		{
			ResourceType: resourceref.KindNgAutomation,
			Kind:         "automation",
			Load: staticGraphSources(
				&GraphSource{ID: 5, Name: "Lead Scoring", Refs: []resourceref.Ref{
					resourceref.MakeIdent(resourceref.KindComposeModule, "lead", resourceref.ReasonTriggerModule),
				}},
			),
		},
	}

	g, err := assembleProjectGraph(context.Background(), 0, registry)
	require.NoError(t, err)

	require.Equal(t, []*types.ProjectGraphEdge{
		{SourceID: 5, TargetID: 1, Reason: resourceref.ReasonTriggerModule},
	}, g.Edges)

	require.Empty(t, g.Missing)
}

func TestAssembleProjectGraphExternalTarget(t *testing.T) {
	registry := []*GraphKind{
		{
			ResourceType: resourceref.KindComposeModule,
			Kind:         "module",
			Load: staticGraphSources(
				&GraphSource{ID: 1, Name: "Lead", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindDalConnection, 40, resourceref.ReasonModuleConnection),
				}},
				// second ref to same connection must reuse the node and dedup the edge
				&GraphSource{ID: 2, Name: "Quote", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindDalConnection, 40, resourceref.ReasonModuleConnection),
				}},
			),
		},
		{
			ResourceType: resourceref.KindDalConnection,
			Kind:         "connection",
			LoadOne: func(_ context.Context, id uint64) (*GraphSource, error) {
				if id != 40 {
					return nil, nil
				}
				return &GraphSource{ID: 40, Name: "Google Sheets", Sensitivity: "confidential"}, nil
			},
		},
	}

	g, err := assembleProjectGraph(context.Background(), 0, registry)
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
	registry := []*GraphKind{
		{
			ResourceType: resourceref.KindComposeModule,
			Kind:         "module",
			Load: staticGraphSources(
				&GraphSource{ID: 1, Name: "Lead", Refs: []resourceref.Ref{
					// ID target absent, kind has no LoadOne
					resourceref.Make(resourceref.KindComposeModule, 99, resourceref.ReasonModuleFieldRef),
					// handle resolves to nothing
					resourceref.MakeIdent(resourceref.KindComposeModule, "ghost", resourceref.ReasonModuleFieldRef),
					// target kind not registered at all
					resourceref.Make(resourceref.KindLlmProvider, 7, resourceref.ReasonAgentLlmProvider),
					// LoadOne returns nil (deleted resource)
					resourceref.Make(resourceref.KindDalConnection, 41, resourceref.ReasonModuleConnection),
				}},
			),
		},
		{
			ResourceType: resourceref.KindDalConnection,
			Kind:         "connection",
			LoadOne: func(context.Context, uint64) (*GraphSource, error) {
				return nil, nil
			},
		},
	}

	g, err := assembleProjectGraph(context.Background(), 0, registry)
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
	registry := []*GraphKind{
		{
			ResourceType: resourceref.KindComposeModule,
			Kind:         "module",
			Load: staticGraphSources(
				&GraphSource{ID: 1, Name: "Lead"},
				// two fields referencing the same module collapse into one edge
				&GraphSource{ID: 2, Name: "Opportunity", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindComposeModule, 1, resourceref.ReasonModuleFieldRef),
					resourceref.Make(resourceref.KindComposeModule, 1, resourceref.ReasonModuleFieldRef),
				}},
			),
		},
	}

	g, err := assembleProjectGraph(context.Background(), 0, registry)
	require.NoError(t, err)
	require.Len(t, g.Edges, 1)
}

func TestAssembleProjectGraphLoadError(t *testing.T) {
	registry := []*GraphKind{
		{
			ResourceType: resourceref.KindComposeModule,
			Kind:         "module",
			Load: func(context.Context, uint64) ([]*GraphSource, error) {
				return nil, fmt.Errorf("store down")
			},
		},
	}

	_, err := assembleProjectGraph(context.Background(), 0, registry)
	require.ErrorContains(t, err, "store down")
}

func TestAssembleProjectGraphDynamicWarnings(t *testing.T) {
	registry := []*GraphKind{
		{
			ResourceType: resourceref.KindNgAutomation,
			Kind:         "automation",
			Load: staticGraphSources(
				&GraphSource{ID: 5, Name: "Scoring", Refs: []resourceref.Ref{
					resourceref.MakeDynamic(resourceref.KindComposeModule, resourceref.ReasonStepArgument),
				}},
			),
		},
	}

	g, err := assembleProjectGraph(context.Background(), 0, registry)
	require.NoError(t, err)
	require.Empty(t, g.Edges)
	require.Empty(t, g.Missing)

	require.Equal(t, []*types.ProjectGraphWarning{
		{SourceID: 5, Kind: resourceref.KindComposeModule, Reason: resourceref.ReasonStepArgument},
	}, g.Warnings)
}

// TestProjectGraphNgAutomationResolution drives a realistic NgAutomation
// config through the real extractor and assembler: trigger constraint by
// handle, constant step argument, connection function, computed argument and
// a dangling target all resolve to their respective graph outputs.
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

	registry := []*GraphKind{
		{
			ResourceType: resourceref.KindComposeModule,
			Kind:         "module",
			Load: staticGraphSources(
				&GraphSource{ID: 1001, Name: "Lead", Handle: "lead"},
				&GraphSource{ID: 1002, Name: "Opportunity", Handle: "opportunity"},
			),
		},
		{
			ResourceType: resourceref.KindNgAutomation,
			Kind:         "automation",
			Load: staticGraphSources(
				// what a real store-backed loader will produce
				&GraphSource{ID: au.ID, Name: "Lead Scoring", Refs: au.ResourceRefs()},
			),
		},
		{
			ResourceType: resourceref.KindConfiguredConnection,
			Kind:         "connection",
			LoadOne: func(_ context.Context, id uint64) (*GraphSource, error) {
				if id != 40 {
					return nil, nil
				}
				return &GraphSource{ID: 40, Name: "Google Sheets", Sensitivity: "confidential"}, nil
			},
		},
	}

	g, err := assembleProjectGraph(context.Background(), 0, registry)
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

func TestProjectGraphMockDataset(t *testing.T) {
	g, err := DefaultProjectGraph.Graph(context.Background(), 0)
	require.NoError(t, err)

	// 12 in-scope nodes + 2 external connections
	require.Len(t, g.Nodes, 14)
	require.Len(t, g.Edges, 11)
	require.Empty(t, g.Missing)

	// dynamic step argument surfaces as warning, not edge
	require.Equal(t, []*types.ProjectGraphWarning{
		{SourceID: 5001, Kind: resourceref.KindComposeModule, Reason: resourceref.ReasonStepArgument},
	}, g.Warnings)

	byID := make(map[uint64]*types.ProjectGraphNode)
	for _, n := range g.Nodes {
		byID[n.ID] = n
	}

	// connections resolved through LoadOne, flagged external
	require.True(t, byID[4001].External)
	require.Equal(t, "confidential", byID[4001].Sensitivity)
	require.True(t, byID[4002].External)

	// handle-based trigger constraint resolved to the Lead module
	require.Contains(t, g.Edges, &types.ProjectGraphEdge{SourceID: 5001, TargetID: 1001, Reason: resourceref.ReasonTriggerModule})
}
