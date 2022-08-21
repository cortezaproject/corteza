package dal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRowEvaluatorTest(t *testing.T) {
	tc := []struct {
		name string
		expr string
		in   map[string]ValueGetter
		out  bool
	}{{
		name: "single row ok",
		expr: `row.test == 10`,
		in: map[string]ValueGetter{
			"row": (&row{}).WithValue("test", 0, 10),
		},
		out: true,
	}, {
		name: "single row nok",
		expr: `row.test == 11`,
		in: map[string]ValueGetter{
			"row": (&row{}).WithValue("test", 0, 10),
		},
		out: false,
	}, {
		name: "two rows ok",
		expr: `local.key == foreign.ref`,
		in: map[string]ValueGetter{
			"local":   (&row{}).WithValue("key", 0, 10),
			"foreign": (&row{}).WithValue("ref", 0, 10),
		},
		out: true,
	}, {
		name: "two rows nok",
		expr: `local.key == foreign.ref`,
		in: map[string]ValueGetter{
			"local":   (&row{}).WithValue("key", 0, 10),
			"foreign": (&row{}).WithValue("ref", 0, 11),
		},
		out: false,
	}}

	ctx := context.Background()

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			evl, err := newRunnerGval(c.expr)
			require.NoError(t, err)
			require.Equal(t, c.out, evl.Test(ctx, c.in))
		})
	}

}
