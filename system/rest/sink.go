package rest

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/auth"
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
