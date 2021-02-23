package service

import (
	"context"
	"strconv"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	role struct {
		actionlog actionlog.Recorder

		ac       roleAccessController
		eventbus eventDispatcher

		user UserService

		store store.Storer
	}

	roleAccessController interface {
		CanCreateRole(context.Context) bool
		CanReadRole(context.Context, *types.Role) bool
		CanUpdateRole(context.Context, *types.Role) bool
		CanDeleteRole(context.Context, *types.Role) bool
		CanManageRoleMembers(context.Context, *types.Role) bool
	}

	RoleService interface {
		FindByID(ctx context.Context, roleID uint64) (*types.Role, error)
		FindByName(ctx context.Context, name string) (*types.Role, error)
		FindByHandle(ctx context.Context, handle string) (*types.Role, error)
		FindByAny(ctx context.Context, identifier interface{}) (*types.Role, error)
		Find(context.Context, types.RoleFilter) (types.RoleSet, types.RoleFilter, error)

		Create(ctx context.Context, role *types.Role) (*types.Role, error)
		Update(ctx context.Context, role *types.Role) (*types.Role, error)

		Archive(ctx context.Context, ID uint64) error
		Unarchive(ctx context.Context, ID uint64) error
		Delete(ctx context.Context, ID uint64) error
		Undelete(ctx context.Context, ID uint64) error

		Membership(ctx context.Context, userID uint64) (types.RoleMemberSet, error)
		MemberList(ctx context.Context, roleID uint64) (types.RoleMemberSet, error)
		MemberAdd(ctx context.Context, roleID, userID uint64) error
		MemberRemove(ctx context.Context, roleID, userID uint64) error
	}
)

func Role(ctx context.Context) RoleService {
	return (&role{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),

		actionlog: DefaultActionlog,

		user:  DefaultUser,
		store: DefaultStore,
	})
}

func (svc role) Find(ctx context.Context, filter types.RoleFilter) (rr types.RoleSet, f types.RoleFilter, err error) {
	var (
		raProps = &roleActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Role) (bool, error) {
		if !svc.ac.CanReadRole(ctx, res) {
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
			//if !svc.ac.CanAccess(ctx) {
			//	return RoleErrNotAllowedToListRoles()
			//}
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Role{}.LabelResourceKind(),
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

		if rr, f, err = store.SearchRoles(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledRoles(rr)...); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, RoleActionSearch, err)
}

func (svc role) FindByID(ctx context.Context, roleID uint64) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() error {
		if r, err = svc.findByID(ctx, roleID); err != nil {
			return err
		}

		raProps.setRole(r)
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, RoleActionLookup, err)
}

func (svc role) findByID(ctx context.Context, roleID uint64) (*types.Role, error) {
	if roleID == 0 {
		return nil, RoleErrInvalidID()
	}

	r, err := store.LookupRoleByID(ctx, svc.store, roleID)
	return svc.proc(ctx, r, err)
}

func (svc role) FindByName(ctx context.Context, name string) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{role: &types.Role{Name: name}}
	)

	err = func() error {
		r, err := store.LookupRoleByName(ctx, svc.store, name)
		if r, err = svc.proc(ctx, r, err); err != nil {
			return err
		}

		raProps.setRole(r)
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, RoleActionLookup, err)
}

func (svc role) FindByHandle(ctx context.Context, h string) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{role: &types.Role{Handle: h}}
	)

	err = func() error {
		r, err = store.LookupRoleByHandle(ctx, svc.store, h)
		if r, err = svc.proc(ctx, r, err); err != nil {
			return err
		}

		raProps.setRole(r)
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, RoleActionLookup, err)
}

// FindByAny finds role by given identifier (id, handle, name)
func (svc role) FindByAny(ctx context.Context, identifier interface{}) (r *types.Role, err error) {
	if ID, ok := identifier.(uint64); ok {
		return svc.FindByID(ctx, ID)
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			return svc.FindByID(ctx, ID)
		} else {
			r, err = svc.FindByHandle(ctx, strIdentifier)

			if (err == nil && r.ID == 0) || errors.IsNotFound(err) {
				return svc.FindByName(ctx, strIdentifier)
			}

			return r, err
		}
	} else {
		return nil, RoleErrInvalidID()
	}
}

func (svc role) proc(ctx context.Context, r *types.Role, err error) (*types.Role, error) {
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, RoleErrNotFound()
		}

		return nil, err
	}

	if err = label.Load(ctx, svc.store, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (svc role) Create(ctx context.Context, new *types.Role) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{new: new}
	)

	err = func() (err error) {
		if !handle.IsValid(new.Handle) {
			return RoleErrInvalidHandle()
		}

		if !svc.ac.CanCreateRole(ctx) {
			return RoleErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(ctx, event.RoleBeforeCreate(new, r)); err != nil {
			return
		}

		if err = svc.UniqueCheck(ctx, new); err != nil {
			return
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateRole(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		r = new

		_ = svc.eventbus.WaitFor(ctx, event.RoleAfterCreate(new, r))
		return
	}()

	return r, svc.recordAction(ctx, raProps, RoleActionCreate, err)

}

func (svc role) Update(ctx context.Context, upd *types.Role) (r *types.Role, err error) {
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

		if !svc.ac.CanUpdateRole(ctx, upd) {
			return RoleErrNotAllowedToUpdate()
		}

		if r, err = store.LookupRoleByID(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setRole(r)

		if err = svc.eventbus.WaitFor(ctx, event.RoleBeforeUpdate(upd, r)); err != nil {
			return
		}

		if err = svc.UniqueCheck(ctx, upd); err != nil {
			return
		}

		r.Handle = upd.Handle
		r.Name = upd.Name
		r.UpdatedAt = now()

		// Assign changed values
		if err = store.UpdateRole(ctx, svc.store, r); err != nil {
			return err
		}

		if label.Changed(r.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}

			r.Labels = upd.Labels
		}

		_ = svc.eventbus.WaitFor(ctx, event.RoleAfterUpdate(upd, r))

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, RoleActionUpdate, err)
}

func (svc role) UniqueCheck(ctx context.Context, r *types.Role) (err error) {
	var (
		raProps = &roleActionProps{role: r}
	)

	if r.Handle != "" {
		if ex, _ := store.LookupRoleByHandle(ctx, svc.store, r.Handle); ex != nil && ex.ID > 0 && ex.ID != r.ID {
			raProps.setExisting(ex)
			return RoleErrHandleNotUnique()
		}
	}

	if r.Name != "" {
		if ex, _ := store.LookupRoleByName(ctx, svc.store, r.Name); ex != nil && ex.ID > 0 && ex.ID != r.ID {
			raProps.setExisting(ex)
			return RoleErrNameNotUnique()
		}
	}

	return nil
}

func (svc role) Delete(ctx context.Context, roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(ctx, roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanDeleteRole(ctx, r) {
			return RoleErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(ctx, event.RoleBeforeDelete(nil, r)); err != nil {
			return
		}

		r.DeletedAt = now()

		if err = store.UpdateRole(ctx, svc.store, r); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.RoleAfterDelete(nil, r))

		return
	}()

	return svc.recordAction(ctx, raProps, RoleActionDelete, err)
}

func (svc role) Undelete(ctx context.Context, roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(ctx, roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanDeleteRole(ctx, r) {
			return RoleErrNotAllowedToDelete()
		}

		r.DeletedAt = nil

		if err = store.UpdateRole(ctx, svc.store, r); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, RoleActionUndelete, err)
}

func (svc role) Archive(ctx context.Context, roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(ctx, roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanUpdateRole(ctx, r) {
			return RoleErrNotAllowedToArchive()
		}

		r.ArchivedAt = now()
		if err = store.UpdateRole(ctx, svc.store, r); err != nil {
			return
		}

		return
	}()

	return svc.recordAction(ctx, raProps, RoleActionArchive, err)
}

func (svc role) Unarchive(ctx context.Context, roleID uint64) (err error) {
	var (
		r       *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(ctx, roleID); err != nil {
			return err
		}

		raProps.setRole(r)

		if !svc.ac.CanDeleteRole(ctx, r) {
			return RoleErrNotAllowedToDelete()
		}

		r.ArchivedAt = nil
		if err = store.UpdateRole(ctx, svc.store, r); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, RoleActionUnarchive, err)
}

func (svc role) Membership(ctx context.Context, userID uint64) (types.RoleMemberSet, error) {
	mm, _, err := store.SearchRoleMembers(ctx, svc.store, types.RoleMemberFilter{UserID: userID})
	return mm, err
}

func (svc role) MemberList(ctx context.Context, roleID uint64) (mm types.RoleMemberSet, err error) {
	var (
		r *types.Role

		raProps = &roleActionProps{
			role: &types.Role{ID: roleID},
		}
	)

	err = func() error {
		if roleID == rbac.EveryoneRoleID || roleID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(ctx, roleID); err != nil {
			return err
		}

		if !svc.ac.CanReadRole(ctx, r) {
			return RoleErrNotAllowedToRead()
		}

		mm, _, err = store.SearchRoleMembers(ctx, svc.store, types.RoleMemberFilter{RoleID: roleID})
		return err
	}()

	return mm, svc.recordAction(ctx, raProps, RoleActionMembers, err)
}

// MemberAdd adds member (user) to a role
func (svc role) MemberAdd(ctx context.Context, roleID, memberID uint64) (err error) {
	var (
		r *types.Role
		m *types.User

		raProps = &roleActionProps{
			role:   &types.Role{ID: roleID},
			member: &types.User{ID: memberID},
		}
	)

	err = func() (err error) {
		if roleID == rbac.EveryoneRoleID || roleID == 0 || memberID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(ctx, roleID); err != nil {
			return
		}

		raProps.setRole(r)

		if m, err = svc.user.FindByID(ctx, memberID); err != nil {
			return
		}

		raProps.setMember(m)

		if err = svc.eventbus.WaitFor(ctx, event.RoleMemberBeforeAdd(m, r)); err != nil {
			return
		}

		if !svc.ac.CanManageRoleMembers(ctx, r) {
			return RoleErrNotAllowedToManageMembers()
		}

		if err = store.CreateRoleMember(ctx, svc.store, &types.RoleMember{RoleID: r.ID, UserID: m.ID}); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.RoleMemberAfterAdd(m, r))
		return nil
	}()

	return svc.recordAction(ctx, raProps, RoleActionMemberAdd, err)
}

// MemberRemove removes member (user) from a role
func (svc role) MemberRemove(ctx context.Context, roleID, memberID uint64) (err error) {
	var (
		r       *types.Role
		m       *types.User
		raProps = &roleActionProps{
			role:   &types.Role{ID: roleID},
			member: &types.User{ID: memberID},
		}
	)

	err = func() (err error) {
		if roleID == rbac.EveryoneRoleID || roleID == 0 || memberID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(ctx, roleID); err != nil {
			return
		}

		raProps.setRole(r)

		if m, err = svc.user.FindByID(ctx, memberID); err != nil {
			return
		}

		raProps.setMember(m)

		if err = svc.eventbus.WaitFor(ctx, event.RoleMemberBeforeRemove(m, r)); err != nil {
			return
		}

		if !svc.ac.CanManageRoleMembers(ctx, r) {
			return RoleErrNotAllowedToManageMembers()
		}

		if err = store.DeleteRoleMember(ctx, svc.store, &types.RoleMember{RoleID: r.ID, UserID: m.ID}); err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.RoleMemberAfterRemove(m, r))
		return nil
	}()

	return svc.recordAction(ctx, raProps, RoleActionMemberRemove, err)
}

// toLabeledRoles converts to []label.LabeledResource
//
// This function is auto-generated.
func toLabeledRoles(set []*types.Role) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
