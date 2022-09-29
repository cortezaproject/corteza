package dal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRowEvaluatorTest(t *testing.T) {
	tc := []struct {
		name string
		expr string
		in   *Row
		out  bool
	}{{
		name: "single row ok",
		expr: `test == 10`,
		in:   (&Row{}).WithValue("test", 0, 10),
		out:  true,
	}, {
		name: "single row nok",
		expr: `test == 11`,
		in:   (&Row{}).WithValue("test", 0, 10),
		out:  false,
	}}

	ctx := context.Background()

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			evl, err := newRunnerGval(c.expr)
			require.NoError(t, err)
			tt, err := evl.Test(ctx, c.in)
			require.NoError(t, err)
			require.Equal(t, c.out, tt)
		})
	}
}

func TestHandlers(t *testing.T) {
	yr := time.Now().Year()
	mnt := int(time.Now().Month())
	qtr := mnt / 4
	dy := time.Now().Day()

	// @todo some of the tests were omitted because they would be too flaky
	//       or can not yet pass (the gval nil handlers primarily)
	tcc := []struct {
		name string
		expr string
		in   simpleRow
		out  any
	}{
		{
			name: "not",
			expr: `!true`,
			in:   simpleRow{},
			out:  false,
		},
		{
			name: "and",
			expr: `true && false`,
			in:   simpleRow{},
			out:  false,
		},
		{
			name: "or",
			expr: `false || true`,
			in:   simpleRow{},
			out:  true,
		},
		{
			name: "xor",
			expr: `'a' xor 'b'`,
			in:   simpleRow{},
			out:  true,
		},
		{
			name: "eq",
			expr: `1 == 2`,
			in:   simpleRow{},
			out:  false,
		},
		{
			name: "ne",
			expr: `3 != 2`,
			in:   simpleRow{},
			out:  true,
		},
		{
			name: "lt",
			expr: `1 < 2`,
			in:   simpleRow{},
			out:  true,
		},
		{
			name: "le",
			expr: `1 <= 2`,
			in:   simpleRow{},
			out:  true,
		},
		{
			name: "gt",
			expr: `2 > 1`,
			in:   simpleRow{},
			out:  true,
		},
		{
			name: "ge",
			expr: `1 >= 2`,
			in:   simpleRow{},
			out:  false,
		},
		{
			name: "add",
			expr: `1+3`,
			in:   simpleRow{},
			out:  float64(4),
		},
		{
			name: "sub",
			expr: `3-5`,
			in:   simpleRow{},
			out:  float64(-2),
		},
		{
			name: "mult",
			expr: `1*2`,
			in:   simpleRow{},
			out:  float64(2),
		},
		{
			name: "div",
			expr: `2/2`,
			in:   simpleRow{},
			out:  float64(1),
		},
		{
			name: "concat",
			expr: `concat('a', 'b')`,
			in:   simpleRow{},
			out:  "ab",
		},
		{
			name: "quarter",
			expr: `quarter(now())`,
			in:   simpleRow{},
			out:  qtr,
		},
		{
			name: "year",
			expr: `year(now())`,
			in:   simpleRow{},
			out:  yr,
		},
		{
			name: "month",
			expr: `month(now())`,
			in:   simpleRow{},
			out:  mnt,
		},
		{
			name: "day",
			expr: `day(now())`,
			in:   simpleRow{},
			out:  dy,
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			evl, err := newRunnerGval(tc.expr)
			require.NoError(t, err)
			out, err := evl.Eval(context.Background(), tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.out, out)
		})
	}
}
