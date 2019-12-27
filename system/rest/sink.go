package rest

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/service"
)

var _ = errors.Wrap

type Sink struct {
	svc interface {
		Process(context.Context, string, *http.Request) error
	}

	sign auth.Signer
}

func (ctrl *Sink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()

		// What are we getting + part of the signature
		contentType = r.URL.Query().Get("content-type")

		//
		sign    = r.URL.Query().Get("sign")
		origin  = r.URL.Query().Get("origin")
		expires = r.URL.Query().Get("expires")
		method  = strings.ToUpper(r.Method)

		unsupported = func() {
			http.Error(w, "unsupported content-type", http.StatusBadRequest)
		}
	)

	if sign == "" {
		http.Error(w, "signature missing", http.StatusUnauthorized)
		return
	}

	if ctrl.sign.Verify(sign, 0, method, "/sink", contentType, origin, expires) {
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	if expires != "" {
		if exp, err := time.Parse("2006-01-02", expires); err != nil {
			http.Error(w, "could not process expiration date", http.StatusInternalServerError)
			return
		} else if exp.Before(time.Now()) {
			http.Error(w, "signature expired", http.StatusGone)
			return
		}
	}

	if contentType == "" {
		// If content-type not explicitly set (via QS),
		// try to get it from the headers
		contentType = r.Header.Get("content-type")
		if i := strings.Index(contentType, ";"); i > 0 {
			// intentionally > 0
			contentType = contentType[0 : i-1]
		}
	}

	if contentType == "" {
		unsupported()
		return
	}

	defer r.Body.Close()

	switch ctrl.svc.Process(ctx, contentType, r) {
	case service.ErrSinkContentProcessingFailed:
		http.Error(w, "sink processing failed", http.StatusInternalServerError)

	case service.ErrSinkContentTypeUnsupported:
		unsupported()
	}
}
