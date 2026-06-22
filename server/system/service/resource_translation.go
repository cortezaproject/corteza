package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (FindByID, Search, Create, Update, DeleteByID, UndeleteByID)
// is generated in resource_translation.gen.go from system/resource_translation.cue.
//
// Every op is bespoke: a single, non-RBAC access check
// (CanManageResourceTranslations) gates all of them and the bodies diverge from
// the standard CRUD scaffold, so each op delegates to a hand-written on<Op>
// handler here (customBodyOps + customAccessOps on every op).

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

// onLookup is the custom body for the generated FindByID. The generated method
// owns the action-log scaffold (+ recordAction with ResourceTranslationActionLookup);
// the ID/access guards and the deleted-inclusive lookup live here.
func (svc *resourceTranslation) onLookup(ctx context.Context, ID uint64, ccProps *resourceTranslationActionProps) (cc *types.ResourceTranslation, err error) {
	if ID == 0 {
		return nil, TemplateErrInvalidID()
	}

	if !svc.ac.CanManageResourceTranslations(ctx) {
		return nil, ResourceTranslationErrNotAllowedToManage()
	}

	if cc, err = store.LookupResourceTranslationByID(ctx, svc.store, ID); err != nil {
		return nil, ResourceTranslationErrInvalidID().Wrap(err)
	}

	ccProps.setResourceTranslation(cc)

	return cc, nil
}

// onSearch is the custom body for the generated Search. The access check is
// non-RBAC (CanManageResourceTranslations), so it lives here rather than in the
// generated scaffold.
func (svc *resourceTranslation) onSearch(ctx context.Context, filter types.ResourceTranslationFilter, aProps *resourceTranslationActionProps) (set types.ResourceTranslationSet, f types.ResourceTranslationFilter, err error) {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return nil, f, ResourceTranslationErrNotAllowedToManage()
	}

	if set, f, err = store.SearchResourceTranslations(ctx, svc.store, filter); err != nil {
		return nil, f, err
	}

	return set, f, nil
}

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold; the non-RBAC access check and the bespoke
// field assignment (incl. CreatedBy) live here.
func (svc *resourceTranslation) onCreate(ctx context.Context, new *types.ResourceTranslation) (err error) {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return ResourceTranslationErrNotAllowedToManage()
	}

	// @todo corredor?

	// Set new values after beforeCreate events are emitted
	new.ID = nextID()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	new.CreatedAt = *now()

	if err = store.CreateResourceTranslation(ctx, svc.store, new); err != nil {
		return
	}

	return nil
}

// onUpdate is the custom body for the generated Update. The non-RBAC access
// check, the lookup of the existing record and the bespoke field copy live here.
func (svc *resourceTranslation) onUpdate(ctx context.Context, s store.Storer, upd, res *types.ResourceTranslation, tplProps *resourceTranslationActionProps, _ func() error, _ func() error) error {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return ResourceTranslationErrNotAllowedToManage()
	}

	res.Lang = upd.Lang
	res.Resource = upd.Resource
	res.K = upd.K
	res.Message = upd.Message
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	res.OwnedBy = upd.OwnedBy
	res.UpdatedAt = now()

	return store.UpdateResourceTranslation(ctx, s, res)
}

// onDelete is the custom body for the generated DeleteByID. The non-RBAC access
// check and the soft-delete via UpdateResourceTranslation live here.
func (svc *resourceTranslation) onDelete(ctx context.Context, s store.Storer, res *types.ResourceTranslation, tplProps *resourceTranslationActionProps) error {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return ResourceTranslationErrNotAllowedToManage()
	}

	res.DeletedAt = now()
	return store.UpdateResourceTranslation(ctx, s, res)
}

// onUndelete is the custom body for the generated UndeleteByID. The non-RBAC
// access check and the clearing of deleted_at live here.
func (svc *resourceTranslation) onUndelete(ctx context.Context, s store.Storer, res *types.ResourceTranslation, tplProps *resourceTranslationActionProps) error {
	if !svc.ac.CanManageResourceTranslations(ctx) {
		return ResourceTranslationErrNotAllowedToManage()
	}

	res.DeletedAt = nil
	return store.UpdateResourceTranslation(ctx, s, res)
}
