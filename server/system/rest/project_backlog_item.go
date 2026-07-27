package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ProjectBacklogItem struct {
		projectBacklogItem projectBacklogItemService
		ac                 projectBacklogItemAccessController
	}

	projectBacklogItemPayload struct {
		*types.ProjectBacklogItem

		CanGrant                    bool `json:"canGrant"`
		CanUpdateProjectBacklogItem bool `json:"canUpdateProjectBacklogItem"`
		CanDeleteProjectBacklogItem bool `json:"canDeleteProjectBacklogItem"`
	}

	projectBacklogItemSetPayload struct {
		Filter types.ProjectBacklogItemFilter `json:"filter"`
		Set    []*projectBacklogItemPayload   `json:"set"`
	}

	projectBacklogItemService interface {
		FindByID(ctx context.Context, ID uint64) (a *types.ProjectBacklogItem, err error)
		Create(ctx context.Context, new *types.ProjectBacklogItem) (a *types.ProjectBacklogItem, err error)
		Update(ctx context.Context, upd *types.ProjectBacklogItem) (a *types.ProjectBacklogItem, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.ProjectBacklogItemFilter) (set types.ProjectBacklogItemSet, f types.ProjectBacklogItemFilter, err error)
	}

	projectBacklogItemAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateProjectBacklogItem(context.Context, *types.ProjectBacklogItem) bool
		CanDeleteProjectBacklogItem(context.Context, *types.ProjectBacklogItem) bool
	}
)

func (ProjectBacklogItem) New() *ProjectBacklogItem {
	return &ProjectBacklogItem{
		projectBacklogItem: service.DefaultProjectBacklogItem,
		ac:                 service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *ProjectBacklogItem) makeFilter(ctx context.Context, r *request.ProjectBacklogItemList) (types.ProjectBacklogItemFilter, error) {
	var (
		err error
		f   = types.ProjectBacklogItemFilter{
			Query:      r.Query,
			ProjectID:  r.ProjectID,
			RevisionID: r.RevisionID,
			EventID:    r.EventID,
			Category:   r.Category,
			Status:     r.Status,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate stamps the owning project; the generated controller already
// mapped the plain-value params.
func (ctrl *ProjectBacklogItem) beforeCreate(ctx context.Context, res *types.ProjectBacklogItem, r *request.ProjectBacklogItemCreate) error {
	res.ProjectID = r.ProjectID
	return nil
}

// beforeUpdate is a no-op; the generated controller mapped all value params.
func (ctrl *ProjectBacklogItem) beforeUpdate(ctx context.Context, res *types.ProjectBacklogItem, r *request.ProjectBacklogItemUpdate) error {
	return nil
}

func (ctrl *ProjectBacklogItem) makePayload(ctx context.Context, a *types.ProjectBacklogItem, err error) (*projectBacklogItemPayload, error) {
	if err != nil || a == nil {
		return nil, err
	}

	return &projectBacklogItemPayload{
		ProjectBacklogItem: a,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateProjectBacklogItem: ctrl.ac.CanUpdateProjectBacklogItem(ctx, a),
		CanDeleteProjectBacklogItem: ctrl.ac.CanDeleteProjectBacklogItem(ctx, a),
	}, nil
}

func (ctrl *ProjectBacklogItem) makeFilterPayload(ctx context.Context, nn types.ProjectBacklogItemSet, f types.ProjectBacklogItemFilter, err error) (*projectBacklogItemSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectBacklogItemSetPayload{Filter: f, Set: make([]*projectBacklogItemPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
