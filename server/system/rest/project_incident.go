package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ProjectIncident struct {
		projectIncident projectIncidentService
		ac              projectIncidentAccessController
	}

	projectIncidentPayload struct {
		*types.ProjectIncident

		CanGrant                 bool `json:"canGrant"`
		CanUpdateProjectIncident bool `json:"canUpdateProjectIncident"`
		CanDeleteProjectIncident bool `json:"canDeleteProjectIncident"`
	}

	projectIncidentSetPayload struct {
		Filter types.ProjectIncidentFilter `json:"filter"`
		Set    []*projectIncidentPayload   `json:"set"`
	}

	projectIncidentService interface {
		FindByID(ctx context.Context, ID uint64) (a *types.ProjectIncident, err error)
		Create(ctx context.Context, new *types.ProjectIncident) (a *types.ProjectIncident, err error)
		Update(ctx context.Context, upd *types.ProjectIncident) (a *types.ProjectIncident, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.ProjectIncidentFilter) (set types.ProjectIncidentSet, f types.ProjectIncidentFilter, err error)
	}

	projectIncidentAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateProjectIncident(context.Context, *types.ProjectIncident) bool
		CanDeleteProjectIncident(context.Context, *types.ProjectIncident) bool
	}
)

func (ProjectIncident) New() *ProjectIncident {
	return &ProjectIncident{
		projectIncident: service.DefaultProjectIncident,
		ac:              service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *ProjectIncident) makeFilter(ctx context.Context, r *request.ProjectIncidentList) (types.ProjectIncidentFilter, error) {
	var (
		err error
		f   = types.ProjectIncidentFilter{
			Query:     r.Query,
			ProjectID: r.ProjectID,
			Status:    r.Status,
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
func (ctrl *ProjectIncident) beforeCreate(ctx context.Context, res *types.ProjectIncident, r *request.ProjectIncidentCreate) error {
	res.ProjectID = r.ProjectID
	return nil
}

// beforeUpdate is a no-op; the generated controller mapped all value params.
func (ctrl *ProjectIncident) beforeUpdate(ctx context.Context, res *types.ProjectIncident, r *request.ProjectIncidentUpdate) error {
	return nil
}

func (ctrl *ProjectIncident) makePayload(ctx context.Context, a *types.ProjectIncident, err error) (*projectIncidentPayload, error) {
	if err != nil || a == nil {
		return nil, err
	}

	return &projectIncidentPayload{
		ProjectIncident: a,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateProjectIncident: ctrl.ac.CanUpdateProjectIncident(ctx, a),
		CanDeleteProjectIncident: ctrl.ac.CanDeleteProjectIncident(ctx, a),
	}, nil
}

func (ctrl *ProjectIncident) makeFilterPayload(ctx context.Context, nn types.ProjectIncidentSet, f types.ProjectIncidentFilter, err error) (*projectIncidentSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectIncidentSetPayload{Filter: f, Set: make([]*projectIncidentPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
