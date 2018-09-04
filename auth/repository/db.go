package repository

import (
	"context"
	"github.com/titpetric/factory"
)

var _db *factory.DB

func DB(ctxs ...context.Context) *factory.DB {
	if _db == nil {
		_db = factory.Database.MustGet()
	}
	for _, ctx := range ctxs {
		_db = _db.With(ctx)
		break
	}
	return _db
}
