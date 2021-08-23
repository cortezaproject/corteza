package locale

import (
	"net/http"

	"golang.org/x/text/language"
)

const AcceptLanguageHeader = "Accept-Language"

func DetectLanguage(ll *Languages) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// parse & ignore errors
			var (
				preferred = ll.Default()
				supported = ll.Tags()

				accepted, _, err = language.ParseAcceptLanguage(r.Header.Get(AcceptLanguageHeader))
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

			next.ServeHTTP(w, r)
		})
	}
}
