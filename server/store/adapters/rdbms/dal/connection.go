package dal

import (
	"context"
	"fmt"
	"sync"

	"github.com/cortezaproject/corteza/server/pkg/id"

	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers"
	"github.com/jmoiron/sqlx"
)

type (
	// connection provides (pkg/dal.Connection) interface to RDBMS implementation
	//
	// In other words: this allows Corteza to read Records from the supported SQL databases
	connection struct {
		mux    sync.RWMutex
		models map[string]*model
		driver dal.Driver

		db      sqlx.ExtContext
		dialect drivers.Dialect

		dataDefiner ddl.DataDefiner
	}
)

var (
	dalDriver dal.Driver
)

func init() {
	dalDriver = dal.Driver{
		Type:       "corteza::dal:driver:rdbms",
		Operations: dal.FullOperations(),
		Connection: dal.NewDSNDriverConnectionConfig(),
	}
	dal.RegisterDriver(dalDriver)
}

func Connection(db sqlx.ExtContext, dialect drivers.Dialect, dd ddl.DataDefiner) *connection {
	return &connection{
		db:          db,
		dialect:     dialect,
		dataDefiner: dd,
		models:      make(map[string]*model),
		driver:      dalDriver,
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
	// @todo somehow (probably operations) bring in the info what can be done
	//       for now, since we're quite rigid on the drivers, this will do.
	//
	// @note this is a temporary hack until we properly address the first point.
	//       No point in complicating it at this stage.
	if c.db.DriverName() == "sqlserver" {
		a = map[string]dal.OpAnalysis{}
	} else {
		a = map[string]dal.OpAnalysis{
			dal.OpAnalysisAggregate: {
				ScanCost:   dal.CostCheep,
				SearchCost: dal.CostCheep,
				FilterCost: dal.CostCheep,
				SortCost:   dal.CostCheep,
			},
		}

	}

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
		_, err = c.dataDefiner.TableLookup(ctx, m.Ident)
		if errors.IsNotFound(err) {
			if err = ddl.CreateModel(ctx, c.dataDefiner, m); err != nil {
				return
			}
		} else if err != nil {
			return
		}
		if err = ddl.EnsureIndexes(ctx, c.dataDefiner, m.Indexes...); err != nil {
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
	c.models[cacheKey(new)] = Model(new, c.db, c.dialect)
	return
}

// AssertSchemaAlterations returns a new set of Alterations based on what the underlying
// schema already provides -- it discards alterations for column additions that already exist, etc.
func (c *connection) AssertSchemaAlterations(ctx context.Context, model *dal.Model, aa ...*dal.Alteration) (out []*dal.Alteration, err error) {
	var aux []*dal.Alteration

	t, err := c.dataDefiner.TableLookup(ctx, model.Ident)
	if err != nil && errors.IsNotFound(err) {
		// Since there is no thing for this model we need to create it and all the
		// alterations are pointless
		return []*dal.Alteration{{
			ID:           id.Next(),
			Resource:     model.Resource,
			ResourceType: model.ResourceType,
			ConnectionID: model.ConnectionID,

			ModelAdd: &dal.ModelAdd{
				Model: model,
			},
		}}, nil
	}
	if err != nil {
		return
	}

	// Index columns by ident for easier lookup
	colIndex := make(map[string]*ddl.Column)
	for _, c := range t.Columns {
		colIndex[c.Ident] = c
	}

	for _, a := range aa {
		switch {
		case a.AttributeAdd != nil:
			aux, err = c.assertAlterationAttributeAdd(t, colIndex, a)
			if err != nil {
				return
			}
			out = append(out, aux...)
		case a.AttributeDelete != nil:
			aux, err = c.assertAlterationAttributeDelete(t, colIndex, a)
			if err != nil {
				return
			}
			out = append(out, aux...)
		case a.AttributeReType != nil:
			aux, err = c.assertAlterationAttributeReType(t, colIndex, a)
			if err != nil {
				return
			}
			out = append(out, aux...)
		case a.AttributeReEncode != nil:
			aux, err = c.assertAlterationAttributeReEncode(t, colIndex, a)
			if err != nil {
				return
			}
			out = append(out, aux...)
		case a.ModelAdd != nil:
			aux, err = c.assertAlterationModelAdd(t, colIndex, a)
			if err != nil {
				return
			}
			out = append(out, aux...)
		case a.ModelDelete != nil:
			aux, err = c.assertAlterationModelDelete(t, colIndex, a)
			if err != nil {
				return
			}
			out = append(out, aux...)
		}
	}

	return
}

// ApplyAhlteration applies the given alterations to the underlying schema
//
// The returned slice of error indicates what alterations failed.
// If the corresponding index is nil, the alteration was successful.
func (c *connection) ApplyAlteration(ctx context.Context, model *dal.Model, alt ...*dal.Alteration) (errs []error) {
	var (
		err    error
		failed = make(map[uint64]bool, len(alt)/2)
	)

	for _, a := range alt {
		// Skip since the alteration we depend on failed
		if a.DependsOn != 0 && failed[a.DependsOn] {
			errs = append(errs, fmt.Errorf("skipping alteration %d: depending alteration %d failed", a.ID, a.DependsOn))
			continue
		}

		switch {
		case a.AttributeAdd != nil:
			err = c.applyAlterationAttributeAdd(ctx, model, a)
		case a.AttributeDelete != nil:
			err = c.applyAlterationAttributeDelete(ctx, model, a)
		case a.AttributeReType != nil:
			err = c.applyAlterationAttributeReType(ctx, model, a)
		case a.AttributeReEncode != nil:
			err = c.applyAlterationAttributeReEncode(ctx, model, a)
		case a.ModelAdd != nil:
			err = c.applyAlterationModelAdd(ctx, model, a)
		case a.ModelDelete != nil:
			err = c.applyAlterationModelDelete(ctx, model, a)
		}

		if err != nil {
			failed[a.ID] = true
		}

		errs = append(errs, err)
	}

	return
}

func (c *connection) applyAlterationAttributeAdd(ctx context.Context, model *dal.Model, alt *dal.Alteration) (err error) {
	col, err := c.dataDefiner.ConvertAttribute(alt.AttributeAdd.Attr)
	if err != nil {
		return
	}

	return c.dataDefiner.ColumnAdd(ctx, model.Ident, col)
}

func (c *connection) applyAlterationAttributeDelete(ctx context.Context, model *dal.Model, alt *dal.Alteration) (err error) {
	return c.dataDefiner.ColumnDrop(ctx, model.Ident, alt.AttributeDelete.Attr.StoreIdent())
}

func (c *connection) applyAlterationAttributeReType(ctx context.Context, model *dal.Model, alt *dal.Alteration) (err error) {
	col, err := c.dataDefiner.ConvertAttribute(alt.AttributeReType.Attr)
	if err != nil {
		return
	}

	return c.dataDefiner.ColumnReType(ctx, model.Ident, col.Ident, col.Type)
}

// @todo might consider droppig this one for now
func (c *connection) applyAlterationAttributeReEncode(ctx context.Context, model *dal.Model, alt *dal.Alteration) (err error) {
	// ...
	return
}

func (c *connection) applyAlterationModelAdd(ctx context.Context, model *dal.Model, alt *dal.Alteration) (err error) {
	t, err := c.dataDefiner.ConvertModel(model)
	if err != nil {
		return
	}

	return c.dataDefiner.TableCreate(ctx, t)
}

func (c *connection) applyAlterationModelDelete(ctx context.Context, model *dal.Model, alt *dal.Alteration) (err error) {
	return c.dataDefiner.TableDrop(ctx, model.Ident)
}

func cacheKey(m *dal.Model) (key string) {
	key = m.ResourceType + "|" + m.Resource + "|" + m.Ident
	if key == "" {
		panic("can not add model without a key (combo of resource type, resource and ident)")
	}

	return
}

func (c *connection) assertAlterationAttributeAdd(table *ddl.Table, colIndex map[string]*ddl.Column, alt *dal.Alteration) (out []*dal.Alteration, err error) {
	if alt.AttributeAdd.Attr.Store.Type() == dal.AttributeCodecRecordValueSetJSON {
		// RecordValue codec needs to be checked a bit differently since we're worried about the column that contains
		// the JSON, not the attribute ident itself
		return c.assertAlterationNestedAttributeAdd(table, colIndex, alt)
	} else {
		return c.assertAlterationStandaloneAttributeAdd(table, colIndex, alt)
	}
}

func (c *connection) assertAlterationNestedAttributeAdd(table *ddl.Table, colIndex map[string]*ddl.Column, alt *dal.Alteration) (out []*dal.Alteration, err error) {
	col := colIndex[alt.AttributeAdd.Attr.StoreIdent()]
	if col != nil {
		return
	}

	a := &dal.Alteration{
		ID:           id.Next(),
		BatchID:      alt.BatchID,
		DependsOn:    alt.DependsOn,
		Resource:     alt.Resource,
		ConnectionID: alt.ConnectionID,

		AttributeAdd: &dal.AttributeAdd{
			Attr: &dal.Attribute{
				Ident: alt.AttributeAdd.Attr.StoreIdent(),
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeJSON{Nullable: false},
			},
		},
	}
	out = append(out, a)

	// Update colIndex so other alterations won't duplicate
	auxCol, err := c.dialect.AttributeToColumn(a.AttributeAdd.Attr)
	if err != nil {
		return nil, err
	}
	colIndex[auxCol.Ident] = auxCol

	return out, nil
}

func (c *connection) assertAlterationStandaloneAttributeAdd(table *ddl.Table, colIndex map[string]*ddl.Column, alt *dal.Alteration) (out []*dal.Alteration, err error) {
	col := colIndex[alt.AttributeAdd.Attr.StoreIdent()]
	if col == nil {
		out = append(out, alt)
		return
	}

	auxCol, err := c.dialect.AttributeToColumn(alt.AttributeAdd.Attr)
	if err != nil {
		return nil, err
	}

	if !c.dialect.ColumnFits(col, auxCol) {
		out = append(out, &dal.Alteration{
			ID:           id.Next(),
			BatchID:      alt.BatchID,
			DependsOn:    alt.DependsOn,
			Resource:     alt.Resource,
			ResourceType: alt.ResourceType,
			ConnectionID: alt.ConnectionID,

			AttributeReType: &dal.AttributeReType{
				Attr: alt.AttributeAdd.Attr,
				To:   alt.AttributeAdd.Attr.Type,
			},
		})
	}

	// Update colIndex so other alterations won't duplicate
	colIndex[auxCol.Ident] = auxCol

	return out, nil
}

func (c *connection) assertAlterationAttributeDelete(table *ddl.Table, colIndex map[string]*ddl.Column, alt *dal.Alteration) (out []*dal.Alteration, err error) {
	col := colIndex[alt.AttributeDelete.Attr.StoreIdent()]
	if col == nil {
		return
	}

	delete(colIndex, alt.AttributeDelete.Attr.StoreIdent())

	out = append(out, alt)
	return
}

func (c *connection) assertAlterationAttributeReType(table *ddl.Table, colIndex map[string]*ddl.Column, alt *dal.Alteration) (out []*dal.Alteration, err error) {
	col := colIndex[alt.AttributeReType.Attr.StoreIdent()]
	if col == nil {
		err = fmt.Errorf("cannot alter %s, column does not exist", alt.AttributeReType.Attr.StoreIdent())
		return
	}

	// Since it's a JSON we don't need to do anything
	// @todo consider adding some migration logic here
	if alt.AttributeReType.Attr.Store.Type() == dal.AttributeCodecRecordValueSetJSON {
		return
	}

	auxCol, err := c.dialect.AttributeToColumn(alt.AttributeReType.Attr)
	if err != nil {
		return
	}

	if c.dialect.ColumnFits(auxCol, col) {
		return
	}

	out = append(out, alt)
	return
}

// @todo for now it just creates a new column; we should add some migration in the future
func (c *connection) assertAlterationAttributeReEncode(table *ddl.Table, colIndex map[string]*ddl.Column, alt *dal.Alteration) (out []*dal.Alteration, err error) {
	auxCodec := *alt.AttributeReEncode.Attr
	auxCodec.Store = alt.AttributeReEncode.To

	out = append(out, &dal.Alteration{
		ID:           alt.ID,
		BatchID:      alt.BatchID,
		DependsOn:    alt.DependsOn,
		Resource:     alt.Resource,
		ResourceType: alt.ResourceType,
		ConnectionID: alt.ConnectionID,

		AttributeAdd: &dal.AttributeAdd{
			Attr: &auxCodec,
		},
	})
	return
}

func (c *connection) assertAlterationModelAdd(table *ddl.Table, colIndex map[string]*ddl.Column, alt *dal.Alteration) (out []*dal.Alteration, err error) {
	if table != nil {
		return
	}

	out = append(out, alt)
	return
}

func (c *connection) assertAlterationModelDelete(table *ddl.Table, colIndex map[string]*ddl.Column, alt *dal.Alteration) (out []*dal.Alteration, err error) {
	if table == nil {
		return
	}

	out = append(out, alt)
	return
}
