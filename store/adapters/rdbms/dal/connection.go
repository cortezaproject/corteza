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
	// connection provides (pkg/dal.Connection) interface to RDBMS implementation
	//
	// In other words: this allows Corteza to read Records from the supported SQL databases
	connection struct {
		mux          sync.RWMutex
		models       map[string]*model
		capabilities capabilities.Set

		db      sqlx.ExtContext
		dialect drivers.Dialect
	}
)

func init() {
	dal.RegisterDriver(dal.Driver{
		Type:         "corteza::dal:driver:rdbms",
		Capabilities: capabilities.FullCapabilities(),
		Connection:   dal.NewDSNDriverConnectionConfig(),
	})
}

func Connection(db sqlx.ExtContext, dialect drivers.Dialect, cc ...capabilities.Capability) *connection {
	return &connection{
		db:           db,
		dialect:      dialect,
		models:       make(map[string]*model),
		capabilities: cc,
	}
}

// model returns rdbms/dal model (converted dal.Model)
//
// It constructs key from res-type + res + ident
// and caches it in the connection
//
// This allows us to have same resource or ident on different res-types
// For example: module's model for revisions has same resouce and ident but different type
func (c *connection) model(m *dal.Model) *model {
	key := m.ResourceType + "|" + m.Resource + "|" + m.Ident
	if key == "" {
		panic("can not add model without a key (combo of resource type, resource and ident)")
	}

	c.mux.RLock()
	if c.models[key] == nil {
		c.mux.RUnlock()
		c.mux.Lock()
		c.models[key] = Model(m, c.db, c.dialect)
		defer c.mux.Unlock()
		return c.models[key]
	}

	defer c.mux.RUnlock()
	return c.models[key]
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
	return nil
	panic("implement me")
}

func (c *connection) UpdateModel(ctx context.Context, old *dal.Model, new *dal.Model) error {
	//TODO implement me
	return nil
	panic("implement me")
}

func (c *connection) UpdateModelAttribute(ctx context.Context, sch *dal.Model, old, new *dal.Attribute, trans ...dal.TransformationFunction) error {
	//TODO implement me
	return nil
	panic("implement me")
}
