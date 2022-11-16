package apigw

import (
	"net/http"

	"github.com/cortezaproject/corteza/server/pkg/apigw/profiler"
	h "github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"go.uber.org/zap"
)

const (
	devHelperResponseBody string = `Hey developer!`
)

func helperDefaultResponse(opt options.ApigwOpt, pr *profiler.Profiler, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addToProfiler(opt, pr, log, r, http.StatusNotFound)

		responseBody := ""

		if opt.LogEnabled {
			// Say something friendly when logging is enabled
			responseBody = devHelperResponseBody
		}

		http.Error(w, responseBody, http.StatusNotFound)
	}
}

func helperMethodNotAllowed(opt options.ApigwOpt, pr *profiler.Profiler, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addToProfiler(opt, pr, log, r, http.StatusMethodNotAllowed)

		if opt.LogEnabled {
			// Say something friendly when logging is enabled
			http.Error(w, devHelperResponseBody, http.StatusTeapot)
		} else {
			// Default 405 response
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	}
}

func addToProfiler(opt options.ApigwOpt, pr *profiler.Profiler, log *zap.Logger, r *http.Request, status int) {
	if !(opt.ProfilerEnabled && opt.ProfilerGlobal) {
		return
	}

	// add to profiler
	ar, err := h.NewRequest(r)

	if err != nil {
		log.Warn("could not create request wrapper, not adding to profiler")
		return
	}

	h := pr.Hit(ar)
	h.Status = status

	pr.Push(h)
}
