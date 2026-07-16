package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
)

type (
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
		Search(context.Context, types.UserGroupFilter) (types.UserGroupSet, types.UserGroupFilter, error)

		Create(ctx context.Context, userGroup *types.UserGroup) (*types.UserGroup, error)
		Update(ctx context.Context, userGroup *types.UserGroup) (*types.UserGroup, error)

		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error

		MemberList(ctx context.Context, userGroupID uint64) (types.UserSet, error)
		MemberAdd(ctx context.Context, userGroupID, userID uint64) error
	}

	eventbusUserGroupChangeRegistry interface {
		Register(eventbus.HandlerFn, ...eventbus.HandlerRegOp) uintptr
	}

	rbacUserGroupService interface {
		UpdateUserGroups(rr ...rbac.GroupMembers) (err error)
		AssignGroupMembers(group id.ID, members ...id.ID) (err error)

		AddNode(id id.ID, handle string, paths ...rbac.GroupNodePath) (err error)
		UpdateNode(id id.ID, handle string, paths ...rbac.GroupNodePath) (err error)
		RemoveNode(id id.ID) (err error)
	}

	userGroupAuth interface {
		RemoveAccessTokens(context.Context, *types.User) error
	}

	userGroupServices struct {
		eventbus      eventDispatcher
		rbac          rbacUserGroupService
		user          UserService
		role          RoleService
		rootUserGroup id.ID
	}
)

func UserGroup(rbac rbacUserGroupService) *userGroup {
	return &userGroup{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		services: &userGroupServices{
			eventbus: eventbus.Service(),
			rbac:     rbac,
			user:     DefaultUser,
			role:     DefaultRole,
		},
	}
}

func (svc *userGroup) onActivate(ctx context.Context, _ *userGroupActionProps) (err error) {
	gMembers := []rbac.GroupMembers{}

	groups, _, err := svc.Search(ctx, types.UserGroupFilter{})
	if err != nil {
		return
	}

	for _, g := range groups {
		if len(g.Config.Paths) == 0 {
			svc.services.rootUserGroup = id.MustNumID(g.ID)
		}

		roles, _, err := svc.services.role.Find(ctx, types.RoleFilter{
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

		pp := []rbac.GroupNodePath{}
		for _, p := range g.Config.Paths {
			pp = append(pp, rbac.GroupNodePath{
				SelfID: id.MustNumID(p.SelfID),
				Name:   p.Name,
			})
		}

		// @todo we need to reload this on specific changes
		gMembers = append(gMembers, rbac.ConvUserGroup(id.MustNumID(g.ID), g.Handle, members, rr, pp))
	}

	err = svc.services.rbac.UpdateUserGroups(gMembers...)
	if err != nil {
		return
	}

	return
}

func (svc *userGroup) Search(ctx context.Context, filter types.UserGroupFilter) (rr types.UserGroupSet, f types.UserGroupFilter, err error) {
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

		for _, r := range rr {
			r.IsRoot = id.MustNumID(r.ID).Equal(svc.services.rootUserGroup)
		}

		if err = label.Load(ctx, svc.store, toLabeledUserGroups(rr)...); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, UserGroupActionSearch, err)
}

func (svc *userGroup) FindByID(ctx context.Context, userGroupID uint64) (r *types.UserGroup, err error) {
	var (
		raProps = &userGroupActionProps{userGroup: &types.UserGroup{ID: userGroupID}}
	)

	err = func() error {
		if r, err = svc.findByID(ctx, userGroupID); err != nil {
			return err
		}

		r.IsRoot = id.MustNumID(r.ID).Equal(svc.services.rootUserGroup)

		raProps.setUserGroup(r)
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, UserGroupActionLookup, err)
}

func (svc *userGroup) findByID(ctx context.Context, userGroupID uint64) (*types.UserGroup, error) {
	r, err := loadUserGroup(ctx, svc.store, userGroupID)
	return svc.proc(ctx, r, err)
}

func (svc *userGroup) FindByHandle(ctx context.Context, h string) (r *types.UserGroup, err error) {
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
func (svc *userGroup) FindByAny(ctx context.Context, identifier interface{}) (r *types.UserGroup, err error) {
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

func (svc *userGroup) proc(ctx context.Context, r *types.UserGroup, err error) (*types.UserGroup, error) {
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

// validate runs at the top of the generated Create, before the access check, on
// the incoming record. checkPaths additionally normalises new.Config (defaults
// it and is the reason it runs before the access check / event in the original).
func (svc *userGroup) validate(ctx context.Context, new *types.UserGroup) error {
	if !handle.IsValid(new.Handle) {
		return UserGroupErrInvalidHandle()
	}

	if !svc.checkPaths(new) {
		return UserGroupErrInvalidSelfID()
	}

	if !svc.isValidStructure(ctx, new) {
		return UserGroupErrInvalidUpdateStructure()
	}

	return nil
}

// beforeCreate runs after the access check and before the generated Create
// assigns the ID / timestamps and persists. It fires the (synchronous)
// before-create event and then enforces handle uniqueness, matching the original
// ordering exactly.
func (svc *userGroup) beforeCreate(ctx context.Context, new *types.UserGroup) error {
	if err := svc.services.eventbus.WaitFor(ctx, event.UserGroupBeforeCreate(new, nil)); err != nil {
		return err
	}

	return svc.UniqueCheck(ctx, new)
}

// afterCreate runs after the generated Create has persisted the record and its
// labels. It registers the rbac node and then dispatches the after-create event
// asynchronously (eventbus.Dispatch), preserving the original's
// "AddNode before AfterCreate" ordering and fire-and-forget semantics.
func (svc *userGroup) afterCreate(ctx context.Context, r *types.UserGroup) error {
	pp := []rbac.GroupNodePath{}
	for _, p := range r.Config.Paths {
		pp = append(pp, rbac.GroupNodePath{
			SelfID: id.MustNumID(p.SelfID),
			Name:   p.Name,
		})
	}

	if err := svc.services.rbac.AddNode(id.MustNumID(r.ID), r.Handle, pp...); err != nil {
		return err
	}

	svc.services.eventbus.Dispatch(ctx, event.UserGroupAfterCreate(r, r))
	return nil
}

func (svc *userGroup) Update(ctx context.Context, upd *types.UserGroup) (r *types.UserGroup, err error) {
	var (
		old     *types.UserGroup
		raProps = &userGroupActionProps{update: upd}
	)

	err = func() (err error) {
		if r, err = loadUserGroup(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = r.Clone()
		raProps.setUserGroup(r)

		if !handle.IsValid(upd.Handle) {
			return UserGroupErrInvalidHandle()
		}

		if !svc.ac.CanUpdateUserGroup(ctx, upd) {
			return UserGroupErrNotAllowedToUpdate()
		}

		if len(upd.Config.Paths) == 0 {
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

		if err = svc.services.eventbus.WaitFor(ctx, event.UserGroupBeforeUpdate(upd, r)); err != nil {
			return
		}

		if err = svc.UniqueCheck(ctx, upd); err != nil {
			return
		}

		r.Handle = upd.Handle
		r.Meta = upd.Meta
		r.UpdatedAt = now()
		r.Config = upd.Config

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

		pp := []rbac.GroupNodePath{}
		for _, p := range r.Config.Paths {
			pp = append(pp, rbac.GroupNodePath{
				SelfID: id.MustNumID(p.SelfID),
				Name:   p.Name,
			})
		}

		err = svc.services.rbac.UpdateNode(id.MustNumID(r.ID), r.Handle, pp...)
		if err != nil {
			return
		}

		svc.services.eventbus.Dispatch(ctx, event.UserGroupAfterUpdate(upd, r))

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, UserGroupActionUpdate, err, old, r)
}

func (svc *userGroup) UniqueCheck(ctx context.Context, r *types.UserGroup) (err error) {
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

func (svc *userGroup) DeleteByID(ctx context.Context, userGroupID uint64) (err error) {
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

		if err = svc.referencedByAuthClient(ctx, r.ID); err != nil {
			return err
		}

		if err = svc.services.eventbus.WaitFor(ctx, event.UserGroupBeforeDelete(nil, r)); err != nil {
			return
		}

		r.DeletedAt = now()

		if err = store.UpdateUserGroup(ctx, svc.store, r); err != nil {
			return
		}

		err = svc.services.rbac.RemoveNode(id.MustNumID(r.ID))
		if err != nil {
			return
		}

		svc.services.eventbus.Dispatch(ctx, event.UserGroupAfterDelete(nil, r))

		return
	}()

	return svc.recordAction(ctx, raProps, UserGroupActionDelete, err)
}

// referencedByAuthClient prevents deleting a user group that is used as the
// security user group on any auth client
func (svc *userGroup) referencedByAuthClient(ctx context.Context, userGroupID uint64) error {
	if userGroupID == 0 {
		return nil
	}

	cc, _, err := store.SearchAuthClients(ctx, svc.store, types.AuthClientFilter{})
	if err != nil {
		return err
	}

	for _, c := range cc {
		if c.Security != nil && c.Security.UserGroup == userGroupID {
			return UserGroupErrNotAllowedToDelete()
		}
	}

	return nil
}

func (svc *userGroup) UndeleteByID(ctx context.Context, userGroupID uint64) (err error) {
	var (
		r, upd  *types.UserGroup
		raProps = &userGroupActionProps{userGroup: &types.UserGroup{ID: userGroupID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(ctx, userGroupID); err != nil {
			return err
		}

		upd = r.Clone()
		if err = svc.services.eventbus.WaitFor(ctx, event.UserGroupBeforeUpdate(upd, r)); err != nil {
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

		pp := []rbac.GroupNodePath{}
		for _, p := range upd.Config.Paths {
			pp = append(pp, rbac.GroupNodePath{
				SelfID: id.MustNumID(p.SelfID),
				Name:   p.Name,
			})
		}

		err = svc.services.rbac.AddNode(id.MustNumID(upd.ID), upd.Handle, pp...)
		if err != nil {
			return
		}

		svc.services.eventbus.Dispatch(ctx, event.UserGroupAfterUpdate(upd, r))
		return nil
	}()

	return svc.recordAction(ctx, raProps, UserGroupActionUndelete, err, r, upd)
}

func (svc *userGroup) onMemberList(ctx context.Context, aProps *userGroupActionProps, userGroupID uint64) (mm types.UserSet, err error) {
	var (
		r *types.UserGroup
	)

	aProps.userGroup = &types.UserGroup{ID: userGroupID}

	if userGroupID == 0 {
		return nil, UserGroupErrInvalidID()
	}

	if r, err = svc.findByID(ctx, userGroupID); err != nil {
		return nil, err
	}

	if !svc.ac.CanReadUserGroup(ctx, r) {
		return nil, UserGroupErrNotAllowedToRead()
	}

	mm, _, err = store.SearchUsers(ctx, svc.store, types.UserFilter{
		UserGroupID: userGroupID,
	})

	return mm, err
}

// onMemberAdd adds member (user) to a userGroup
func (svc *userGroup) onMemberAdd(ctx context.Context, aProps *userGroupActionProps, userGroupID, memberID uint64) (err error) {
	var (
		g *types.UserGroup
		m *types.User
	)

	aProps.userGroup = &types.UserGroup{ID: userGroupID}
	aProps.member = &types.User{ID: memberID}

	if userGroupID == 0 || memberID == 0 {
		return UserGroupErrInvalidID()
	}

	if g, err = svc.findByID(ctx, userGroupID); err != nil {
		return
	}

	aProps.setUserGroup(g)

	if m, err = svc.services.user.FindByID(ctx, memberID); err != nil {
		return
	}

	aProps.setMember(m)

	m.UserGroupID = g.ID

	if err = svc.services.eventbus.WaitFor(ctx, event.UserGroupBeforeMemberAdd(g, g)); err != nil {
		return
	}

	if !svc.ac.CanManageMembersOnUserGroup(ctx, g) {
		return UserGroupErrNotAllowedToManageMembers()
	}

	if err = store.UpdateUser(ctx, svc.store, m); err != nil {
		return
	}

	err = svc.services.rbac.AssignGroupMembers(id.MustNumID(m.UserGroupID), id.MustNumID(m.ID))
	if err != nil {
		return
	}

	_ = svc.services.eventbus.WaitFor(ctx, event.UserGroupAfterMemberAdd(g, g))
	return nil
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

func (svc *userGroup) isValidStructure(ctx context.Context, g *types.UserGroup) (ok bool) {
	// @todo :)
	return true
}

func (svc *userGroup) checkSelfID(ctx context.Context, g *types.UserGroup) bool {
	for _, p := range g.Config.Paths {
		// Can't point to itself
		if p.SelfID == g.ID {
			return false
		}

		// The pointed to selfID exists
		_, err := svc.FindByID(ctx, p.SelfID)
		if err != nil {
			return false
		}
	}

	return true
}

func (svc *userGroup) checkPaths(g *types.UserGroup) (ok bool) {
	if id.MustNumID(g.ID).Equal(svc.services.rootUserGroup) {
		return len(g.Config.Paths) == 0
	}

	if g.Config == nil {
		g.Config = &types.UserGroupConfig{}
	}

	if len(g.Config.Paths) == 0 {
		return false
	}

	names := make(map[string]bool, len(g.Config.Paths)/2)
	for _, p := range g.Config.Paths {
		if p.SelfID == 0 {
			return false
		}

		// Do not allow duplicate paths
		if names[p.Name] {
			return false
		}

		names[p.Name] = true
	}

	return true
}
