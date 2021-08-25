package locale

import (
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/text/language"
)

const AcceptLanguageHeader = "Accept-Language"

func DetectLanguage(ll *service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, upgradeRequest(ll, r))
		})
	}
}

func upgradeRequest(ll *service, r *http.Request) *http.Request {
	return r.WithContext(SetLanguageToContext(r.Context(), detectLanguage(ll, r)))
}

func detectLanguage(ll *service, r *http.Request) (tag language.Tag) {
	if ll.opt.DevelopmentMode {
		if err := ll.Reload(); err != nil {
			// when in development mode, refresh languages for every request
			ll.log.Error("failed to load locales", zap.Error(err))
			return
		}
	}

	if ll.Default() == nil {
		// locale service does not have anything loaded...
		return
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
		return
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
