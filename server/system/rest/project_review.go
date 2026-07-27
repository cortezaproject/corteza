package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ProjectReview struct {
		projectReview projectReviewService
		ac            projectReviewAccessController
	}

	projectReviewPayload struct {
		*types.ProjectReview

		CanGrant               bool `json:"canGrant"`
		CanUpdateProjectReview bool `json:"canUpdateProjectReview"`
		CanDeleteProjectReview bool `json:"canDeleteProjectReview"`
	}

	projectReviewSetPayload struct {
		Filter types.ProjectReviewFilter `json:"filter"`
		Set    []*projectReviewPayload   `json:"set"`
	}

	projectReviewService interface {
		FindByID(ctx context.Context, ID uint64) (a *types.ProjectReview, err error)
		Create(ctx context.Context, new *types.ProjectReview) (a *types.ProjectReview, err error)
		Update(ctx context.Context, upd *types.ProjectReview) (a *types.ProjectReview, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.ProjectReviewFilter) (set types.ProjectReviewSet, f types.ProjectReviewFilter, err error)
	}

	projectReviewAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateProjectReview(context.Context, *types.ProjectReview) bool
		CanDeleteProjectReview(context.Context, *types.ProjectReview) bool
	}
)

func (ProjectReview) New() *ProjectReview {
	return &ProjectReview{
		projectReview: service.DefaultProjectReview,
		ac:            service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *ProjectReview) makeFilter(ctx context.Context, r *request.ProjectReviewList) (types.ProjectReviewFilter, error) {
	var (
		err error
		f   = types.ProjectReviewFilter{
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
func (ctrl *ProjectReview) beforeCreate(ctx context.Context, res *types.ProjectReview, r *request.ProjectReviewCreate) error {
	res.ProjectID = r.ProjectID
	return nil
}

// beforeUpdate is a no-op; the generated controller mapped all value params.
func (ctrl *ProjectReview) beforeUpdate(ctx context.Context, res *types.ProjectReview, r *request.ProjectReviewUpdate) error {
	return nil
}

func (ctrl *ProjectReview) makePayload(ctx context.Context, a *types.ProjectReview, err error) (*projectReviewPayload, error) {
	if err != nil || a == nil {
		return nil, err
	}

	return &projectReviewPayload{
		ProjectReview: a,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateProjectReview: ctrl.ac.CanUpdateProjectReview(ctx, a),
		CanDeleteProjectReview: ctrl.ac.CanDeleteProjectReview(ctx, a),
	}, nil
}

func (ctrl *ProjectReview) makeFilterPayload(ctx context.Context, nn types.ProjectReviewSet, f types.ProjectReviewFilter, err error) (*projectReviewSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectReviewSetPayload{Filter: f, Set: make([]*projectReviewPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
