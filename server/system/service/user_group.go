package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/label"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/service/event"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	userGroup struct {
		actionlog actionlog.Recorder

		ac       userGroupAccessController
		eventbus eventDispatcher
		rbac     rbacUserGroupService

		user UserService
		role RoleService

		store store.Storer
	}

	userGroupAccessController interface {
		CanGrant(context.Context) bool

		CanSearchUserGroups(context.Context) bool
		CanCreateUserGroup(context.Context) bool
		CanReadUserGroup(context.Context, *types.UserGroup) bool
		CanUpdateUserGroup(context.Context, *types.UserGroup) bool
		CanDeleteUserGroup(context.Context, *types.UserGroup) bool
		CanManageMembersOnUserGroup(context.Context, *types.UserGroup) bool
	}

	UserGroupService interface {
		FindByID(ctx context.Context, userGroupID uint64) (*types.UserGroup, error)
		FindByHandle(ctx context.Context, handle string) (*types.UserGroup, error)
		Find(context.Context, types.UserGroupFilter) (types.UserGroupSet, types.UserGroupFilter, error)

		Create(ctx context.Context, userGroup *types.UserGroup) (*types.UserGroup, error)
		Update(ctx context.Context, userGroup *types.UserGroup) (*types.UserGroup, error)

		Delete(ctx context.Context, ID uint64) error
		Undelete(ctx context.Context, ID uint64) error

		MemberList(ctx context.Context, userGroupID uint64) (types.UserSet, error)
		MemberAdd(ctx context.Context, userGroupID, userID uint64) error
	}

	eventbusUserGroupChangeRegistry interface {
		Register(eventbus.HandlerFn, ...eventbus.HandlerRegOp) uintptr
	}

	rbacUserGroupService interface {
		UpdateUserGroups(rr ...rbac.GroupMembers) (err error)
		AssignGroupMembers(group id.ID, members ...id.ID) (err error)

		AddNode(id id.ID, handle string, selfID id.ID) (err error)
		UpdateNode(id id.ID, handle string, selfID id.ID) (err error)
		RemoveNode(id id.ID) (err error)
	}

	userGroupAuth interface {
		RemoveAccessTokens(context.Context, *types.User) error
	}
)

func UserGroup(rbac rbacUserGroupService) *userGroup {
	return &userGroup{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
		rbac:     rbac,

		actionlog: DefaultActionlog,

		user:  DefaultUser,
		role:  DefaultRole,
		store: DefaultStore,
	}
}

func (svc userGroup) Activate(ctx context.Context) (err error) {
	gMembers := []rbac.GroupMembers{}

	groups, _, err := svc.Find(ctx, types.UserGroupFilter{})
	if err != nil {
		return
	}

	for _, g := range groups {
		roles, _, err := svc.role.Find(ctx, types.RoleFilter{
			Resource: fmt.Sprintf("corteza::system:user-group/%d", g.ID),
		})
		if err != nil {
			return err
		}

		var mm types.UserSet
		mm, err = svc.MemberList(ctx, g.ID)
		if err != nil {
			return err
		}

		var rr []id.ID
		for _, r := range roles {
			rr = append(rr, id.MustNumID(r.ID))
		}

		members := make([]id.ID, len(mm))
		for i, m := range mm {
			members[i] = id.MustNumID(m.ID)
		}

		// @todo we need to reload this on specific changes
		gMembers = append(gMembers, rbac.ConvUserGroup(id.MustNumID(g.ID), g.Handle, id.MustNumID(g.SelfID), members, rr))
	}

	err = svc.rbac.UpdateUserGroups(gMembers...)
	if err != nil {
		return
	}

	return
}

func (svc userGroup) Find(ctx context.Context, filter types.UserGroupFilter) (rr types.UserGroupSet, f types.UserGroupFilter, err error) {
	var (
		raProps = &userGroupActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.UserGroup) (bool, error) {
		if !svc.ac.CanReadUserGroup(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchUserGroups(ctx) {
			return UserGroupErrNotAllowedToSearch()
		}

		if filter.Deleted > 0 {
			// If list with deleted or suspended users is requested
			// user must have access permissions to system (ie: is admin)
			//
			// not the best solution but ATM it allows us to have at least
			// some kind of control over who can see deleted or archived userGroups
			//if !svc.ac.CanAccess(ctx) {
			//	return UserGroupErrNotAllowedToListUserGroups()
			//}
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.UserGroup{}.LabelResourceKind(),
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

		if rr, f, err = store.SearchUserGroups(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledUserGroups(rr)...); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, UserGroupActionSearch, err)
}

func (svc userGroup) FindByID(ctx context.Context, userGroupID uint64) (r *types.UserGroup, err error) {
	var (
		raProps = &userGroupActionProps{userGroup: &types.UserGroup{ID: userGroupID}}
	)

	err = func() error {
		if r, err = svc.findByID(ctx, userGroupID); err != nil {
			return err
		}

		raProps.setUserGroup(r)
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, UserGroupActionLookup, err)
}

func (svc userGroup) findByID(ctx context.Context, userGroupID uint64) (*types.UserGroup, error) {
	r, err := loadUserGroup(ctx, svc.store, userGroupID)
	return svc.proc(ctx, r, err)
}

func (svc userGroup) FindByHandle(ctx context.Context, h string) (r *types.UserGroup, err error) {
	var (
		raProps = &userGroupActionProps{userGroup: &types.UserGroup{Handle: h}}
	)

	err = func() error {
		r, err = store.LookupUserGroupByHandle(ctx, svc.store, h)
		if r, err = svc.proc(ctx, r, err); err != nil {
			return err
		}

		raProps.setUserGroup(r)
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, UserGroupActionLookup, err)
}

// FindByAny finds userGroup by given identifier (id, handle, name)
func (svc userGroup) FindByAny(ctx context.Context, identifier interface{}) (r *types.UserGroup, err error) {
	if ID, ok := identifier.(uint64); ok {
		return svc.FindByID(ctx, ID)
	} else if strIdentifier, ok := identifier.(string); ok {
		if ID, _ := strconv.ParseUint(strIdentifier, 10, 64); ID > 0 {
			return svc.FindByID(ctx, ID)
		} else {
			r, err = svc.FindByHandle(ctx, strIdentifier)
			return r, err
		}
	} else {
		return nil, UserGroupErrInvalidID()
	}
}

func (svc userGroup) proc(ctx context.Context, r *types.UserGroup, err error) (*types.UserGroup, error) {
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, UserGroupErrNotFound()
		}

		return nil, err
	}

	if err = label.Load(ctx, svc.store, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (svc userGroup) Create(ctx context.Context, new *types.UserGroup) (r *types.UserGroup, err error) {
	var (
		raProps = &userGroupActionProps{userGroup: new}
	)

	err = func() (err error) {
		if !handle.IsValid(new.Handle) {
			return UserGroupErrInvalidHandle()
		}

		// @todo !!!

		// if new.SelfID == 0 {
		// 	return UserGroupErrMissingSelfID()
		// }

		// if !svc.checkSelfID(ctx, new) {
		// 	return UserGroupErrInvalidSelfID()
		// }

		if !svc.isValidStructure(ctx, new) {
			return UserGroupErrInvalidUpdateStructure()
		}

		if !svc.ac.CanCreateUserGroup(ctx) {
			return UserGroupErrNotAllowedToCreate()
		}

		if err = svc.eventbus.WaitFor(ctx, event.UserGroupBeforeCreate(new, r)); err != nil {
			return
		}

		if err = svc.UniqueCheck(ctx, new); err != nil {
			return
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		raProps.setNew(new)

		if err = store.CreateUserGroup(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		r = new

		err = svc.rbac.AddNode(id.MustNumID(new.ID), new.Handle, id.MustNumID(new.SelfID))
		if err != nil {
			return
		}

		svc.eventbus.Dispatch(ctx, event.UserGroupAfterCreate(new, r))
		return
	}()

	return r, svc.recordAction(ctx, raProps, UserGroupActionCreate, err)

}

func (svc userGroup) Update(ctx context.Context, upd *types.UserGroup) (r *types.UserGroup, err error) {
	var (
		raProps = &userGroupActionProps{update: upd}
	)

	err = func() (err error) {
		if r, err = loadUserGroup(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setUserGroup(r)

		if !handle.IsValid(upd.Handle) {
			return UserGroupErrInvalidHandle()
		}

		if !svc.ac.CanUpdateUserGroup(ctx, upd) {
			return UserGroupErrNotAllowedToUpdate()
		}

		if upd.SelfID == 0 {
			return UserGroupErrMissingSelfID()
		}

		if !svc.checkSelfID(ctx, upd) {
			return UserGroupErrInvalidSelfID()
		}

		if !svc.isValidStructure(ctx, upd) {
			return UserGroupErrInvalidUpdateStructure()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, r.UpdatedAt, r.CreatedAt) {
			return UserGroupErrStaleData()
		}

		if err = svc.eventbus.WaitFor(ctx, event.UserGroupBeforeUpdate(upd, r)); err != nil {
			return
		}

		if err = svc.UniqueCheck(ctx, upd); err != nil {
			return
		}

		r.Handle = upd.Handle
		r.Meta = upd.Meta
		r.UpdatedAt = now()

		// Assign changed values
		if err = store.UpdateUserGroup(ctx, svc.store, r); err != nil {
			return err
		}

		if label.Changed(r.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}

			r.Labels = upd.Labels
		}

		err = svc.rbac.UpdateNode(id.MustNumID(r.ID), r.Handle, id.MustNumID(r.SelfID))
		if err != nil {
			return
		}

		svc.eventbus.Dispatch(ctx, event.UserGroupAfterUpdate(upd, r))

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, UserGroupActionUpdate, err)
}

func (svc userGroup) UniqueCheck(ctx context.Context, r *types.UserGroup) (err error) {
	var (
		raProps = &userGroupActionProps{userGroup: r}
	)

	if r.Handle != "" {
		if ex, _ := store.LookupUserGroupByHandle(ctx, svc.store, r.Handle); ex != nil && ex.ID > 0 && ex.ID != r.ID {
			raProps.setExisting(ex)
			return UserGroupErrHandleNotUnique()
		}
	}

	return nil
}

func (svc userGroup) Delete(ctx context.Context, userGroupID uint64) (err error) {
	var (
		r       *types.UserGroup
		raProps = &userGroupActionProps{userGroup: &types.UserGroup{ID: userGroupID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(ctx, userGroupID); err != nil {
			return err
		}

		raProps.setUserGroup(r)

		if !svc.ac.CanDeleteUserGroup(ctx, r) {
			return UserGroupErrNotAllowedToDelete()
		}

		if err = svc.eventbus.WaitFor(ctx, event.UserGroupBeforeDelete(nil, r)); err != nil {
			return
		}

		r.DeletedAt = now()

		if err = store.UpdateUserGroup(ctx, svc.store, r); err != nil {
			return
		}

		err = svc.rbac.RemoveNode(id.MustNumID(r.ID))
		if err != nil {
			return
		}

		svc.eventbus.Dispatch(ctx, event.UserGroupAfterDelete(nil, r))

		return
	}()

	return svc.recordAction(ctx, raProps, UserGroupActionDelete, err)
}

func (svc userGroup) Undelete(ctx context.Context, userGroupID uint64) (err error) {
	var (
		r, upd  *types.UserGroup
		raProps = &userGroupActionProps{userGroup: &types.UserGroup{ID: userGroupID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(ctx, userGroupID); err != nil {
			return err
		}

		upd = r.Clone()
		if err = svc.eventbus.WaitFor(ctx, event.UserGroupBeforeUpdate(upd, r)); err != nil {
			return
		}

		raProps.setUserGroup(upd)

		if !svc.ac.CanDeleteUserGroup(ctx, upd) {
			return UserGroupErrNotAllowedToDelete()
		}

		upd.DeletedAt = nil
		if err = store.UpdateUserGroup(ctx, svc.store, upd); err != nil {
			return
		}

		err = svc.rbac.AddNode(id.MustNumID(upd.ID), upd.Handle, id.MustNumID(upd.SelfID))
		if err != nil {
			return
		}

		svc.eventbus.Dispatch(ctx, event.UserGroupAfterUpdate(upd, r))
		return nil
	}()

	return svc.recordAction(ctx, raProps, UserGroupActionUndelete, err)
}

func (svc userGroup) MemberList(ctx context.Context, userGroupID uint64) (mm types.UserSet, err error) {
	var (
		r *types.UserGroup

		raProps = &userGroupActionProps{
			userGroup: &types.UserGroup{ID: userGroupID},
		}
	)

	err = func() error {
		if userGroupID == 0 {
			return UserGroupErrInvalidID()
		}

		if r, err = svc.findByID(ctx, userGroupID); err != nil {
			return err
		}

		if !svc.ac.CanReadUserGroup(ctx, r) {
			return UserGroupErrNotAllowedToRead()
		}

		mm, _, err = store.SearchUsers(ctx, svc.store, types.UserFilter{
			UserGroupID: userGroupID,
		})

		return err
	}()

	return mm, svc.recordAction(ctx, raProps, UserGroupActionMembers, err)
}

// MemberAdd adds member (user) to a userGroup
func (svc userGroup) MemberAdd(ctx context.Context, userGroupID, memberID uint64) (err error) {
	var (
		g *types.UserGroup
		m *types.User

		raProps = &userGroupActionProps{
			userGroup: &types.UserGroup{ID: userGroupID},
			member:    &types.User{ID: memberID},
		}
	)

	err = func() (err error) {
		if userGroupID == 0 || memberID == 0 {
			return UserGroupErrInvalidID()
		}

		if g, err = svc.findByID(ctx, userGroupID); err != nil {
			return
		}

		raProps.setUserGroup(g)

		if m, err = svc.user.FindByID(ctx, memberID); err != nil {
			return
		}

		raProps.setMember(m)

		m.UserGroupID = g.ID

		if err = svc.eventbus.WaitFor(ctx, event.UserGroupBeforeMemberAdd(g, g)); err != nil {
			return
		}

		if !svc.ac.CanManageMembersOnUserGroup(ctx, g) {
			return UserGroupErrNotAllowedToManageMembers()
		}

		if err = store.UpdateUser(ctx, svc.store, m); err != nil {
			return
		}

		err = svc.rbac.AssignGroupMembers(id.MustNumID(m.UserGroupID), id.MustNumID(m.ID))
		if err != nil {
			return
		}

		_ = svc.eventbus.WaitFor(ctx, event.UserGroupAfterMemberAdd(g, g))
		return nil
	}()

	return svc.recordAction(ctx, raProps, UserGroupActionMemberAdd, err)
}

func loadUserGroup(ctx context.Context, s store.UserGroups, ID uint64) (res *types.UserGroup, err error) {
	if ID == 0 {
		return nil, UserGroupErrInvalidID()
	}

	if res, err = store.LookupUserGroupByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, UserGroupErrNotFound()
	}

	return
}

// toLabeledUserGroups converts to []label.LabeledResource
//
// This function is auto-generated.
func toLabeledUserGroups(set []*types.UserGroup) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

func (svc userGroup) isValidStructure(ctx context.Context, g *types.UserGroup) (ok bool) {
	// @todo :)
	return true
}

func (svc userGroup) checkSelfID(ctx context.Context, g *types.UserGroup) bool {
	if g.SelfID == g.ID {
		return false
	}

	_, err := svc.FindByID(ctx, g.SelfID)
	return err == nil
}
