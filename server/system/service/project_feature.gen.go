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

type projectFeatureAccessController interface {
	CanCreateProjectFeature(context.Context) bool
	CanSearchProjectFeatures(context.Context) bool
	CanReadProjectFeature(context.Context, *types.ProjectFeature) bool
	CanUpdateProjectFeature(context.Context, *types.ProjectFeature) bool
	CanDeleteProjectFeature(context.Context, *types.ProjectFeature) bool
}

type projectFeature struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectFeatureAccessController
}

type projectFeatureServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *projectFeature) FindByID(ctx context.Context, ID uint64) (res *types.ProjectFeature, err error) {
	var (
		aProps = &projectFeatureActionProps{projectFeature: &types.ProjectFeature{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProjectFeature(ctx, svc.store, ID); err != nil {
			return ProjectFeatureErrInvalidID().Wrap(err)
		}

		aProps.setProjectFeature(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanReadProjectFeature(ctx, res) {
			return ProjectFeatureErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectFeatureActionLookup, err)
}

func (svc *projectFeature) Search(ctx context.Context, filter types.ProjectFeatureFilter) (set types.ProjectFeatureSet, f types.ProjectFeatureFilter, err error) {
	var (
		aProps = &projectFeatureActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ProjectFeature) (bool, error) {
		if !svc.ac.CanReadProjectFeature(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if err = svc.beforeSearch(ctx, &filter); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjectFeatures(ctx) {
			return ProjectFeatureErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchProjectFeatures(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectFeatureActionSearch, err)
}

func (svc *projectFeature) Create(ctx context.Context, new *types.ProjectFeature) (res *types.ProjectFeature, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectFeatureActionProps{projectFeature: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProjectFeature(ctx) {
			return ProjectFeatureErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProjectFeature(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectFeatureActionCreate, err)
}

func (svc *projectFeature) Update(ctx context.Context, upd *types.ProjectFeature) (res *types.ProjectFeature, err error) {
	var (
		aProps = &projectFeatureActionProps{update: upd}
		old    *types.ProjectFeature
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProjectFeature(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProjectFeature(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanUpdateProjectFeature(ctx, res) {
			return ProjectFeatureErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectFeatureErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.RevisionID = upd.RevisionID
		res.Title = upd.Title
		res.Description = upd.Description
		res.FeatureType = upd.FeatureType
		res.Status = upd.Status
		res.Severity = upd.Severity
		res.Risk = upd.Risk
		res.FeatureOwner = upd.FeatureOwner
		res.ChangeOwner = upd.ChangeOwner
		res.ChangeApprovedBy = upd.ChangeApprovedBy
		res.RiskFeature = upd.RiskFeature
		res.ChangeRequired = upd.ChangeRequired
		res.RiskChange = upd.RiskChange
		res.DateDue = upd.DateDue
		res.UpdatedAt = now()

		if err = store.UpdateProjectFeature(ctx, svc.store, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectFeatureActionUpdate, err, old, res)
}

func (svc *projectFeature) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectFeatureActionProps{}
		res    *types.ProjectFeature
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadProjectFeature(ctx, svc.store, ID); err != nil {
			return
		}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		aProps.setProjectFeature(res)

		if !svc.ac.CanDeleteProjectFeature(ctx, res) {
			return ProjectFeatureErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateProjectFeature(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectFeatureActionDelete, err)
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
func (svc *projectFeature) guard(_ context.Context, _ *types.ProjectFeature) error { return nil }

func (svc *projectFeature) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
func (svc *projectFeature) scopeServices(ctx context.Context) *projectFeatureServices {
	return &projectFeatureServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
