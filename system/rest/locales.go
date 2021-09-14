package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"golang.org/x/text/language"
)

type (
	localeResponse struct {
		Name string       `json:"name"`
		Self string       `json:"self"`
		Tag  language.Tag `json:"code"`
	}

	Locale struct {
		svc service.ResourceTranslationService
	}

	resourceTranslationSetPayload struct {
		Filter types.ResourceTranslationFilter `json:"filter"`
		Set    []*resourceTranslationPayload   `json:"set"`
	}

	resourceTranslationPayload struct {
		*types.ResourceTranslation
	}
)

func (Locale) New() *Locale {
	return &Locale{
		svc: service.DefaultResourceTranslation,
	}
}

func (ctrl Locale) List(ctx context.Context, r *request.LocaleList) (interface{}, error) {
	return locale.Global().LocalizedList(ctx), nil
}

func (ctrl Locale) Get(ctx context.Context, r *request.LocaleGet) (interface{}, error) {
	svc := locale.Global()

	// We're using + as a language delimiter
	// because webapp client i18n lib does this by default
	// (url encoded space becomes a +)
	const langSplit = "+"

	return func(w http.ResponseWriter, req *http.Request) {
		if len(svc.List()) == 0 {
			errors.ProperlyServeHTTP(w, req, errors.New(
				errors.KindNotFound,
				"no languages found",
				errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
			), true)
			return
		}

		// default to 1st language
		//var def = svc.List()[0].Tag
		var ll = []language.Tag{}

		for _, candidate := range strings.Split(r.Lang, langSplit) {
			ll = append(ll, language.Make(candidate))
			//if svc.HasLanguage(tmp) {
			//	ll = append(ll, lang)
			//	break
			//}
		}

		//if !svc.HasLanguage(lang) {
		//	errors.ProperlyServeHTTP(w, req, errors.New(
		//		errors.KindNotFound,
		//		"no such language",
		//		errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
		//	), true)
		//	return
		//}

		//if !svc.HasApplication(lang, r.Application) {
		//	errors.ProperlyServeHTTP(w, req, errors.New(
		//		errors.KindNotFound,
		//		"no such application",
		//		errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
		//	), true)
		//	return
		//}

		if err := locale.Global().EncodeExternal(w, r.Application, ll...); err != nil {
			errors.ProperlyServeHTTP(w, req, err, true)
		}

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	}, nil
}

func (ctrl Locale) ListResource(ctx context.Context, r *request.LocaleListResource) (interface{}, error) {
	var (
		err error
		f   = types.ResourceTranslationFilter{
			Lang:         r.Lang,
			Resource:     r.Resource,
			ResourceType: r.ResourceType,
			OwnerID:      r.OwnerID,
			Deleted:      filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.List(ctx, f)
	return ctrl.makeResourceTranslationSetPayload(ctx, set, filter, err)
}

func (ctrl Locale) CreateResource(ctx context.Context, r *request.LocaleCreateResource) (interface{}, error) {
	var (
		err error
		app = &types.ResourceTranslation{
			Lang:     types.Lang{Tag: language.Make(r.Lang)},
			Resource: r.Resource,
			K:        r.Key,
			Message:  r.Message,
			OwnedBy:  r.OwnerID,
		}
	)

	app, err = ctrl.svc.Create(ctx, app)
	return ctrl.makeResourceTranslationPayload(ctx, app, err)
}

func (ctrl Locale) UpdateResource(ctx context.Context, r *request.LocaleUpdateResource) (interface{}, error) {
	var (
		err error
		app = &types.ResourceTranslation{
			ID:       r.TranslationID,
			Lang:     types.Lang{Tag: language.Make(r.Lang)},
			Resource: r.Resource,
			K:        r.Key,
			Message:  r.Message,
			OwnedBy:  r.OwnerID,
		}
	)

	app, err = ctrl.svc.Update(ctx, app)
	return ctrl.makeResourceTranslationPayload(ctx, app, err)
}

func (ctrl Locale) ReadResource(ctx context.Context, r *request.LocaleReadResource) (interface{}, error) {
	tpl, err := ctrl.svc.Read(ctx, r.TranslationID)
	return ctrl.makeResourceTranslationPayload(ctx, tpl, err)
}

func (ctrl Locale) DeleteResource(ctx context.Context, r *request.LocaleDeleteResource) (interface{}, error) {
	return api.OK(), ctrl.svc.Delete(ctx, r.TranslationID)
}

func (ctrl Locale) UndeleteResource(ctx context.Context, r *request.LocaleUndeleteResource) (interface{}, error) {
	return api.OK(), ctrl.svc.Undelete(ctx, r.TranslationID)
}

// Utilities

func (ctrl Locale) makeResourceTranslationSetPayload(ctx context.Context, tt types.ResourceTranslationSet, f types.ResourceTranslationFilter, err error) (*resourceTranslationSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &resourceTranslationSetPayload{Filter: f, Set: make([]*resourceTranslationPayload, len(tt))}

	for i := range tt {
		msp.Set[i], _ = ctrl.makeResourceTranslationPayload(ctx, tt[i], nil)
	}

	return msp, nil
}

func (ctrl Locale) makeResourceTranslationPayload(ctx context.Context, rt *types.ResourceTranslation, err error) (*resourceTranslationPayload, error) {
	if err != nil || rt == nil {
		return nil, err
	}

	pl := &resourceTranslationPayload{
		ResourceTranslation: rt,
	}

	return pl, nil
}
