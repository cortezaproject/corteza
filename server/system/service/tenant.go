package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	tenantAccessController interface {
		CanCreateTenant(ctx context.Context) bool
		CanSearchTenants(ctx context.Context) bool
		CanReadTenant(ctx context.Context, t *types.Tenant) bool
		CanUpdateTenant(ctx context.Context, t *types.Tenant) bool
		CanDeleteTenant(ctx context.Context, t *types.Tenant) bool
		CanSuspendTenant(ctx context.Context, t *types.Tenant) bool
		CanManageMembersOnTenant(ctx context.Context, t *types.Tenant) bool
	}

	TenantService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Tenant, error)
		FindByHandle(ctx context.Context, handle string) (*types.Tenant, error)
		Search(ctx context.Context, filter types.TenantFilter) (types.TenantSet, types.TenantFilter, error)

		Create(ctx context.Context, new *types.Tenant) (*types.Tenant, error)
		Update(ctx context.Context, upd *types.Tenant) (*types.Tenant, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error

		Suspend(ctx context.Context, ID uint64) error
		Activate(ctx context.Context, ID uint64) error
		Archive(ctx context.Context, ID uint64) error

		SearchMembers(ctx context.Context, filter types.TenantMembershipFilter) (types.TenantMembershipSet, types.TenantMembershipFilter, error)
		Invite(ctx context.Context, m *types.TenantMembership) (*types.TenantMembership, error)
		AcceptInvite(ctx context.Context, tenantID, userID uint64) error
		UpdateMember(ctx context.Context, m *types.TenantMembership) (*types.TenantMembership, error)
		RemoveMember(ctx context.Context, tenantID, userID uint64) error
		SuspendMember(ctx context.Context, tenantID, userID uint64) error
		ActivateMember(ctx context.Context, tenantID, userID uint64) error
	}
)

func (svc *tenant) FindByID(ctx context.Context, ID uint64) (t *types.Tenant, err error) {
	var taProps = &tenantActionProps{tenant: &types.Tenant{ID: ID}}

	err = func() error {
		if t, err = loadTenant(ctx, svc.store, ID); err != nil {
			return err
		}

		taProps.setTenant(t)

		if !svc.ac.CanReadTenant(ctx, t) {
			return TenantErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, t); err != nil {
			return err
		}

		return nil
	}()

	return t, svc.recordAction(ctx, taProps, TenantActionLookup, err)
}

func (svc *tenant) FindByHandle(ctx context.Context, h string) (t *types.Tenant, err error) {
	var taProps = &tenantActionProps{tenant: &types.Tenant{Handle: h}}

	err = func() error {
		if t, err = store.LookupTenantByHandle(ctx, svc.store, h); errors.IsNotFound(err) {
			return TenantErrNotFound()
		} else if err != nil {
			return err
		}

		taProps.setTenant(t)

		if !svc.ac.CanReadTenant(ctx, t) {
			return TenantErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, t); err != nil {
			return err
		}

		return nil
	}()

	return t, svc.recordAction(ctx, taProps, TenantActionLookup, err)
}

// beforeCreate runs after the create RBAC check and before id/timestamps are
// assigned. It validates the incoming tenant, defaults + validates the status,
// enforces handle uniqueness and stamps CreatedBy (the generated Create body
// only sets ID + CreatedAt).
func (svc *tenant) beforeCreate(ctx context.Context, new *types.Tenant) (err error) {
	if !handle.IsValid(new.Handle) {
		return TenantErrInvalidHandle()
	}

	if new.Status == "" {
		new.Status = types.TenantStatusActive
	}
	if err = validateTenantStatus(new.Status); err != nil {
		return err
	}

	if err = svc.uniqueCheck(ctx, new); err != nil {
		return err
	}

	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

	return nil
}

func (svc *tenant) Update(ctx context.Context, upd *types.Tenant) (t *types.Tenant, err error) {
	var (
		taProps  = &tenantActionProps{update: upd}
		existing *types.Tenant
	)

	err = func() (err error) {
		if existing, err = loadTenant(ctx, svc.store, upd.ID); err != nil {
			return
		}

		taProps.setTenant(existing)

		if !svc.ac.CanUpdateTenant(ctx, existing) {
			return TenantErrNotAllowedToUpdate()
		}

		if isStale(upd.UpdatedAt, existing.UpdatedAt, existing.CreatedAt) {
			return TenantErrStaleData()
		}

		if !handle.IsValid(upd.Handle) {
			return TenantErrInvalidHandle()
		}

		if upd.Status == "" {
			upd.Status = existing.Status
		}
		if err = validateTenantStatus(upd.Status); err != nil {
			return err
		}

		if upd.Handle != existing.Handle {
			if err = svc.uniqueCheck(ctx, upd); err != nil {
				return err
			}
		}

		upd.UpdatedAt = now()
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
		upd.CreatedAt = existing.CreatedAt
		upd.CreatedBy = existing.CreatedBy
		upd.SuspendedAt = existing.SuspendedAt
		upd.DeletedAt = existing.DeletedAt

		if err = store.UpdateTenant(ctx, svc.store, upd); err != nil {
			return
		}

		if err = label.Update(ctx, svc.store, upd); err != nil {
			return
		}

		t = upd
		return nil
	}()

	return t, svc.recordAction(ctx, taProps, TenantActionUpdate, err, existing, upd)
}

// beforeDelete runs after the delete RBAC check and before the soft-delete
// store update. It stamps DeletedBy from the context identity (the generated
// delete body only sets DeletedAt).
func (svc *tenant) beforeDelete(ctx context.Context, t *types.Tenant) error {
	t.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *tenant) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		taProps = &tenantActionProps{tenant: &types.Tenant{ID: ID}}
		t       *types.Tenant
		old     *types.Tenant
	)

	err = func() (err error) {
		if t, err = loadTenant(ctx, svc.store, ID); err != nil {
			return
		}

		taProps.setTenant(t)

		if !svc.ac.CanDeleteTenant(ctx, t) {
			return TenantErrNotAllowedToDelete()
		}

		old = t.Clone()

		t.DeletedAt = nil
		t.DeletedBy = 0
		if err = store.UpdateTenant(ctx, svc.store, t); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, taProps, TenantActionUndelete, err, old, t)
}

func (svc *tenant) onSuspend(ctx context.Context, _ *tenantActionProps, ID uint64) error {
	return svc.setStatus(ctx, ID, types.TenantStatusSuspended, TenantActionSuspend)
}

func (svc *tenant) onActivate(ctx context.Context, _ *tenantActionProps, ID uint64) error {
	return svc.setStatus(ctx, ID, types.TenantStatusActive, TenantActionActivate)
}

func (svc *tenant) onArchive(ctx context.Context, _ *tenantActionProps, ID uint64) error {
	return svc.setStatus(ctx, ID, types.TenantStatusArchived, TenantActionArchive)
}

// setStatus transitions a tenant to the given status, stamping SuspendedAt when
// suspending and clearing it otherwise.
func (svc *tenant) setStatus(ctx context.Context, ID uint64, status types.TenantStatus, action func(...*tenantActionProps) *tenantAction) (err error) {
	var (
		taProps = &tenantActionProps{tenant: &types.Tenant{ID: ID}}
		t       *types.Tenant
		old     *types.Tenant
	)

	err = func() (err error) {
		if t, err = loadTenant(ctx, svc.store, ID); err != nil {
			return
		}

		taProps.setTenant(t)

		if !svc.ac.CanSuspendTenant(ctx, t) {
			return TenantErrNotAllowedToSuspend()
		}

		old = t.Clone()

		t.Status = status
		if status == types.TenantStatusSuspended {
			t.SuspendedAt = now()
		} else {
			t.SuspendedAt = nil
		}
		t.UpdatedAt = now()
		t.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateTenant(ctx, svc.store, t); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, taProps, action, err, old, t)
}

// --- members ---

func (svc *tenant) onSearchMembers(ctx context.Context, _ *tenantActionProps, filter types.TenantMembershipFilter) (set types.TenantMembershipSet, f types.TenantMembershipFilter, err error) {
	var t *types.Tenant
	if t, err = loadTenant(ctx, svc.store, filter.TenantID); err != nil {
		return set, f, err
	}

	if !svc.ac.CanReadTenant(ctx, t) {
		return set, f, TenantErrNotAllowedToRead()
	}

	if set, f, err = store.SearchTenantMemberships(ctx, svc.store, filter); err != nil {
		return set, f, err
	}

	return set, f, nil
}

// Invite adds a user to a tenant.
//
// Constraint: one user = one tenant. The user must not already hold an active
// membership in any tenant.
func (svc *tenant) onInvite(ctx context.Context, _ *tenantActionProps, m *types.TenantMembership) (res *types.TenantMembership, err error) {
	var t *types.Tenant
	if t, err = loadTenant(ctx, svc.store, m.TenantID); err != nil {
		return res, err
	}

	if !svc.ac.CanManageMembersOnTenant(ctx, t) {
		return res, TenantErrNotAllowedToManageMembers()
	}

	if m.UserID == 0 {
		return res, TenantErrMemberNotFound()
	}

	if m.Role == "" {
		m.Role = types.TenantRoleMember
	}
	if !m.Role.Valid() {
		return res, TenantErrInvalidRole()
	}

	// one membership record per user per tenant
	if existing, e := store.LookupTenantMembershipByTenantIDUserID(ctx, svc.store, m.TenantID, m.UserID); e == nil && existing != nil {
		return res, TenantErrMemberAlreadyExists()
	} else if e != nil && !errors.IsNotFound(e) {
		return res, e
	}

	// one user = one tenant: refuse if user has an active membership anywhere
	if err = svc.assertUserHasNoActiveTenant(ctx, m.UserID); err != nil {
		return res, err
	}

	if m.Status == "" {
		m.Status = types.TenantMemberStatusInvited
	}
	if !m.Status.Valid() {
		return res, TenantErrInvalidMemberStatus()
	}

	m.ID = nextID()
	m.CreatedAt = *now()
	m.InvitedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.CreateTenantMembership(ctx, svc.store, m); err != nil {
		return res, err
	}

	res = m
	return res, nil
}

func (svc *tenant) onAcceptInvite(ctx context.Context, _ *tenantActionProps, tenantID, userID uint64) (err error) {
	var existing *types.TenantMembership
	if existing, err = store.LookupTenantMembershipByTenantIDUserID(ctx, svc.store, tenantID, userID); errors.IsNotFound(err) {
		return TenantErrMemberNotFound()
	} else if err != nil {
		return err
	}

	existing.Status = types.TenantMemberStatusActive
	existing.UpdatedAt = now()
	return store.UpdateTenantMembership(ctx, svc.store, existing)
}

func (svc *tenant) onUpdateMember(ctx context.Context, _ *tenantActionProps, m *types.TenantMembership) (res *types.TenantMembership, err error) {
	var t *types.Tenant
	if t, err = loadTenant(ctx, svc.store, m.TenantID); err != nil {
		return res, err
	}

	if !svc.ac.CanManageMembersOnTenant(ctx, t) {
		return res, TenantErrNotAllowedToManageMembers()
	}

	var existing *types.TenantMembership
	if existing, err = store.LookupTenantMembershipByTenantIDUserID(ctx, svc.store, m.TenantID, m.UserID); errors.IsNotFound(err) {
		return res, TenantErrMemberNotFound()
	} else if err != nil {
		return res, err
	}

	if !m.Role.Valid() {
		return res, TenantErrInvalidRole()
	}

	existing.Role = m.Role
	if m.Status != "" {
		if !m.Status.Valid() {
			return res, TenantErrInvalidMemberStatus()
		}
		existing.Status = m.Status
	}
	existing.UpdatedAt = now()

	if err = store.UpdateTenantMembership(ctx, svc.store, existing); err != nil {
		return res, err
	}

	res = existing
	return res, nil
}

func (svc *tenant) onRemoveMember(ctx context.Context, _ *tenantActionProps, tenantID, userID uint64) (err error) {
	var t *types.Tenant
	if t, err = loadTenant(ctx, svc.store, tenantID); err != nil {
		return err
	}

	if !svc.ac.CanManageMembersOnTenant(ctx, t) {
		return TenantErrNotAllowedToManageMembers()
	}

	var existing *types.TenantMembership
	if existing, err = store.LookupTenantMembershipByTenantIDUserID(ctx, svc.store, tenantID, userID); errors.IsNotFound(err) {
		return TenantErrMemberNotFound()
	} else if err != nil {
		return err
	}

	return store.DeleteTenantMembership(ctx, svc.store, existing)
}

func (svc *tenant) onSuspendMember(ctx context.Context, _ *tenantActionProps, tenantID, userID uint64) error {
	return svc.setMemberStatus(ctx, tenantID, userID, types.TenantMemberStatusSuspended)
}

func (svc *tenant) onActivateMember(ctx context.Context, _ *tenantActionProps, tenantID, userID uint64) error {
	return svc.setMemberStatus(ctx, tenantID, userID, types.TenantMemberStatusActive)
}

func (svc *tenant) setMemberStatus(ctx context.Context, tenantID, userID uint64, status types.TenantMemberStatus) (err error) {
	return func() (err error) {
		var t *types.Tenant
		if t, err = loadTenant(ctx, svc.store, tenantID); err != nil {
			return
		}

		if !svc.ac.CanManageMembersOnTenant(ctx, t) {
			return TenantErrNotAllowedToManageMembers()
		}

		var existing *types.TenantMembership
		if existing, err = store.LookupTenantMembershipByTenantIDUserID(ctx, svc.store, tenantID, userID); errors.IsNotFound(err) {
			return TenantErrMemberNotFound()
		} else if err != nil {
			return err
		}

		existing.Status = status
		existing.UpdatedAt = now()
		return store.UpdateTenantMembership(ctx, svc.store, existing)
	}()
}

// --- helpers ---

func (svc *tenant) uniqueCheck(ctx context.Context, t *types.Tenant) error {
	if t.Handle == "" {
		return nil
	}
	if e, err := store.LookupTenantByHandle(ctx, svc.store, t.Handle); err == nil && e != nil && e.ID != t.ID {
		return TenantErrHandleNotUnique()
	} else if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// assertUserHasNoActiveTenant enforces the one-user-one-tenant rule. It returns
// TenantErrUserAlreadyInTenant when the user holds an active membership in any
// tenant.
func (svc *tenant) assertUserHasNoActiveTenant(ctx context.Context, userID uint64) error {
	f := types.TenantMembershipFilter{
		UserID: userID,
		Status: types.TenantMemberStatusActive,
	}
	set, _, err := store.SearchTenantMemberships(ctx, svc.store, f)
	if err != nil {
		return err
	}
	if len(set) > 0 {
		return TenantErrUserAlreadyInTenant()
	}
	return nil
}

func validateTenantStatus(s types.TenantStatus) error {
	switch s {
	case types.TenantStatusActive, types.TenantStatusArchived, types.TenantStatusSuspended:
		return nil
	}
	return TenantErrInvalidStatus()
}
