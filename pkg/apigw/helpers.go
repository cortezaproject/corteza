package apigw

import (
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/options"
)

const (
	devHelperResponseBody string = `Hey developer!`
)

func helperDefaultResponse(opt *options.ApigwOpt) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if opt.LogEnabled {
			// Say something friendly when logging is enabled
			http.Error(w, devHelperResponseBody, http.StatusTeapot)
		} else {
			// Default 404 response
			http.Error(w, "", http.StatusNotFound)
		}
	}
}
