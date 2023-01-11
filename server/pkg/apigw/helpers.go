package apigw

import (
	"net/http"

	"github.com/cortezaproject/corteza/server/pkg/apigw/profiler"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	h "github.com/cortezaproject/corteza/server/pkg/http"
	"go.uber.org/zap"
)

func helperDefaultResponse(cfg types.Config, pr *profiler.Profiler, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addToProfiler(cfg, pr, log, r, http.StatusNotFound)
		http.Error(w, "", http.StatusNotFound)
	}
}

func helperMethodNotAllowed(cfg types.Config, pr *profiler.Profiler, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addToProfiler(cfg, pr, log, r, http.StatusMethodNotAllowed)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func addToProfiler(cfg types.Config, pr *profiler.Profiler, log *zap.Logger, r *http.Request, status int) {
	if !(cfg.Profiler.Enabled && cfg.Profiler.Global) {
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
