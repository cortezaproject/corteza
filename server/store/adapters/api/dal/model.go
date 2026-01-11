package dal

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"
	"github.com/spf13/cast"
)

type (
	runner interface {
		Run(ctx context.Context, method string, endpoint string, payload []byte, headers map[string][]string) (rsp []byte, err error)
	}

	model struct {
		model *dal.Model
		conn  runner

		dialect drivers.Dialect

		table drivers.TableCodec
	}
)

func validate(m *dal.Model) error {
	return nil
}

// Model returns fully initialized model store
//
// It abstracts database table and its columns and provides unified interface
// for fetching and storing records.
func Model(m *dal.Model, c runner, d drivers.Dialect) *model {
	var (
		ms = &model{
			model:   m,
			conn:    c,
			dialect: d,

			table: drivers.NewTableCodec(m, d),
		}
	)

	return ms
}

func (d *model) Create(ctx context.Context, rr ...dal.ValueGetter) error {
	for _, r := range rr {
		method, endpoint, payload, err := d.makeCreatePayload(r)
		if err != nil {
			return err
		}

		_, err = d.conn.Run(ctx, method, endpoint, payload, nil)
		if err != nil {
			return err
		}

	}

	return nil
}

func (d *model) makeCreatePayload(r dal.ValueGetter) (method string, endpoint string, payload []byte, err error) {
	method, endpointTpl, err := d.dialect.MapOpParams("create")
	if err != nil {
		return
	}

	payload, err = d.table.Encode(r)
	if err != nil {
		return
	}

	endpoint = d.procEndpointM(endpointTpl, d.table.Ident())
	return
}

func (d *model) Update(ctx context.Context, r dal.ValueGetter) error {
	method, endpoint, payload, err := d.makeUpdatePayload(r)
	if err != nil {
		return err
	}

	_, err = d.conn.Run(ctx, method, endpoint, payload, nil)
	if err != nil {
		return err
	}

	return nil
}

func (d *model) makeUpdatePayload(r dal.ValueGetter) (method string, endpoint string, payload []byte, err error) {
	method, endpointTpl, err := d.dialect.MapOpParams("update")
	if err != nil {
		return
	}

	var (
		key any
	)

	payload, err = d.table.Encode(r)
	if err != nil {
		return
	}

	for _, c := range d.table.Columns() {
		if c.IsPrimaryKey() {
			if key != nil {
				err = fmt.Errorf("composite keys not supported")
				return
			}

			key, err = c.Encode(r)
			if err != nil {
				return
			}
		}
	}

	endpoint = d.procEndpointMR(endpointTpl, d.table.Ident(), cast.ToString(key))
	return
}

func (d *model) Delete(ctx context.Context, r dal.ValueGetter) error {
	method, endpoint, err := d.makeDeletePayload(r)
	if err != nil {
		return err
	}

	_, err = d.conn.Run(ctx, method, endpoint, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (d *model) makeDeletePayload(r dal.ValueGetter) (method string, endpoint string, err error) {
	method, endpointTpl, err := d.dialect.MapOpParams("delete")
	if err != nil {
		return
	}

	var (
		key any
	)

	for _, c := range d.table.Columns() {
		if c.IsPrimaryKey() {
			if key != nil {
				err = fmt.Errorf("composite keys not supported")
				return
			}

			key, err = c.Encode(r)
			if err != nil {
				return
			}
		}
	}

	endpoint = d.procEndpointMR(endpointTpl, d.table.Ident(), cast.ToString(key))
	return
}

func (d *model) Lookup(ctx context.Context, pkv dal.ValueGetter, r dal.ValueSetter) (err error) {
	method, endpoint, err := d.makeLookupPayload(pkv)
	if err != nil {
		return
	}

	rsp, err := d.conn.Run(ctx, method, endpoint, nil, nil)
	if err != nil {
		return
	}

	err = d.table.Decode(rsp, r)
	if err != nil {
		return
	}

	return
}

func (d *model) makeLookupPayload(pkv dal.ValueGetter) (method string, endpoint string, err error) {
	method, endpointTpl, err := d.dialect.MapOpParams("read")
	if err != nil {
		return
	}

	cc := pkv.CountValues()
	if len(cc) > 1 {
		err = fmt.Errorf("can not use composite keys at this point")
		return
	}

	var val any
	for k := range cc {
		val, err = pkv.GetValue(k, 0)
		if err != nil {
			return
		}

		break
	}

	endpoint = d.procEndpointMR(endpointTpl, d.table.Ident(), cast.ToString(val))
	return
}

func (d *model) Search(f filter.Filter) (i *iterator, err error) {
	method, endpoint, query, err := d.makeSearchPayload(f)
	if err != nil {
		return
	}

	// @todo can we stream this?

	i = &iterator{
		src: d,
		dst: d,

		method:   method,
		endpoint: endpoint,
		query:    query,

		sorting: f.OrderBy(),
		cursor:  f.Cursor(),
		limit:   f.Limit(),
	}

	return
}

func (d *model) makeSearchPayload(f filter.Filter) (method string, endpoint string, query map[string]string, err error) {
	method, endpointTpl, err := d.dialect.MapOpParams("search")
	if err != nil {
		return
	}

	// @todo preparing and processing things

	endpoint = d.procEndpointM(endpointTpl, d.table.Ident())
	return
}

func (d *model) Truncate(ctx context.Context) error {
	return fmt.Errorf("operation not supported: truncate")
}

func (d *model) Count(ctx context.Context, f filter.Filter) (c uint, err error) {
	err = fmt.Errorf("operation not supported: count")
	return
}

// Aggregate constructs SELECT sql with group-by and an optional having CLAUSE
//
// All group-by attributes are prepended to aggregation
// expressions when constructing expressions & columns to select from.
//
// Passing in filter with cursor, empty groupBy or aggrExpr slice will result in an error
func (d *model) Aggregate(f filter.Filter, groupBy []dal.AggregateAttr, aggrExpr []dal.AggregateAttr, having *ql.ASTNode) (i *iterator, err error) {
	err = fmt.Errorf("operation not supported: aggregate")
	return
}

func (d *model) procEndpointM(tpl, module string) string {
	r := strings.NewReplacer(
		"{{moduleID}}",
		module,
	)

	return r.Replace(tpl)
}

func (d *model) procEndpointMR(tpl, module string, record string) string {
	r := strings.NewReplacer(
		"{{moduleID}}",
		module,

		"{{recordID}}",
		record,
	)

	return r.Replace(tpl)
}
