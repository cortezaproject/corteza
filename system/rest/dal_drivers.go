package rest

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	DalDriver struct{}

	dalDriverSetPayload struct {
		Set []dal.Driver `json:"set"`
	}
)

func (DalDriver) New() *DalDriver {
	return &DalDriver{}
}

func (ctrl DalDriver) List(ctx context.Context, r *request.DalDriverList) (interface{}, error) {
	return ctrl.makeFilterPayload(ctx, dal.Service().Drivers())
}

func (ctrl DalDriver) makeFilterPayload(ctx context.Context, drivers []dal.Driver) (out *dalDriverSetPayload, err error) {
	out = &dalDriverSetPayload{
		Set: drivers,
	}

	return
}

func (ctrl DalDriver) serve(ctx context.Context, fn string, archive io.ReadSeeker, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Disposition", "attachment; filename="+fn)

		http.ServeContent(w, req, fn, time.Now(), archive)
	}, nil
}
