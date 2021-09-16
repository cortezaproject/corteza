package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"golang.org/x/text/language"
)

type (
	localeResponse struct {
		Name string       `json:"name"`
		Self string       `json:"self"`
		Tag  language.Tag `json:"code"`
	}

	Locale struct{}
)

func (Locale) New() *Locale {
	ctrl := &Locale{}
	return ctrl
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
