package dal

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"sync"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/jmoiron/sqlx"
)

type (
	// connection provides (pkg/dal.Connection) interface to RDBMS implementation
	//
	// In other words: this allows Corteza to read Records from the supported SQL databases
	connection struct {
		mux        sync.RWMutex
		models     map[string]*model
		operations dal.OperationSet

		db      sqlx.ExtContext
		dialect drivers.Dialect

		dataDefiner ddl.DataDefiner
	}
)

func init() {
	dal.RegisterDriver(dal.Driver{
		Type:       "corteza::dal:driver:rdbms",
		Operations: dal.FullOperations(),
		Connection: dal.NewDSNDriverConnectionConfig(),
	})
}

func Connection(db sqlx.ExtContext, dialect drivers.Dialect, dd ddl.DataDefiner, cc ...dal.Operation) *connection {
	return &connection{
		db:          db,
		dialect:     dialect,
		dataDefiner: dd,
		models:      make(map[string]*model),
		operations:  cc,
	}
}

// model returns rdbms/dal model (converted dal.Model)
//
// It constructs key from res-type + res + ident
// and caches it in the connection
//
// This allows us to have same resource or ident on different res-types
// For example: module's model for revisions has same resouce and ident but different type
func (c *connection) withModel(m *dal.Model, fn func(m *model) error) error {
	var (
		key = cacheKey(m)
	)
	c.mux.RLock()
	defer c.mux.RUnlock()
	if cached, ok := c.models[cacheKey(m)]; ok {
		return fn(cached)
	}

	return fmt.Errorf("model %q (%d) not loaded", key, m.ResourceID)
}

func (c *connection) Operations() dal.OperationSet {
	return c.operations
}

func (c *connection) Can(operations ...dal.Operation) bool {
	return c.operations.IsSuperset(operations...)
}

func (c *connection) Create(ctx context.Context, m *dal.Model, rr ...dal.ValueGetter) (err error) {
	return c.withModel(m, func(m *model) error {
		return m.Create(ctx, rr...)
	})
}

func (c *connection) Update(ctx context.Context, m *dal.Model, r dal.ValueGetter) (err error) {
	return c.withModel(m, func(m *model) error {
		return m.Update(ctx, r)
	})
}

func (c *connection) Lookup(ctx context.Context, m *dal.Model, pkv dal.ValueGetter, r dal.ValueSetter) (err error) {
	return c.withModel(m, func(m *model) error {
		return m.Lookup(ctx, pkv, r)
	})
}

func (c *connection) Search(ctx context.Context, m *dal.Model, f filter.Filter) (i dal.Iterator, _ error) {
	return i, c.withModel(m, func(m *model) (err error) {
		i, err = m.Search(f)
		return
	})
}

func (c *connection) Delete(ctx context.Context, m *dal.Model, pkv dal.ValueGetter) (err error) {
	return c.withModel(m, func(m *model) error {
		return m.Delete(ctx, pkv)
	})
}

func (c *connection) Truncate(ctx context.Context, m *dal.Model) (err error) {
	return c.withModel(m, func(m *model) error {
		return m.Truncate(ctx)
	})
}

func (c *connection) Models(ctx context.Context) (dal.ModelSet, error) {
	// not raising not-supported error
	// because we do not want to break
	// DAL service model adding procedure
	return nil, nil
}

// CreateModel checks/creates db tables in the database and catches the processed model
//
// @todo DDL operations
func (c *connection) CreateModel(ctx context.Context, mm ...*dal.Model) (err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, m := range mm {
		err = ddl.UpdateModel(ctx, c.dataDefiner, m)
		if !errors.IsNotFound(err) && err != nil {
			return
		}

		if err = ddl.CreateModel(ctx, c.dataDefiner, m); err != nil {
			return
		}

		// cache the model
		c.models[cacheKey(m)] = Model(m, c.db, c.dialect)
	}

	return
}

// DeleteModel removes db tables from the database and removes the processed model from cache
//
// @todo DDL operations
// @todo some tables should not be removed (like compose_record on primary connection)
func (c *connection) DeleteModel(ctx context.Context, mm ...*dal.Model) (err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, m := range mm {
		// @todo check if table exists and if it can be removed
		if err = c.DeleteModel(ctx, m); err != nil {
			return
		}

		// remove from cache
		delete(c.models, cacheKey(m))
	}

	return
}

// UpdateModel alters db tables from the database and refreshes the processed model in the cache
//
// @todo DDL operations
// @todo some tables should not be removed (like compose_record on primary connection)
func (c *connection) UpdateModel(ctx context.Context, old *dal.Model, new *dal.Model) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	// remove from cache
	delete(c.models, cacheKey(old))

	// @todo check if column exists and if it can be removed

	// update the cache
	c.models[cacheKey(new)] = Model(new, c.db, c.dialect)
	return nil
}

// UpdateModelAttribute alters column on a db table and runs data transformations
func (c *connection) UpdateModelAttribute(ctx context.Context, sch *dal.Model, old, new *dal.Attribute, trans ...dal.TransformationFunction) error {
	// not raising not-supported error
	// because we do not want to break
	// DAL service model adding procedure

	// @todo implement model column altering
	return nil
}

func cacheKey(m *dal.Model) (key string) {
	key = m.ResourceType + "|" + m.Resource + "|" + m.Ident
	if key == "" {
		panic("can not add model without a key (combo of resource type, resource and ident)")
	}

	return
}
