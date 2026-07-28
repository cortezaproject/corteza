package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type projectReviewAccessController interface {
	CanCreateProjectReview(context.Context) bool
	CanSearchProjectReviews(context.Context) bool
	CanReadProjectReview(context.Context, *types.ProjectReview) bool
	CanUpdateProjectReview(context.Context, *types.ProjectReview) bool
	CanDeleteProjectReview(context.Context, *types.ProjectReview) bool
}

type projectReview struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectReviewAccessController
}

type projectReviewServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *projectReview) FindByID(ctx context.Context, ID uint64) (res *types.ProjectReview, err error) {
	var (
		aProps = &projectReviewActionProps{projectReview: &types.ProjectReview{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProjectReview(ctx, svc.store, ID); err != nil {
			return ProjectReviewErrInvalidID().Wrap(err)
		}

		aProps.setProjectReview(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanReadProjectReview(ctx, res) {
			return ProjectReviewErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectReviewActionLookup, err)
}

func (svc *projectReview) Search(ctx context.Context, filter types.ProjectReviewFilter) (set types.ProjectReviewSet, f types.ProjectReviewFilter, err error) {
	var (
		aProps = &projectReviewActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ProjectReview) (bool, error) {
		if !svc.ac.CanReadProjectReview(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjectReviews(ctx) {
			return ProjectReviewErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchProjectReviews(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectReviewActionSearch, err)
}

func (svc *projectReview) Create(ctx context.Context, new *types.ProjectReview) (res *types.ProjectReview, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectReviewActionProps{projectReview: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProjectReview(ctx) {
			return ProjectReviewErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProjectReview(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectReviewActionCreate, err)
}

func (svc *projectReview) Update(ctx context.Context, upd *types.ProjectReview) (res *types.ProjectReview, err error) {
	var (
		aProps = &projectReviewActionProps{update: upd}
		old    *types.ProjectReview
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProjectReview(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProjectReview(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanUpdateProjectReview(ctx, res) {
			return ProjectReviewErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectReviewErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.RevisionID = upd.RevisionID
		res.Title = upd.Title
		res.Description = upd.Description
		res.ReviewType = upd.ReviewType
		res.ReviewFrequency = upd.ReviewFrequency
		res.Scope = upd.Scope
		res.Reviewer = upd.Reviewer
		res.ApprovedBy = upd.ApprovedBy
		res.Status = upd.Status
		res.DateDue = upd.DateDue
		res.UpdatedAt = now()

		if err = store.UpdateProjectReview(ctx, svc.store, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectReviewActionUpdate, err, old, res)
}

func (svc *projectReview) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectReviewActionProps{}
		res    *types.ProjectReview
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadProjectReview(ctx, svc.store, ID); err != nil {
			return
		}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		aProps.setProjectReview(res)

		if !svc.ac.CanDeleteProjectReview(ctx, res) {
			return ProjectReviewErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateProjectReview(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectReviewActionDelete, err)
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
func (svc *projectReview) guard(_ context.Context, _ *types.ProjectReview) error { return nil }

func (svc *projectReview) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
func (svc *projectReview) scopeServices(ctx context.Context) *projectReviewServices {
	return &projectReviewServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
