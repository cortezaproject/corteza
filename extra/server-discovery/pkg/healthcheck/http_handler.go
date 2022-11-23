package healthcheck

import "net/http"

func HttpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := Defaults().Run(r.Context())
		if results.Healthy() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		results.WriteTo(w)
	}
}
