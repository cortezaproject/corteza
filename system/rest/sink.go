package rest

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/system/internal/service"
)

var _ = errors.Wrap

type Sink struct {
	// xxx service.XXXService
	svc interface {
		Process(context.Context, string, io.Reader) error
	}
}

func (ctrl *Sink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		cType = r.URL.Query().Get("content-type")

		unsupported = func() {
			http.Error(w, "unsupported content-type", http.StatusBadRequest)
		}
	)

	if cType == "" {
		// If content-type not explicitly set (via QS),
		// try to get it from the headers
		cType = r.Header.Get("content-type")
		if i := strings.Index(cType, ";"); i > 0 {
			// intentionally > 0
			cType = cType[0 : i-1]
		}
	}

	if cType == "" {
		unsupported()
		return
	}

	defer r.Body.Close()

	switch ctrl.svc.Process(ctx, cType, r.Body) {
	case service.ErrSinkContentProcessingFailed:
		http.Error(w, "sink processing failed", http.StatusInternalServerError)

	case service.ErrSinkContentTypeUnsupported:
		unsupported()
	}
}
