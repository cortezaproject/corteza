package dml

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// dalConnectionTypeDSN is the Config.DAL.Type value identifying an RDBMS
// (DSN-backed) connection — the only kind DML can introspect in v1.
const dalConnectionTypeDSN = "corteza::dal:connection:dsn"

type (
	dalReader interface {
		GetConnectionByID(connectionID uint64) *dal.ConnectionWrap
		SearchModels(ctx context.Context) (dal.ModelSet, error)
		SearchExternalModels(ctx context.Context, connectionID uint64) (dal.ModelSet, error)
		SearchConnectionIssues(connectionID uint64) []dal.Issue
	}

	// Connection projects pkg/dal connection + model state into API-facing DML types.
	Connection struct {
		store store.Storer
		dal   dalReader
	}
)

func NewConnection(s store.Storer, d dalReader) *Connection {
	return &Connection{store: s, dal: d}
}

// Find returns all RDBMS DalConnections, projected into DML connections,
// matching the filter.
func (r *Connection) Find(ctx context.Context, f types.DmlConnectionFilter) ([]*types.DmlConnection, error) {
	cc, err := r.listRdbmsConnections(ctx)
	if err != nil {
		return nil, err
	}
	return filterConnections(cc, f), nil
}

// FindByID returns the DML connection with the given ID, backed by the
// matching DalConnection record.
func (r *Connection) FindByID(ctx context.Context, connectionID uint64) (*types.DmlConnection, error) {
	c, err := store.LookupDalConnectionByID(ctx, r.store, connectionID)
	if err != nil {
		return nil, err
	}
	// Match the listing contract: only live RDBMS/DSN connections are visible.
	// LookupDalConnectionByID applies no deleted/type predicate, so guard here
	// — otherwise soft-deleted or non-DSN (REST/primary) connections would be
	// resolvable by ID even though Find never lists them.
	if c == nil || c.DeletedAt != nil || c.Config.DAL == nil || c.Config.DAL.Type != dalConnectionTypeDSN {
		return nil, fmt.Errorf("dml: connection %d not found", connectionID)
	}
	return dmlConnectionFromDal(c), nil
}

// listRdbmsConnections reads DalConnections straight from the store (so
// Config.DAL stays populated) and keeps only the RDBMS/DSN ones.
func (r *Connection) listRdbmsConnections(ctx context.Context) ([]*types.DmlConnection, error) {
	set, _, err := store.SearchDalConnections(ctx, r.store, types.DalConnectionFilter{
		Deleted: filter.StateExcluded,
	})
	if err != nil {
		return nil, err
	}

	out := make([]*types.DmlConnection, 0, len(set))
	for _, c := range set {
		if c.Config.DAL == nil || c.Config.DAL.Type != dalConnectionTypeDSN {
			continue
		}
		out = append(out, dmlConnectionFromDal(c))
	}
	return out, nil
}

// dmlConnectionFromDal projects a persisted DalConnection into a DML connection.
func dmlConnectionFromDal(c *types.DalConnection) *types.DmlConnection {
	out := &types.DmlConnection{
		ID:              c.ID,
		DalConnectionID: c.ID,
		Handle:          c.Handle,
		Type:            c.Type,
		Label:           c.Handle,
	}
	if c.Config.DAL != nil {
		// Driver here is the DAL connection-kind (dsn/rest), not the concrete
		// RDBMS driver — the real driver is parsed from the DSN at runtime.
		out.Driver = c.Config.DAL.Type
		out.ModelIdent = c.Config.DAL.ModelIdent
	}
	return out
}

// FindModels introspects the external connection's live schema and returns
// the resulting models. connectionID 0 resolves to the primary DAL connection.
// The connection must already be registered in the DAL pool, otherwise the
// lookup fails with "connection not found".
func (r *Connection) FindModels(ctx context.Context, connectionID uint64) ([]*types.DmlModel, error) {
	models, err := r.dal.SearchExternalModels(ctx, connectionID)
	if err != nil {
		return nil, err
	}

	out := make([]*types.DmlModel, 0, len(models))
	for _, m := range models {
		out = append(out, fromDalModel(m))
	}
	return out, nil
}

// FindModelByIdent resolves a single model on a connection by its ident.
func (r *Connection) FindModelByIdent(ctx context.Context, connectionID uint64, ident string) (*types.DmlModel, error) {
	models, err := r.FindModels(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	for _, m := range models {
		if m.Ident == ident {
			return m, nil
		}
	}
	return nil, fmt.Errorf("dml: model %q not found on connection %d", ident, connectionID)
}

func filterConnections(in []*types.DmlConnection, f types.DmlConnectionFilter) []*types.DmlConnection {
	if len(f.ConnectionID) == 0 && f.Handle == "" && f.Type == "" {
		return in
	}
	ids := map[string]struct{}{}
	for _, id := range f.ConnectionID {
		if id != "" {
			ids[id] = struct{}{}
		}
	}
	out := make([]*types.DmlConnection, 0, len(in))
	for _, c := range in {
		if f.Handle != "" && c.Handle != f.Handle {
			continue
		}
		if f.Type != "" && c.Type != f.Type {
			continue
		}
		if len(ids) > 0 {
			if _, ok := ids[fmt.Sprintf("%d", c.ID)]; !ok {
				continue
			}
		}
		out = append(out, c)
	}
	return out
}
