package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ProjectTask struct {
		projectTask projectTaskService
		ac          projectTaskAccessController
	}

	projectTaskPayload struct {
		*types.ProjectTask

		CanGrant             bool `json:"canGrant"`
		CanUpdateProjectTask bool `json:"canUpdateProjectTask"`
		CanDeleteProjectTask bool `json:"canDeleteProjectTask"`
	}

	projectTaskSetPayload struct {
		Filter types.ProjectTaskFilter `json:"filter"`
		Set    []*projectTaskPayload   `json:"set"`
	}

	projectTaskService interface {
		FindByID(ctx context.Context, ID uint64) (a *types.ProjectTask, err error)
		Create(ctx context.Context, new *types.ProjectTask) (a *types.ProjectTask, err error)
		Update(ctx context.Context, upd *types.ProjectTask) (a *types.ProjectTask, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.ProjectTaskFilter) (set types.ProjectTaskSet, f types.ProjectTaskFilter, err error)
	}

	projectTaskAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateProjectTask(context.Context, *types.ProjectTask) bool
		CanDeleteProjectTask(context.Context, *types.ProjectTask) bool
	}
)

func (ProjectTask) New() *ProjectTask {
	return &ProjectTask{
		projectTask: service.DefaultProjectTask,
		ac:          service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *ProjectTask) makeFilter(ctx context.Context, r *request.ProjectTaskList) (types.ProjectTaskFilter, error) {
	var (
		err error
		f   = types.ProjectTaskFilter{
			Query:      r.Query,
			ProjectID:  r.ProjectID,
			RevisionID: r.RevisionID,
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
func (ctrl *ProjectTask) beforeCreate(ctx context.Context, res *types.ProjectTask, r *request.ProjectTaskCreate) error {
	res.ProjectID = r.ProjectID
	return nil
}

// beforeUpdate is a no-op; the generated controller mapped all value params.
func (ctrl *ProjectTask) beforeUpdate(ctx context.Context, res *types.ProjectTask, r *request.ProjectTaskUpdate) error {
	return nil
}

func (ctrl *ProjectTask) makePayload(ctx context.Context, a *types.ProjectTask, err error) (*projectTaskPayload, error) {
	if err != nil || a == nil {
		return nil, err
	}

	return &projectTaskPayload{
		ProjectTask: a,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateProjectTask: ctrl.ac.CanUpdateProjectTask(ctx, a),
		CanDeleteProjectTask: ctrl.ac.CanDeleteProjectTask(ctx, a),
	}, nil
}

func (ctrl *ProjectTask) makeFilterPayload(ctx context.Context, nn types.ProjectTaskSet, f types.ProjectTaskFilter, err error) (*projectTaskSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectTaskSetPayload{Filter: f, Set: make([]*projectTaskPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
