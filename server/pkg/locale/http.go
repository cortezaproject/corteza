package locale

import (
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/text/language"
)

const AcceptLanguageHeader = "Accept-Language"
const ContentLanguageHeader = "Content-Language"

func DetectLanguage(ll *service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			// resolve accept-language header
			// Accept-Language specifies the language of the response payload.
			ctx = SetAcceptLanguageToContext(ctx, resolveAcceptLanguageHeaders(r, ll))

			// resolve content-language header
			// Content-Language specifies the language of the request payload.
			ctx = SetContentLanguageToContext(ctx, resolveContentLanguageHeaders(r.Header, ll))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// reads Content-Language headers from the request and returns
// parsed value as language.Tag
//
// There are 4 valid scenarios for :
//  - lang == 'skip':   returns und & services will (likely) ignore all translatable content
//  - invalid language: (same as skip)
//  - valid language:   returns valid language; services will treat translatable content from the payload as translations
//  - no header:        returns default language; (same as valid language)
func resolveContentLanguageHeaders(h http.Header, ll *service) language.Tag {
	if ll == nil || ll.Default() == nil {
		return language.Und
	}

	var cLang = h.Get(ContentLanguageHeader)

	if cLang == "" {
		return ll.Default().Tag
	}

	if cLang == "skip" {
		// more than 1 header or value equal to skip
		return language.Und
	}

	if tag, err := language.Parse(cLang); err != nil {
		return language.Und
	} else {
		return tag
	}
}

func resolveAcceptLanguageHeaders(r *http.Request, ll *service) (tag language.Tag) {
	if ll == nil || ll.Default() == nil {
		// locale service does not have anything loaded...
		return language.Und
	}

	if ll.opt.DevelopmentMode {
		if err := ll.ReloadStatic(); err != nil {
			// when in development mode, refresh languages for every request
			ll.log.Error("failed to load locales", zap.Error(err))
			return
		}
	}

	var (
		raw string
	)

	// try to detect the language from the request's query string:
	if ll.opt.QueryStringParam != "" {
		raw = r.URL.Query().Get(ll.opt.QueryStringParam)
	}

	// try to detect the language from the request's headers:
	if len(raw) == 0 {
		raw = r.Header.Get(AcceptLanguageHeader)
	}

	if len(raw) == 0 {
		// no need for lang detection
		return ll.Default().Tag
	}

	// parse & ignore errors
	var (
		preferred = ll.Default().Tag
		supported = ll.Tags()

		accepted, _, err = language.ParseAcceptLanguage(raw)
	)

	if err == nil {
		// ignoring index & confidence
		preferred, _, _ = language.NewMatcher(supported).Match(accepted...)

		var match bool
		for _, s := range supported {
			if s == preferred {
				match = true
				break
			}
		}

		if !match {
			base, _ := preferred.Base()
			preferred = language.MustParse(base.String())
		}
	}

	if ll.opt.DevelopmentMode {
		ll.log.Debug(
			"language detected",
			zap.String("preferred", preferred.String()),
			zap.String("raw", raw),
		)
	}

	// new request with new context
	return preferred

}
