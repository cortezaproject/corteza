package apigw

import (
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/apigw/profiler"
	h "github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/pkg/options"
)

const (
	devHelperResponseBody string = `Hey developer!`
)

func helperDefaultResponse(opt *options.ApigwOpt, pr *profiler.Profiler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addToProfiler(opt, pr, r, http.StatusNotFound)

		if opt.LogEnabled {
			// Say something friendly when logging is enabled
			http.Error(w, devHelperResponseBody, http.StatusTeapot)
		} else {
			// Default 404 response
			http.Error(w, "", http.StatusNotFound)
		}
	}
}

func helperMethodNotAllowed(opt *options.ApigwOpt, pr *profiler.Profiler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addToProfiler(opt, pr, r, http.StatusMethodNotAllowed)

		if opt.LogEnabled {
			// Say something friendly when logging is enabled
			http.Error(w, devHelperResponseBody, http.StatusTeapot)
		} else {
			// Default 405 response
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	}
}

func addToProfiler(opt *options.ApigwOpt, pr *profiler.Profiler, r *http.Request, status int) {
	if !(opt.ProfilerEnabled && opt.ProfilerGlobal) {
		return
	}

	// add to profiler
	ar, err := h.NewRequest(r)

	if err != nil {
		panic(err)
	}

	h := pr.Hit(ar)
	h.Status = status

	pr.Push(h)
}
