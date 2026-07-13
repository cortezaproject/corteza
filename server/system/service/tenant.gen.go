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
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type tenant struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        tenantAccessController
}

type tenantServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *tenant) Search(ctx context.Context, filter types.TenantFilter) (set types.TenantSet, f types.TenantFilter, err error) {
	var (
		aProps = &tenantActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Tenant) (bool, error) {
		if !svc.ac.CanReadTenant(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchTenants(ctx) {
			return TenantErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Tenant{}.LabelResourceKind(),
				filter.Labels,
			)
			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchTenants(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledTenants(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, TenantActionSearch, err)
}

func (svc *tenant) Create(ctx context.Context, new *types.Tenant) (res *types.Tenant, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &tenantActionProps{tenant: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateTenant(ctx) {
			return TenantErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateTenant(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, TenantActionCreate, err)
}

func (svc *tenant) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &tenantActionProps{}
		res    *types.Tenant
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadTenant(ctx, svc.store, ID); err != nil {
			return
		}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		aProps.setTenant(res)

		if !svc.ac.CanDeleteTenant(ctx, res) {
			return TenantErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateTenant(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, TenantActionDelete, err)
}

func loadTenant(ctx context.Context, s store.Tenants, ID uint64) (res *types.Tenant, err error) {
	if ID == 0 {
		return nil, TenantErrInvalidID()
	}

	if res, err = store.LookupTenantByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, TenantErrNotFound()
	}

	return
}

// toLabeledTenants converts to []label.LabeledResource
func toLabeledTenants(set []*types.Tenant) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
func (svc *tenant) guard(_ context.Context, _ *types.Tenant) error { return nil }

func (svc *tenant) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
func (svc *tenant) scopeServices(ctx context.Context) *tenantServices {
	return &tenantServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}

func (svc *tenant) Suspend(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onSuspend(ctx, aProps, ID)
		return err
	}()

	return svc.recordAction(ctx, aProps, TenantActionSuspend, err)
}

func (svc *tenant) Activate(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onActivate(ctx, aProps, ID)
		return err
	}()

	return svc.recordAction(ctx, aProps, TenantActionActivate, err)
}

func (svc *tenant) Archive(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onArchive(ctx, aProps, ID)
		return err
	}()

	return svc.recordAction(ctx, aProps, TenantActionArchive, err)
}

func (svc *tenant) SearchMembers(ctx context.Context, filter types.TenantMembershipFilter) (set types.TenantMembershipSet, f types.TenantMembershipFilter, err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		set, f, err = svc.onSearchMembers(ctx, aProps, filter)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, TenantActionSearchMembers, err)
}

func (svc *tenant) Invite(ctx context.Context, m *types.TenantMembership) (res *types.TenantMembership, err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		res, err = svc.onInvite(ctx, aProps, m)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, TenantActionInvite, err)
}

func (svc *tenant) AcceptInvite(ctx context.Context, tenantID uint64, userID uint64) (err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onAcceptInvite(ctx, aProps, tenantID, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, TenantActionAcceptInvite, err)
}

func (svc *tenant) UpdateMember(ctx context.Context, m *types.TenantMembership) (res *types.TenantMembership, err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		res, err = svc.onUpdateMember(ctx, aProps, m)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, TenantActionUpdateMember, err)
}

func (svc *tenant) RemoveMember(ctx context.Context, tenantID uint64, userID uint64) (err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onRemoveMember(ctx, aProps, tenantID, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, TenantActionRemoveMember, err)
}

func (svc *tenant) SuspendMember(ctx context.Context, tenantID uint64, userID uint64) (err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onSuspendMember(ctx, aProps, tenantID, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, TenantActionSuspendMember, err)
}

func (svc *tenant) ActivateMember(ctx context.Context, tenantID uint64, userID uint64) (err error) {
	var (
		aProps = &tenantActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onActivateMember(ctx, aProps, tenantID, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, TenantActionActivateMember, err)
}
