package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	role struct {
		db  *factory.DB
		ctx context.Context

		actionlog actionlog.Recorder

		ac       roleAccessController
		eventbus eventDispatcher

		user UserService

		role repository.RoleRepository
	}

	roleAccessController interface {
		CanAccess(context.Context) bool

		CanCreateRole(context.Context) bool
		CanReadRole(context.Context, *types.Role) bool
		CanUpdateRole(context.Context, *types.Role) bool
		CanDeleteRole(context.Context, *types.Role) bool
		CanManageRoleMembers(context.Context, *types.Role) bool

		FilterReadableRoles(ctx context.Context) *permissions.ResourceFilter
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
		Merge(roleID, targetroleID uint64) error
		Move(roleID, organisationID uint64) error

		Archive(ID uint64) error
		Unarchive(ID uint64) error
		Delete(ID uint64) error
		Undelete(ID uint64) error

		Membership(userID uint64) ([]*types.RoleMember, error)
		MemberList(roleID uint64) ([]*types.RoleMember, error)
		MemberAdd(roleID, userID uint64) error
		MemberRemove(roleID, userID uint64) error
	}
)

func Role(ctx context.Context) RoleService {
	return (&role{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),

		actionlog: DefaultActionlog,

		user: DefaultUser.With(ctx),
	}).With(ctx)
}

func (svc role) With(ctx context.Context) RoleService {
	db := repository.DB(ctx)
	return &role{
		db:  db,
		ctx: ctx,

		actionlog: svc.actionlog,

		ac:       svc.ac,
		eventbus: svc.eventbus,
		user:     svc.user,

		role: repository.Role(ctx, db),
	}
}

func (svc role) Find(filter types.RoleFilter) (rr types.RoleSet, f types.RoleFilter, err error) {
	var (
		raProps = &roleActionProps{filter: &filter}
	)

	err = svc.db.Transaction(func() error {
		filter.IsReadable = svc.ac.FilterReadableRoles(svc.ctx)

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

		rr, f, err = svc.role.Find(filter)
		return err
	})

	return rr, f, svc.recordAction(svc.ctx, raProps, RoleActionSearch, err)
}

func (svc role) FindByID(roleID uint64) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = svc.db.Transaction(func() error {
		r, err = svc.findByID(roleID)
		raProps.setRole(r)
		return err
	})

	return r, svc.recordAction(svc.ctx, raProps, RoleActionLookup, err)
}

func (svc role) findByID(roleID uint64) (*types.Role, error) {
	if roleID == 0 {
		return nil, RoleErrInvalidID()
	}

	return svc.role.FindByID(roleID)
}

func (svc role) FindByName(name string) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{role: &types.Role{Name: name}}
	)

	err = svc.db.Transaction(func() error {
		r, err = svc.role.FindByName(name)
		raProps.setRole(r)
		return err
	})

	return r, svc.recordAction(svc.ctx, raProps, RoleActionLookup, err)
}

func (svc role) FindByHandle(h string) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{role: &types.Role{Handle: h}}
	)

	err = svc.db.Transaction(func() error {
		r, err = svc.role.FindByName(h)
		raProps.setRole(r)
		return err
	})

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

	err = svc.db.Transaction(func() (err error) {
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

		if r, err = svc.role.Create(new); err != nil {
			return
		}

		raProps.setRole(r)

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleAfterCreate(new, r))
		return
	})

	return r, svc.recordAction(svc.ctx, raProps, RoleActionCreate, err)

}

func (svc role) Update(upd *types.Role) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{update: upd}
	)

	err = svc.db.Transaction(func() (err error) {
		if upd.ID == 0 {
			return RoleErrInvalidID()
		}

		if !handle.IsValid(upd.Handle) {
			return RoleErrInvalidHandle()
		}

		if !svc.ac.CanUpdateRole(svc.ctx, upd) {
			return RoleErrNotAllowedToUpdate()
		}

		if r, err = svc.role.FindByID(upd.ID); err != nil {
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

		// Assign changed values
		if r, err = svc.role.Update(r); err != nil {
			return err
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleAfterUpdate(upd, r))

		return nil
	})

	return r, svc.recordAction(svc.ctx, raProps, RoleActionUpdate, err)
}

func (svc role) UniqueCheck(r *types.Role) (err error) {
	var (
		raProps = &roleActionProps{role: r}
	)

	if r.Handle != "" {
		if ex, _ := svc.role.FindByHandle(r.Handle); ex != nil && ex.ID > 0 && ex.ID != r.ID {
			raProps.setExisting(ex)
			return RoleErrHandleNotUnique()
		}
	}

	if r.Name != "" {
		if ex, _ := svc.role.FindByName(r.Name); ex != nil && ex.ID > 0 && ex.ID != r.ID {
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

	err = svc.db.Transaction(func() (err error) {
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

		if err = svc.role.DeleteByID(roleID); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleAfterDelete(nil, r))

		return
	})

	return svc.recordAction(svc.ctx, raProps, RoleActionDelete, err)
}

func (svc role) Undelete(roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanDeleteRole(svc.ctx, r) {
			return RoleErrNotAllowedToDelete()
		}

		if err = svc.role.UndeleteByID(roleID); err != nil {
			return
		}

		return nil
	})

	return svc.recordAction(svc.ctx, raProps, RoleActionUndelete, err)
}

func (svc role) Archive(roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanUpdateRole(svc.ctx, r) {
			return RoleErrNotAllowedToArchive()
		}

		if err = svc.role.ArchiveByID(roleID); err != nil {
			return
		}

		return
	})

	return svc.recordAction(svc.ctx, raProps, RoleActionArchive, err)
}

func (svc role) Unarchive(roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = svc.db.Transaction(func() (err error) {
		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanDeleteRole(svc.ctx, r) {
			return RoleErrNotAllowedToDelete()
		}

		if err = svc.role.UndeleteByID(roleID); err != nil {
			return
		}

		return nil
	})

	return svc.recordAction(svc.ctx, raProps, RoleActionUnarchive, err)
}

func (svc role) Merge(roleID, targetRoleID uint64) (err error) {
	var (
		r *types.Role
		t *types.Role

		raProps = &roleActionProps{
			role:   &types.Role{ID: roleID},
			target: &types.Role{ID: targetRoleID},
		}
	)

	err = svc.db.Transaction(func() (err error) {
		if roleID == 0 || targetRoleID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanUpdateRole(svc.ctx, r) {
			return RoleErrNotAllowedToUpdate()
		}

		if t, err = svc.findByID(targetRoleID); err != nil {
			return err
		}

		raProps.setTarget(t)

		if !svc.ac.CanUpdateRole(svc.ctx, t) {
			return RoleErrNotAllowedToUpdate()
		}

		if err = svc.role.MergeByID(roleID, targetRoleID); err != nil {
			return
		}

		return nil
	})

	return svc.recordAction(svc.ctx, raProps, RoleActionMerge, err)
}

// Move
//
// @obsolete
func (svc role) Move(roleID, targetOrganisationID uint64) error {
	return RoleErrGeneric().Wrap(fmt.Errorf("obsolete"))
}

func (svc role) Membership(userID uint64) ([]*types.RoleMember, error) {
	return svc.role.MembershipsFindByUserID(userID)
}

func (svc role) MemberList(roleID uint64) (mm []*types.RoleMember, err error) {
	var (
		r *types.Role

		raProps = &roleActionProps{
			role: &types.Role{ID: roleID},
		}
	)

	err = svc.db.Transaction(func() error {
		if roleID == permissions.EveryoneRoleID || roleID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(roleID); err != nil {
			return err
		}

		if !svc.ac.CanReadRole(svc.ctx, r) {
			return RoleErrNotAllowedToRead()
		}

		if mm, err = svc.role.MemberFindByRoleID(roleID); err != nil {
			return err
		}

		return nil
	})

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

	err = svc.db.Transaction(func() (err error) {
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

		if err = svc.role.MemberAddByID(r.ID, m.ID); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleMemberAfterAdd(m, r))
		return nil
	})

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

	err = svc.db.Transaction(func() (err error) {
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

		if err = svc.role.MemberRemoveByID(r.ID, m.ID); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(svc.ctx, event.RoleMemberAfterRemove(m, r))
		return nil
	})

	return svc.recordAction(svc.ctx, raProps, RoleActionMemberRemove, err)
}

var _ RoleService = &role{}
