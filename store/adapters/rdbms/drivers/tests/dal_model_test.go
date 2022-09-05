package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModelManagement(t *testing.T) {
	var (
		ctx = context.Background()
		log = logger.Default()
	)

	eachDB(t, func(t *testing.T, c *conn) error {
		var (
			req = require.New(t)
		)

		if c.isSQLite {
			t.Skip("sqlite does not support index management")
		}

		const (
			dalConnID    = 1
			dalTableName = "lil_dal_test"
		)

		_, err := c.db.Exec("DROP TABLE IF EXISTS " + dalTableName)
		req.NoError(err)

		cw := dal.MakeConnection(
			dalConnID,
			c.store.ToDalConn(),
			dal.ConnectionParams{},
			dal.ConnectionConfig{},
		)

		svc, err := dal.New(log, true)
		req.NoError(err)
		req.NoError(svc.ReplaceConnection(ctx, cw, true))

		req.NoError(svc.ReplaceModel(ctx, &dal.Model{
			Ident:        dalTableName,
			ConnectionID: dalConnID,
			Attributes: dal.AttributeSet{
				&dal.Attribute{
					Ident:      "ID",
					PrimaryKey: true,
					Type:       &dal.TypeID{},
					Store:      &dal.CodecAlias{Ident: "id"},
				},

				&dal.Attribute{
					Ident:    "OwnerID",
					Sortable: true,
					Type:     &dal.TypeText{},
					Store:    &dal.CodecAlias{Ident: "rel_owner"},
				},
			},
		}))

		_, err = c.db.Exec("DROP TABLE " + dalTableName)
		req.NoError(err)

		return nil
	})
}
