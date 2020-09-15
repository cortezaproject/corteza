package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"strconv"
)

type (
	role struct {
		ctx context.Context

		actionlog actionlog.Recorder

		ac       roleAccessController
		eventbus eventDispatcher

		user UserService

		store store.Storer
	}

	roleAccessController interface {
		CanAccess(context.Context) bool

		CanCreateRole(context.Context) bool
		CanReadRole(context.Context, *types.Role) bool
		CanUpdateRole(context.Context, *types.Role) bool
		CanDeleteRole(context.Context, *types.Role) bool
		CanManageRoleMembers(context.Context, *types.Role) bool
	}

	RoleService interface {
		With(ctx context.Context) RoleService

		FindByID(roleID uint64) (*types.Role, error)
		FindByName(name string) (*types.Role, error)
		FindByHandle(handle string) (*types.Role, error)
		FindByAny(ctx context.Context, identifier interface{}) (*types.Role, error)
		Find(types.RoleFilter) (types.RoleSet, types.RoleFilter, error)

		Create(role *types.Role) (*types.Role, error)
		Update(role *types.Role) (*types.Role, error)

		Archive(ID uint64) error
		Unarchive(ID uint64) error
		Delete(ID uint64) error
		Undelete(ID uint64) error

		Membership(userID uint64) (types.RoleMemberSet, error)
		MemberList(roleID uint64) (types.RoleMemberSet, error)
		MemberAdd(roleID, userID uint64) error
		MemberRemove(roleID, userID uint64) error
	}
)

func Role(ctx context.Context) RoleService {
	return (&role{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),

		actionlog: DefaultActionlog,

		user:  DefaultUser.With(ctx),
		store: DefaultStore,
	}).With(ctx)
}

func (svc role) With(ctx context.Context) RoleService {
	return &role{
		ctx: ctx,

		actionlog: svc.actionlog,

		ac:       svc.ac,
		eventbus: svc.eventbus,
		user:     svc.user,

		store: DefaultStore,
	}
}

func (svc role) Find(filter types.RoleFilter) (rr types.RoleSet, f types.RoleFilter, err error) {
	var (
		raProps = &roleActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Role) (bool, error) {
		if !svc.ac.CanReadRole(svc.ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if filter.Deleted > 0 {
			// If list with deleted or suspended users is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted or archived roles
			if !svc.ac.CanAccess(svc.ctx) {
				return RoleErrNotAllowedToListRoles()
			}
		}

		rr, f, err = store.SearchRoles(svc.ctx, svc.store, filter)
		return err
	}()

	return rr, f, svc.recordAction(svc.ctx, raProps, RoleActionSearch, err)
}

func (svc role) FindByID(roleID uint64) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() error {
		r, err = svc.findByID(roleID)
		raProps.setRole(r)
		return err
	}()

	return r, svc.recordAction(svc.ctx, raProps, RoleActionLookup, err)
}

func (svc role) findByID(roleID uint64) (*types.Role, error) {
	if roleID == 0 {
		return nil, RoleErrInvalidID()
	}

	return store.LookupRoleByID(svc.ctx, svc.store, roleID)
}

func (svc role) FindByName(name string) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{role: &types.Role{Name: name}}
	)

	err = func() error {
		r, err = store.LookupRoleByName(svc.ctx, svc.store, name)
		raProps.setRole(r)
		return err
	}()

	return r, svc.recordAction(svc.ctx, raProps, RoleActionLookup, err)
}

func (svc role) FindByHandle(h string) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{role: &types.Role{Handle: h}}
	)

	err = func() error {
		r, err = store.LookupRoleByName(svc.ctx, svc.store, h)
		raProps.setRole(r)
		return err
	}()

	return r, svc.recordAction(svc.ctx, raProps, RoleActionLookup, err)
}

// FindByAny finds role by given identifier (id, handle, name)
func (svc role) FindByAny(ctx context.Context, identifier interface{}) (r *types.Role, err error) {
	if ID, ok := identifier.(uint64); ok {
		return svc.With(ctx).FindByID(ID)
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			return svc.With(ctx).FindByID(ID)
		} else {
			r, err = svc.With(ctx).FindByHandle(strIdentifier)
			if err == nil && r.ID == 0 {
				return svc.With(ctx).FindByName(strIdentifier)
			}

			return r, err
		}
	} else {
		return nil, RoleErrInvalidID()
	}
}

func (svc role) Create(new *types.Role) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{new: new}
	)

	err = func() (err error) {
		if !handle.IsValid(new.Handle) {
			return RoleErrInvalidHandle()
		}

		if !svc.ac.CanCreateRole(svc.ctx) {
			return RoleErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.RoleBeforeCreate(new, r)); err != nil {
			return
		}

		if err = svc.UniqueCheck(new); err != nil {
			return
		}

		new.ID = id.Next()
		new.CreatedAt = now()

		if err = store.CreateRole(svc.ctx, svc.store, new); err != nil {
			return
		}

		raProps.setRole(r)

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleAfterCreate(new, r))
		return
	}()

	return r, svc.recordAction(svc.ctx, raProps, RoleActionCreate, err)

}

func (svc role) Update(upd *types.Role) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{update: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return RoleErrInvalidID()
		}

		if !handle.IsValid(upd.Handle) {
			return RoleErrInvalidHandle()
		}

		if !svc.ac.CanUpdateRole(svc.ctx, upd) {
			return RoleErrNotAllowedToUpdate()
		}

		if r, err = store.LookupRoleByID(svc.ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setRole(r)

		if err = svc.eventbus.WaitFor(svc.ctx, event.RoleBeforeUpdate(upd, r)); err != nil {
			return
		}

		if err = svc.UniqueCheck(upd); err != nil {
			return
		}

		r.Handle = upd.Handle
		r.Name = upd.Name
		r.UpdatedAt = nowPtr()

		// Assign changed values
		if err = store.UpdateRole(svc.ctx, svc.store, r); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleAfterUpdate(upd, r))

		return nil
	}()

	return r, svc.recordAction(svc.ctx, raProps, RoleActionUpdate, err)
}

func (svc role) UniqueCheck(r *types.Role) (err error) {
	var (
		raProps = &roleActionProps{role: r}
	)

	if r.Handle != "" {
		if ex, _ := store.LookupRoleByHandle(svc.ctx, svc.store, r.Handle); ex != nil && ex.ID > 0 && ex.ID != r.ID {
			raProps.setExisting(ex)
			return RoleErrHandleNotUnique()
		}
	}

	if r.Name != "" {
		if ex, _ := store.LookupRoleByName(svc.ctx, svc.store, r.Name); ex != nil && ex.ID > 0 && ex.ID != r.ID {
			raProps.setExisting(ex)
			return RoleErrNameNotUnique()
		}
	}

	return nil
}

func (svc role) Delete(roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanDeleteRole(svc.ctx, r) {
			return RoleErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.RoleBeforeDelete(nil, r)); err != nil {
			return
		}

		r.DeletedAt = nowPtr()

		if err = store.UpdateRole(svc.ctx, svc.store, r); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleAfterDelete(nil, r))

		return
	}()

	return svc.recordAction(svc.ctx, raProps, RoleActionDelete, err)
}

func (svc role) Undelete(roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanDeleteRole(svc.ctx, r) {
			return RoleErrNotAllowedToDelete()
		}

		r.DeletedAt = nil

		if err = store.UpdateRole(svc.ctx, svc.store, r); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(svc.ctx, raProps, RoleActionUndelete, err)
}

func (svc role) Archive(roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanUpdateRole(svc.ctx, r) {
			return RoleErrNotAllowedToArchive()
		}

		r.ArchivedAt = nowPtr()
		if err = store.UpdateRole(svc.ctx, svc.store, r); err != nil {
			return
		}

		return
	}()

	return svc.recordAction(svc.ctx, raProps, RoleActionArchive, err)
}

func (svc role) Unarchive(roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanDeleteRole(svc.ctx, r) {
			return RoleErrNotAllowedToDelete()
		}

		r.ArchivedAt = nil
		if err = store.UpdateRole(svc.ctx, svc.store, r); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(svc.ctx, raProps, RoleActionUnarchive, err)
}

func (svc role) Membership(userID uint64) (types.RoleMemberSet, error) {
	mm, _, err := store.SearchRoleMembers(svc.ctx, svc.store, types.RoleMemberFilter{UserID: userID})
	return mm, err
}

func (svc role) MemberList(roleID uint64) (mm types.RoleMemberSet, err error) {
	var (
		r *types.Role

		raProps = &roleActionProps{
			role: &types.Role{ID: roleID},
		}
	)

	err = func() error {
		if roleID == permissions.EveryoneRoleID || roleID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		if !svc.ac.CanReadRole(svc.ctx, r) {
			return RoleErrNotAllowedToRead()
		}

		mm, _, err = store.SearchRoleMembers(svc.ctx, svc.store, types.RoleMemberFilter{RoleID: roleID})
		return err
	}()

	return mm, svc.recordAction(svc.ctx, raProps, RoleActionMembers, err)
}

// MemberAdd adds member (user) to a role
func (svc role) MemberAdd(roleID, memberID uint64) (err error) {
	var (
		r *types.Role
		m *types.User

		raProps = &roleActionProps{
			role:   &types.Role{ID: roleID},
			member: &types.User{ID: memberID},
		}
	)

	err = func() (err error) {
		if roleID == permissions.EveryoneRoleID || roleID == 0 || memberID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(roleID); err != nil {
			return
		}

		raProps.setRole(r)

		if m, err = svc.user.FindByID(memberID); err != nil {
			return
		}

		raProps.setMember(m)

		if err = svc.eventbus.WaitFor(svc.ctx, event.RoleMemberBeforeAdd(m, r)); err != nil {
			return
		}

		if !svc.ac.CanManageRoleMembers(svc.ctx, r) {
			return RoleErrNotAllowedToManageMembers()
		}

		if err = store.CreateRoleMember(svc.ctx, svc.store, &types.RoleMember{RoleID: r.ID, UserID: m.ID}); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleMemberAfterAdd(m, r))
		return nil
	}()

	return svc.recordAction(svc.ctx, raProps, RoleActionMemberAdd, err)
}

// MemberRemove removes member (user) from a role
func (svc role) MemberRemove(roleID, memberID uint64) (err error) {
	var (
		r       *types.Role
		m       *types.User
		raProps = &roleActionProps{
			role:   &types.Role{ID: roleID},
			member: &types.User{ID: memberID},
		}
	)

	err = func() (err error) {
		if roleID == permissions.EveryoneRoleID || roleID == 0 || memberID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(roleID); err != nil {
			return
		}

		raProps.setRole(r)

		if m, err = svc.user.FindByID(memberID); err != nil {
			return
		}

		raProps.setMember(m)

		if err = svc.eventbus.WaitFor(svc.ctx, event.RoleMemberBeforeRemove(m, r)); err != nil {
			return
		}

		if !svc.ac.CanManageRoleMembers(svc.ctx, r) {
			return RoleErrNotAllowedToManageMembers()
		}

		if err = store.DeleteRoleMember(svc.ctx, svc.store, &types.RoleMember{RoleID: r.ID, UserID: m.ID}); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleMemberAfterRemove(m, r))
		return nil
	}()

	return svc.recordAction(svc.ctx, raProps, RoleActionMemberRemove, err)
}
