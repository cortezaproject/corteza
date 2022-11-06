package dal_test

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	. "github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestModel_Search(t *testing.T) {
	_ = logger.Default()

	const (
		items = 1000
	)

	var (
		req = require.New(t)

		ctx = context.Background()

		baseModel = &dal.Model{
			Ident: t.Name(),
			Attributes: []*dal.Attribute{
				{Ident: "item", Type: &dal.TypeText{}},
				{Ident: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true},
				{Ident: "price", Type: &dal.TypeNumber{}, Filterable: true},
				{Ident: "published", Type: &dal.TypeBoolean{}, Filterable: true},
			},
		}

		m = Model(baseModel, s.DB, s.Dialect)

		i dal.Iterator

		table, err = s.DataDefiner.ConvertModel(baseModel)
		row        kv
	)

	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	t.Logf("Creating temporary table %q", table.Ident)
	table.Temporary = true
	req.NoError(s.DataDefiner.TableCreate(ctx, table))

	bm := time.Now()
	ctx = context.Background() // no need to log inserts
	count := 0
	for i := 1; i <= items; i++ {
		req.NoError(m.Create(ctx, &kv{
			"item":      fmt.Sprintf("i%d", i),
			"group":     fmt.Sprintf("g%d", i%1000),
			"price":     i,
			"published": i%2 == 0,
		}))
		count++
	}

	t.Logf("Inserted %d entries in %v", count, time.Now().Sub(bm))

	t.Log("Search through records, filter out published")
	i, err = m.Search(filter.Generic(
		filter.WithExpression("published"),
		filter.WithOrderBy(filter.SortExprSet{
			&filter.SortExpr{Column: "group"},
		}),
	))

	req.NoError(err)
	req.NotNil(i)

	defer req.NoError(i.Close())

	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	t.Log("Iterating over results")
	rows := make([]kv, 0, 5)
	for i.Next(ctx) {
		row = kv{}
		req.NoError(i.Scan(row))
		rows = append(rows, row)
	}

	req.NoError(i.Err())
	req.Len(rows, items/2)
	req.Equal("group=g0 item=i1000 price=1000 published=1", rows[0].String())
}

// Should be part of general DAL testing (not only RDBMS)
func TestModel_Search2(t *testing.T) {
	var (
		req = require.New(t)

		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

		baseModel = &dal.Model{
			Ident: t.Name(),
			Attributes: func() (aa []*dal.Attribute) {
				s := &dal.CodecRecordValueSetJSON{Ident: "values"}
				aa = []*dal.Attribute{
					{Ident: "phyTxt", Type: &dal.TypeText{}},
					{Ident: "phyNum", Type: &dal.TypeNumber{}},
					{Ident: "phyBool", Type: &dal.TypeBoolean{}},
					{Ident: "jsonEncTxt", Type: &dal.TypeText{}, Store: s},
					{Ident: "jsonEncTxtMV", Type: &dal.TypeText{}, Store: s, MultiValue: true},
					{Ident: "jsonEncNum", Type: &dal.TypeNumber{}, Store: s},
					{Ident: "jsonEncNumMV", Type: &dal.TypeNumber{}, Store: s, MultiValue: true},
					{Ident: "jsonEncBool", Type: &dal.TypeBoolean{}, Store: s},
					{Ident: "jsonEncBoolMV", Type: &dal.TypeBoolean{}, Store: s, MultiValue: true},
				}

				// iterate through all attributsand set Filterable flag to true
				for _, a := range aa {
					a.Filterable = true
					a.Sortable = true
				}

				return
			}(),
		}

		m = Model(baseModel, s.DB, s.Dialect)

		search = func(t *testing.T, f filter.Filter) []kv {
			req := require.New(t)
			i, err := m.Search(f)

			req.NoError(err)
			req.NotNil(i)

			defer req.NoError(i.Close())

			ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

			rows := make([]kv, 0, 5)
			for i.Next(ctx) {
				row := kv{}
				req.NoError(i.Scan(row))
				rows = append(rows, row)
			}

			req.NoError(i.Err())
			return rows
		}
	)

	table, err := s.DataDefiner.ConvertModel(baseModel)
	req.NoError(err)

	t.Logf("Creating temporary table %q", table.Ident)
	table.Temporary = true
	req.NoError(s.DataDefiner.TableCreate(ctx, table))

	{
		noLogCtx := context.Background() // no need to log inserts
		req.NoError(m.Create(noLogCtx, (&kvv{}).
			Set("phyTxt", "bar").
			Set("phyNum", 42).
			Set("phyBool", false).
			Set("jsonEncTxt", "bar").
			Set("jsonEncTxtMV", "bar", "foo").
			Set("jsonEncNum", 42).
			Set("jsonEncNumMV", 21, 42).
			Set("jsonEncBool", false).
			Set("jsonEncBoolMV", false, true),
		))
	}

	req.Len(search(t, filter.Generic(filter.WithExpression("phyTxt = 'bar'"))), 1)
	req.Len(search(t, filter.Generic(filter.WithExpression("phyTxt = 'baz'"))), 0)
	req.Len(search(t, filter.Generic(filter.WithExpression("phyNum = 42"))), 1)
	req.Len(search(t, filter.Generic(filter.WithExpression("phyNum = 21"))), 0)
	req.Len(search(t, filter.Generic(filter.WithExpression("!phyBool"))), 1)
	req.Len(search(t, filter.Generic(filter.WithExpression("phyBool"))), 0)
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncTxt = 'bar'"))), 1)
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncTxt = 'baz'"))), 0)
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncTxtMV = 'bar'"))), 1)
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncTxtMV = 'baz'"))), 0)
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncTxtMV = 'foo'"))), 0, "should not match the second value")
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncNum = 42"))), 1)
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncNum = 21"))), 0)
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncNumMV = 21"))), 1)
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncNumMV = 22"))), 0)
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncNumMV = 42"))), 0, "should not match the second value")
	req.Len(search(t, filter.Generic(filter.WithExpression("jsonEncBoolMV"))), 0, "should not match the second value")
}

func TestModel_Aggregate(t *testing.T) {
	_ = logger.Default()

	var (
		req = require.New(t)

		ctx = context.Background()

		baseModel = &dal.Model{
			Ident: t.Name(),
			Attributes: []*dal.Attribute{
				{Ident: "item", Type: &dal.TypeText{}},
				{Ident: "date", Type: &dal.TypeDate{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true},
				{Ident: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true},
				{Ident: "quantity", Type: &dal.TypeNumber{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Filterable: true},
				{Ident: "price", Type: &dal.TypeNumber{}, Filterable: true},
				{Ident: "published", Type: &dal.TypeBoolean{}, Filterable: true},
			},
		}

		m = Model(baseModel, s.DB, s.Dialect)

		i dal.Iterator

		table, err = s.DataDefiner.ConvertModel(baseModel)
		row        kv
	)

	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	t.Logf("Creating temporary table %q", table.Ident)
	table.Temporary = true
	req.NoError(s.DataDefiner.TableCreate(ctx, table))

	req.NoError(m.Create(ctx, &kv{"item": "i1", "date": "2022-10-06", "group": "g1", "price": "1000", "quantity": "10", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i2", "date": "2022-10-06", "group": "g1", "price": "3000", "quantity": "30", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i3", "date": "2022-10-06", "group": "g2", "price": "4000", "quantity": "40", "published": false}))
	req.NoError(m.Create(ctx, &kv{"item": "i4", "date": "2022-10-06", "group": "g2", "price": "1000", "quantity": "10", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i5", "date": "2022-10-07", "group": "g2", "price": "1000", "quantity": "10", "published": false}))
	req.NoError(m.Create(ctx, &kv{"item": "i6", "date": "2022-10-07", "group": "g2", "price": "5000", "quantity": "50", "published": true}))

	t.Log("Aggregating all records, calculating min & max price per group")
	i, err = m.Aggregate(
		filter.Generic(
			filter.WithExpression("published"),
			filter.WithOrderBy(filter.SortExprSet{
				&filter.SortExpr{Column: "group", Descending: true},
				&filter.SortExpr{Column: "date", Descending: false},
			}),
		),
		// group-by
		[]dal.AggregateAttr{
			{Identifier: "date", Type: &dal.TypeDate{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
			{Identifier: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
		},
		// aggregation expressions
		[]dal.AggregateAttr{
			{Identifier: "count", RawExpr: "COUNT(*)", Type: &dal.TypeNumber{}},
			{Identifier: "max", RawExpr: "MAX(price)", Type: &dal.TypeNumber{}},
			{Identifier: "min", RawExpr: "MIN(price)", Type: &dal.TypeNumber{}},
			{Identifier: "avg", RawExpr: "AVG(price)", Type: &dal.TypeNumber{}},
			{Identifier: "stock", RawExpr: "SUM(quantity)", Type: &dal.TypeNumber{}},
		},
		nil, // <== here be having condition
	)
	req.NoError(err)
	req.NotNil(i)

	defer req.NoError(i.Close())

	t.Log("Iterating over results")
	rows := make([]kv, 0, 3)
	for i.Next(ctx) {
		row = kv{}
		req.NoError(i.Scan(row))

		// due to difference of number of decimal digits in different DBs, we need to do this
		// to make sure we get the same result
		row["avg"] = fmt.Sprintf("%.2f", cast.ToFloat64(row["avg"]))
		row["stock"] = fmt.Sprintf("%.2f", cast.ToFloat64(row["stock"]))

		rows = append(rows, row)
	}

	req.NoError(i.Err())
	req.Len(rows, 3)
	req.Equal("avg=1000.00 count=1 date=2022-10-06 group=g2 max=1000 min=1000 stock=10.00", rows[0].String())
	req.Equal("avg=5000.00 count=1 date=2022-10-07 group=g2 max=5000 min=5000 stock=50.00", rows[1].String())
	req.Equal("avg=2000.00 count=2 date=2022-10-06 group=g1 max=3000 min=1000 stock=40.00", rows[2].String())
}

func TestModel_AggregationModelIdentifiers(t *testing.T) {
	_ = logger.Default()

	var (
		req = require.New(t)

		ctx = context.Background()

		baseModel = &dal.Model{
			Ident: "t",
			Attributes: []*dal.Attribute{
				{Ident: "group", Type: &dal.TypeText{}, Filterable: true},
			},
		}

		m = Model(baseModel, s.DB, s.Dialect)

		i dal.Iterator

		table, err = s.DataDefiner.ConvertModel(baseModel)
		row        kv
	)

	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	t.Logf("Creating temporary table %q", table.Ident)
	table.Temporary = true
	req.NoError(s.DataDefiner.TableCreate(ctx, table))

	req.NoError(m.Create(ctx, &kv{"group": "g1"}))
	req.NoError(m.Create(ctx, &kv{"group": "g1"}))
	req.NoError(m.Create(ctx, &kv{"group": "g2"}))
	req.NoError(m.Create(ctx, &kv{"group": "g2"}))
	req.NoError(m.Create(ctx, &kv{"group": "g2"}))

	t.Log("Aggregating all records, calculating min & max price per group")
	i, err = m.Aggregate(
		filter.Generic(
			filter.WithOrderBy(filter.SortExprSet{
				&filter.SortExpr{Column: "grp", Descending: true},
			}),
		),
		// group-by
		[]dal.AggregateAttr{
			{Identifier: "grp", RawExpr: "CONCAT(group, 'x')", Type: &dal.TypeText{}},
		},
		// aggregation expressions
		[]dal.AggregateAttr{
			{Identifier: "total", RawExpr: "COUNT(*)", Type: &dal.TypeNumber{}},
		},
		// having
		qlParse(req, "total = 3"),
	)
	req.NoError(err)
	req.NotNil(i)

	defer req.NoError(i.Close())

	t.Log("Iterating over results")
	rows := make([]kv, 0)
	for i.Next(ctx) {
		row = kv{}
		req.NoError(i.Scan(row))

		rows = append(rows, row)
	}

	req.NoError(i.Err())
	req.Len(rows, 1)
	req.Equal("grp=g2x total=3", rows[0].String())
}

func TestModel_AggregateWithHaving(t *testing.T) {
	_ = logger.Default()

	var (
		req = require.New(t)

		ctx = context.Background()

		baseModel = &dal.Model{
			Ident: t.Name(),
			Attributes: []*dal.Attribute{
				{Ident: "item", Type: &dal.TypeText{}},
				{Ident: "date", Type: &dal.TypeDate{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true},
				{Ident: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true},
				{Ident: "quantity", Type: &dal.TypeNumber{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Filterable: true},
				{Ident: "price", Type: &dal.TypeNumber{}, Filterable: true},
				{Ident: "published", Type: &dal.TypeBoolean{}, Filterable: true},
			},
		}

		m = Model(baseModel, s.DB, s.Dialect)

		i dal.Iterator

		table, err = s.DataDefiner.ConvertModel(baseModel)
		row        kv
	)

	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	t.Logf("Creating temporary table %q", table.Ident)
	table.Temporary = true
	req.NoError(s.DataDefiner.TableCreate(ctx, table))

	req.NoError(m.Create(ctx, &kv{"item": "i1", "date": "2022-10-06", "group": "g1", "price": "1000", "quantity": "0", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i2", "date": "2022-10-06", "group": "g1", "price": "3000", "quantity": "0", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i3", "date": "2022-10-06", "group": "g2", "price": "4000", "quantity": "40", "published": false}))
	req.NoError(m.Create(ctx, &kv{"item": "i4", "date": "2022-10-06", "group": "g2", "price": "1000", "quantity": "10", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i5", "date": "2022-10-07", "group": "g2", "price": "1000", "quantity": "10", "published": false}))
	req.NoError(m.Create(ctx, &kv{"item": "i6", "date": "2022-10-07", "group": "g2", "price": "5000", "quantity": "50", "published": true}))

	t.Log("Aggregating all records, calculating min & max price per group, ignoring empty quantity")
	i, err = m.Aggregate(
		filter.Generic(
			filter.WithExpression("published"),
			filter.WithOrderBy(filter.SortExprSet{
				&filter.SortExpr{Column: "group", Descending: true},
				&filter.SortExpr{Column: "date", Descending: false},
			}),
		),
		// group-by
		[]dal.AggregateAttr{
			{Identifier: "date", Type: &dal.TypeDate{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
			{Identifier: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
		},
		// aggregation expressions
		[]dal.AggregateAttr{
			{Identifier: "count", RawExpr: "COUNT(*)", Type: &dal.TypeNumber{}},
			{Identifier: "max", RawExpr: "MAX(price)", Type: &dal.TypeNumber{}},
			{Identifier: "min", RawExpr: "MIN(price)", Type: &dal.TypeNumber{}},
			{Identifier: "avg", RawExpr: "AVG(price)", Type: &dal.TypeNumber{}},
			{Identifier: "stock", RawExpr: "SUM(quantity)", Type: &dal.TypeNumber{}},
		},
		// having
		qlParse(req, "SUM(quantity) > 0"),
	)
	req.NoError(err)
	req.NotNil(i)

	defer req.NoError(i.Close())

	t.Log("Iterating over results")
	rows := make([]kv, 0, 3)
	for i.Next(ctx) {
		row = kv{}
		req.NoError(i.Scan(row))

		// due to difference of number of decimal digits in different DBs, we need to do this
		// to make sure we get the same result
		row["avg"] = fmt.Sprintf("%.2f", cast.ToFloat64(row["avg"]))
		row["stock"] = fmt.Sprintf("%.2f", cast.ToFloat64(row["stock"]))

		rows = append(rows, row)
	}

	req.NoError(i.Err())
	req.Len(rows, 2)
	req.Equal("avg=1000.00 count=1 date=2022-10-06 group=g2 max=1000 min=1000 stock=10.00", rows[0].String())
	req.Equal("avg=5000.00 count=1 date=2022-10-07 group=g2 max=5000 min=5000 stock=50.00", rows[1].String())
}

func TestModel_AggregateHavingGroup(t *testing.T) {
	_ = logger.Default()
	var (
		req = require.New(t)

		ctx = context.Background()

		baseModel = &dal.Model{
			Ident: t.Name(),
			Attributes: []*dal.Attribute{
				{Ident: "item", Type: &dal.TypeText{}},
				{Ident: "date", Type: &dal.TypeDate{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true},
				{Ident: "grp", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true, Filterable: true},
				{Ident: "quantity", Type: &dal.TypeNumber{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Filterable: true},
				{Ident: "price", Type: &dal.TypeNumber{}, Filterable: true},
				{Ident: "published", Type: &dal.TypeBoolean{}, Filterable: true},
			},
		}

		m = Model(baseModel, s.DB, s.Dialect)

		i dal.Iterator

		table, err = s.DataDefiner.ConvertModel(baseModel)
		row        kv
	)

	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	t.Logf("Creating temporary table %q", table.Ident)
	table.Temporary = true
	req.NoError(s.DataDefiner.TableCreate(ctx, table))

	req.NoError(m.Create(ctx, &kv{"item": "i1", "date": "2022-10-06", "grp": "g1", "price": "1000", "quantity": "10", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i2", "date": "2022-10-06", "grp": "g1", "price": "3000", "quantity": "30", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i3", "date": "2022-10-06", "grp": "g2", "price": "4000", "quantity": "40", "published": false}))
	req.NoError(m.Create(ctx, &kv{"item": "i4", "date": "2022-10-06", "grp": "g2", "price": "1000", "quantity": "10", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i5", "date": "2022-10-07", "grp": "g2", "price": "1000", "quantity": "10", "published": false}))
	req.NoError(m.Create(ctx, &kv{"item": "i6", "date": "2022-10-07", "grp": "g2", "price": "5000", "quantity": "50", "published": true}))

	t.Log("Aggregating all records, calculating min & max price per group")
	i, err = m.Aggregate(
		filter.Generic(),
		// group-by
		[]dal.AggregateAttr{
			{Identifier: "agg", Expression: &ql.ASTNode{Symbol: "grp"}, Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
		},
		// aggregation expressions
		[]dal.AggregateAttr{
			{Identifier: "count", RawExpr: "COUNT(*)", Type: &dal.TypeNumber{}},
		},
		qlParse(req, "agg = 'g1'"),
	)
	req.NoError(err)
	req.NotNil(i)

	defer req.NoError(i.Close())

	t.Log("Iterating over results")
	rows := make([]kv, 0, 1)
	for i.Next(ctx) {
		row = kv{}
		req.NoError(i.Scan(row))

		rows = append(rows, row)
	}

	req.NoError(i.Err())
	req.Len(rows, 1)
	req.Equal("agg=g1 count=2", rows[0].String())
}

func TestModel_Distinct(t *testing.T) {
	_ = logger.Default()

	var (
		req = require.New(t)

		ctx = context.Background()

		baseModel = &dal.Model{
			Ident: t.Name(),
			Attributes: []*dal.Attribute{
				{Ident: "item", Type: &dal.TypeText{}},
				{Ident: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true},
				{Ident: "published", Type: &dal.TypeBoolean{}, Filterable: true},
			},
		}

		m = Model(baseModel, s.DB, s.Dialect)

		i dal.Iterator

		table, err = s.DataDefiner.ConvertModel(baseModel)
		row        kv
	)

	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	t.Logf("Creating temporary table %q", table.Ident)
	table.Temporary = true
	req.NoError(s.DataDefiner.TableCreate(ctx, table))

	req.NoError(m.Create(ctx, &kv{"item": "i1", "group": "g1", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i2", "group": "g1", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i3", "group": "g2", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i4", "group": "g2", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i5", "group": "g3", "published": false}))
	req.NoError(m.Create(ctx, &kv{"item": "i6", "group": "g3", "published": false}))

	t.Log("Aggregating all records, returning distinct groups")
	i, err = m.Aggregate(
		filter.Generic(
			filter.WithExpression("published"),
			filter.WithOrderBy(filter.SortExprSet{
				&filter.SortExpr{Column: "group"},
			})),
		// group-by
		[]dal.AggregateAttr{
			{Identifier: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
		},
		nil,
		nil, // <== here be having condition
	)
	req.NoError(err)
	req.NotNil(i)

	defer req.NoError(i.Close())

	t.Log("Iterating over results")
	rows := make([]kv, 0, 2)
	for i.Next(ctx) {
		row = kv{}
		req.NoError(i.Scan(row))

		rows = append(rows, row)
	}

	req.NoError(i.Err())
	req.Len(rows, 2)
	req.Equal("group=g1", rows[0].String())
	req.Equal("group=g2", rows[1].String())
}

func TestModel_AggregateWithCursors(t *testing.T) {
	_ = logger.Default()

	var (
		req = require.New(t)

		ctx = context.Background()

		baseModel = &dal.Model{
			Ident: t.Name(),
			Attributes: []*dal.Attribute{
				{Ident: "item", Type: &dal.TypeText{}},
				{Ident: "date", Type: &dal.TypeDate{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true},
				{Ident: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Sortable: true},
				{Ident: "quantity", Type: &dal.TypeNumber{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}, Filterable: true},
				{Ident: "price", Type: &dal.TypeNumber{}, Filterable: true},
				{Ident: "published", Type: &dal.TypeBoolean{}, Filterable: true},
			},
		}

		m = Model(baseModel, s.DB, s.Dialect)

		i dal.Iterator

		table, err = s.DataDefiner.ConvertModel(baseModel)
		row        kv

		cur *filter.PagingCursor

		results = []string{
			"avg=1000.00 count=1 date=2022-10-06 group=g2 max=1000 min=1000 stock=10.00",
			"avg=5000.00 count=1 date=2022-10-07 group=g2 max=5000 min=5000 stock=50.00",
			"",
		}
	)

	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	t.Logf("Creating temporary table %q", table.Ident)
	table.Temporary = true
	req.NoError(s.DataDefiner.TableCreate(ctx, table))

	req.NoError(m.Create(ctx, &kv{"item": "i1", "date": "2022-10-06", "group": "g1", "price": "1000", "quantity": "0", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i2", "date": "2022-10-06", "group": "g1", "price": "3000", "quantity": "0", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i3", "date": "2022-10-06", "group": "g2", "price": "4000", "quantity": "40", "published": false}))
	req.NoError(m.Create(ctx, &kv{"item": "i4", "date": "2022-10-06", "group": "g2", "price": "1000", "quantity": "10", "published": true}))
	req.NoError(m.Create(ctx, &kv{"item": "i5", "date": "2022-10-07", "group": "g2", "price": "1000", "quantity": "10", "published": false}))
	req.NoError(m.Create(ctx, &kv{"item": "i6", "date": "2022-10-07", "group": "g2", "price": "5000", "quantity": "50", "published": true}))

	// fetching two pages
	for p, result := range results {
		t.Logf("#%d Aggregating all records, calculating min & max price per group, ignoring empty quantity", p+1)
		i, err = m.Aggregate(
			filter.Generic(
				filter.WithCursor(cur),
				filter.WithLimit(1),
				filter.WithExpression("published"),
				filter.WithOrderBy(filter.SortExprSet{
					&filter.SortExpr{Column: "group", Descending: true},
					&filter.SortExpr{Column: "date", Descending: false},
				}),
			),
			// group-by
			[]dal.AggregateAttr{
				{Identifier: "date", Type: &dal.TypeDate{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
				{Identifier: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
			},
			// aggregation expressions
			[]dal.AggregateAttr{
				{Identifier: "count", RawExpr: "COUNT(*)", Type: &dal.TypeNumber{}},
				{Identifier: "max", RawExpr: "MAX(price)", Type: &dal.TypeNumber{}},
				{Identifier: "min", RawExpr: "MIN(price)", Type: &dal.TypeNumber{}},
				{Identifier: "avg", RawExpr: "AVG(price)", Type: &dal.TypeNumber{}},
				{Identifier: "stock", RawExpr: "SUM(quantity)", Type: &dal.TypeNumber{}},
			},
			// having
			qlParse(req, "SUM(quantity) > 0"),
		)
		req.NoError(err)
		req.NotNil(i)

		defer req.NoError(i.Close())

		// uncomment to se generated query
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

		t.Logf("#%d Iterating over results", p+1)
		rows := make([]kv, 0, 3)
		for i.Next(ctx) {
			row = kv{}
			req.NoError(i.Scan(row))

			// due to difference of number of decimal digits in different DBs, we need to do this
			// to make sure we get the same result
			row["avg"] = fmt.Sprintf("%.2f", cast.ToFloat64(row["avg"]))
			row["stock"] = fmt.Sprintf("%.2f", cast.ToFloat64(row["stock"]))

			rows = append(rows, row)
		}

		req.NoError(i.Err())
		if len(result) == 0 {
			req.Len(rows, 0)
		} else {
			req.Len(rows, 1)
			req.Equal(result, rows[0].String())

			cur, err = i.ForwardCursor(row)
			req.NotNil(cur)
			req.NoError(err)
		}

	}

}
