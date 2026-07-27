package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ProjectFeature struct {
		projectFeature projectFeatureService
		ac             projectFeatureAccessController
	}

	projectFeaturePayload struct {
		*types.ProjectFeature

		CanGrant                bool `json:"canGrant"`
		CanUpdateProjectFeature bool `json:"canUpdateProjectFeature"`
		CanDeleteProjectFeature bool `json:"canDeleteProjectFeature"`
	}

	projectFeatureSetPayload struct {
		Filter types.ProjectFeatureFilter `json:"filter"`
		Set    []*projectFeaturePayload   `json:"set"`
	}

	projectFeatureService interface {
		FindByID(ctx context.Context, ID uint64) (a *types.ProjectFeature, err error)
		Create(ctx context.Context, new *types.ProjectFeature) (a *types.ProjectFeature, err error)
		Update(ctx context.Context, upd *types.ProjectFeature) (a *types.ProjectFeature, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.ProjectFeatureFilter) (set types.ProjectFeatureSet, f types.ProjectFeatureFilter, err error)
	}

	projectFeatureAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateProjectFeature(context.Context, *types.ProjectFeature) bool
		CanDeleteProjectFeature(context.Context, *types.ProjectFeature) bool
	}
)

func (ProjectFeature) New() *ProjectFeature {
	return &ProjectFeature{
		projectFeature: service.DefaultProjectFeature,
		ac:             service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *ProjectFeature) makeFilter(ctx context.Context, r *request.ProjectFeatureList) (types.ProjectFeatureFilter, error) {
	var (
		err error
		f   = types.ProjectFeatureFilter{
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
func (ctrl *ProjectFeature) beforeCreate(ctx context.Context, res *types.ProjectFeature, r *request.ProjectFeatureCreate) error {
	res.ProjectID = r.ProjectID
	return nil
}

// beforeUpdate is a no-op; the generated controller mapped all value params.
func (ctrl *ProjectFeature) beforeUpdate(ctx context.Context, res *types.ProjectFeature, r *request.ProjectFeatureUpdate) error {
	return nil
}

func (ctrl *ProjectFeature) makePayload(ctx context.Context, a *types.ProjectFeature, err error) (*projectFeaturePayload, error) {
	if err != nil || a == nil {
		return nil, err
	}

	return &projectFeaturePayload{
		ProjectFeature: a,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateProjectFeature: ctrl.ac.CanUpdateProjectFeature(ctx, a),
		CanDeleteProjectFeature: ctrl.ac.CanDeleteProjectFeature(ctx, a),
	}, nil
}

func (ctrl *ProjectFeature) makeFilterPayload(ctx context.Context, nn types.ProjectFeatureSet, f types.ProjectFeatureFilter, err error) (*projectFeatureSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectFeatureSetPayload{Filter: f, Set: make([]*projectFeaturePayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
