package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *role) FindByID(ctx context.Context, ID uint64) (res *types.Role, err error) {
	var (
		aProps = &roleActionProps{role: &types.Role{ID: ID}}
	)

	err = func() error {
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
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, RoleActionUpdate, err)
}

func (svc *role) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, RoleActionDelete, err)
}

func (svc *role) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &roleActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, RoleActionUndelete, err)
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
