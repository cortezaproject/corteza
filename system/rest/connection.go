package rest

import (
	"context"
	"io"
	"net/http"
	"time"

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
	Connection struct {
		svc          connectionService
		connectionAc connectionAccessController
	}

	connectionSetPayload struct {
		Filter types.ConnectionFilter `json:"filter"`
		Set    types.ConnectionSet    `json:"set"`
	}

	connectionAccessController interface {
		CanCreateConnection(context.Context) bool
		CanUpdateConnection(context.Context, *types.Connection) bool
	}

	connectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Connection, error)
		Create(ctx context.Context, new *types.Connection) (*types.Connection, error)
		Update(ctx context.Context, upd *types.Connection) (*types.Connection, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.ConnectionFilter) (types.ConnectionSet, types.ConnectionFilter, error)
	}
)

func (Connection) New() *Connection {
	return &Connection{
		svc:          service.DefaultConnection,
		connectionAc: service.DefaultAccessControl,
	}
}

func (ctrl Connection) List(ctx context.Context, r *request.ConnectionList) (interface{}, error) {
	var (
		err error
		set types.ConnectionSet
		f   = types.ConnectionFilter{
			ConnectionID: payload.ParseUint64s(r.ConnectionID),
			Handle:       r.Handle,
			Location:     r.Location,
			Ownership:    r.Ownership,

			Labels:  r.Labels,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	if f.Deleted == 0 {
		f.Deleted = filter.StateInclusive
	}

	set, f, err = ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl Connection) Create(ctx context.Context, r *request.ConnectionCreate) (interface{}, error) {
	connection := &types.Connection{
		Handle: r.Handle,

		DSN:       r.Dsn,
		Location:  r.Location,
		Ownership: r.Ownership,
		Sensitive: r.Sensitive,

		Config:       r.Config,
		Capabilities: r.Capabilities,

		Labels: r.Labels,
	}

	return ctrl.svc.Create(ctx, connection)
}

func (ctrl Connection) Update(ctx context.Context, r *request.ConnectionUpdate) (interface{}, error) {
	connection := &types.Connection{
		ID:     r.ConnectionID,
		Handle: r.Handle,

		DSN:       r.Dsn,
		Location:  r.Location,
		Ownership: r.Ownership,
		Sensitive: r.Sensitive,

		Config:       r.Config,
		Capabilities: r.Capabilities,

		Labels: r.Labels,
	}

	return ctrl.svc.Update(ctx, connection)
}

func (ctrl Connection) ReadPrimary(ctx context.Context, r *request.ConnectionReadPrimary) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, 0)
}

func (ctrl Connection) Read(ctx context.Context, r *request.ConnectionRead) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.ConnectionID)
}

func (ctrl Connection) Delete(ctx context.Context, r *request.ConnectionDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ConnectionID)
}

func (ctrl Connection) Undelete(ctx context.Context, r *request.ConnectionUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.ConnectionID)
}

func (ctrl Connection) makeFilterPayload(ctx context.Context, uu types.ConnectionSet, f types.ConnectionFilter, err error) (*connectionSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(uu) == 0 {
		uu = make([]*types.Connection, 0)
	}

	return &connectionSetPayload{Filter: f, Set: uu}, nil
}

func (ctrl Connection) serve(ctx context.Context, fn string, archive io.ReadSeeker, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Disposition", "attachment; filename="+fn)

		http.ServeContent(w, req, fn, time.Now(), archive)
	}, nil
}
