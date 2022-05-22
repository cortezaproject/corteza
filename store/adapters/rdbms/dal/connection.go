package dal

import (
	"context"
	"sync"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/jmoiron/sqlx"
)

type (
	connection struct {
		mux          sync.RWMutex
		models       map[string]*model
		capabilities capabilities.Set

		db      *sqlx.DB
		dialect drivers.Dialect
	}
)

func Connection(db *sqlx.DB, dialect drivers.Dialect, cc ...capabilities.Capability) *connection {
	return &connection{
		db:           db,
		dialect:      dialect,
		models:       make(map[string]*model),
		capabilities: cc,
	}
}

func (c *connection) model(m *dal.Model) *model {
	if m.Resource == "" {
		// if resource is empty, use ident
		m.Resource = m.Ident
	}

	if m.Resource == "" {
		panic("can not add model with empty resource")
	}

	c.mux.RLock()
	if c.models[m.Resource] == nil {
		c.mux.RUnlock()
		c.mux.Lock()
		c.models[m.Resource] = Model(m, c.db, c.dialect)
		defer c.mux.Unlock()
		return c.models[m.Resource]
	}

	defer c.mux.RUnlock()
	return c.models[m.Resource]
}

func (c *connection) Capabilities() capabilities.Set {
	return c.capabilities
}

func (c *connection) Can(capabilities ...capabilities.Capability) bool {
	return c.capabilities.IsSuperset(capabilities...)
}

func (c *connection) Create(ctx context.Context, m *dal.Model, rr ...dal.ValueGetter) error {
	return c.model(m).Create(ctx, rr...)
}

func (c *connection) Update(ctx context.Context, m *dal.Model, r dal.ValueGetter) error {
	return c.model(m).Update(ctx, r)
}

func (c *connection) Lookup(ctx context.Context, m *dal.Model, pkv dal.ValueGetter, r dal.ValueSetter) error {
	return c.model(m).Lookup(ctx, pkv, r)
}

func (c *connection) Search(ctx context.Context, m *dal.Model, f filter.Filter) (dal.Iterator, error) {
	return c.model(m).Search(f)
}

func (c *connection) Delete(ctx context.Context, m *dal.Model, pkv dal.ValueGetter) error {
	return c.model(m).Delete(ctx, pkv)
}

func (c *connection) Truncate(ctx context.Context, m *dal.Model) error {
	return c.model(m).Truncate(ctx)
}

func (c *connection) Models(ctx context.Context) (dal.ModelSet, error) {
	//TODO implement me
	return nil, nil
	panic("implement me")
}

func (c *connection) CreateModel(ctx context.Context, model *dal.Model, model2 ...*dal.Model) error {
	//TODO implement me
	return nil
	panic("implement me")
}

func (c *connection) DeleteModel(ctx context.Context, model *dal.Model, model2 ...*dal.Model) error {
	//TODO implement me
	panic("implement me")
}

func (c *connection) UpdateModel(ctx context.Context, old *dal.Model, new *dal.Model) error {
	//TODO implement me
	panic("implement me")
}

func (c *connection) UpdateModelAttribute(ctx context.Context, sch *dal.Model, old dal.Attribute, new dal.Attribute, trans ...dal.TransformationFunction) error {
	//TODO implement me
	panic("implement me")
}
