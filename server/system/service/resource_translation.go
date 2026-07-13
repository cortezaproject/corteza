package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	resourceTranslationAccessController interface {
		CanManageResourceTranslations(context.Context) bool
	}

	ResourceTranslationService interface {
		Search(context.Context, types.ResourceTranslationFilter) (types.ResourceTranslationSet, types.ResourceTranslationFilter, error)
		Create(context.Context, *types.ResourceTranslation) (*types.ResourceTranslation, error)
		Update(context.Context, *types.ResourceTranslation) (*types.ResourceTranslation, error)
		FindByID(context.Context, uint64) (*types.ResourceTranslation, error)
		DeleteByID(context.Context, uint64) error
		UndeleteByID(context.Context, uint64) error
	}
)

func ResourceTranslation() *resourceTranslation {
	return &resourceTranslation{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

func (svc *resourceTranslation) onLookup(ctx context.Context, ID uint64, aProps *resourceTranslationActionProps) (*types.ResourceTranslation, error) {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return nil, ResourceTranslationErrNotAllowedToManage()
	}

	res, err := store.LookupResourceTranslationByID(ctx, svc.store, ID)
	if err != nil {
		return nil, ResourceTranslationErrInvalidID().Wrap(err)
	}

	aProps.setResourceTranslation(res)
	return res, nil
}

func (svc *resourceTranslation) onSearch(ctx context.Context, filter types.ResourceTranslationFilter, aProps *resourceTranslationActionProps) (types.ResourceTranslationSet, types.ResourceTranslationFilter, error) {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return nil, filter, ResourceTranslationErrNotAllowedToManage()
	}

	set, f, err := store.SearchResourceTranslations(ctx, svc.store, filter)
	return set, f, err
}

func (svc *resourceTranslation) onCreate(ctx context.Context, new *types.ResourceTranslation) error {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return ResourceTranslationErrNotAllowedToManage()
	}

	new.ID = nextID()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	new.CreatedAt = *now()

	return store.CreateResourceTranslation(ctx, svc.store, new)
}

func (svc *resourceTranslation) onUpdate(ctx context.Context, s store.Storer, upd *types.ResourceTranslation, res *types.ResourceTranslation, aProps *resourceTranslationActionProps, before func() error, after func() error) error {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return ResourceTranslationErrNotAllowedToManage()
	}

	// Lang is not copied by gen; set it on upd so gen's field-copy picks it up via res.
	// Gen copies Resource/K/Message/OwnedBy/UpdatedBy but not Lang — copy it here.
	res.Lang = upd.Lang
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	return nil
}

func (svc *resourceTranslation) onDelete(ctx context.Context, s store.Storer, res *types.ResourceTranslation, aProps *resourceTranslationActionProps) error {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return ResourceTranslationErrNotAllowedToManage()
	}

	res.DeletedAt = now()
	return store.UpdateResourceTranslation(ctx, s, res)
}

func (svc *resourceTranslation) onUndelete(ctx context.Context, s store.Storer, res *types.ResourceTranslation, aProps *resourceTranslationActionProps) error {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return ResourceTranslationErrNotAllowedToManage()
	}

	res.DeletedAt = nil
	return store.UpdateResourceTranslation(ctx, s, res)
}
