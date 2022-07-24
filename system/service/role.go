package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

type (
	role struct {
		actionlog actionlog.Recorder

		ac       roleAccessController
		eventbus eventDispatcher
		rbac     rbacRuleService

		user UserService
		auth roleAuth

		store store.Storer

		// list of all system roles
		system map[string]bool

		// list of all closed roles
		closed map[string]bool
	}

	roleAccessController interface {
		CanGrant(context.Context) bool

		CanSearchRoles(context.Context) bool
		CanCreateRole(context.Context) bool
		CanReadRole(context.Context, *types.Role) bool
		CanUpdateRole(context.Context, *types.Role) bool
		CanDeleteRole(context.Context, *types.Role) bool
		CanManageMembersOnRole(context.Context, *types.Role) bool
	}

	RoleService interface {
		FindByID(ctx context.Context, roleID uint64) (*types.Role, error)
		FindByName(ctx context.Context, name string) (*types.Role, error)
		FindByHandle(ctx context.Context, handle string) (*types.Role, error)
		FindByAny(ctx context.Context, identifier interface{}) (*types.Role, error)
		Find(context.Context, types.RoleFilter) (types.RoleSet, types.RoleFilter, error)

		IsSystem(r *types.Role) bool
		IsClosed(r *types.Role) bool

		Create(ctx context.Context, role *types.Role) (*types.Role, error)
		Update(ctx context.Context, role *types.Role) (*types.Role, error)

		Archive(ctx context.Context, ID uint64) error
		Unarchive(ctx context.Context, ID uint64) error
		Delete(ctx context.Context, ID uint64) error
		Undelete(ctx context.Context, ID uint64) error
		CloneRules(ctx context.Context, ID uint64, cloneToRoleID ...uint64) error

		Membership(ctx context.Context, userID uint64) (types.RoleMemberSet, error)
		MemberList(ctx context.Context, roleID uint64) (types.RoleMemberSet, error)
		MemberAdd(ctx context.Context, roleID, userID uint64) error
		MemberRemove(ctx context.Context, roleID, userID uint64) error
	}

	eventbusRoleChangeRegistry interface {
		Register(eventbus.HandlerFn, ...eventbus.HandlerRegOp) uintptr
	}

	rbacRoleUpdater interface {
		UpdateRoles(rr ...*rbac.Role)
	}

	rbacRuleService interface {
		CloneRulesByRoleID(ctx context.Context, roleID uint64, toRoleID ...uint64) error
	}

	roleAuth interface {
		RemoveAccessTokens(context.Context, *types.User) error
	}
)

func Role(rbac rbacRuleService) *role {
	return &role{
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
		rbac:     rbac,

		actionlog: DefaultActionlog,

		user:  DefaultUser,
		auth:  DefaultAuth,
		store: DefaultStore,

		system: make(map[string]bool),
		closed: make(map[string]bool),
	}
}

// SetSystem sets list of handles for all system roles
//
// System roles can not be changed or deleted
func (svc *role) SetSystem(hh ...string) {
	svc.system = slice.ToStringBoolMap(hh)
	delete(svc.system, "")
}

func (svc role) IsSystem(r *types.Role) bool {
	return len(r.Handle) > 0 && svc.system[r.Handle]
}

// SetClosed sets list of handles for all closed roles
//
// Closed roles can not have members
func (svc *role) SetClosed(hh ...string) {
	svc.closed = slice.ToStringBoolMap(hh)
	delete(svc.closed, "")
}

func (svc role) IsClosed(r *types.Role) bool {
	return len(r.Handle) > 0 && svc.closed[r.Handle]
}

func (svc role) IsContextual(r *types.Role) bool {
	return r.Meta != nil && r.Meta.Context != nil && len(r.Meta.Context.Expr) > 0
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
		if !svc.ac.CanSearchRoles(ctx) {
			return RoleErrNotAllowedToSearch()
		}

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
	r, err := loadRole(ctx, svc.store, roleID)
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

		if new.Meta != nil && new.Meta.Context != nil {
			if err = svc.validateContext(ctx, new.Meta.Context); err != nil {
				return
			}
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

		svc.eventbus.Dispatch(ctx, event.RoleAfterCreate(new, r))
		return
	}()

	return r, svc.recordAction(ctx, raProps, RoleActionCreate, err)

}

func (svc role) Update(ctx context.Context, upd *types.Role) (r *types.Role, err error) {
	var (
		raProps = &roleActionProps{update: upd}
	)

	err = func() (err error) {
		if !handle.IsValid(upd.Handle) {
			return RoleErrInvalidHandle()
		}

		if !svc.ac.CanUpdateRole(ctx, upd) {
			return RoleErrNotAllowedToUpdate()
		}

		if r, err = loadRole(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setRole(r)

		if svc.IsSystem(r) {
			// prevent system role updates
			// we need this here because of the clumsy way
			// how rest endpoint handler is implemented ATM
			if r.Handle == upd.Handle && r.Name == upd.Name {
				// no change.
				return nil
			}
			return RoleErrNotAllowedToUpdate()
		}

		if err = svc.eventbus.WaitFor(ctx, event.RoleBeforeUpdate(upd, r)); err != nil {
			return
		}

		if upd.Meta != nil && upd.Meta.Context != nil {
			if err = svc.validateContext(ctx, upd.Meta.Context); err != nil {
				return
			}
		}

		if err = svc.UniqueCheck(ctx, upd); err != nil {
			return
		}

		r.Handle = upd.Handle
		r.Name = upd.Name
		r.Meta = upd.Meta
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

		svc.eventbus.Dispatch(ctx, event.RoleAfterUpdate(upd, r))

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

// validateContext validates role context expression
func (svc role) validateContext(ctx context.Context, r *types.RoleContext) error {
	if len(strings.TrimSpace(r.Expr)) == 0 {
		return nil
	}

	_, err := expr.NewParser().Parse(r.Expr)
	if err != nil {
		return err
	}

	// this is really as much as we can validate at this point
	// any further validation of the expression through actual execution
	// would require us to bring in information about resources from non-system components
	// and figuring out the exact module structure when using this with record values
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

		if svc.IsSystem(r) {
			return RoleErrNotAllowedToDelete()
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

		svc.eventbus.Dispatch(ctx, event.RoleAfterDelete(nil, r))

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

		if svc.IsSystem(r) {
			return RoleErrNotAllowedToUndelete()
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
		r, upd  *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(ctx, roleID); err != nil {
			return err
		}

		if svc.IsSystem(r) {
			return RoleErrNotAllowedToArchive()
		}

		upd = r.Clone()
		if err = svc.eventbus.WaitFor(ctx, event.RoleBeforeUpdate(upd, r)); err != nil {
			return
		}

		raProps.setRole(upd)

		if !svc.ac.CanUpdateRole(ctx, upd) {
			return RoleErrNotAllowedToArchive()
		}

		upd.ArchivedAt = now()
		if err = store.UpdateRole(ctx, svc.store, upd); err != nil {
			return
		}

		svc.eventbus.Dispatch(ctx, event.RoleAfterUpdate(upd, r))
		return
	}()

	return svc.recordAction(ctx, raProps, RoleActionArchive, err)
}

func (svc role) Unarchive(ctx context.Context, roleID uint64) (err error) {
	var (
		r, upd  *types.Role
		raProps = &roleActionProps{role: &types.Role{ID: roleID}}
	)

	err = func() (err error) {
		if r, err = svc.findByID(ctx, roleID); err != nil {
			return err
		}

		if svc.IsSystem(r) {
			return RoleErrNotAllowedToUnarchive()
		}

		upd = r.Clone()
		if err = svc.eventbus.WaitFor(ctx, event.RoleBeforeUpdate(upd, r)); err != nil {
			return
		}

		raProps.setRole(upd)

		if !svc.ac.CanDeleteRole(ctx, upd) {
			return RoleErrNotAllowedToUndelete()
		}

		r.ArchivedAt = nil
		if err = store.UpdateRole(ctx, svc.store, upd); err != nil {
			return
		}

		svc.eventbus.Dispatch(ctx, event.RoleAfterUpdate(upd, r))
		return nil
	}()

	return svc.recordAction(ctx, raProps, RoleActionUnarchive, err)
}

func (svc role) CloneRules(ctx context.Context, roleID uint64, cloneToRoleID ...uint64) (err error) {
	if !svc.ac.CanGrant(ctx) {
		return RoleErrNotAllowedToCloneRules()
	}

	return svc.rbac.CloneRulesByRoleID(ctx, roleID, cloneToRoleID...)
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
		if roleID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(ctx, roleID); err != nil {
			return err
		}

		if svc.IsClosed(r) || svc.IsContextual(r) {
			return RoleErrNotAllowedToManageMembers()
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
		if roleID == 0 || memberID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(ctx, roleID); err != nil {
			return
		}

		raProps.setRole(r)

		if svc.IsClosed(r) || svc.IsContextual(r) {
			return RoleErrNotAllowedToManageMembers()
		}

		if m, err = svc.user.FindByID(ctx, memberID); err != nil {
			return
		}

		raProps.setMember(m)

		if err = svc.eventbus.WaitFor(ctx, event.RoleMemberBeforeAdd(m, r)); err != nil {
			return
		}

		if !svc.ac.CanManageMembersOnRole(ctx, r) {
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
		if roleID == 0 || memberID == 0 {
			return RoleErrInvalidID()
		}

		if r, err = svc.findByID(ctx, roleID); err != nil {
			return
		}

		if svc.IsClosed(r) || svc.IsContextual(r) {
			return RoleErrNotAllowedToManageMembers()
		}

		raProps.setRole(r)

		if m, err = svc.user.FindByID(ctx, memberID); err != nil {
			return
		}

		raProps.setMember(m)

		if err = svc.eventbus.WaitFor(ctx, event.RoleMemberBeforeRemove(m, r)); err != nil {
			return
		}

		if !svc.ac.CanManageMembersOnRole(ctx, r) {
			return RoleErrNotAllowedToManageMembers()
		}

		if err = store.DeleteRoleMember(ctx, svc.store, &types.RoleMember{RoleID: r.ID, UserID: m.ID}); err != nil {
			return
		}

		// @todo skipping this for now,
		//			AS per now PUT `/role/{roleID}` endpoint updates role and it's member with it but
		//			only if one or more members are included in request otherwise we ignore it,
		//			which causes issue in admin when we remove all the members from role.
		//			we have to rework the role membership management logic in admin,
		//			and use the dedicated endpoints for POST, DELETE role member.
		// if err = svc.auth.RemoveAccessTokens(ctx, m); err != nil {
		//	 return
		// }

		_ = svc.eventbus.WaitFor(ctx, event.RoleMemberAfterRemove(m, r))
		return nil
	}()

	return svc.recordAction(ctx, raProps, RoleActionMemberRemove, err)
}

func loadRole(ctx context.Context, s store.Roles, ID uint64) (res *types.Role, err error) {
	if ID == 0 {
		return nil, RoleErrInvalidID()
	}

	if res, err = store.LookupRoleByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, RoleErrNotFound()
	}

	return
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

// Initializes roles to RBAC and default role service
//
// Sets all closed & system roles
func initRoles(ctx context.Context, log *zap.Logger, opt options.RbacOpt, eb eventbusRoleChangeRegistry, ru rbacRoleUpdater) (err error) {
	var (
		// splits space separated string into map
		s = func(l string) (map[string]bool, error) {
			m := make(map[string]bool)
			for _, r := range strings.Split(l, " ") {
				if r = strings.TrimSpace(r); len(r) == 0 {
					continue
				}

				if !handle.IsValid(r) {
					return nil, fmt.Errorf("invalid handle '%s'", r)
				}

				m[r] = true
			}

			return m, nil
		}

		// joins map keys into string slice
		j = func(mm ...map[string]bool) []string {
			o := make([]string, 0)

			for _, m := range mm {
				for r := range m {
					o = append(o, r)
				}
			}

			return o
		}

		bypass, authenticated, anonymous map[string]bool
	)

	if bypass, err = s(opt.BypassRoles); err != nil {
		return fmt.Errorf("failed to process list of bypass roles (RBAC_BYPASS_ROLES): %w", err)
	} else if !bypass[intAuth.BypassRoleHandle] {
		return fmt.Errorf("default bypass role %s (DefaultBypassRole) not in bypass role list (RBAC_BYPASS_ROLES)", intAuth.BypassRoleHandle)
	}

	if authenticated, err = s(opt.AuthenticatedRoles); err != nil {
		return fmt.Errorf("failed to process list of authenticated roles (RBAC_AUTHENTICATED_ROLES): %w", err)
	} else if !authenticated[intAuth.AuthenticatedRoleHandle] {
		return fmt.Errorf("default authenticated role %s (DefaultAuthenticatedRole) not in authenticated role list (RBAC_AUTHENTICATED_ROLES)", intAuth.AuthenticatedRoleHandle)
	}

	if anonymous, err = s(opt.AnonymousRoles); err != nil {
		return fmt.Errorf("failed to process list of anonymous roles (RBAC_ANONYMOUS_ROLES): %w", err)
	} else if !anonymous[intAuth.AnonymousRoleHandle] {
		return fmt.Errorf("default anonymous role %s (DefaultAnonymousRole) not in anonymous role list (RBAC_ANONYMOUS_ROLES)", intAuth.AnonymousRoleHandle)
	}

	for r := range authenticated {
		if bypass[r] {
			return fmt.Errorf("role %s used for authenticated users must not be used as bypass role", r)
		}
	}

	for r := range anonymous {
		if bypass[r] {
			return fmt.Errorf("role %s used for anonymous users must not be used as bypass role", r)
		}

		if authenticated[r] {
			return fmt.Errorf("role %s used for anonymous users must not be used as bypass role", r)
		}
	}

	tmp := j(bypass, authenticated, anonymous)
	log.Debug("setting system roles", zap.Strings("roles", tmp))
	DefaultRole.SetSystem(tmp...)

	tmp = j(authenticated, anonymous)
	log.Debug("setting closed roles", zap.Strings("roles", tmp))
	DefaultRole.SetClosed(tmp...)

	// Initial RBAC update
	if err = UpdateRbacRoles(ctx, log, ru, bypass, authenticated, anonymous); err != nil {
		return
	}

	// Hook to role create, update & delete events and
	// re-apply all roles to RBAC
	eb.Register(
		func(_ context.Context, ev eventbus.Event) error {
			log.Debug("role changed, updating RBAC")
			return UpdateRbacRoles(ctx, log, ru, bypass, authenticated, anonymous)
		},
		eventbus.For("system:role"),
		eventbus.On("afterUpdate", "afterCreate", "afterDelete"),
	)

	return nil
}

func UpdateRbacRoles(ctx context.Context, log *zap.Logger, ru rbacRoleUpdater, bypass, authenticated, anonymous map[string]bool) error {
	var (
		p  = expr.NewParser()
		f  = types.RoleFilter{}
		rr []*rbac.Role

		countBypass, countAuth, countAnony int
	)

	f.Paging.Total = 0
	roles, _, err := DefaultStore.SearchRoles(ctx, f)
	if err != nil {
		log.Error("failed to read roles", zap.Error(err))
		return nil
	}

	for _, r := range roles {
		log := log.With(
			zap.Uint64("ID", r.ID),
			zap.String("handle", r.Handle),
		)

		switch {
		case bypass[r.Handle]:
			countBypass++
			rr = append(rr, rbac.BypassRole.Make(r.ID, r.Handle))
			log.Debug("bypass role added")

		case anonymous[r.Handle]:
			countAnony++
			rr = append(rr, rbac.AnonymousRole.Make(r.ID, r.Handle))
			log.Debug("anonymous role added")

		case authenticated[r.Handle]:
			countAuth++
			rr = append(rr, rbac.AuthenticatedRole.Make(r.ID, r.Handle))
			log.Debug("authenticated role added")

		case r.Meta != nil && r.Meta.Context != nil && len(r.Meta.Context.Expr) > 0:
			log := log.With(zap.String("expr", r.Meta.Context.Expr))
			eval, err := p.Parse(r.Meta.Context.Expr)
			if err != nil {
				log.Warn("failed to parse role context expression, defaulting to deny", zap.Error(err))

				rr = append(rr, rbac.MakeContextRole(r.ID, r.Handle, func(_ map[string]interface{}) bool {
					log.Warn("role context expression not parsed, fallback to deny", zap.Error(err))
					return false
				}))
				continue
			}

			check := func(scope map[string]interface{}) bool {
				vars, err := expr.NewVars(scope)
				if err != nil {
					log.Warn("failed to convert check scope to expr.Vars, fallback to deny", zap.Error(err))
					return false
				}

				test, err := eval.Test(ctx, vars)
				if err != nil {
					log.Warn("failed to evaluate role context expression, fallback to deny", zap.Error(err))
					return false
				}

				return test
			}

			rr = append(rr, rbac.MakeContextRole(r.ID, r.Handle, check, r.Meta.Context.Resource...))
			log.Debug("context role added")

		default:
			rr = append(rr, rbac.CommonRole.Make(r.ID, r.Handle))
			log.Debug("common role added")
		}
	}

	if countBypass == 0 {
		log.Warn("no bypass roles registered, Corteza might not work as expected")
	}

	if countAuth == 0 {
		log.Warn("no roles for authentication users registered, Corteza might not work as expected")
	}

	if countAnony == 0 {
		log.Warn("no roles for anonymous users registered, Corteza might not work as expected")
	}

	ru.UpdateRoles(rr...)
	return nil
}
