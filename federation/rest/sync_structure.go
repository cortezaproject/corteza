package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	SyncStructure struct{}

	listResponse struct {
		Filter *types.ExposedModuleFilter `json:"filter"`
		Set    *types.ExposedModuleSet    `json:"set"`
	}
)

func (SyncStructure) New() *SyncStructure {
	return &SyncStructure{}
}

func (ctrl SyncStructure) ReadExposedAll(ctx context.Context, r *request.SyncStructureReadExposedAll) (interface{}, error) {
	var (
		err  error
		node *types.Node
	)

	if node, err = service.DefaultNode.FindBySharedNodeID(ctx, r.NodeID); err != nil {
		return nil, err
	}

	f := types.ExposedModuleFilter{
		NodeID: node.ID,
	}

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	list, f, err := (service.ExposedModule()).Find(context.Background(), f)

	return listResponse{
		Set:    &list,
		Filter: &f,
	}, nil
}
