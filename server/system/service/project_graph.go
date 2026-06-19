package service

import (
	"context"
	"fmt"

	automationTypes "github.com/crusttech/human/server/automation/types"
	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	GraphSource struct {
		ResourceType string
		Kind         string

		ID          uint64
		Name        string
		Handle      string
		Sensitivity string

		Refs []resourceref.Ref
	}

	externalLoader func(ctx context.Context, resourceType string, id uint64) (*GraphSource, error)

	projectGraphService struct {
		store store.Storer
	}
)

var DefaultProjectGraph = &projectGraphService{store: DefaultStore}

func (s *projectGraphService) Graph(ctx context.Context, projectID uint64) (*types.ProjectGraph, error) {
	g, err := s.build(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return filterProjectGraph(g), nil
}

func (s *projectGraphService) build(ctx context.Context, projectID uint64) (*types.ProjectGraph, error) {
	proj, err := store.LookupProjectByID(ctx, s.store, projectID)
	if err != nil {
		return nil, err
	}

	ctx = scope.SetScopeToContext(ctx, scope.Scope{TenantID: proj.TenantID, ProjectID: projectID})

	sens, err := s.loadSensitivityLevels(ctx)
	if err != nil {
		return nil, err
	}

	sources, err := s.loadSources(ctx, sens)
	if err != nil {
		return nil, err
	}

	return assembleProjectGraph(ctx, sources, s.loadExternal)
}

func filterProjectGraph(g *types.ProjectGraph) *types.ProjectGraph {
	return g
}

func (s *projectGraphService) loadSources(ctx context.Context, sens map[uint64]string) ([]*GraphSource, error) {
	var out []*GraphSource

	mm, _, err := store.SearchComposeModules(ctx, s.store, composeTypes.ModuleFilter{})
	if err != nil {
		return nil, err
	}
	if len(mm) > 0 {
		ff, _, err := store.SearchComposeModuleFields(ctx, s.store, composeTypes.ModuleFieldFilter{ModuleID: mm.IDs()})
		if err != nil {
			return nil, err
		}
		for _, m := range mm {
			m.Fields = ff.FilterByModule(m.ID)
		}
	}
	for _, m := range mm {
		out = append(out, &GraphSource{
			ResourceType: resourceref.KindComposeModule,
			Kind:         "module",
			ID:           m.ID,
			Name:         m.Name,
			Handle:       m.Handle,
			Sensitivity:  sens[m.Config.Privacy.SensitivityLevelID],
			Refs:         m.ResourceRefs(),
		})
	}

	pp, _, err := store.SearchComposePages(ctx, s.store, composeTypes.PageFilter{})
	if err != nil {
		return nil, err
	}
	for _, p := range pp {
		out = append(out, &GraphSource{
			ResourceType: resourceref.KindComposePage,
			Kind:         "page",
			ID:           p.ID,
			Name:         firstNonEmpty(p.Title, p.Handle),
			Handle:       p.Handle,
			Refs:         p.ResourceRefs(),
		})
	}

	cc, _, err := store.SearchComposeCharts(ctx, s.store, composeTypes.ChartFilter{})
	if err != nil {
		return nil, err
	}
	for _, c := range cc {
		out = append(out, &GraphSource{
			ResourceType: resourceref.KindComposeChart,
			Kind:         "chart",
			ID:           c.ID,
			Name:         firstNonEmpty(c.Name, c.Handle),
			Handle:       c.Handle,
			Refs:         c.ResourceRefs(),
		})
	}

	aa, _, err := store.SearchAutomationNgAutomations(ctx, s.store, automationTypes.NgAutomationFilter{})
	if err != nil {
		return nil, err
	}
	for _, a := range aa {
		name := a.Handle
		if a.Meta != nil && a.Meta.Short != "" {
			name = a.Meta.Short
		}
		out = append(out, &GraphSource{
			ResourceType: resourceref.KindNgAutomation,
			Kind:         "automation",
			ID:           a.ID,
			Name:         firstNonEmpty(name, a.Handle),
			Handle:       a.Handle,
			Refs:         a.ResourceRefs(),
		})
	}

	ag, _, err := store.SearchAgents(ctx, s.store, types.AgentFilter{})
	if err != nil {
		return nil, err
	}
	for _, a := range ag {
		out = append(out, &GraphSource{
			ResourceType: resourceref.KindAgent,
			Kind:         "agent",
			ID:           a.ID,
			Name:         firstNonEmpty(a.Meta.Short, a.Handle),
			Handle:       a.Handle,
			Refs:         a.ResourceRefs(),
		})
	}

	cb, _, err := store.SearchChatbots(ctx, s.store, types.ChatbotFilter{})
	if err != nil {
		return nil, err
	}
	for _, c := range cb {
		out = append(out, &GraphSource{
			ResourceType: types.ChatbotResourceType,
			Kind:         "chatbot",
			ID:           c.ID,
			Name:         firstNonEmpty(c.Name, c.Handle),
			Handle:       c.Handle,
			Refs:         c.ResourceRefs(),
		})
	}

	kk, _, err := store.SearchKnowledgeBases(ctx, s.store, types.KnowledgeBaseFilter{})
	if err != nil {
		return nil, err
	}
	for _, k := range kk {
		out = append(out, &GraphSource{
			ResourceType: resourceref.KindKnowledgeBase,
			Kind:         "knowledge-base",
			ID:           k.ID,
			Name:         firstNonEmpty(k.Title, k.Handle),
			Handle:       k.Handle,
			Refs:         k.ResourceRefs(),
		})
	}

	return out, nil
}

func (s *projectGraphService) loadExternal(ctx context.Context, resourceType string, id uint64) (*GraphSource, error) {
	switch resourceType {
	case resourceref.KindDalConnection:
		c, err := store.LookupConnectionByID(ctx, s.store, id)
		if err != nil {
			return nil, ignoreNotFound(err)
		}
		return &GraphSource{ResourceType: resourceType, Kind: "connection", ID: c.ID, Name: firstNonEmpty(c.Meta.Short, c.Handle), Handle: c.Handle}, nil

	case resourceref.KindLlmProvider:
		p, err := store.LookupLlmProviderByID(ctx, s.store, id)
		if err != nil {
			return nil, ignoreNotFound(err)
		}
		return &GraphSource{ResourceType: resourceType, Kind: "llm-provider", ID: p.ID, Name: p.Handle, Handle: p.Handle}, nil

	case resourceref.KindAutomationWorkflow:
		w, err := store.LookupAutomationWorkflowByID(ctx, s.store, id)
		if err != nil {
			return nil, ignoreNotFound(err)
		}
		return &GraphSource{ResourceType: resourceType, Kind: "workflow", ID: w.ID, Name: firstNonEmpty(w.Meta.Name, w.Handle), Handle: w.Handle}, nil

	case resourceref.KindRole:
		r, err := store.LookupRoleByID(ctx, s.store, id)
		if err != nil {
			return nil, ignoreNotFound(err)
		}
		return &GraphSource{ResourceType: resourceType, Kind: "role", ID: r.ID, Name: firstNonEmpty(r.Name, r.Handle), Handle: r.Handle}, nil

	case resourceref.KindTemplate:
		t, err := store.LookupTemplateByID(ctx, s.store, id)
		if err != nil {
			return nil, ignoreNotFound(err)
		}
		return &GraphSource{ResourceType: resourceType, Kind: "template", ID: t.ID, Name: t.Handle, Handle: t.Handle}, nil

	case resourceref.KindComposeNamespace:
		ns, err := store.LookupComposeNamespaceByID(ctx, s.store, id)
		if err != nil {
			return nil, ignoreNotFound(err)
		}
		return &GraphSource{ResourceType: resourceType, Kind: "namespace", ID: ns.ID, Name: firstNonEmpty(ns.Name, ns.Slug), Handle: ns.Slug}, nil
	}

	return nil, nil
}

func assembleProjectGraph(ctx context.Context, sources []*GraphSource, loadExternal externalLoader) (g *types.ProjectGraph, err error) {
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
	)

	g = &types.ProjectGraph{
		Nodes: []*types.ProjectGraphNode{},
		Edges: []*types.ProjectGraphEdge{},
	}

	var (
		nodes   = make(map[nodeKey]bool)
		handles = make(map[string]map[string]uint64)
	)

	addNode := func(src *GraphSource, external bool) {
		key := nodeKey{src.ResourceType, src.ID}
		if nodes[key] {
			return
		}
		nodes[key] = true

		g.Nodes = append(g.Nodes, &types.ProjectGraphNode{
			ID:          src.ID,
			Kind:        src.Kind,
			Name:        src.Name,
			Sensitivity: src.Sensitivity,
			External:    external,
		})

		if src.Handle != "" {
			if handles[src.ResourceType] == nil {
				handles[src.ResourceType] = make(map[string]uint64)
			}
			handles[src.ResourceType][src.Handle] = src.ID
		}
	}

	for _, src := range sources {
		addNode(src, false)
	}

	seen := make(map[edgeKey]bool)
	for _, src := range sources {
		for _, ref := range src.Refs {
			if ref.Unresolved {
				g.Warnings = append(g.Warnings, &types.ProjectGraphWarning{
					SourceID: src.ID,
					Kind:     ref.Kind(),
					Reason:   ref.Reason,
				})
				continue
			}

			miss := func() {
				g.Missing = append(g.Missing, &types.ProjectGraphMissingRef{
					SourceID:    src.ID,
					Kind:        ref.Kind(),
					TargetID:    ref.ID(),
					TargetIdent: ref.Label,
					Reason:      ref.Reason,
				})
			}

			targetID := ref.ID()
			if targetID == 0 {
				targetID = handles[ref.Kind()][ref.Label]
			}
			if targetID == 0 {
				miss()
				continue
			}

			if !nodes[nodeKey{ref.Kind(), targetID}] {
				ext, err := loadExternal(ctx, ref.Kind(), targetID)
				if err != nil {
					return nil, fmt.Errorf("load external %s %d: %w", ref.Kind(), targetID, err)
				}
				if ext == nil {
					miss()
					continue
				}

				addNode(ext, true)
			}

			ek := edgeKey{src.ID, targetID, ref.Reason}
			if seen[ek] {
				continue
			}
			seen[ek] = true

			g.Edges = append(g.Edges, &types.ProjectGraphEdge{
				SourceID: src.ID,
				TargetID: targetID,
				Reason:   ref.Reason,
			})
		}
	}

	return
}

func (s *projectGraphService) loadSensitivityLevels(ctx context.Context) (map[uint64]string, error) {
	ll, _, err := store.SearchDalSensitivityLevels(ctx, s.store, types.DalSensitivityLevelFilter{})
	if err != nil {
		return nil, err
	}

	out := make(map[uint64]string, len(ll))
	for _, l := range ll {
		out[l.ID] = l.Handle
	}
	return out, nil
}

func ignoreNotFound(err error) error {
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
