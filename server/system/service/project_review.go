package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	projectReview struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        projectReviewAccessController
	}

	projectReviewAccessController interface {
		CanSearchProjectReviews(context.Context) bool
		CanCreateProjectReview(context.Context) bool
		CanReadProjectReview(context.Context, *types.ProjectReview) bool
		CanUpdateProjectReview(context.Context, *types.ProjectReview) bool
		CanDeleteProjectReview(context.Context, *types.ProjectReview) bool
	}
)

func ProjectReview() *projectReview {
	return &projectReview{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc projectReview) FindByID(ctx context.Context, ID uint64) (r *types.ProjectReview, err error) {
	var (
		raProps = &projectReviewActionProps{projectReview: &types.ProjectReview{ID: ID}}
	)

	err = func() error {
		if r, err = loadProjectReview(ctx, svc.store, ID); err != nil {
			return err
		}

		raProps.setProjectReview(r)

		if !svc.ac.CanReadProjectReview(ctx, r) {
			return ProjectReviewErrNotAllowedToRead()
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectReviewActionLookup, err)
}

func (svc projectReview) Search(ctx context.Context, filter types.ProjectReviewFilter) (rr types.ProjectReviewSet, f types.ProjectReviewFilter, err error) {
	var (
		raProps = &projectReviewActionProps{filter: &filter}
	)

	filter.Check = func(res *types.ProjectReview) (bool, error) {
		if !svc.ac.CanReadProjectReview(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchProjectReviews(ctx) {
			return ProjectReviewErrNotAllowedToSearch()
		}

		if rr, f, err = store.SearchProjectReviews(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, ProjectReviewActionSearch, err)
}

func (svc projectReview) Create(ctx context.Context, new *types.ProjectReview) (r *types.ProjectReview, err error) {
	var (
		raProps = &projectReviewActionProps{projectReview: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateProjectReview(ctx) {
			return ProjectReviewErrNotAllowedToCreate()
		}

		if len(new.Title) == 0 {
			new.Title = firstNonEmpty(new.Description, "(untitled)")
		}
		if len(new.Status) == 0 {
			new.Status = "Open"
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.CreateProjectReview(ctx, svc.store, new); err != nil {
			return
		}

		r = new
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectReviewActionCreate, err)
}

func (svc projectReview) Update(ctx context.Context, upd *types.ProjectReview) (r *types.ProjectReview, err error) {
	var (
		raProps = &projectReviewActionProps{projectReview: upd, update: upd}
	)

	err = func() (err error) {
		if r, err = loadProjectReview(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setProjectReview(r)

		if !svc.ac.CanUpdateProjectReview(ctx, r) {
			return ProjectReviewErrNotAllowedToUpdate()
		}

		r.Title = upd.Title
		r.Description = upd.Description
		r.ReviewType = upd.ReviewType
		r.ReviewFrequency = upd.ReviewFrequency
		r.Scope = upd.Scope
		r.Reviewer = upd.Reviewer
		r.ApprovedBy = upd.ApprovedBy
		r.Status = upd.Status
		r.DateDue = upd.DateDue
		r.UpdatedAt = now()
		r.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateProjectReview(ctx, svc.store, r); err != nil {
			return err
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectReviewActionUpdate, err)
}

func (svc projectReview) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		raProps = &projectReviewActionProps{}
		r       *types.ProjectReview
	)

	err = func() (err error) {
		if r, err = loadProjectReview(ctx, svc.store, ID); err != nil {
			return
		}

		raProps.setProjectReview(r)

		if !svc.ac.CanDeleteProjectReview(ctx, r) {
			return ProjectReviewErrNotAllowedToDelete()
		}

		r.DeletedAt = now()
		if err = store.UpdateProjectReview(ctx, svc.store, r); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, ProjectReviewActionDelete, err)
}

func loadProjectReview(ctx context.Context, s store.ProjectReviews, ID uint64) (res *types.ProjectReview, err error) {
	if ID == 0 {
		return nil, ProjectReviewErrInvalidID()
	}

	if res, err = store.LookupProjectReviewByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectReviewErrNotFound()
	}

	return
}
