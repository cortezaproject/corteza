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

	TempStruct struct {
		Filter types.ExposedModuleFilter `json:"filter"`
		Set    types.ExposedModuleSet    `json:"set"`
	}
)

func (SyncStructure) New() *SyncStructure {
	return &SyncStructure{}
}

func (ctrl SyncStructure) ReadExposedAll(ctx context.Context, r *request.SyncStructureReadExposedAll) (interface{}, error) {
	// TODO - fixed values for now

	var (
		err error
		f   = types.ExposedModuleFilter{
			NodeID: 276342359342989444,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	list, f, err := (service.ExposedModule()).Find(context.Background(), f)

	return TempStruct{
		Set:    list,
		Filter: f,
	}, err
}
