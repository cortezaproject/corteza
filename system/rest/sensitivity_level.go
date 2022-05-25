package rest

import (
	"context"
	"io"
	"net/http"
	"time"

	federationService "github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	SensitivityLevel struct {
		svc           sensitivityLevelService
		federationSvc federationNodeService
	}

	sensitivityLevelSetPayload struct {
		Filter types.DalSensitivityLevelFilter `json:"filter"`
		Set    types.DalSensitivityLevelSet    `json:"set"`
	}

	sensitivityLevelService interface {
		FindByID(ctx context.Context, ID uint64) (*types.DalSensitivityLevel, error)
		Create(ctx context.Context, new *types.DalSensitivityLevel) (*types.DalSensitivityLevel, error)
		Update(ctx context.Context, upd *types.DalSensitivityLevel) (*types.DalSensitivityLevel, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.DalSensitivityLevelFilter) (types.DalSensitivityLevelSet, types.DalSensitivityLevelFilter, error)
	}
)

func (SensitivityLevel) New() *SensitivityLevel {
	return &SensitivityLevel{
		svc:           service.DefaultDalSensitivityLevel,
		federationSvc: federationService.DefaultNode,
	}
}

func (ctrl SensitivityLevel) List(ctx context.Context, r *request.DalSensitivityLevelList) (interface{}, error) {
	var (
		err error
		set types.DalSensitivityLevelSet

		f = types.DalSensitivityLevelFilter{
			SensitivityLevelID: payload.ParseUint64s(r.SensitivityLevelID),

			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Deleted == 0 {
		f.Deleted = filter.StateInclusive
	}

	set, f, err = ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl SensitivityLevel) Create(ctx context.Context, r *request.DalSensitivityLevelCreate) (interface{}, error) {
	sensitivityLevel := &types.DalSensitivityLevel{
		Handle: r.Handle,
		Level:  r.Level,
		Meta:   r.Meta,
	}

	return ctrl.svc.Create(ctx, sensitivityLevel)
}

func (ctrl SensitivityLevel) Update(ctx context.Context, r *request.DalSensitivityLevelUpdate) (interface{}, error) {
	sensitivityLevel := &types.DalSensitivityLevel{
		ID:     r.SensitivityLevelID,
		Handle: r.Handle,
		Level:  r.Level,
		Meta:   r.Meta,
	}

	return ctrl.svc.Update(ctx, sensitivityLevel)
}

func (ctrl SensitivityLevel) Read(ctx context.Context, r *request.DalSensitivityLevelRead) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.SensitivityLevelID)
}

func (ctrl SensitivityLevel) Delete(ctx context.Context, r *request.DalSensitivityLevelDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.SensitivityLevelID)
}

func (ctrl SensitivityLevel) Undelete(ctx context.Context, r *request.DalSensitivityLevelUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.SensitivityLevelID)
}

func (ctrl SensitivityLevel) makeFilterPayload(ctx context.Context, ll types.DalSensitivityLevelSet, f types.DalSensitivityLevelFilter, err error) (*sensitivityLevelSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(ll) == 0 {
		ll = make([]*types.DalSensitivityLevel, 0)
	}

	return &sensitivityLevelSetPayload{Filter: f, Set: ll}, nil
}

func (ctrl SensitivityLevel) serve(ctx context.Context, fn string, archive io.ReadSeeker, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Disposition", "attachment; filename="+fn)

		http.ServeContent(w, req, fn, time.Now(), archive)
	}, nil
}
