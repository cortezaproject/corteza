package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIndexManagement(t *testing.T) {
	var (
		ctx = context.Background()
	)

	eachDB(t, func(t *testing.T, c *conn) error {
		if c.isSQLite {
			t.Skip("sqlite does not support index management")
		}

		var (
			req     = require.New(t)
			dd      = c.store.DataDefiner(c.db)
			ii, err = dd.IndexLookup(ctx, ddl.PRIMARY_KEY, "users")
		)

		req.NoError(err)
		req.NotNil(ii)

		return nil
	})
}
