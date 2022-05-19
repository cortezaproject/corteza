package mysql

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	rdbmsdal "github.com/cortezaproject/corteza-server/store/adapters/rdbms/dal"
)

func init() {
	dal.Register(dalConnector, baseSchema, debugSchema)
}

func dalConnector(ctx context.Context, dsn string, cc ...capabilities.Capability) (_ dal.Connection, err error) {
	db, _, err := connectBase(ctx, dsn)
	if err != nil {
		return
	}
	return rdbmsdal.Connection(db, Dialect(), cc...), nil
}
