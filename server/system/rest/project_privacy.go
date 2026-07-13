package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ProjectPrivacy struct {
		projectPrivacy projectPrivacyService
		ac             projectPrivacyAccessController
	}

	projectPrivacyPayload struct {
		*types.ProjectPrivacy

		CanGrant                bool `json:"canGrant"`
		CanUpdateProjectPrivacy bool `json:"canUpdateProjectPrivacy"`
		CanDeleteProjectPrivacy bool `json:"canDeleteProjectPrivacy"`
	}

	projectPrivacySetPayload struct {
		Filter types.ProjectPrivacyFilter `json:"filter"`
		Set    []*projectPrivacyPayload   `json:"set"`
	}

	projectPrivacyService interface {
		FindByID(ctx context.Context, ID uint64) (a *types.ProjectPrivacy, err error)
		Create(ctx context.Context, new *types.ProjectPrivacy) (a *types.ProjectPrivacy, err error)
		Update(ctx context.Context, upd *types.ProjectPrivacy) (a *types.ProjectPrivacy, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.ProjectPrivacyFilter) (set types.ProjectPrivacySet, f types.ProjectPrivacyFilter, err error)
	}

	projectPrivacyAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateProjectPrivacy(context.Context, *types.ProjectPrivacy) bool
		CanDeleteProjectPrivacy(context.Context, *types.ProjectPrivacy) bool
	}
)

func (ProjectPrivacy) New() *ProjectPrivacy {
	return &ProjectPrivacy{
		projectPrivacy: service.DefaultProjectPrivacy,
		ac:             service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *ProjectPrivacy) makeFilter(ctx context.Context, r *request.ProjectPrivacyList) (types.ProjectPrivacyFilter, error) {
	var (
		err error
		f   = types.ProjectPrivacyFilter{
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
func (ctrl *ProjectPrivacy) beforeCreate(ctx context.Context, res *types.ProjectPrivacy, r *request.ProjectPrivacyCreate) error {
	res.ProjectID = r.ProjectID
	return nil
}

// beforeUpdate is a no-op; the generated controller mapped all value params.
func (ctrl *ProjectPrivacy) beforeUpdate(ctx context.Context, res *types.ProjectPrivacy, r *request.ProjectPrivacyUpdate) error {
	return nil
}

func (ctrl *ProjectPrivacy) makePayload(ctx context.Context, a *types.ProjectPrivacy, err error) (*projectPrivacyPayload, error) {
	if err != nil || a == nil {
		return nil, err
	}

	return &projectPrivacyPayload{
		ProjectPrivacy: a,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateProjectPrivacy: ctrl.ac.CanUpdateProjectPrivacy(ctx, a),
		CanDeleteProjectPrivacy: ctrl.ac.CanDeleteProjectPrivacy(ctx, a),
	}, nil
}

func (ctrl *ProjectPrivacy) makeFilterPayload(ctx context.Context, nn types.ProjectPrivacySet, f types.ProjectPrivacyFilter, err error) (*projectPrivacySetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectPrivacySetPayload{Filter: f, Set: make([]*projectPrivacyPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
