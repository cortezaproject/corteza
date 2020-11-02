package rest

import (
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/pkg/errors"
	"net/http"
)

var _ = errors.Wrap

type Sink struct {
	svc interface {
		ProcessRequest(w http.ResponseWriter, r *http.Request)
	}

	sign auth.Signer
}

func (ctrl *Sink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctrl.svc.ProcessRequest(w, r)
}
