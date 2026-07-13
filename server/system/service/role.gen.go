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
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type role struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        roleAccessController
	services  *roleServices
}

func (svc *role) FindByID(ctx context.Context, ID uint64) (res *types.Role, err error) {
	var (
		aProps = &roleActionProps{role: &types.Role{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, RoleActionLookup, err)
}

func (svc *role) Search(ctx context.Context, filter types.RoleFilter) (set types.RoleSet, f types.RoleFilter, err error) {
	var (
		aProps = &roleActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Role) (bool, error) {
		if !svc.ac.CanReadRole(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if err = svc.beforeSearch(ctx, &filter); err != nil {
			return err
		}
		if !svc.ac.CanSearchRoles(ctx) {
			return RoleErrNotAllowedToSearch()
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

		if set, f, err = store.SearchRoles(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledRoles(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, RoleActionSearch, err)
}

func (svc *role) Create(ctx context.Context, new *types.Role) (res *types.Role, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &roleActionProps{role: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if err = svc.validate(ctx, new); err != nil {
			return err
		}
		if !svc.ac.CanCreateRole(ctx) {
			return RoleErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateRole(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new

		if err = svc.afterCreate(ctx, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, RoleActionCreate, err)
}

func (svc *role) Update(ctx context.Context, upd *types.Role) (res *types.Role, err error) {
	var (
		aProps = &roleActionProps{update: upd}
		old    *types.Role
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadRole(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setRole(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return RoleErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return RoleErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Handle = upd.Handle
		res.Name = upd.Name
		res.Meta = upd.Meta
		res.UpdatedAt = now()

		if err = store.UpdateRole(ctx, s, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, RoleActionUpdate, err, old, res)
}

func (svc *role) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
		res    *types.Role
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadRole(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setRole(res)

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, RoleActionDelete, err)
}

func (svc *role) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
		res    *types.Role
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadRole(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setRole(res)

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, RoleActionUndelete, err)
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
func (svc *role) guard(_ context.Context, _ *types.Role) error { return nil }

func (svc *role) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *role) Archive(ctx context.Context, roleID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onArchive(ctx, aProps, roleID)
		return err
	}()

	return svc.recordAction(ctx, aProps, RoleActionArchive, err)
}

func (svc *role) Unarchive(ctx context.Context, roleID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onUnarchive(ctx, aProps, roleID)
		return err
	}()

	return svc.recordAction(ctx, aProps, RoleActionUnarchive, err)
}

func (svc *role) CloneRules(ctx context.Context, roleID uint64, cloneToRoleID ...uint64) (err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if !svc.ac.CanGrant(ctx) {
			return RoleErrNotAllowedToCloneRules()
		}

		err = svc.onCloneRules(ctx, aProps, roleID, cloneToRoleID...)
		return err
	}()

	return svc.recordAction(ctx, aProps, RoleActionCloneRules, err)
}

func (svc *role) Membership(ctx context.Context, userID uint64) (mm types.RoleMemberSet, err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		mm, err = svc.onMembership(ctx, aProps, userID)
		return err
	}()

	return mm, svc.recordAction(ctx, aProps, RoleActionMembership, err)
}

func (svc *role) MemberList(ctx context.Context, roleID uint64) (mm types.RoleMemberSet, err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		mm, err = svc.onMemberList(ctx, aProps, roleID)
		return err
	}()

	return mm, svc.recordAction(ctx, aProps, RoleActionMembers, err)
}

func (svc *role) MemberAdd(ctx context.Context, roleID uint64, memberID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onMemberAdd(ctx, aProps, roleID, memberID)
		return err
	}()

	return svc.recordAction(ctx, aProps, RoleActionMemberAdd, err)
}

func (svc *role) MemberAddGroup(ctx context.Context, roleID uint64, userGroupID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onMemberAddGroup(ctx, aProps, roleID, userGroupID)
		return err
	}()

	return svc.recordAction(ctx, aProps, RoleActionMemberAdd, err)
}

func (svc *role) MemberRemove(ctx context.Context, roleID uint64, memberID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onMemberRemove(ctx, aProps, roleID, memberID)
		return err
	}()

	return svc.recordAction(ctx, aProps, RoleActionMemberRemove, err)
}

func (svc *role) MemberRemoveGroup(ctx context.Context, roleID uint64, userGroupID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onMemberRemoveGroup(ctx, aProps, roleID, userGroupID)
		return err
	}()

	return svc.recordAction(ctx, aProps, RoleActionMemberRemove, err)
}
