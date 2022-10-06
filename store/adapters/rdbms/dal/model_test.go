package dal_test

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	. "github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestModel_Aggregate(t *testing.T) {
	var (
		req = require.New(t)

		ctx = context.Background()

		baseModel = &dal.Model{
			Ident: "test_dal_Aggregation",
			Attributes: []*dal.Attribute{
				{Ident: "item", Type: &dal.TypeText{}},
				{Ident: "group", Type: &dal.TypeText{}},
				{Ident: "price", Type: &dal.TypeNumber{}, Filterable: true},
				{Ident: "published", Type: &dal.TypeBoolean{}},
			},
		}

		m = Model(baseModel, s.DB, s.Dialect)

		i dal.Iterator

		table, err = s.DataDefiner.ConvertModel(baseModel)
		row        *kv
	)

	ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	table.Temporary = true
	req.NoError(s.DataDefiner.TableCreate(ctx, table))

	t.Log("Inserting test data")
	bm := time.Now()
	ctx = context.Background() // no need to log inserts
	count := 0
	for i := 1; i <= 1000; i++ {
		for g := 1; g <= 5; g++ {
			req.NoError(m.Create(ctx, &kv{
				"item":      fmt.Sprintf("i%d", i),
				"group":     fmt.Sprintf("g%d", g),
				"price":     (100000 * g) + i,
				"published": i%2 == 0,
			}))
			count++
		}
	}
	t.Logf("inserted %d entries in %v", count, time.Now().Sub(bm))

	t.Log("Aggregating all records, calculating min & max price per group")
	i, err = m.Aggregate(
		filter.Generic(
			filter.WithOrderBy(filter.SortExprSet{
				&filter.SortExpr{Column: "group", Descending: true},
			}),
		),
		// group-by
		[]*dal.AggregateAttr{
			{Identifier: "group", Type: &dal.TypeText{}},
		},
		// aggregation expressions
		[]*dal.AggregateAttr{
			{Identifier: "max", RawExpr: "MAX(price)", Type: &dal.TypeNumber{}},
			{Identifier: "min", RawExpr: "MIN(price)", Type: &dal.TypeNumber{}},
		},
		"", // <== here be having condition
	)
	req.NoError(err)
	req.NotNil(i)

	defer req.NoError(i.Close())

	// ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

	t.Log("Iterating over results")
	rows := make([]*kv, 0, 5)
	for i.Next(ctx) {
		row = &kv{}
		req.NoError(i.Scan(row))
		rows = append(rows, row)
	}

	req.Len(rows, 5)
	req.Equal(&kv{"group": "g5", "max": "501000", "min": "500001"}, rows[0])
	req.Equal(&kv{"group": "g4", "max": "401000", "min": "400001"}, rows[1])
	req.Equal(&kv{"group": "g3", "max": "301000", "min": "300001"}, rows[2])
	req.Equal(&kv{"group": "g2", "max": "201000", "min": "200001"}, rows[3])
	req.Equal(&kv{"group": "g1", "max": "101000", "min": "100001"}, rows[4])
}
