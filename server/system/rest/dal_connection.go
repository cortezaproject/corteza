package rest

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	DalConnection struct {
		svc dalConnectionService

		connectionAc dalConnectionAccessController
	}

	dalConnectionPayload struct {
		*types.DalConnection

		CanGrant            bool `json:"canGrant"`
		CanUpdateConnection bool `json:"canUpdateConnection"`
		CanDeleteConnection bool `json:"canDeleteConnection"`
		CanManageDalConfig  bool `json:"canManageDalConfig"`
	}

	dalConnectionSetPayload struct {
		Filter types.DalConnectionFilter `json:"filter"`
		Set    []*dalConnectionPayload   `json:"set"`
	}

	dalConnectionAccessController interface {
		CanGrant(context.Context) bool
		CanCreateDalConnection(context.Context) bool
		CanUpdateDalConnection(context.Context, *types.DalConnection) bool
		CanDeleteDalConnection(context.Context, *types.DalConnection) bool
		CanManageDalConfigOnDalConnection(context.Context, *types.DalConnection) bool
	}

	dalConnectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.DalConnection, error)
		Create(ctx context.Context, new *types.DalConnection) (*types.DalConnection, error)
		Update(ctx context.Context, upd *types.DalConnection) (*types.DalConnection, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.DalConnectionFilter) (types.DalConnectionSet, types.DalConnectionFilter, error)
	}
)

func (DalConnection) New() *DalConnection {
	return &DalConnection{
		svc: service.DefaultDalConnection,

		connectionAc: service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl DalConnection) makeFilter(ctx context.Context, r *request.DalConnectionList) (f types.DalConnectionFilter, err error) {
	f = types.DalConnectionFilter{
		DalConnectionID: r.ConnectionID,
		Handle:          r.Handle,
		Type:            r.Type,
		Query:           r.Query,

		Deleted: filter.State(r.Deleted),
	}

	if f.Deleted == 0 {
		f.Deleted = filter.StateExcluded
	}

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl DalConnection) beforeCreate(ctx context.Context, res *types.DalConnection, r *request.DalConnectionCreate) error {
	res.Meta = r.Meta
	res.Config = r.Config
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl DalConnection) beforeUpdate(ctx context.Context, res *types.DalConnection, r *request.DalConnectionUpdate) error {
	res.Meta = r.Meta
	res.Config = r.Config
	return nil
}

func (ctrl DalConnection) makeFilterPayload(ctx context.Context, connections types.DalConnectionSet, f types.DalConnectionFilter, err error) (*dalConnectionSetPayload, error) {
	if err != nil {
		return nil, err
	}

	out := &dalConnectionSetPayload{
		Filter: f,
		Set:    make([]*dalConnectionPayload, 0, len(connections)),
	}

	for _, c := range connections {
		p, _ := ctrl.makePayload(ctx, c, nil)
		out.Set = append(out.Set, p)
	}

	return out, nil
}

// Make payload for dal-connection
//
// # Payload is without connection params on the config prop
//
// An explicit call to /params
func (ctrl DalConnection) makePayload(ctx context.Context, c *types.DalConnection, err error) (*dalConnectionPayload, error) {
	if err != nil || c == nil {
		return nil, err
	}

	return &dalConnectionPayload{
		DalConnection: c,

		CanGrant:            ctrl.connectionAc.CanGrant(ctx),
		CanUpdateConnection: ctrl.connectionAc.CanUpdateDalConnection(ctx, c),
		CanDeleteConnection: ctrl.connectionAc.CanDeleteDalConnection(ctx, c),
		CanManageDalConfig:  ctrl.connectionAc.CanManageDalConfigOnDalConnection(ctx, c),
	}, nil
}

func (ctrl DalConnection) serve(ctx context.Context, fn string, archive io.ReadSeeker, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Disposition", "attachment; filename="+fn)

		http.ServeContent(w, req, fn, time.Now(), archive)
	}, nil
}
