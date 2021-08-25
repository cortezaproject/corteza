package locale

import (
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/text/language"
)

const AcceptLanguageHeader = "Accept-Language"

func DetectLanguage(ll *Languages) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				rawLanguageTag string
			)

			// try to detect the language from the request's query string:
			if ll.Options().QueryStringParam != "" {
				rawLanguageTag = r.URL.Query().Get(ll.Options().QueryStringParam)
			}

			// try to detect the language from the request's headers:
			if len(rawLanguageTag) == 0 {
				rawLanguageTag = r.Header.Get(AcceptLanguageHeader)
			}

			if len(rawLanguageTag) > 0 {
				// parse & ignore errors
				var (
					preferred = ll.Default()
					supported = ll.Tags()

					accepted, _, err = language.ParseAcceptLanguage(rawLanguageTag)
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

				// new request with new context
				r = r.WithContext(SetLanguageToContext(r.Context(), preferred))

				ll.log.Debug(
					"language detected",
					zap.String("preferred", preferred.String()),
					zap.String("raw", rawLanguageTag),
				)
			}

			next.ServeHTTP(w, r)
		})
	}
}
