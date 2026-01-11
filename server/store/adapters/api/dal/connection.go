package dal

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	// connection provides (pkg/dal.Connection) interface to HTTP implementation
	//
	// In other words: this allows Corteza to read Records from the REST API
	connection struct {
		mux    sync.RWMutex
		models map[string]*model
		driver dal.Driver

		runner  runner
		dialect drivers.Dialect
	}
)

var (
	dalDriver dal.Driver
)

func init() {
	dalDriver = dal.Driver{
		Type:       "corteza::dal:driver:api",
		Operations: dal.FullOperations(),
		Connection: dal.NewRESTriverConnectionConfig(),
	}
	dal.RegisterDriver(dalDriver)
}

func Connection(runner runner, dialect drivers.Dialect) *connection {
	return &connection{
		runner:  runner,
		models:  make(map[string]*model),
		driver:  dalDriver,
		dialect: dialect,
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
	if cached, ok := c.models[key]; ok {
		return fn(cached)
	}

	return fmt.Errorf("model %q (%d) not loaded", key, m.ResourceID)
}

func (c *connection) Operations() dal.OperationSet {
	return c.driver.Operations
}

func (c *connection) Can(operations ...dal.Operation) bool {
	return c.Operations().IsSuperset(operations...)
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

func (c *connection) Count(ctx context.Context, m *dal.Model, f filter.Filter) (i uint, _ error) {
	return i, c.withModel(m, func(m *model) (err error) {
		i, err = m.Count(ctx, f)
		return
	})
}

func (c *connection) Analyze(ctx context.Context, m *dal.Model) (a map[string]dal.OpAnalysis, err error) {
	a = map[string]dal.OpAnalysis{}

	// // @todo somehow (probably operations) bring in the info what can be done
	// //       for now, since we're quite rigid on the drivers, this will do.
	// //
	// // @note this is a temporary hack until we properly address the first point.
	// //       No point in complicating it at this stage.
	// if c.db.DriverName() == "sqlserver" {
	// 	a = map[string]dal.OpAnalysis{}
	// } else {
	// 	a = map[string]dal.OpAnalysis{
	// 		dal.OpAnalysisAggregate: {
	// 			ScanCost:   dal.CostCheep,
	// 			SearchCost: dal.CostCheep,
	// 			FilterCost: dal.CostCheep,
	// 			SortCost:   dal.CostCheep,
	// 		},
	// 	}

	// }

	return
}

func (c *connection) Aggregate(ctx context.Context, m *dal.Model, f filter.Filter, groupBy []dal.AggregateAttr, aggrExpr []dal.AggregateAttr, having *ql.ASTNode) (i dal.Iterator, _ error) {
	return i, c.withModel(m, func(m *model) (err error) {
		i, err = m.Aggregate(f, groupBy, aggrExpr, having)
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
	for _, m := range mm {
		if err = validate(m); err != nil {
			return
		}
	}

	c.mux.Lock()
	defer c.mux.Unlock()
	for _, m := range mm {
		_ = m
		// @todo can we upsert the model?

		// cache the model
		c.models[cacheKey(m)] = Model(m, c.runner, c.dialect)
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
func (c *connection) UpdateModel(ctx context.Context, old *dal.Model, new *dal.Model) (err error) {
	if err = validate(new); err != nil {
		return
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	// remove from cache
	delete(c.models, cacheKey(old))

	// @todo check if column exists and if it can be removed

	// update the cache
	// @todo dialects??
	c.models[cacheKey(new)] = Model(new, c.runner, c.dialect)
	return
}

// AssertSchemaAlterations returns a new set of Alterations based on what the underlying
// schema already provides -- it discards alterations for column additions that already exist, etc.
func (c *connection) AssertSchemaAlterations(ctx context.Context, model *dal.Model, aa ...*dal.Alteration) (out []*dal.Alteration, err error) {
	// @todo can we do alterations even?
	return
}

// ApplyAhlteration applies the given alterations to the underlying schema
//
// The returned slice of error indicates what alterations failed.
// If the corresponding index is nil, the alteration was successful.
func (c *connection) ApplyAlteration(ctx context.Context, model *dal.Model, alt ...*dal.Alteration) (errs []error) {
	errs = []error{errors.New("operations not supported by the API driver")}
	return
}

func cacheKey(m *dal.Model) (key string) {
	key = m.ResourceType + "|" + m.Resource + "|" + m.Ident
	if key == "" {
		panic("can not add model without a key (combo of resource type, resource and ident)")
	}

	return
}
