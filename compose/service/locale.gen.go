package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - compose.module.yaml
// - compose.namespace.yaml
// - compose.page.yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	resourceTranslation struct {
		actionlog actionlog.Recorder
		locale    locale.Resource
		store     store.Storer
		ac        localeAccessController
	}

	localeAccessController interface {
		// CanManageResourceTranslation(context.Context) bool
	}

	ResourceTranslationService interface {
		Module(ctx context.Context, namespaceID uint64, ID uint64) (locale.ResourceTranslationSet, error)
		Namespace(ctx context.Context, ID uint64) (locale.ResourceTranslationSet, error)
		Page(ctx context.Context, namespaceID uint64, ID uint64) (locale.ResourceTranslationSet, error)

		Upsert(context.Context, locale.ResourceTranslationSet) error
		Locale() locale.Resource
	}
)

func ResourceTranslation(ls locale.Resource) *resourceTranslation {
	return &resourceTranslation{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		locale:    ls,
	}
}

func (svc resourceTranslation) Upsert(ctx context.Context, rr locale.ResourceTranslationSet) (err error) {
	// @todo AC
	// @todo validation

	defer locale.Global().ReloadResourceTranslations(ctx)
	me := auth.GetIdentityFromContext(ctx)

	// - group by resource
	localeByRes := make(map[string]locale.ResourceTranslationSet)
	for _, r := range rr {
		localeByRes[r.Resource] = append(localeByRes[r.Resource], r)
	}

	// - for each resource, fetch the current state
	sysLocale := make(systemTypes.ResourceTranslationSet, 0, len(rr))
	for res, rr := range localeByRes {
		current, _, err := store.SearchResourceTranslations(ctx, svc.store, systemTypes.ResourceTranslationFilter{
			Resource: res,
		})
		if err != nil {
			return err
		}

		// get deltas and prepare upsert accordingly
		aux := current.New(rr)
		aux.Walk(func(cc *systemTypes.ResourceTranslation) error {
			cc.ID = nextID()
			cc.CreatedAt = *now()
			cc.CreatedBy = me.Identity()

			return nil
		})
		sysLocale = append(sysLocale, aux...)

		aux = current.Old(rr)
		aux.Walk(func(cc *systemTypes.ResourceTranslation) error {
			cc.UpdatedAt = now()
			cc.UpdatedBy = me.Identity()
			return nil
		})
		sysLocale = append(sysLocale, aux...)
	}

	err = store.UpsertResourceTranslation(ctx, svc.store, sysLocale...)
	if err != nil {
		return err
	}
	return nil
}

func (svc resourceTranslation) Locale() locale.Resource {
	return svc.locale
}

func (svc resourceTranslation) Module(ctx context.Context, namespaceID uint64, ID uint64) (locale.ResourceTranslationSet, error) {
	var (
		err error
		out locale.ResourceTranslationSet
		res *types.Module
		k   types.LocaleKey
	)

	res, err = svc.loadModule(ctx, svc.store, namespaceID, ID)
	if err != nil {
		return nil, err
	}

	for _, tag := range svc.locale.Tags() {
		k = types.LocaleKeyModuleName
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), k.Path),
		})

	}

	tmp, err := svc.moduleExtended(ctx, res)
	return append(out, tmp...), err
}

func (svc resourceTranslation) Namespace(ctx context.Context, ID uint64) (locale.ResourceTranslationSet, error) {
	var (
		err error
		out locale.ResourceTranslationSet
		res *types.Namespace
		k   types.LocaleKey
	)

	res, err = svc.loadNamespace(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	for _, tag := range svc.locale.Tags() {
		k = types.LocaleKeyNamespaceName
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), k.Path),
		})

		k = types.LocaleKeyNamespaceSubtitle
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), k.Path),
		})

		k = types.LocaleKeyNamespaceDescription
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), k.Path),
		})

	}
	return out, nil
}

func (svc resourceTranslation) Page(ctx context.Context, namespaceID uint64, ID uint64) (locale.ResourceTranslationSet, error) {
	var (
		err error
		out locale.ResourceTranslationSet
		res *types.Page
		k   types.LocaleKey
	)

	res, err = svc.loadPage(ctx, svc.store, namespaceID, ID)
	if err != nil {
		return nil, err
	}

	for _, tag := range svc.locale.Tags() {
		k = types.LocaleKeyPageTitle
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), k.Path),
		})

		k = types.LocaleKeyPageDescription
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), k.Path),
		})

	}

	tmp, err := svc.pageExtended(ctx, res)
	return append(out, tmp...), err
}
