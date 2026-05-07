package rest

import (
	"context"
	"encoding/json"

	"github.com/crusttech/human/server/pkg/api"
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

func (ctrl Connection) List(ctx context.Context, r *request.ConnectionList) (interface{}, error) {
	var (
		err error
		set types.ConnectionSet

		f = types.ConnectionFilter{
			Handle: r.Handle,
			Status: r.Status,
			Query:  r.Query,
			Tags:   r.Tags,

			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Deleted == 0 {
		f.Deleted = filter.StateExcluded
	}

	f.IncTotal = r.IncTotal

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	set, f, err = ctrl.svc.Search(ctx, f)
	if err != nil {
		return nil, err
	}

	return ctrl.makeFilterPayload(ctx, set, f)
}

func (ctrl Connection) Create(ctx context.Context, r *request.ConnectionCreate) (interface{}, error) {
	connection := &types.Connection{
		Handle:  r.Handle,
		Meta:    r.Meta,
		Service: r.Service,
		Status:  "draft",
	}

	if r.Resources != nil {
		_ = json.Unmarshal(r.Resources, &connection.Resources)
	}
	if r.Operations != nil {
		_ = json.Unmarshal(r.Operations, &connection.Operations)
	}

	res, err := ctrl.svc.Create(ctx, connection)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res), nil
}

func (ctrl Connection) Update(ctx context.Context, r *request.ConnectionUpdate) (interface{}, error) {
	connection := &types.Connection{
		ID:        r.ConnectionID,
		Handle:    r.Handle,
		Meta:      r.Meta,
		Service:   r.Service,
		UpdatedAt: r.UpdatedAt,
	}

	if r.Resources != nil {
		_ = json.Unmarshal(r.Resources, &connection.Resources)
	}
	if r.Operations != nil {
		_ = json.Unmarshal(r.Operations, &connection.Operations)
	}

	res, err := ctrl.svc.Update(ctx, connection)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res), nil
}

func (ctrl Connection) Read(ctx context.Context, r *request.ConnectionRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ConnectionID)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, res), nil
}

func (ctrl Connection) Delete(ctx context.Context, r *request.ConnectionDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ConnectionID)
}

func (ctrl Connection) Undelete(ctx context.Context, r *request.ConnectionUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.ConnectionID)
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
	return ctrl.makePayload(ctx, res), nil
}

func (ctrl Connection) Configure(ctx context.Context, r *request.ConnectionConfigure) (interface{}, error) {
	connectionID := r.ConnectionID

	// Auto-import catalog connection on first configure.
	if connectionID == 0 && r.CatalogID != "" {
		imported, err := ctrl.svc.Import(ctx, r.CatalogID)
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

func (ctrl Connection) makeFilterPayload(ctx context.Context, set types.ConnectionSet, f types.ConnectionFilter) (*connectionSetPayload, error) {
	out := &connectionSetPayload{
		Filter: f,
		Set:    make([]*connectionPayload, 0, len(set)),
	}

	for _, c := range set {
		out.Set = append(out.Set, ctrl.makePayload(ctx, c))
	}

	return out, nil
}

func (ctrl Connection) makePayload(ctx context.Context, c *types.Connection) *connectionPayload {
	return &connectionPayload{
		Connection:          c,
		CanUpdateConnection: ctrl.ac.CanUpdateConnection(ctx, c),
		CanDeleteConnection: ctrl.ac.CanDeleteConnection(ctx, c),
	}
}

func (ctrl Connection) makeConfigurePayload(ctx context.Context, c *types.ConfiguredConnection) *configuredConnectionPayload {
	return &configuredConnectionPayload{
		ConfiguredConnection:          c,
		CanUpdateConfiguredConnection: ctrl.ac.CanUpdateConfiguredConnection(ctx, c),
		CanDeleteConfiguredConnection: ctrl.ac.CanDeleteConfiguredConnection(ctx, c),
	}
}
