package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type userGroupServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *userGroup) Create(ctx context.Context, new *types.UserGroup) (res *types.UserGroup, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &userGroupActionProps{userGroup: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if err = svc.validate(ctx, new); err != nil {
			return err
		}
		if !svc.ac.CanCreateUserGroup(ctx) {
			return UserGroupErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateUserGroup(ctx, svc.store, new); err != nil {
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

	return res, svc.recordAction(ctx, aProps, UserGroupActionCreate, err)
}

// toLabeledUserGroups converts to []label.LabeledResource
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

func (svc *userGroup) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *userGroup) scopeServices(ctx context.Context) *userGroupServices {
	return &userGroupServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}

func (svc *userGroup) Activate(ctx context.Context) (err error) {
	var (
		aProps = &userGroupActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onActivate(ctx, aProps)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserGroupActionActivate, err)
}

func (svc *userGroup) MemberList(ctx context.Context, userGroupID uint64) (mm types.UserSet, err error) {
	var (
		aProps = &userGroupActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		mm, err = svc.onMemberList(ctx, aProps, userGroupID)
		return err
	}()

	return mm, svc.recordAction(ctx, aProps, UserGroupActionMembers, err)
}

func (svc *userGroup) MemberAdd(ctx context.Context, userGroupID uint64, memberID uint64) (err error) {
	var (
		aProps = &userGroupActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onMemberAdd(ctx, aProps, userGroupID, memberID)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserGroupActionMemberAdd, err)
}
