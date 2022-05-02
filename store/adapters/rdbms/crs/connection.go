package crs

import (
	"context"
	"sync"

	"github.com/cortezaproject/corteza-server/compose/crs"
	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/jmoiron/sqlx"
)

type (
	connection struct {
		mux    sync.RWMutex
		models map[string]*model

		db      *sqlx.DB
		dialect drivers.Dialect
	}
)

func Connection(db *sqlx.DB, dialect drivers.Dialect) *connection {
	return &connection{
		db:      db,
		dialect: dialect,
		models:  make(map[string]*model),
	}
}

func (c *connection) model(m *data.Model) *model {
	c.mux.RLock()
	if c.models[m.Ident] == nil {
		c.mux.RUnlock()
		c.mux.Lock()
		c.models[m.Ident] = Model(m, c.db, c.dialect)
		defer c.mux.Unlock()
		return c.models[m.Ident]
	}

	defer c.mux.RUnlock()
	return c.models[m.Ident]
}

func (c *connection) Capabilities() capabilities.Set {
	//TODO implement me
	panic("implement me")
}

func (c *connection) Can(capabilities ...capabilities.Capability) bool {
	//TODO implement me
	panic("implement me")
}

func (c *connection) Close(ctx context.Context) error {
	//return c.db.Close() // <<= should we really?
	return nil
}

func (c *connection) CreateRecords(ctx context.Context, m *data.Model, rr ...crs.ValueGetter) error {
	return c.model(m).Create(ctx, rr...)
}

func (c *connection) LookupRecord(ctx context.Context, m *data.Model, pkv crs.ValueGetter, r crs.ValueSetter) error {
	return c.model(m).Lookup(ctx, pkv, r)
}

func (c *connection) SearchRecords(ctx context.Context, m *data.Model, filter any) (crs.Iterator, error) {
	return c.model(m).Search(filter.(types.RecordFilter))
}

func (c *connection) Models(ctx context.Context) (data.ModelSet, error) {
	//TODO implement me
	panic("implement me")
}

func (c *connection) AddModel(ctx context.Context, model *data.Model, model2 ...*data.Model) error {
	//TODO implement me
	panic("implement me")
}

func (c *connection) RemoveModel(ctx context.Context, model *data.Model, model2 ...*data.Model) error {
	//TODO implement me
	panic("implement me")
}

func (c *connection) AlterModel(ctx context.Context, old *data.Model, new *data.Model) error {
	//TODO implement me
	panic("implement me")
}

func (c *connection) AlterModelAttribute(ctx context.Context, sch *data.Model, old data.Attribute, new data.Attribute, trans ...crs.TransformationFunction) error {
	//TODO implement me
	panic("implement me")
}
