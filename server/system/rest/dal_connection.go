package rest

import (
	"context"
	"io"
	"net/http"
	"time"

	federationService "github.com/cortezaproject/corteza/server/federation/service"
	federationTypes "github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/modern-go/reflect2"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	DalConnection struct {
		svc           connectionService
		federationSvc federationNodeService

		connectionAc connectionAccessController
	}

	connectionPayload struct {
		*types.DalConnection

		CanGrant            bool `json:"canGrant"`
		CanUpdateConnection bool `json:"canUpdateConnection"`
		CanDeleteConnection bool `json:"canDeleteConnection"`
		CanManageDalConfig  bool `json:"canManageDalConfig"`
	}

	connectionSetPayload struct {
		Filter types.DalConnectionFilter `json:"filter"`
		Set    []*connectionPayload      `json:"set"`
	}

	connectionAccessController interface {
		CanGrant(context.Context) bool
		CanCreateDalConnection(context.Context) bool
		CanUpdateDalConnection(context.Context, *types.DalConnection) bool
		CanDeleteDalConnection(context.Context, *types.DalConnection) bool
		CanManageDalConfigOnDalConnection(context.Context, *types.DalConnection) bool
	}

	connectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.DalConnection, error)
		Create(ctx context.Context, new *types.DalConnection) (*types.DalConnection, error)
		Update(ctx context.Context, upd *types.DalConnection) (*types.DalConnection, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.DalConnectionFilter) (types.DalConnectionSet, types.DalConnectionFilter, error)
	}

	federationNodeService interface {
		Search(ctx context.Context, filter federationTypes.NodeFilter) (set federationTypes.NodeSet, f federationTypes.NodeFilter, err error)
	}
)

func (DalConnection) New() *DalConnection {
	return &DalConnection{
		svc:           service.DefaultDalConnection,
		federationSvc: federationService.DefaultNode,

		connectionAc: service.DefaultAccessControl,
	}
}

func (ctrl DalConnection) List(ctx context.Context, r *request.DalConnectionList) (interface{}, error) {
	var (
		err            error
		dalConnections types.DalConnectionSet

		f = types.DalConnectionFilter{
			ConnectionID: payload.ParseUint64s(r.ConnectionID),
			Handle:       r.Handle,
			Type:         r.Type,

			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Deleted == 0 {
		f.Deleted = filter.StateExcluded
	}

	f.IncTotal = r.IncTotal

	dalConnections, f, err = ctrl.collectConnections(ctx, f)
	if err != nil {
		return nil, err
	}

	return ctrl.makeFilterPayload(ctx, dalConnections, f)
}

func (ctrl DalConnection) Create(ctx context.Context, r *request.DalConnectionCreate) (interface{}, error) {
	connection := &types.DalConnection{
		Handle: r.Handle,
		Type:   r.Type,
		Meta:   r.Meta,
		Config: r.Config,
	}

	res, err := ctrl.svc.Create(ctx, connection)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res), nil
}

func (ctrl DalConnection) Update(ctx context.Context, r *request.DalConnectionUpdate) (interface{}, error) {
	connection := &types.DalConnection{
		ID:     r.ConnectionID,
		Handle: r.Handle,
		Type:   r.Type,
		Meta:   r.Meta,
		Config: r.Config,
	}

	res, err := ctrl.svc.Update(ctx, connection)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res), nil
}

func (ctrl DalConnection) Read(ctx context.Context, r *request.DalConnectionRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ConnectionID)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res), nil
}

func (ctrl DalConnection) Delete(ctx context.Context, r *request.DalConnectionDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ConnectionID)
}

func (ctrl DalConnection) Undelete(ctx context.Context, r *request.DalConnectionUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.ConnectionID)
}

func (ctrl DalConnection) makeFilterPayload(ctx context.Context, connections types.DalConnectionSet, f types.DalConnectionFilter) (out *connectionSetPayload, err error) {
	out = &connectionSetPayload{
		Filter: f,
		Set:    make([]*connectionPayload, 0, len(connections)),
	}

	for _, c := range connections {
		out.Set = append(out.Set, ctrl.makePayload(ctx, c))
	}

	return
}

// Make payload for dal-connection
//
// Payload is without connection params on the config prop
//
// An explicit call to /params
func (ctrl DalConnection) makePayload(ctx context.Context, c *types.DalConnection) *connectionPayload {
	return &connectionPayload{
		DalConnection: c,

		CanGrant:            ctrl.connectionAc.CanGrant(ctx),
		CanUpdateConnection: ctrl.connectionAc.CanUpdateDalConnection(ctx, c),
		CanDeleteConnection: ctrl.connectionAc.CanDeleteDalConnection(ctx, c),
		CanManageDalConfig:  ctrl.connectionAc.CanManageDalConfigOnDalConnection(ctx, c),
	}
}

func (ctrl DalConnection) federatedNodeToConnection(f *federationTypes.Node) *types.DalConnection {
	h, _ := handle.Cast(nil, f.Name)

	return &types.DalConnection{
		ID: f.ID,

		Meta: types.ConnectionMeta{
			Name:      f.Name,
			Ownership: f.Contact,
		},

		Handle: h,
		Type:   federationTypes.NodeResourceType,

		//Config: types.ConnectionConfig{
		//	Connection: dal.NewFederatedNodeConnection(f.BaseURL, f.PairToken, f.AuthToken),
		//},

		CreatedAt: f.CreatedAt,
		CreatedBy: f.CreatedBy,
		UpdatedAt: f.UpdatedAt,
		UpdatedBy: f.UpdatedBy,
		DeletedAt: f.DeletedAt,
		DeletedBy: f.DeletedBy,
	}
}

func (ctrl DalConnection) collectConnections(ctx context.Context, filter types.DalConnectionFilter) (out types.DalConnectionSet, f types.DalConnectionFilter, err error) {
	var (
		dalConnections types.DalConnectionSet
		federatedNodes federationTypes.NodeSet
	)

	if dalConnections, f, err = ctrl.svc.Search(ctx, filter); err != nil {
		return nil, f, err
	}

	if !reflect2.IsNil(ctrl.federationSvc) {
		if federatedNodes, _, err = ctrl.federationSvc.Search(ctx, federationTypes.NodeFilter{
			// @todo IDs?
			Deleted: f.Deleted,
		}); err != nil {
			return nil, f, err
		}
	}

	out = append(out, dalConnections...)

	// We're converting federation nodes to DAL connection structs so that we have
	// a unified output.
	//
	// Eventually federation nodes will become connections, so this is ok
	for _, nn := range federatedNodes {
		out = append(out, ctrl.federatedNodeToConnection(nn))
	}

	out = ctrl.filterConnections(out, f)

	return
}

func (ctrl DalConnection) filterConnections(baseConnections types.DalConnectionSet, f types.DalConnectionFilter) (out types.DalConnectionSet) {
	for _, conn := range baseConnections {
		include := true

		if len(f.ConnectionID) > 0 {
			include = include && ctrl.inIDSet(f.ConnectionID, conn.ID)
		}

		if f.Handle != "" {
			include = include && f.Handle == conn.Handle
		}

		if f.Type != "" {
			include = include && f.Type == conn.Type
		}

		{
			if f.Deleted == filter.StateExcluded {
				include = include && conn.DeletedAt == nil
			}

			if f.Deleted == filter.StateExclusive {
				include = include && conn.DeletedAt != nil
			}
		}

		if include {
			out = append(out, conn)
		}
	}

	return
}

func (ctrl DalConnection) inIDSet(set []uint64, target uint64) (out bool) {
	for _, id := range set {
		out = out || id == target
	}

	return
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
