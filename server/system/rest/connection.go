package rest

import (
	"context"
	"encoding/json"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	Connection struct {
		svc connectionService
		ac  connectionAccessController
	}

	connectionPayload struct {
		*types.Connection

		CanUpdateConnection bool `json:"canUpdateConnection"`
		CanDeleteConnection bool `json:"canDeleteConnection"`
	}

	connectionSetPayload struct {
		Filter types.ConnectionFilter `json:"filter"`
		Set    []*connectionPayload   `json:"set"`
	}

	connectionAccessController interface {
		CanCreateConnection(context.Context) bool
		CanUpdateConnection(context.Context, *types.Connection) bool
		CanDeleteConnection(context.Context, *types.Connection) bool

		CanUpdateConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
		CanDeleteConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
	}

	connectionService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Connection, error)
		Create(ctx context.Context, new *types.Connection) (*types.Connection, error)
		Update(ctx context.Context, upd *types.Connection) (*types.Connection, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.ConnectionFilter) (types.ConnectionSet, types.ConnectionFilter, error)
		Import(ctx context.Context, catalogID string) (*types.Connection, error)
		Enable(ctx context.Context, ID uint64) (*types.Connection, error)
		Configure(ctx context.Context, new *types.ConfiguredConnection) (*types.ConfiguredConnection, error)
		UpdateConfiguration(ctx context.Context, upd *types.ConfiguredConnection) (*types.ConfiguredConnection, error)
	}
)

func (Connection) New() *Connection {
	return &Connection{
		svc: service.DefaultConnection,
		ac:  service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl Connection) makeFilter(ctx context.Context, r *request.ConnectionList) (types.ConnectionFilter, error) {
	var (
		err error

		f = types.ConnectionFilter{
			Handle: r.Handle,
			Status: r.Status,
			Query:  r.Query,
			Tags:   r.Tags,
			Source: r.Source,

			Deleted: filter.State(r.Deleted),
		}
	)

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
func (ctrl Connection) beforeCreate(ctx context.Context, res *types.Connection, r *request.ConnectionCreate) error {
	res.Meta = r.Meta
	res.Service = r.Service
	res.Status = "draft"

	if r.Resources != nil {
		_ = json.Unmarshal(r.Resources, &res.Resources)
	}
	if r.Operations != nil {
		_ = json.Unmarshal(r.Operations, &res.Operations)
	}

	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl Connection) beforeUpdate(ctx context.Context, res *types.Connection, r *request.ConnectionUpdate) error {
	res.Meta = r.Meta
	res.Service = r.Service

	if r.Resources != nil {
		_ = json.Unmarshal(r.Resources, &res.Resources)
	}
	if r.Operations != nil {
		_ = json.Unmarshal(r.Operations, &res.Operations)
	}

	return nil
}

func (ctrl Connection) Enable(ctx context.Context, r *request.ConnectionEnable) (interface{}, error) {
	res, err := ctrl.svc.Enable(ctx, r.ConnectionID)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res, nil)
}

func (ctrl Connection) Generate(ctx context.Context, r *request.ConnectionGenerate) (interface{}, error) {
	// TODO: wire to connection-builder agent
	return nil, nil
}

func (ctrl Connection) Import(ctx context.Context, r *request.ConnectionImport) (interface{}, error) {
	res, err := ctrl.svc.Import(ctx, r.CatalogID)
	if err != nil {
		return nil, err
	}
	return ctrl.makePayload(ctx, res, nil)
}

func (ctrl Connection) Configure(ctx context.Context, r *request.ConnectionConfigure) (interface{}, error) {
	connectionID := r.ConnectionID

	// Auto-import catalog connection on first configure.
	// Two paths:
	//   1. Explicit catalog ID + connectionID == 0 (frontend already chose to import)
	//   2. Synthetic high-bit ID (frontend passed a catalog list entry's ID directly)
	if connectionID == 0 && r.CatalogID != "" {
		imported, err := ctrl.svc.Import(ctx, r.CatalogID)
		if err != nil {
			return nil, err
		}
		connectionID = imported.ID
	} else if connectionID&(1<<63) != 0 {
		found, err := ctrl.svc.FindByID(ctx, connectionID)
		if err != nil {
			return nil, err
		}
		imported, err := ctrl.svc.Import(ctx, found.CatalogID)
		if err != nil {
			return nil, err
		}
		connectionID = imported.ID
	}

	conn := &types.ConfiguredConnection{
		ConnectionID: connectionID,
		Name:         r.Name,
		Config:       r.Config,
		Status:       "draft",
		Labels:       r.Labels,
	}

	res, err := ctrl.svc.Configure(ctx, conn)
	if err != nil {
		return nil, err
	}

	return ctrl.makeConfigurePayload(ctx, res), nil
}

func (ctrl Connection) UpdateConfiguration(ctx context.Context, r *request.ConnectionUpdateConfiguration) (interface{}, error) {
	conn := &types.ConfiguredConnection{
		ID:           r.ConfiguredConnectionID,
		ConnectionID: r.ConnectionID,
		Name:         r.Name,
		Config:       r.Config,
		Status:       "draft",
		Labels:       r.Labels,
	}

	res, err := ctrl.svc.UpdateConfiguration(ctx, conn)
	if err != nil {
		return nil, err
	}

	return ctrl.makeConfigurePayload(ctx, res), nil
}

func (ctrl Connection) makeFilterPayload(ctx context.Context, set types.ConnectionSet, f types.ConnectionFilter, err error) (*connectionSetPayload, error) {
	if err != nil {
		return nil, err
	}

	out := &connectionSetPayload{
		Filter: f,
		Set:    make([]*connectionPayload, 0, len(set)),
	}

	for _, c := range set {
		p, _ := ctrl.makePayload(ctx, c, nil)
		out.Set = append(out.Set, p)
	}

	return out, nil
}

func (ctrl Connection) makePayload(ctx context.Context, c *types.Connection, err error) (*connectionPayload, error) {
	if err != nil || c == nil {
		return nil, err
	}

	return &connectionPayload{
		Connection:          c,
		CanUpdateConnection: ctrl.ac.CanUpdateConnection(ctx, c),
		CanDeleteConnection: ctrl.ac.CanDeleteConnection(ctx, c),
	}, nil
}

func (ctrl Connection) makeConfigurePayload(ctx context.Context, c *types.ConfiguredConnection) *configuredConnectionPayload {
	return &configuredConnectionPayload{
		ConfiguredConnection:          c,
		CanUpdateConfiguredConnection: ctrl.ac.CanUpdateConfiguredConnection(ctx, c),
		CanDeleteConfiguredConnection: ctrl.ac.CanDeleteConfiguredConnection(ctx, c),
	}
}
