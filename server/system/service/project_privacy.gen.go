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

type projectPrivacyAccessController interface {
	CanCreateProjectPrivacy(context.Context) bool
	CanSearchProjectPrivacys(context.Context) bool
	CanReadProjectPrivacy(context.Context, *types.ProjectPrivacy) bool
	CanUpdateProjectPrivacy(context.Context, *types.ProjectPrivacy) bool
	CanDeleteProjectPrivacy(context.Context, *types.ProjectPrivacy) bool
}

type projectPrivacy struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectPrivacyAccessController
}

type projectPrivacyServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *projectPrivacy) FindByID(ctx context.Context, ID uint64) (res *types.ProjectPrivacy, err error) {
	var (
		aProps = &projectPrivacyActionProps{projectPrivacy: &types.ProjectPrivacy{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProjectPrivacy(ctx, svc.store, ID); err != nil {
			return ProjectPrivacyErrInvalidID().Wrap(err)
		}

		aProps.setProjectPrivacy(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanReadProjectPrivacy(ctx, res) {
			return ProjectPrivacyErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectPrivacyActionLookup, err)
}

func (svc *projectPrivacy) Search(ctx context.Context, filter types.ProjectPrivacyFilter) (set types.ProjectPrivacySet, f types.ProjectPrivacyFilter, err error) {
	var (
		aProps = &projectPrivacyActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ProjectPrivacy) (bool, error) {
		if !svc.ac.CanReadProjectPrivacy(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjectPrivacys(ctx) {
			return ProjectPrivacyErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchProjectPrivacys(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectPrivacyActionSearch, err)
}

func (svc *projectPrivacy) Create(ctx context.Context, new *types.ProjectPrivacy) (res *types.ProjectPrivacy, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectPrivacyActionProps{projectPrivacy: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProjectPrivacy(ctx) {
			return ProjectPrivacyErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProjectPrivacy(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectPrivacyActionCreate, err)
}

func (svc *projectPrivacy) Update(ctx context.Context, upd *types.ProjectPrivacy) (res *types.ProjectPrivacy, err error) {
	var (
		aProps = &projectPrivacyActionProps{update: upd}
		old    *types.ProjectPrivacy
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProjectPrivacy(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProjectPrivacy(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanUpdateProjectPrivacy(ctx, res) {
			return ProjectPrivacyErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectPrivacyErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.Title = upd.Title
		res.Description = upd.Description
		res.RequestType = upd.RequestType
		res.Status = upd.Status
		res.Severity = upd.Severity
		res.Risk = upd.Risk
		res.RequestOwner = upd.RequestOwner
		res.ChangeOwner = upd.ChangeOwner
		res.ChangeApprovedBy = upd.ChangeApprovedBy
		res.RiskAssessment = upd.RiskAssessment
		res.ChangeRequired = upd.ChangeRequired
		res.RiskChange = upd.RiskChange
		res.DateDue = upd.DateDue
		res.UpdatedAt = now()

		if err = store.UpdateProjectPrivacy(ctx, svc.store, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectPrivacyActionUpdate, err, old, res)
}

func (svc *projectPrivacy) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectPrivacyActionProps{}
		res    *types.ProjectPrivacy
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadProjectPrivacy(ctx, svc.store, ID); err != nil {
			return
		}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		aProps.setProjectPrivacy(res)

		if !svc.ac.CanDeleteProjectPrivacy(ctx, res) {
			return ProjectPrivacyErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateProjectPrivacy(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectPrivacyActionDelete, err)
}

func loadProjectPrivacy(ctx context.Context, s store.ProjectPrivacys, ID uint64) (res *types.ProjectPrivacy, err error) {
	if ID == 0 {
		return nil, ProjectPrivacyErrInvalidID()
	}

	if res, err = store.LookupProjectPrivacyByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectPrivacyErrNotFound()
	}

	return
}
func (svc *projectPrivacy) guard(_ context.Context, _ *types.ProjectPrivacy) error { return nil }

func (svc *projectPrivacy) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
func (svc *projectPrivacy) scopeServices(ctx context.Context) *projectPrivacyServices {
	return &projectPrivacyServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
