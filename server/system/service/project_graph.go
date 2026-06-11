package service

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/system/types"
)

type (
	// GraphSource is one resource as seen by the project graph assembler
	GraphSource struct {
		ID          uint64
		Name        string
		Handle      string
		Sensitivity string

		// Refs as returned by the resource's ResourceRefs() extractor
		Refs []resourceref.Ref
	}

	// GraphKind describes one resource kind participating in the project graph
	GraphKind struct {
		// ResourceType is matched against resourceref.Ref.Kind when resolving
		// edge targets (see pkg/resourceref kind constants)
		ResourceType string

		// Kind of the produced graph nodes
		Kind string

		// Load enumerates all sources of this kind within the project scope;
		// nil for kinds that only appear as edge targets (tenant-level resources)
		Load func(ctx context.Context, projectID uint64) ([]*GraphSource, error)

		// LoadOne fetches a single source referenced from the project but
		// living outside its scope; resulting nodes are flagged as external
		LoadOne func(ctx context.Context, id uint64) (*GraphSource, error)
	}

	projectGraphService struct {
		registry []*GraphKind
	}
)

// DefaultProjectGraph assembles project graphs from a mock registry until
// resource scoping (project_id columns on compose resources) is in place;
// swap in store-backed loaders to go live.
var DefaultProjectGraph = &projectGraphService{registry: mockProjectGraphRegistry()}

func (s *projectGraphService) Graph(ctx context.Context, projectID uint64) (*types.ProjectGraph, error) {
	return assembleProjectGraph(ctx, projectID, s.registry)
}

// assembleProjectGraph builds a dependency graph from the given kind registry
//
// Phases:
//  1. enumerate in-scope sources per kind; every source becomes a node and is
//     indexed by ID and handle
//  2. resolve each source's refs: ID refs directly, ident refs through the
//     handle index; targets outside the scope are fetched via LoadOne and
//     flagged external; unresolvable refs are reported as missing
//
// Edges are deduplicated on (source, target, reason); output order follows
// registry and source order, so results are deterministic.
func assembleProjectGraph(ctx context.Context, projectID uint64, registry []*GraphKind) (g *types.ProjectGraph, err error) {
	type (
		nodeKey struct {
			resourceType string
			id           uint64
		}
		edgeKey struct {
			sourceID uint64
			targetID uint64
			reason   string
		}
		loadedSource struct {
			kind *GraphKind
			src  *GraphSource
		}
	)

	g = &types.ProjectGraph{
		Nodes: []*types.ProjectGraphNode{},
		Edges: []*types.ProjectGraphEdge{},
	}

	var (
		kinds   = make(map[string]*GraphKind, len(registry))
		nodes   = make(map[nodeKey]bool)
		handles = make(map[string]map[string]uint64)
		sources []loadedSource
	)

	addNode := func(k *GraphKind, s *GraphSource, external bool) {
		key := nodeKey{k.ResourceType, s.ID}
		if nodes[key] {
			return
		}
		nodes[key] = true

		g.Nodes = append(g.Nodes, &types.ProjectGraphNode{
			ID:          s.ID,
			Kind:        k.Kind,
			Name:        s.Name,
			Sensitivity: s.Sensitivity,
			External:    external,
		})

		if s.Handle != "" {
			if handles[k.ResourceType] == nil {
				handles[k.ResourceType] = make(map[string]uint64)
			}
			handles[k.ResourceType][s.Handle] = s.ID
		}
	}

	for _, k := range registry {
		kinds[k.ResourceType] = k

		if k.Load == nil {
			continue
		}

		ss, err := k.Load(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("load %s graph sources: %w", k.Kind, err)
		}

		for _, s := range ss {
			addNode(k, s, false)
			sources = append(sources, loadedSource{kind: k, src: s})
		}
	}

	seen := make(map[edgeKey]bool)
	for _, ls := range sources {
		for _, ref := range ls.src.Refs {
			if ref.Dynamic {
				g.Warnings = append(g.Warnings, &types.ProjectGraphWarning{
					SourceID: ls.src.ID,
					Kind:     ref.Kind,
					Reason:   ref.Reason,
					Path:     ref.Path,
				})
				continue
			}

			miss := func() {
				g.Missing = append(g.Missing, &types.ProjectGraphMissingRef{
					SourceID:    ls.src.ID,
					Kind:        ref.Kind,
					TargetID:    ref.ID,
					TargetIdent: ref.Ident,
					Reason:      ref.Reason,
				})
			}

			targetID := ref.ID
			if targetID == 0 {
				targetID = handles[ref.Kind][ref.Ident]
			}
			if targetID == 0 {
				miss()
				continue
			}

			if !nodes[nodeKey{ref.Kind, targetID}] {
				k := kinds[ref.Kind]
				if k == nil || k.LoadOne == nil {
					miss()
					continue
				}

				s, err := k.LoadOne(ctx, targetID)
				if err != nil {
					return nil, fmt.Errorf("load external %s %d: %w", k.Kind, targetID, err)
				}
				if s == nil {
					miss()
					continue
				}

				addNode(k, s, true)
			}

			ek := edgeKey{ls.src.ID, targetID, ref.Reason}
			if seen[ek] {
				continue
			}
			seen[ek] = true

			g.Edges = append(g.Edges, &types.ProjectGraphEdge{
				SourceID: ls.src.ID,
				TargetID: targetID,
				Reason:   ref.Reason,
			})
		}
	}

	return
}

// mockProjectGraphRegistry serves a static dataset through the real assembler
// until store-backed loaders land (blocked on compose project scoping)
func mockProjectGraphRegistry() []*GraphKind {
	connections := map[uint64]*GraphSource{
		4001: {ID: 4001, Name: "Google Sheets", Sensitivity: "confidential"},
		4002: {ID: 4002, Name: "Stripe", Sensitivity: "restricted"},
	}

	return []*GraphKind{
		{
			ResourceType: resourceref.KindComposeModule,
			Kind:         "module",
			Load: staticGraphSources(
				&GraphSource{ID: 1001, Name: "Lead", Handle: "lead", Sensitivity: "internal", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindDalConnection, 4001, resourceref.ReasonModuleConnection, "Config.DAL.ConnectionID"),
				}},
				&GraphSource{ID: 1002, Name: "Opportunity", Handle: "opportunity", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonModuleFieldRef, "Fields.lead.Options.ModuleID"),
				}},
				&GraphSource{ID: 1003, Name: "Quote", Handle: "quote", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonModuleFieldRef, "Fields.opportunity.Options.ModuleID"),
					resourceref.Make(resourceref.KindDalConnection, 4002, resourceref.ReasonModuleConnection, "Config.DAL.ConnectionID"),
				}},
			),
		},
		{
			ResourceType: resourceref.KindComposePage,
			Kind:         "page",
			Load: staticGraphSources(
				&GraphSource{ID: 2001, Name: "Lead List", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonPageModule, "ModuleID"),
				}},
				&GraphSource{ID: 2002, Name: "Opportunity Board", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonPageModule, "ModuleID"),
				}},
			),
		},
		{
			ResourceType: resourceref.KindComposeChart,
			Kind:         "chart",
			Load: staticGraphSources(
				&GraphSource{ID: 3001, Name: "Pipeline Forecast", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonChartModule, "Config.Reports.0.ModuleID"),
				}},
			),
		},
		{
			ResourceType: resourceref.KindDalConnection,
			Kind:         "connection",
			// tenant-level; appears only as an external edge target
			LoadOne: func(_ context.Context, id uint64) (*GraphSource, error) {
				return connections[id], nil
			},
		},
		{
			ResourceType: resourceref.KindNgAutomation,
			Kind:         "automation",
			Load: staticGraphSources(
				// handle-based ref exercises ident resolution (trigger constraints
				// reference modules by handle); dynamic ref exercises the
				// runtime-resolved-argument warning
				&GraphSource{ID: 5001, Name: "Lead Scoring", Refs: []resourceref.Ref{
					resourceref.MakeIdent(resourceref.KindComposeModule, "lead", resourceref.ReasonTriggerModule, "Triggers.0.Constraints.0.Values.0"),
					resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonStepArgument, "Steps.2.Arguments.0"),
					resourceref.MakeDynamic(resourceref.KindComposeModule, resourceref.ReasonStepArgument, "Steps.4.Arguments.0"),
				}},
			),
		},
		{
			ResourceType: resourceref.KindAgent,
			Kind:         "agent",
			Load: staticGraphSources(
				&GraphSource{ID: 6001, Name: "Sales Assistant", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonAgentModule, "Access.Tools.0.Allow.0.ModuleIDs.0"),
				}},
			),
		},
		{
			ResourceType: types.ChatbotResourceType,
			Kind:         "chatbot",
			Load: staticGraphSources(
				&GraphSource{ID: 7001, Name: "Support Bot", Refs: []resourceref.Ref{
					resourceref.Make(resourceref.KindAgent, 6001, resourceref.ReasonChatbotAgent, "Scenarios.0.AgentID"),
				}},
			),
		},
		{
			ResourceType: types.RoleResourceType,
			Kind:         "role",
			// RBAC-derived edges deferred; roles render as standalone nodes
			Load: staticGraphSources(
				&GraphSource{ID: 8001, Name: "Sales Rep"},
				&GraphSource{ID: 8002, Name: "Sales Manager"},
				&GraphSource{ID: 8003, Name: "Customer"},
			),
		},
	}
}

func staticGraphSources(ss ...*GraphSource) func(context.Context, uint64) ([]*GraphSource, error) {
	return func(context.Context, uint64) ([]*GraphSource, error) { return ss, nil }
}
