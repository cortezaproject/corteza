package dal_test

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/logger"
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
			Ident: "test_dal_select",
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

	//ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

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

func TestModel_Aggregate(t *testing.T) {
	_ = logger.Default()

	var (
		req = require.New(t)

		ctx = context.Background()

		baseModel = &dal.Model{
			Ident: "test_dal_aggregation",
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
		[]*dal.AggregateAttr{
			{Identifier: "date", Type: &dal.TypeDate{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
			{Identifier: "group", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "values"}},
		},
		// aggregation expressions
		[]*dal.AggregateAttr{
			{Identifier: "count", RawExpr: "COUNT(*)", Type: &dal.TypeNumber{}},
			{Identifier: "max", RawExpr: "MAX(price)", Type: &dal.TypeNumber{}},
			{Identifier: "min", RawExpr: "MIN(price)", Type: &dal.TypeNumber{}},
			{Identifier: "avg", RawExpr: "AVG(price)", Type: &dal.TypeNumber{}},
			{Identifier: "stock", RawExpr: "SUM(quantity)", Type: &dal.TypeNumber{}},
		},
		"", // <== here be having condition
	)
	req.NoError(err)
	req.NotNil(i)

	defer req.NoError(i.Close())

	// uncomment to se generated query
	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

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
