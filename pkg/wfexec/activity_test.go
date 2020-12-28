package wfexec

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestActivities_NoArgs(t *testing.T) {
	var (
		req = require.New(t)
	)

	{
		a := Activity(
			func(ctx context.Context, variables Variables) (Variables, error) {
				return Variables{"foo": "42"}, nil
			},
			nil,
			nil,
		)

		r, err := a.Exec(context.Background(), &ExecRequest{})
		req.NoError(err)
		req.IsType(Variables{}, r)
		req.Empty(r)
	}

	{
		results := NewExpressions(expr.Parser())
		results.Set("bar", "foo")

		a := Activity(
			func(ctx context.Context, variables Variables) (Variables, error) {
				return Variables{"foo": 42}, nil
			},
			nil,
			results,
		)

		r, err := a.Exec(context.Background(), &ExecRequest{})
		req.NoError(err)
		req.NotNil(r)
		req.Contains(r, "bar")
		req.NotContains(r, "foo")
		req.Equal(42, r.(Variables)["bar"])
	}
}
