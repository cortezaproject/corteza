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
	projectFeature struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        projectFeatureAccessController
	}

	projectFeatureAccessController interface {
		CanSearchProjectFeatures(context.Context) bool
		CanCreateProjectFeature(context.Context) bool
		CanReadProjectFeature(context.Context, *types.ProjectFeature) bool
		CanUpdateProjectFeature(context.Context, *types.ProjectFeature) bool
		CanDeleteProjectFeature(context.Context, *types.ProjectFeature) bool
	}
)

func ProjectFeature() *projectFeature {
	return &projectFeature{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc projectFeature) FindByID(ctx context.Context, ID uint64) (r *types.ProjectFeature, err error) {
	var (
		raProps = &projectFeatureActionProps{projectFeature: &types.ProjectFeature{ID: ID}}
	)

	err = func() error {
		if r, err = loadProjectFeature(ctx, svc.store, ID); err != nil {
			return err
		}

		raProps.setProjectFeature(r)

		if !svc.ac.CanReadProjectFeature(ctx, r) {
			return ProjectFeatureErrNotAllowedToRead()
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectFeatureActionLookup, err)
}

func (svc projectFeature) Search(ctx context.Context, filter types.ProjectFeatureFilter) (rr types.ProjectFeatureSet, f types.ProjectFeatureFilter, err error) {
	var (
		raProps = &projectFeatureActionProps{filter: &filter}
	)

	filter.Check = func(res *types.ProjectFeature) (bool, error) {
		if !svc.ac.CanReadProjectFeature(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchProjectFeatures(ctx) {
			return ProjectFeatureErrNotAllowedToSearch()
		}

		if rr, f, err = store.SearchProjectFeatures(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, ProjectFeatureActionSearch, err)
}

func (svc projectFeature) Create(ctx context.Context, new *types.ProjectFeature) (r *types.ProjectFeature, err error) {
	var (
		raProps = &projectFeatureActionProps{projectFeature: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateProjectFeature(ctx) {
			return ProjectFeatureErrNotAllowedToCreate()
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

		if err = store.CreateProjectFeature(ctx, svc.store, new); err != nil {
			return
		}

		r = new
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectFeatureActionCreate, err)
}

func (svc projectFeature) Update(ctx context.Context, upd *types.ProjectFeature) (r *types.ProjectFeature, err error) {
	var (
		raProps = &projectFeatureActionProps{projectFeature: upd, update: upd}
	)

	err = func() (err error) {
		if r, err = loadProjectFeature(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setProjectFeature(r)

		if !svc.ac.CanUpdateProjectFeature(ctx, r) {
			return ProjectFeatureErrNotAllowedToUpdate()
		}

		r.Title = upd.Title
		r.Description = upd.Description
		r.FeatureType = upd.FeatureType
		r.Status = upd.Status
		r.Severity = upd.Severity
		r.Risk = upd.Risk
		r.FeatureOwner = upd.FeatureOwner
		r.ChangeOwner = upd.ChangeOwner
		r.ChangeApprovedBy = upd.ChangeApprovedBy
		r.RiskFeature = upd.RiskFeature
		r.ChangeRequired = upd.ChangeRequired
		r.RiskChange = upd.RiskChange
		r.Backlog = upd.Backlog
		r.DateDue = upd.DateDue
		r.UpdatedAt = now()
		r.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateProjectFeature(ctx, svc.store, r); err != nil {
			return err
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectFeatureActionUpdate, err)
}

func (svc projectFeature) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		raProps = &projectFeatureActionProps{}
		r       *types.ProjectFeature
	)

	err = func() (err error) {
		if r, err = loadProjectFeature(ctx, svc.store, ID); err != nil {
			return
		}

		raProps.setProjectFeature(r)

		if !svc.ac.CanDeleteProjectFeature(ctx, r) {
			return ProjectFeatureErrNotAllowedToDelete()
		}

		r.DeletedAt = now()
		if err = store.UpdateProjectFeature(ctx, svc.store, r); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, ProjectFeatureActionDelete, err)
}

func loadProjectFeature(ctx context.Context, s store.ProjectFeatures, ID uint64) (res *types.ProjectFeature, err error) {
	if ID == 0 {
		return nil, ProjectFeatureErrInvalidID()
	}

	if res, err = store.LookupProjectFeatureByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectFeatureErrNotFound()
	}

	return
}
