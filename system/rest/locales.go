package rest

import (
	"context"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
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
	var (
		rval        = make([]*localeResponse, 0, 16)
		currentLang = locale.GetLanguageFromContext(ctx)
	)

	for _, l := range locale.Global().List() {
		rval = append(rval, &localeResponse{
			Name: display.Languages(l.Tag).Name(l.Tag),
			Self: display.Languages(currentLang).Name(l.Tag),
			Tag:  l.Tag,
		})
	}

	return rval, nil
}

func (ctrl Locale) Get(ctx context.Context, r *request.LocaleGet) (interface{}, error) {
	svc := locale.Global()

	return func(w http.ResponseWriter, req *http.Request) {
		if !svc.HasLanguage(language.Make(r.Lang)) {
			// @todo temp workaround until frontend knows what languages it can use
			r.Lang = svc.List()[0].Tag.String()
		}

		if !svc.HasLanguage(language.Make(r.Lang)) {
			errors.ProperlyServeHTTP(w, req, errors.New(
				errors.KindNotFound,
				"no such language",
				errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
			), true)
			return
		}

		if !svc.HasApplication(language.Make(r.Lang), r.Application) {
			errors.ProperlyServeHTTP(w, req, errors.New(
				errors.KindNotFound,
				"no such application",
				errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
			), true)
			return
		}

		if err := locale.Global().EncodeExternal(w, language.Make(r.Lang), r.Application); err != nil {
			errors.ProperlyServeHTTP(w, req, err, true)
		}

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	}, nil
}
