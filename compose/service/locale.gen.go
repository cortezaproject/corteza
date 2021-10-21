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
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"golang.org/x/text/language"
)

type (
	localeAccessControl interface {
		CanManageResourceTranslations(ctx context.Context) bool
	}

	resourceTranslationsManager struct {
		actionlog actionlog.Recorder
		locale    locale.Resource
		store     store.Storer
		ac        localeAccessController
	}

	localeAccessController interface {
		CanManageResourceTranslations(context.Context) bool
	}

	ResourceTranslationsManagerService interface {
		Module(ctx context.Context, namespaceID uint64, ID uint64) (locale.ResourceTranslationSet, error)
		Namespace(ctx context.Context, ID uint64) (locale.ResourceTranslationSet, error)
		Page(ctx context.Context, namespaceID uint64, ID uint64) (locale.ResourceTranslationSet, error)

		Upsert(context.Context, locale.ResourceTranslationSet) error
		Locale() locale.Resource
	}
)

var ErrNotAllowedToManageResourceTranslations = errors.Unauthorized("not allowed to manage resource translations")

func ResourceTranslationsManager(ls locale.Resource) *resourceTranslationsManager {
	return &resourceTranslationsManager{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		locale:    ls,
	}
}

func (svc resourceTranslationsManager) Upsert(ctx context.Context, rr locale.ResourceTranslationSet) (err error) {
	// User is allowed to manage resource translations when:
	//  - managed resource translation strings are all for default language
	//  or
	//  - user is allowed to manage resource translations
	if rr.ContainsForeign(svc.Locale().Default().Tag) {
		if !svc.ac.CanManageResourceTranslations(ctx) {
			return ErrNotAllowedToManageResourceTranslations
		}
	}

	for _, r := range rr {
		r.Msg = locale.SanitizeMessage(r.Msg)
	}

	// @todo validation

	me := intAuth.GetIdentityFromContext(ctx)

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

	// Reload ALL resource translations
	// @todo we could probably do this more selectively and refresh only updated resources?
	_ = locale.Global().ReloadResourceTranslations(ctx)

	return nil
}

func (svc resourceTranslationsManager) Locale() locale.Resource {
	return svc.locale
}

func (svc resourceTranslationsManager) Module(ctx context.Context, namespaceID uint64, ID uint64) (locale.ResourceTranslationSet, error) {
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
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
		})

	}

	tmp, err := svc.moduleExtended(ctx, res)
	return append(out, tmp...), err
}

func (svc resourceTranslationsManager) Namespace(ctx context.Context, ID uint64) (locale.ResourceTranslationSet, error) {
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
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
		})

		k = types.LocaleKeyNamespaceSubtitle
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
		})

		k = types.LocaleKeyNamespaceDescription
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
		})

	}
	return out, nil
}

func (svc resourceTranslationsManager) Page(ctx context.Context, namespaceID uint64, ID uint64) (locale.ResourceTranslationSet, error) {
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
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
		})

		k = types.LocaleKeyPageDescription
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
		})

	}

	tmp, err := svc.pageExtended(ctx, res)
	return append(out, tmp...), err
}

func updateTranslations(ctx context.Context, ac localeAccessControl, lsvc ResourceTranslationsManagerService, tt ...*locale.ResourceTranslation) error {
	if lsvc == nil || lsvc.Locale() == nil || lsvc.Locale().Default() == nil {
		// gracefully handle partial initializations
		return nil
	}

	var (
		// assuming options will not change after start
		contentLang = lsvc.Locale().Default().Tag
	)

	if options.Locale().ResourceTranslationsEnabled {
		contentLang = locale.GetContentLanguageFromContext(ctx)
		// Resource translations enabled
		if contentLang == language.Und {
			// If no content-language meta (HTTP header) info was
			// used, do not run update translations - we do not know
			// what is the language that we're sending in
			return nil
		}

		if !lsvc.Locale().SupportedLang(contentLang) {
			// unsupported language
			return errors.InvalidData("unsupported language")
		}

		if !ac.CanManageResourceTranslations(ctx) {
			return errors.Unauthorized("not allowed to manage resource translations")
		}
	}

	locale.ResourceTranslationSet(tt).SetLanguage(contentLang)
	if err := lsvc.Upsert(ctx, tt); err != nil {
		return err
	}

	return nil
}
