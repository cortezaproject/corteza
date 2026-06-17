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

func (svc *resourceTranslation) FindByID(ctx context.Context, ID uint64) (res *types.ResourceTranslation, err error) {
	var (
		aProps = &resourceTranslationActionProps{resourceTranslation: &types.ResourceTranslation{ID: ID}}
	)

	err = func() error {
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ResourceTranslationActionLookup, err)
}

func (svc *resourceTranslation) Search(ctx context.Context, filter types.ResourceTranslationFilter) (set types.ResourceTranslationSet, f types.ResourceTranslationFilter, err error) {
	var (
		aProps = &resourceTranslationActionProps{filter: &filter}
	)

	err = func() error {
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ResourceTranslationActionSearch, err)
}

func (svc *resourceTranslation) Create(ctx context.Context, new *types.ResourceTranslation) (res *types.ResourceTranslation, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &resourceTranslationActionProps{resourceTranslation: new, new: new}
	)

	err = func() (err error) {
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ResourceTranslationActionCreate, err)
}

func (svc *resourceTranslation) Update(ctx context.Context, upd *types.ResourceTranslation) (res *types.ResourceTranslation, err error) {
	var (
		aProps = &resourceTranslationActionProps{update: upd}
	)

	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ResourceTranslationActionUpdate, err)
}

func (svc *resourceTranslation) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &resourceTranslationActionProps{}
	)

	err = func() (err error) {
		return svc.onDelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ResourceTranslationActionDelete, err)
}

func (svc *resourceTranslation) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &resourceTranslationActionProps{}
	)

	err = func() (err error) {
		return svc.onUndelete(ctx, ID, aProps)
	}()

	return svc.recordAction(ctx, aProps, ResourceTranslationActionUndelete, err)
}
