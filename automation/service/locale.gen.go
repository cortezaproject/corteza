package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - automation.workflow.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/types"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	resourceTranslationsManager struct {
		actionlog actionlog.Recorder
		locale    locale.Resource
		store     store.Storer
		ac        localeAccessController
	}

	localeAccessController interface {
		// CanManageResourceTranslation(context.Context) bool
	}

	ResourceTranslationsManagerService interface {
		Workflow(ctx context.Context, ID uint64) (locale.ResourceTranslationSet, error)

		Upsert(context.Context, locale.ResourceTranslationSet) error
		Locale() locale.Resource
	}
)

func ResourceTranslationsManager(ls locale.Resource) *resourceTranslationsManager {
	return &resourceTranslationsManager{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		locale:    ls,
	}
}

func (svc resourceTranslationsManager) Upsert(ctx context.Context, rr locale.ResourceTranslationSet) (err error) {
	// @todo AC
	//if (!svc.ac.CanManageResourceTranslation(ctx)) {
	//	return *****ErrNotAllowedToCreate()
	//}

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

func (svc resourceTranslationsManager) Workflow(ctx context.Context, ID uint64) (locale.ResourceTranslationSet, error) {
	var (
		err error
		out locale.ResourceTranslationSet
		res *types.Workflow
		k   types.LocaleKey
	)

	res, err = svc.loadWorkflow(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	for _, tag := range svc.locale.Tags() {
		k = types.LocaleKeyWorkflowName
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
		})

		k = types.LocaleKeyWorkflowDescription
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
		})

	}
	return out, nil
}
