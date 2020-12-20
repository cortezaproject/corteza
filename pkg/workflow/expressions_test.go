package workflow

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpressions(t *testing.T) {
	var (
		req = require.New(t)
		ee  = NewExpressions(expr.Parser())
	)

	req.NoError(ee.Set("foo", `40`))
	req.NoError(ee.Set("bar", `foo + 2`))

	r, err := ee.Exec(context.Background(), &ExecRequest{})
	req.NoError(err)
	req.NotNil(r)
	req.Contains(r, "foo")
	req.Contains(r, "bar")
	req.Equal(40, r.(Variables).Int("foo", 0))
	req.Equal(42, r.(Variables).Int("bar", 0))
}
