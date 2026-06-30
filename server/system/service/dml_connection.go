package service

import (
	"context"
	"fmt"
	"time"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	dmlDAL interface {
		SearchExternalModels(ctx context.Context, connectionID uint64) (dal.ModelSet, error)
		SearchExternalData(ctx context.Context, connectionID uint64, model *dal.Model, f filter.Filter) (dal.Iterator, error)
		ReplaceConnection(ctx context.Context, cw *dal.ConnectionWrap, isDefault bool) error
		RemoveConnection(ctx context.Context, ID uint64) error
	}

	dmlAC interface {
		CanSearchDalConnections(ctx context.Context) bool
		CanCreateDalConnection(ctx context.Context) bool
	}

	DmlConnectionSvc struct {
		store store.Storer
		dal   dmlDAL
		ac    dmlAC
	}
)

func NewDmlConnectionSvc(s store.Storer, d dmlDAL, ac dmlAC) *DmlConnectionSvc {
	return &DmlConnectionSvc{store: s, dal: d, ac: ac}
}

func (svc *DmlConnectionSvc) Create(ctx context.Context, in types.DmlConnectionInput) (*types.DmlConnection, error) {
	if !svc.ac.CanCreateDalConnection(ctx) {
		return nil, errors.Unauthorized("not allowed to manage DML connections")
	}
	if in.Handle == "" {
		return nil, fmt.Errorf("dml: connection handle required")
	}
	if in.Params == nil {
		return nil, fmt.Errorf("dml: connection params required")
	}

	label := in.Label
	if label == "" {
		label = in.Handle
	}
	if in.Params.ModelIdent == "" {
		in.Params.ModelIdent = "compose_records_{{namespace}}_{{module}}"
	}

	conn := &types.DmlConnection{
		ID:        id.Next(),
		Handle:    in.Handle,
		Label:     label,
		Params:    *in.Params,
		CreatedAt: time.Now(),
	}
	if err := store.CreateDmlConnection(ctx, svc.store, conn); err != nil {
		return nil, fmt.Errorf("dml: store connection: %w", err)
	}
	if err := svc.registerInDAL(ctx, conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func (svc *DmlConnectionSvc) Find(ctx context.Context, f types.DmlConnectionFilter) ([]*types.DmlConnection, error) {
	if !svc.ac.CanSearchDalConnections(ctx) {
		return nil, errors.Unauthorized("not allowed to access DML connections")
	}
	set, _, err := store.SearchDmlConnections(ctx, svc.store, f)
	return set, err
}

func (svc *DmlConnectionSvc) FindByID(ctx context.Context, connectionID uint64) (*types.DmlConnection, error) {
	if !svc.ac.CanSearchDalConnections(ctx) {
		return nil, errors.Unauthorized("not allowed to access DML connections")
	}
	conn, err := store.LookupDmlConnectionByID(ctx, svc.store, connectionID)
	if err != nil {
		return nil, err
	}
	if conn == nil || conn.DeletedAt != nil {
		return nil, fmt.Errorf("dml: connection %d not found", connectionID)
	}
	return conn, nil
}

// LoadAll registers all stored DML connections with the DAL. Called at startup, no AC check.
func (svc *DmlConnectionSvc) LoadAll(ctx context.Context) error {
	set, _, err := store.SearchDmlConnections(ctx, svc.store, types.DmlConnectionFilter{
		Deleted: filter.StateExcluded,
	})
	if err != nil {
		return err
	}
	for _, conn := range set {
		if err := svc.registerInDAL(ctx, conn); err != nil {
			return err
		}
	}
	return nil
}

func (svc *DmlConnectionSvc) FindModels(ctx context.Context, connectionID uint64, f types.DmlModelFilter) ([]*types.DmlModel, error) {
	if !svc.ac.CanSearchDalConnections(ctx) {
		return nil, errors.Unauthorized("not allowed to access DML connections")
	}
	models, err := svc.dal.SearchExternalModels(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	out := make([]*types.DmlModel, 0, len(models))
	for _, m := range models {
		out = append(out, dmlModelFromDAL(m))
	}
	return filterDmlModels(out, f), nil
}

func (svc *DmlConnectionSvc) FindModelByIdent(ctx context.Context, connectionID uint64, ident string) (*types.DmlModel, error) {
	models, err := svc.FindModels(ctx, connectionID, types.DmlModelFilter{Ident: []string{ident}})
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, fmt.Errorf("dml: model %q not found on connection %d", ident, connectionID)
	}
	return models[0], nil
}

func (svc *DmlConnectionSvc) registerInDAL(ctx context.Context, conn *types.DmlConnection) error {
	cw := dal.MakeConnection(conn.ID, nil, dal.ConnectionParams{
		Type:   conn.Params.Type,
		Params: conn.Params.Params,
	}, dal.ConnectionConfig{})
	if err := svc.dal.ReplaceConnection(ctx, cw, false); err != nil {
		return fmt.Errorf("dml: register connection in dal: %w", err)
	}
	return nil
}

func dmlModelFromDAL(m *dal.Model) *types.DmlModel {
	dm := &types.DmlModel{
		ConnectionID: m.ConnectionID,
		Ident:        m.Ident,
		Label:        m.Label,
		ResourceType: m.ResourceType,
	}
	for _, a := range m.Attributes {
		dm.Attributes = append(dm.Attributes, dmlAttributeFromDAL(a))
	}
	return dm
}

func dmlAttributeFromDAL(a *dal.Attribute) *types.DmlAttribute {
	da := &types.DmlAttribute{
		Ident:      a.Ident,
		Label:      a.Label,
		PrimaryKey: a.PrimaryKey,
		Sortable:   a.Sortable,
		Filterable: a.Filterable,
	}
	if a.Type != nil {
		da.Type = string(a.Type.Type())
	}
	if a.Store != nil {
		da.Store = string(a.Store.Type())
	}
	return da
}

func filterDmlModels(in []*types.DmlModel, f types.DmlModelFilter) []*types.DmlModel {
	if len(f.Ident) == 0 {
		return in
	}
	allowed := make(map[string]struct{}, len(f.Ident))
	for _, ident := range f.Ident {
		allowed[ident] = struct{}{}
	}
	out := make([]*types.DmlModel, 0, len(f.Ident))
	for _, m := range in {
		if _, ok := allowed[m.Ident]; ok {
			out = append(out, m)
		}
	}
	return out
}
