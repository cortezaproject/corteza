package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	intAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/pkg/slice"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
)

type (
	roleServices struct {
		eventbus eventDispatcher
		rbac     rbacRuleService
		user     UserService
		auth     roleAuth
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

		// Membership answers "which roles does this user hold", so it is gated
		// on reading the subject rather than on each role.
		CanReadUser(context.Context, *types.User) bool
	}

	RoleService interface {
		FindByID(ctx context.Context, roleID uint64) (*types.Role, error)
		FindByName(ctx context.Context, name string) (*types.Role, error)
		FindByHandle(ctx context.Context, handle string) (*types.Role, error)
		FindByAny(ctx context.Context, identifier interface{}) (*types.Role, error)
		Find(context.Context, types.RoleFilter) (types.RoleSet, types.RoleFilter, error)
		Search(context.Context, types.RoleFilter) (types.RoleSet, types.RoleFilter, error)
		DeleteByID(ctx context.Context, ID uint64) error

		IsSystem(r *types.Role) bool
		IsClosed(r *types.Role) bool

		Create(ctx context.Context, role *types.Role) (*types.Role, error)
		Update(ctx context.Context, role *types.Role) (*types.Role, error)

		Archive(ctx context.Context, ID uint64) error
		Unarchive(ctx context.Context, ID uint64) error
		Delete(ctx context.Context, ID uint64) error
		Undelete(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		CloneRules(ctx context.Context, ID uint64, cloneToRoleID ...uint64) error

		Membership(ctx context.Context, userID uint64) (types.RoleMemberSet, error)
		MemberList(ctx context.Context, roleID uint64) (types.RoleMemberSet, error)
		MemberAdd(ctx context.Context, roleID, userID uint64) error
		MemberAddGroup(ctx context.Context, roleID, userGroupID uint64) error
		MemberRemove(ctx context.Context, roleID, userID uint64) error
		MemberRemoveGroup(ctx context.Context, roleID, userGroupID uint64) error
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
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		services: &roleServices{
			eventbus: eventbus.Service(),
			rbac:     rbac,
			user:     DefaultUser,
			auth:     DefaultAuth,
			system:   make(map[string]bool),
			closed:   make(map[string]bool),
		},
	}
}

// SetSystem sets list of handles for all system roles
//
// System roles can not be changed or deleted
func (svc *role) SetSystem(hh ...string) {
	svc.services.system = slice.ToStringBoolMap(hh)
	delete(svc.services.system, "")
}

func (svc role) IsSystem(r *types.Role) bool {
	return len(r.Handle) > 0 && svc.services.system[r.Handle]
}

// SetClosed sets list of handles for all closed roles
//
// Closed roles can not have members
func (svc *role) SetClosed(hh ...string) {
	svc.services.closed = slice.ToStringBoolMap(hh)
	delete(svc.services.closed, "")
}

func (svc role) IsClosed(r *types.Role) bool {
	return len(r.Handle) > 0 && svc.services.closed[r.Handle]
}

func (svc role) IsContextual(r *types.Role) bool {
	return r.Meta != nil && r.Meta.Context != nil && len(r.Meta.Context.Expr) > 0
}

func (svc role) Delete(ctx context.Context, ID uint64) error   { return svc.DeleteByID(ctx, ID) }
func (svc role) Undelete(ctx context.Context, ID uint64) error { return svc.UndeleteByID(ctx, ID) }

func (svc role) Find(ctx context.Context, filter types.RoleFilter) (types.RoleSet, types.RoleFilter, error) {
	return svc.Search(ctx, filter)
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

func (svc role) UniqueCheck(ctx context.Context, r *types.Role) (err error) {
	var (
		raProps = &roleActionProps{role: r}
	)

	// Scoped to the project, matching the store's constraint check. A role
	// belongs to a project, so its handle and name only have to be unique within
	// one -- two projects may each define a "Customer" role, and a project
	// revision branch deliberately puts the same role in two revisions at once
	// (see system/service/project_revision_clone.go). Global roles all carry
	// ProjectID 0 and so keep being checked against each other.
	if r.Handle != "" {
		if ex, _ := store.LookupRoleByProjectIDHandle(ctx, svc.store, r.ProjectID, r.Handle); ex != nil && ex.ID > 0 && ex.ID != r.ID {
			raProps.setExisting(ex)
			return RoleErrHandleNotUnique()
		}
	}

	if r.Name != "" {
		if ex, _ := store.LookupRoleByProjectIDName(ctx, svc.store, r.ProjectID, r.Name); ex != nil && ex.ID > 0 && ex.ID != r.ID {
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

// -- hooks --

func (svc *role) validate(ctx context.Context, new *types.Role) error {
	if !handle.IsValid(new.Handle) {
		return RoleErrInvalidHandle()
	}

	if new.Meta != nil && new.Meta.Context != nil {
		if err := svc.validateContext(ctx, new.Meta.Context); err != nil {
			return err
		}
	}

	return svc.UniqueCheck(ctx, new)
}

func (svc *role) beforeSearch(_ context.Context, _ *types.RoleFilter) error {
	return nil
}

func (svc *role) beforeCreate(ctx context.Context, new *types.Role) error {
	return svc.services.eventbus.WaitFor(ctx, event.RoleBeforeCreate(new, nil))
}

func (svc *role) afterCreate(ctx context.Context, res *types.Role) error {
	svc.services.eventbus.Dispatch(ctx, event.RoleAfterCreate(res, res))
	return nil
}

func (svc *role) onLookup(ctx context.Context, roleID uint64, aProps *roleActionProps) (*types.Role, error) {
	r, err := svc.findByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadRole(ctx, r) {
		return nil, RoleErrNotAllowedToRead()
	}

	aProps.setRole(r)
	return r, nil
}

func (svc *role) onUpdate(ctx context.Context, _ store.Storer, upd *types.Role, res *types.Role, _ *roleActionProps, _ func() error, _ func() error) error {
	// Permission first, identity second. This check used to sit below the
	// system-role branch, and that branch returns nil when handle and name are
	// unchanged — so an edit to a system role's description or meta succeeded
	// without update permission at all. Whether a system role may be touched is
	// a separate question from whether this caller may touch any role.
	if !svc.ac.CanUpdateRole(ctx, res) {
		return RoleErrNotAllowedToUpdate()
	}

	if svc.IsSystem(res) {
		// prevent system role updates unless handle and name are unchanged
		if res.Handle == upd.Handle && res.Name == upd.Name {
			return nil
		}
		return RoleErrNotAllowedToUpdate()
	}

	if err := svc.services.eventbus.WaitFor(ctx, event.RoleBeforeUpdate(upd, res)); err != nil {
		return err
	}

	if upd.Meta != nil && upd.Meta.Context != nil {
		if err := svc.validateContext(ctx, upd.Meta.Context); err != nil {
			return err
		}
	}

	if err := svc.UniqueCheck(ctx, upd); err != nil {
		return err
	}

	svc.services.eventbus.Dispatch(ctx, event.RoleAfterUpdate(upd, res))
	return nil
}

func (svc *role) onDelete(ctx context.Context, s store.Storer, res *types.Role, _ *roleActionProps) error {
	if svc.IsSystem(res) {
		return RoleErrNotAllowedToDelete()
	}

	if !svc.ac.CanDeleteRole(ctx, res) {
		return RoleErrNotAllowedToDelete()
	}

	if err := svc.services.eventbus.WaitFor(ctx, event.RoleBeforeDelete(nil, res)); err != nil {
		return err
	}

	res.DeletedAt = now()
	if err := store.UpdateRole(ctx, s, res); err != nil {
		return err
	}

	svc.services.eventbus.Dispatch(ctx, event.RoleAfterDelete(nil, res))
	return nil
}

func (svc *role) onUndelete(ctx context.Context, s store.Storer, res *types.Role, _ *roleActionProps) error {
	if svc.IsSystem(res) {
		return RoleErrNotAllowedToUndelete()
	}

	upd := res.Clone()
	if err := svc.services.eventbus.WaitFor(ctx, event.RoleBeforeUpdate(upd, res)); err != nil {
		return err
	}

	if !svc.ac.CanDeleteRole(ctx, upd) {
		return RoleErrNotAllowedToDelete()
	}

	upd.DeletedAt = nil
	if err := store.UpdateRole(ctx, s, upd); err != nil {
		return err
	}

	svc.services.eventbus.Dispatch(ctx, event.RoleAfterUpdate(upd, res))
	return nil
}

func (svc *role) onArchive(ctx context.Context, _ *roleActionProps, roleID uint64) error {
	r, err := svc.findByID(ctx, roleID)
	if err != nil {
		return err
	}

	if svc.IsSystem(r) {
		return RoleErrNotAllowedToArchive()
	}

	upd := r.Clone()
	if err = svc.services.eventbus.WaitFor(ctx, event.RoleBeforeUpdate(upd, r)); err != nil {
		return err
	}

	if !svc.ac.CanUpdateRole(ctx, upd) {
		return RoleErrNotAllowedToArchive()
	}

	upd.ArchivedAt = now()
	if err = store.UpdateRole(ctx, svc.store, upd); err != nil {
		return err
	}

	svc.services.eventbus.Dispatch(ctx, event.RoleAfterUpdate(upd, r))
	return nil
}

func (svc *role) onUnarchive(ctx context.Context, _ *roleActionProps, roleID uint64) error {
	r, err := svc.findByID(ctx, roleID)
	if err != nil {
		return err
	}

	if svc.IsSystem(r) {
		return RoleErrNotAllowedToUnarchive()
	}

	upd := r.Clone()
	if err = svc.services.eventbus.WaitFor(ctx, event.RoleBeforeUpdate(upd, r)); err != nil {
		return err
	}

	if !svc.ac.CanDeleteRole(ctx, upd) {
		return RoleErrNotAllowedToUndelete()
	}

	upd.ArchivedAt = nil
	if err = store.UpdateRole(ctx, svc.store, upd); err != nil {
		return err
	}

	svc.services.eventbus.Dispatch(ctx, event.RoleAfterUpdate(upd, r))
	return nil
}

func (svc *role) onCloneRules(ctx context.Context, _ *roleActionProps, roleID uint64, cloneToRoleID ...uint64) error {
	return svc.services.rbac.CloneRulesByRoleID(ctx, roleID, cloneToRoleID...)
}

// onMembership returns every role a user holds.
//
// The read check is on the subject, not on each role: this answers "what does
// this user have", so the caller must be allowed to read that user. Previously
// it checked nothing at all beyond checkScope, so any authenticated caller
// could enumerate anyone's roles.
//
// A caller reading their own memberships is always permitted — that is the
// self-service case the webapp relies on.
func (svc *role) onMembership(ctx context.Context, _ *roleActionProps, userID uint64) (types.RoleMemberSet, error) {
	if userID == 0 {
		return nil, RoleErrInvalidID()
	}

	if userID != intAuth.GetIdentityFromContext(ctx).Identity() {
		u, err := store.LookupUserByID(ctx, svc.store, userID)
		if err != nil {
			return nil, err
		}

		if !svc.ac.CanReadUser(ctx, u) {
			return nil, RoleErrNotAllowedToRead()
		}
	}

	resource := fmt.Sprintf("corteza::system:user/%d", userID)
	mm, _, err := store.SearchRoleMembers(ctx, svc.store, types.RoleMemberFilter{Resource: resource})
	return mm, err
}

func (svc *role) onMemberList(ctx context.Context, _ *roleActionProps, roleID uint64) (mm types.RoleMemberSet, err error) {
	if roleID == 0 {
		return nil, RoleErrInvalidID()
	}

	r, err := svc.findByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	if svc.IsClosed(r) || svc.IsContextual(r) {
		return nil, RoleErrNotAllowedToManageMembers()
	}

	if !svc.ac.CanReadRole(ctx, r) {
		return nil, RoleErrNotAllowedToRead()
	}

	mm, _, err = store.SearchRoleMembers(ctx, svc.store, types.RoleMemberFilter{RoleID: roleID})
	return mm, err
}

func (svc *role) onMemberAdd(ctx context.Context, _ *roleActionProps, roleID, memberID uint64) error {
	if roleID == 0 || memberID == 0 {
		return RoleErrInvalidID()
	}

	r, err := svc.findByID(ctx, roleID)
	if err != nil {
		return err
	}

	if svc.IsClosed(r) || svc.IsContextual(r) {
		return RoleErrNotAllowedToManageMembers()
	}

	m, err := svc.services.user.FindByID(ctx, memberID)
	if err != nil {
		return err
	}

	if err = svc.services.eventbus.WaitFor(ctx, event.RoleMemberBeforeAdd(m, r)); err != nil {
		return err
	}

	if !svc.ac.CanManageMembersOnRole(ctx, r) {
		return RoleErrNotAllowedToManageMembers()
	}

	resource := fmt.Sprintf("corteza::system:user/%d", m.ID)
	if err = store.CreateRoleMember(ctx, svc.store, &types.RoleMember{RoleID: r.ID, Resource: resource}); err != nil {
		return err
	}

	_ = svc.services.eventbus.WaitFor(ctx, event.RoleMemberAfterAdd(m, r))
	return nil
}

func (svc *role) onMemberAddGroup(ctx context.Context, _ *roleActionProps, roleID, userGroupID uint64) error {
	if roleID == 0 || userGroupID == 0 {
		return RoleErrInvalidID()
	}

	r, err := svc.findByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Closed and contextual roles refuse membership changes. onMemberAdd and
	// onMemberRemove have always enforced this for users; the group variants
	// did not, so a closed role could gain members by adding a group instead
	// of a user — a bypass of the same rule, by a different door.
	if svc.IsClosed(r) || svc.IsContextual(r) {
		return RoleErrNotAllowedToManageMembers()
	}

	if !svc.ac.CanManageMembersOnRole(ctx, r) {
		return RoleErrNotAllowedToManageMembers()
	}

	resource := fmt.Sprintf("corteza::system:user-group/%d", userGroupID)
	return store.CreateRoleMember(ctx, svc.store, &types.RoleMember{RoleID: r.ID, Resource: resource})
}

func (svc *role) onMemberRemove(ctx context.Context, _ *roleActionProps, roleID, memberID uint64) error {
	if roleID == 0 || memberID == 0 {
		return RoleErrInvalidID()
	}

	r, err := svc.findByID(ctx, roleID)
	if err != nil {
		return err
	}

	if svc.IsClosed(r) || svc.IsContextual(r) {
		return RoleErrNotAllowedToManageMembers()
	}

	m, err := svc.services.user.FindByID(ctx, memberID)
	if err != nil {
		return err
	}

	if err = svc.services.eventbus.WaitFor(ctx, event.RoleMemberBeforeRemove(m, r)); err != nil {
		return err
	}

	if !svc.ac.CanManageMembersOnRole(ctx, r) {
		return RoleErrNotAllowedToManageMembers()
	}

	resource := fmt.Sprintf("corteza::system:user/%d", m.ID)
	if err = store.DeleteRoleMember(ctx, svc.store, &types.RoleMember{RoleID: r.ID, Resource: resource}); err != nil {
		return err
	}

	_ = svc.services.eventbus.WaitFor(ctx, event.RoleMemberAfterRemove(m, r))
	return nil
}

func (svc *role) onMemberRemoveGroup(ctx context.Context, _ *roleActionProps, roleID, userGroupID uint64) error {
	if roleID == 0 || userGroupID == 0 {
		return RoleErrInvalidID()
	}

	r, err := svc.findByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Closed and contextual roles refuse membership changes. onMemberAdd and
	// onMemberRemove have always enforced this for users; the group variants
	// did not, so a closed role could gain members by adding a group instead
	// of a user — a bypass of the same rule, by a different door.
	if svc.IsClosed(r) || svc.IsContextual(r) {
		return RoleErrNotAllowedToManageMembers()
	}

	if !svc.ac.CanManageMembersOnRole(ctx, r) {
		return RoleErrNotAllowedToManageMembers()
	}

	resource := fmt.Sprintf("corteza::system:user-group/%d", userGroupID)
	return store.DeleteRoleMember(ctx, svc.store, &types.RoleMember{RoleID: r.ID, Resource: resource})
}

// -- init helpers --

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
			logger.Uint64("ID", r.ID),
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
