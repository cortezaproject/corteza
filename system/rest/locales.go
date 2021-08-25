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
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		err := locale.Global().EncodeExternal(w, language.Make(r.Lang), r.Application)
		if err != nil {
			errors.ProperlyServeHTTP(w, req, err, false)
		}
	}, nil
}
