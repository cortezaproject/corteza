package service

import (
	"context"
	"fmt"

	automationTypes "github.com/crusttech/human/server/automation/types"
	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/rbac"
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
	ctx, sources, err := s.fetch(ctx, projectID)
	if err != nil {
		return nil, err
	}

	sources = filterSources(sources)

	return s.transform(ctx, sources)
}

func (s *projectGraphService) fetch(ctx context.Context, projectID uint64) (context.Context, []*GraphSource, error) {
	proj, err := store.LookupProjectByID(ctx, s.store, projectID)
	if err != nil {
		return ctx, nil, err
	}

	ctx = scope.SetScopeToContext(ctx, scope.Scope{TenantID: proj.TenantID, ProjectID: projectID})

	sens, err := s.loadSensitivityLevels(ctx)
	if err != nil {
		return ctx, nil, err
	}

	sources, err := s.loadSources(ctx, sens)
	if err != nil {
		return ctx, nil, err
	}

	return ctx, sources, nil
}

// graphNodeKinds is the set of GraphSource.Kind values shown on the project
// graph. Sources of any other kind (knowledge-base, role, …) are dropped whole.
var graphNodeKinds = map[string]bool{
	"module":     true,
	"page":       true,
	"chart":      true,
	"connection": true,
	"automation": true,
	"agent":      true,
	"chatbot":    true,
}

// graphTargetKind maps a ref's target resource type to its displayed node kind.
// Refs whose target type is absent (workflow, llm-provider, knowledge-base,
// namespace, page-layout, template, role) are dropped.
var graphTargetKind = map[string]string{
	resourceref.KindComposeModule:        "module",
	resourceref.KindComposePage:          "page",
	resourceref.KindComposeChart:         "chart",
	resourceref.KindNgAutomation:         "automation",
	resourceref.KindAgent:                "agent",
	resourceref.KindChatbot:              "chatbot",
	resourceref.KindDalConnection:        "connection",
	resourceref.KindConfiguredConnection: "connection",
}

// graphEdgeKinds is the allowed set of unordered node-kind pairs, mirroring the
// resource-relationship matrix. Keys are sorted so direction does not matter.
var graphEdgeKinds = map[[2]string]bool{
	{"module", "module"}:         true,
	{"module", "page"}:           true,
	{"chart", "module"}:          true,
	{"automation", "module"}:     true,
	{"agent", "module"}:          true,
	{"chart", "page"}:            true,
	{"automation", "page"}:       true,
	{"agent", "page"}:            true,
	{"chatbot", "page"}:          true,
	{"automation", "connection"}: true,
	{"agent", "automation"}:      true,
	{"automation", "chatbot"}:    true,
	{"agent", "chatbot"}:         true,
}

func filterSources(sources []*GraphSource) []*GraphSource {
	out := sources[:0]
	for _, src := range sources {
		if !graphNodeKinds[src.Kind] {
			continue
		}
		src.Refs = filterRefs(src.Kind, src.Refs)
		out = append(out, src)
	}
	return out
}

func filterRefs(srcKind string, refs []resourceref.Ref) []resourceref.Ref {
	out := refs[:0]
	for _, ref := range refs {
		targetKind, ok := graphTargetKind[ref.Kind()]
		if !ok {
			continue
		}
		if graphEdgeKinds[graphEdgeKey(srcKind, targetKind)] {
			out = append(out, ref)
		}
	}
	return out
}

func graphEdgeKey(a, b string) [2]string {
	if a > b {
		a, b = b, a
	}
	return [2]string{a, b}
}

func (s *projectGraphService) transform(ctx context.Context, sources []*GraphSource) (*types.ProjectGraph, error) {
	return assembleProjectGraph(ctx, sources, s.loadExternal)
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

	// roles: their edges come from RBAC rules, not a ResourceRefs() extractor.
	// Roles are project-scoped so the scoped search returns this project's
	// roles; rules carry no project scope and the store filter can't narrow by
	// role, so load all rules once and group them by role in memory.
	rr, _, err := store.SearchRoles(ctx, s.store, types.RoleFilter{})
	if err != nil {
		return nil, err
	}
	if len(rr) > 0 {
		rules, _, err := store.SearchRbacRules(ctx, s.store, rbac.RuleFilter{})
		if err != nil {
			return nil, err
		}

		byRole := make(map[uint64]rbac.RuleSet, len(rr))
		for _, rule := range rules {
			byRole[rule.RoleID] = append(byRole[rule.RoleID], rule)
		}

		for _, r := range rr {
			refs := roleRuleRefs(byRole[r.ID])
			if len(refs) == 0 {
				// role grants nothing on graph resources; omit the node
				continue
			}
			out = append(out, &GraphSource{
				ResourceType: resourceref.KindRole,
				Kind:         "role",
				ID:           r.ID,
				Name:         firstNonEmpty(r.Name, r.Handle),
				Handle:       r.Handle,
				Refs:         refs,
			})
		}
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

	case resourceref.KindComposePageLayout:
		l, err := store.LookupComposePageLayoutByID(ctx, s.store, id)
		if err != nil {
			return nil, ignoreNotFound(err)
		}
		return &GraphSource{ResourceType: resourceType, Kind: "page-layout", ID: l.ID, Name: firstNonEmpty(l.Meta.Title, l.Handle), Handle: l.Handle}, nil
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
			a uint64
			b uint64
		}
	)

	g = &types.ProjectGraph{
		Nodes: []*types.ProjectGraphNode{},
		Edges: []*types.ProjectGraphEdge{},
	}

	var (
		nodes       = make(map[nodeKey]bool)
		nodesByKind = make(map[string][]uint64)
		handles     = make(map[string]map[string]uint64)
	)

	addNode := func(src *GraphSource, external bool) {
		key := nodeKey{src.ResourceType, src.ID}
		if nodes[key] {
			return
		}
		nodes[key] = true
		nodesByKind[src.ResourceType] = append(nodesByKind[src.ResourceType], src.ID)

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
	addEdge := func(sourceID, targetID uint64, reason string) {
		// normalize to unordered pair so A→B, B→A and duplicate reasons collapse to one edge
		ek := edgeKey{sourceID, targetID}
		if ek.a > ek.b {
			ek.a, ek.b = ek.b, ek.a
		}
		if seen[ek] {
			return
		}
		seen[ek] = true

		g.Edges = append(g.Edges, &types.ProjectGraphEdge{
			SourceID: sourceID,
			TargetID: targetID,
			Reason:   reason,
		})
	}

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

			// wildcard ref (RBAC role grant on a whole kind): fan out to every
			// in-scope node of the kind, skipping the source's self-edge.
			if ref.Wildcard {
				for _, targetID := range nodesByKind[ref.Kind()] {
					if targetID == src.ID && ref.Kind() == src.ResourceType {
						continue
					}
					addEdge(src.ID, targetID, ref.Reason)
				}
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

			addEdge(src.ID, targetID, ref.Reason)
		}
	}

	return
}

// graphResourceKinds is the set of rbac/resourceref kinds that correspond to
// project-graph nodes. Role RBAC rules on any other resource (records, grants,
// settings, users, …) are not graph edges and are ignored.
var graphResourceKinds = map[string]bool{
	resourceref.KindComposeNamespace:   true,
	resourceref.KindComposeModule:      true,
	resourceref.KindComposeChart:       true,
	resourceref.KindComposePage:        true,
	resourceref.KindComposePageLayout:  true,
	resourceref.KindAutomationWorkflow: true,
	resourceref.KindNgAutomation:       true,
	resourceref.KindDalConnection:      true,
	resourceref.KindLlmProvider:        true,
	resourceref.KindKnowledgeBase:      true,
	resourceref.KindAgent:              true,
	resourceref.KindChatbot:            true,
	resourceref.KindRole:               true,
	resourceref.KindTemplate:           true,
}

// roleRuleRefs converts a role's RBAC rules into graph refs. A rule on a
// specific resource yields a direct ref; a rule on a whole kind (trailing
// wildcard, e.g. "corteza::compose:module/*/*") yields a wildcard ref the
// assembler fans out to every in-scope node of the kind. Rules on non-graph
// resources are skipped; allow and deny both count as a relationship. Refs are
// deduplicated so a role's many per-operation rules collapse to one edge.
func roleRuleRefs(rules rbac.RuleSet) []resourceref.Ref {
	var (
		out  = make([]resourceref.Ref, 0, len(rules))
		seen = make(map[string]bool)
	)

	for _, rule := range rules {
		kind := rbac.ResourceType(rule.Resource)
		if !graphResourceKinds[kind] {
			continue
		}

		var ref resourceref.Ref
		if id := rbac.ResourceID(rule.Resource); id > 0 {
			ref = resourceref.Make(kind, id, resourceref.ReasonRoleRbac)
		} else {
			ref = resourceref.MakeWildcard(kind, resourceref.ReasonRoleRbac)
		}
		if ref.IsEmpty() {
			continue
		}

		key := ref.Resource
		if ref.Wildcard {
			key = "*" + ref.Resource
		}
		if seen[key] {
			continue
		}
		seen[key] = true

		out = append(out, ref)
	}

	return out
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
