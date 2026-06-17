package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	types "github.com/crusttech/human/server/system/types"
)

func (svc *dalSensitivityLevel) FindByID(ctx context.Context, ID uint64) (res *types.DalSensitivityLevel, err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{sensitivityLevel: &types.DalSensitivityLevel{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, DalSensitivityLevelActionLookup, err)
}

func (svc *dalSensitivityLevel) Search(ctx context.Context, filter types.DalSensitivityLevelFilter) (set types.DalSensitivityLevelSet, f types.DalSensitivityLevelFilter, err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{search: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, DalSensitivityLevelActionSearch, err)
}

func (svc *dalSensitivityLevel) Create(ctx context.Context, new *types.DalSensitivityLevel) (res *types.DalSensitivityLevel, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &dalSensitivityLevelActionProps{sensitivityLevel: new, new: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, DalSensitivityLevelActionCreate, err)
}

func (svc *dalSensitivityLevel) Update(ctx context.Context, upd *types.DalSensitivityLevel) (res *types.DalSensitivityLevel, err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{update: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, DalSensitivityLevelActionUpdate, err)
}

func (svc *dalSensitivityLevel) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, DalSensitivityLevelActionDelete, err)
}

func (svc *dalSensitivityLevel) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &dalSensitivityLevelActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, DalSensitivityLevelActionUndelete, err)
}
