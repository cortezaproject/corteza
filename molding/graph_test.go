package molding

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	testEndEvent struct{ scope Variables }
)

var (
	_ Finalizer = &testEndEvent{}
)

func NewTestEndEvent() *testEndEvent                                  { return &testEndEvent{} }
func (t *testEndEvent) Finalize(_ context.Context, s Variables) error { t.scope = s; return nil }
func (testEndEvent) NodeRef() string                                  { return "end-event" }

func MustMakeNode(n Node, err error) Node {
	if err != nil {
		panic(err)
	}
	return n
}

func TestRun_Set(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()
		ee  = NewTestEndEvent()
		sa  = MustMakeNode(NewSetActivity(ee, Expr("foo", `"bar"`)))
	)

	req.NoError(Workflow(ctx, sa, Variables{}))
	req.Equal(ee.scope["foo"], "bar")
}

func TestRun_Condition(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()

		ee  = NewTestEndEvent()
		bar = MustMakeNode(NewSetActivity(ee, Expr("foo", `"bar"`)))
		baz = MustMakeNode(NewSetActivity(ee, Expr("foo", `"baz"`)))

		cnd = MustMakeNode(NewExclGateway(
			NewGatewayCondition("setFooToBar", bar),
			NewGatewayNoCondition(baz),
		))
	)

	req.NoError(Workflow(ctx, cnd, Variables{"setFooToBar": true}))
	req.Equal(ee.scope["foo"], "bar")
	req.NoError(Workflow(ctx, cnd, Variables{"setFooToBar": false}))
	req.Equal(ee.scope["foo"], "baz")
}

func TestRun_Loop(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()

		ee  = NewTestEndEvent()
		inc = MustMakeNode(NewSetActivity(nil, Expr("counter", `counter + 1`)))

		cnd = MustMakeNode(NewExclGateway(
			NewGatewayCondition("counter < 5", inc),
			NewGatewayNoCondition(ee),
		))
	)

	inc.(*setActivity).SetNext(cnd)

	req.NoError(Workflow(ctx, cnd, Variables{"counter": 0}))
	req.Equal(float64(5), ee.scope["counter"])
}

func TestRun_Join(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()

		ee  = NewTestEndEvent()
		foo = MustMakeNode(NewSetActivity(nil, Expr("foo", `"set"`), Expr("count", "count + 1"), Expr("order", "1")))
		bar = MustMakeNode(NewSetActivity(nil, Expr("bar", `"set"`), Expr("count", "count + 1"), Expr("order", "2")))
		baz = MustMakeNode(NewSetActivity(nil, Expr("baz", `"set"`), Expr("count", "count + 1"), Expr("order", "3")))

		fork = MustMakeNode(NewForkGateway(foo, bar, baz))
		join = MustMakeNode(NewJoinGateway(foo.(Iterator), bar.(Iterator), baz.(Iterator)))
	)

	join.(*joinGateway).SetNext(ee)

	req.NoError(Workflow(ctx, fork, Variables{
		"setFooToBar": true,
		"count":       0,
		"order":       0,
	}))

	// test if all paths are executed
	// all there should be set
	req.Equal(ee.scope["foo"], "set")
	req.Equal(ee.scope["bar"], "set")
	req.Equal(ee.scope["baz"], "set")

	// tests if paths are isolated;
	// expecting 1, all setters start with scope where count==0
	req.Equal(ee.scope["count"], float64(1))
	// tests if scope from paths merged in proper order
	// expecting 3 - value from the node last added to the join
	req.Equal(ee.scope["order"], float64(3))
}
