package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	a "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	resourceTranslation struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        resourceTranslationAccessController
	}

	resourceTranslationAccessController interface {
		CanManageResourceTranslation(context.Context) bool
	}

	ResourceTranslationService interface {
		List(context.Context, types.ResourceTranslationFilter) (types.ResourceTranslationSet, types.ResourceTranslationFilter, error)
		Create(context.Context, *types.ResourceTranslation) (*types.ResourceTranslation, error)
		Update(context.Context, *types.ResourceTranslation) (*types.ResourceTranslation, error)
		Read(context.Context, uint64) (*types.ResourceTranslation, error)
		Delete(context.Context, uint64) error
		Undelete(context.Context, uint64) error
	}
)

func ResourceTranslation() *resourceTranslation {
	return &resourceTranslation{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

func (svc resourceTranslation) Read(ctx context.Context, ID uint64) (cc *types.ResourceTranslation, err error) {
	var (
		ccProps = &resourceTranslationActionProps{resourceTranslation: &types.ResourceTranslation{ID: ID}}
	)

	err = func() error {
		if ID == 0 {
			return TemplateErrInvalidID()
		}

		if !svc.ac.CanManageResourceTranslation(ctx) {
			return ResourceTranslationErrNotAllowedToManage()
		}

		if cc, err = store.LookupResourceTranslationByID(ctx, svc.store, ID); err != nil {
			return ResourceTranslationErrInvalidID().Wrap(err)
		}

		ccProps.setResourceTranslation(cc)

		return nil
	}()

	return cc, svc.recordAction(ctx, ccProps, ResourceTranslationActionLookup, err)
}

func (svc resourceTranslation) List(ctx context.Context, filter types.ResourceTranslationFilter) (set types.ResourceTranslationSet, f types.ResourceTranslationFilter, err error) {
	var (
		aProps = &resourceTranslationActionProps{filter: &filter}
	)

	err = func() error {
		if !svc.ac.CanManageResourceTranslation(ctx) {
			return ResourceTranslationErrNotAllowedToManage()
		}

		if set, f, err = store.SearchResourceTranslations(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ResourceTranslationActionSearch, err)
}

func (svc resourceTranslation) Create(ctx context.Context, new *types.ResourceTranslation) (cc *types.ResourceTranslation, err error) {
	var (
		tplProps = &resourceTranslationActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanManageResourceTranslation(ctx) {
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

		cc = new

		return nil
	}()

	return cc, svc.recordAction(ctx, tplProps, ResourceTranslationActionCreate, err)
}

func (svc resourceTranslation) Update(ctx context.Context, upd *types.ResourceTranslation) (cc *types.ResourceTranslation, err error) {
	var (
		tplProps = &resourceTranslationActionProps{update: upd}
	)

	err = func() (err error) {
		if upd.ID == 0 {
			return ResourceTranslationErrInvalidID()
		}

		if !svc.ac.CanManageResourceTranslation(ctx) {
			return ResourceTranslationErrNotAllowedToManage()
		}

		if cc, err = store.LookupResourceTranslationByID(ctx, svc.store, upd.ID); err != nil {
			return
		}

		tplProps.setResourceTranslation(cc)

		// @todo corredor?
		cc.Lang = upd.Lang
		cc.Resource = upd.Resource
		cc.K = upd.K
		cc.Message = upd.Message
		cc.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
		cc.OwnedBy = upd.OwnedBy

		cc.UpdatedAt = now()

		if err = store.UpdateResourceTranslation(ctx, svc.store, cc); err != nil {
			return err
		}

		return nil
	}()

	return cc, svc.recordAction(ctx, tplProps, ResourceTranslationActionUpdate, err)
}

func (svc resourceTranslation) Delete(ctx context.Context, ID uint64) (err error) {
	var (
		tplProps = &resourceTranslationActionProps{}
		cc       *types.ResourceTranslation
	)

	err = func() (err error) {
		if ID == 0 {
			return ResourceTranslationErrInvalidID()
		}

		if !svc.ac.CanManageResourceTranslation(ctx) {
			return ResourceTranslationErrNotAllowedToManage()
		}

		if cc, err = store.LookupResourceTranslationByID(ctx, svc.store, ID); err != nil {
			return
		}

		tplProps.setResourceTranslation(cc)

		// @todo corredor?

		cc.DeletedAt = now()
		if err = store.UpdateResourceTranslation(ctx, svc.store, cc); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, tplProps, ResourceTranslationActionDelete, err)
}

func (svc resourceTranslation) Undelete(ctx context.Context, ID uint64) (err error) {
	var (
		tplProps = &resourceTranslationActionProps{}
		cc       *types.ResourceTranslation
	)

	err = func() (err error) {
		if ID == 0 {
			return ResourceTranslationErrInvalidID()
		}

		if !svc.ac.CanManageResourceTranslation(ctx) {
			return ResourceTranslationErrNotAllowedToManage()
		}

		if cc, err = store.LookupResourceTranslationByID(ctx, svc.store, ID); err != nil {
			return
		}

		tplProps.setResourceTranslation(cc)

		// @todo corredor?
		cc.DeletedAt = nil
		if err = store.UpdateResourceTranslation(ctx, svc.store, cc); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, tplProps, ResourceTranslationActionUndelete, err)
}
