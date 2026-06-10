package service

import (
	"context"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type projectGraphService struct {
	store store.Storer
}

var DefaultProjectGraph = &projectGraphService{}

func ProjectGraph() *projectGraphService {
	return &projectGraphService{store: DefaultStore}
}

// Graph derives the project resource graph from the real resources in the
// project's compose namespace. Modules are the only project-scoped resource
// kind so far; other kinds (pages, connections, automations, agents, chatbots)
// join as they gain project scoping.
//
// Edges:
//   - module-field-ref: a Record field on one module referencing another
func (s *projectGraphService) Graph(ctx context.Context, projectID uint64) (*types.ProjectGraph, error) {
	str := s.store
	if str == nil {
		str = DefaultStore
	}

	p, err := loadProject(ctx, str, projectID)
	if err != nil {
		return nil, err
	}

	g := &types.ProjectGraph{
		Nodes: []*types.ProjectGraphNode{},
		Edges: []*types.ProjectGraphEdge{},
	}

	if p.Config.NamespaceID == 0 {
		return g, nil
	}

	mm, _, err := store.SearchComposeModules(ctx, str, composeTypes.ModuleFilter{
		NamespaceID: p.Config.NamespaceID,
	})
	if err != nil {
		return nil, err
	}
	if len(mm) == 0 {
		return g, nil
	}

	var (
		moduleIDs = make([]uint64, len(mm))
		inProject = make(map[uint64]bool, len(mm))
	)
	for i, m := range mm {
		moduleIDs[i] = m.ID
		inProject[m.ID] = true
		g.Nodes = append(g.Nodes, &types.ProjectGraphNode{
			ID:          m.ID,
			Kind:        "module",
			Name:        m.Name,
			Sensitivity: s.sensitivityHandle(ctx, str, m.Config.Privacy.SensitivityLevelID),
		})
	}

	ff, _, err := store.SearchComposeModuleFields(ctx, str, composeTypes.ModuleFieldFilter{
		ModuleID: moduleIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, f := range ff {
		if f.Kind != "Record" {
			continue
		}
		target := f.Options.UInt64("moduleID")
		if target == 0 || target == f.ModuleID || !inProject[target] {
			continue
		}
		g.Edges = append(g.Edges, &types.ProjectGraphEdge{
			SourceID: f.ModuleID,
			TargetID: target,
			Reason:   "module-field-ref",
		})
	}

	return g, nil
}

// sensitivityHandle resolves a sensitivity level ID to its handle; empty when
// unset or unresolvable.
func (s *projectGraphService) sensitivityHandle(ctx context.Context, str store.Storer, id uint64) string {
	if id == 0 {
		return ""
	}
	lvl, err := store.LookupDalSensitivityLevelByID(ctx, str, id)
	if err != nil || lvl == nil {
		return ""
	}
	return lvl.Handle
}
