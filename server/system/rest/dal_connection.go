package rest

import (
	"context"
	"io"
	"net/http"
	"time"

	federationService "github.com/crusttech/human/server/federation/service"
	federationTypes "github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
	"github.com/modern-go/reflect2"
)

type (
	DalConnection struct {
		svc           dalConnectionService
		federationSvc federationNodeService

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

// makeFilter builds the search filter for the generated List controller.
func (ctrl DalConnection) makeFilter(ctx context.Context, r *request.DalConnectionList) (types.DalConnectionFilter, error) {
	f := types.DalConnectionFilter{
		DalConnectionID: r.ConnectionID,
		Handle:          r.Handle,
		Type:            r.Type,

		Deleted: filter.State(r.Deleted),
	}

	if f.Deleted == 0 {
		f.Deleted = filter.StateExcluded
	}

	f.IncTotal = r.IncTotal

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

	// Merge federation nodes into the base (service-searched) set and apply
	// the in-memory filtering, preserving the original List behavior.
	if connections, f, err = ctrl.collectConnections(ctx, connections, f); err != nil {
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

func (ctrl DalConnection) federatedNodeToConnection(f *federationTypes.Node) *types.DalConnection {
	h, _ := handle.Cast(nil, f.Name)

	return &types.DalConnection{
		ID: f.ID,

		Meta: types.DalConnectionMeta{
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

// collectConnections merges federation nodes into the base (already
// service-searched) connection set and applies the in-memory filtering.
func (ctrl DalConnection) collectConnections(ctx context.Context, dalConnections types.DalConnectionSet, f types.DalConnectionFilter) (out types.DalConnectionSet, _ types.DalConnectionFilter, err error) {
	var (
		federatedNodes federationTypes.NodeSet
	)

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

	return out, f, nil
}

func (ctrl DalConnection) filterConnections(baseConnections types.DalConnectionSet, f types.DalConnectionFilter) (out types.DalConnectionSet) {
	for _, conn := range baseConnections {
		include := true

		if len(f.DalConnectionID) > 0 {
			include = include && ctrl.inIDSet(id.Uints(f.DalConnectionID...), conn.ID)
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
