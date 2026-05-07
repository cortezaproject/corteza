package rdbms

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

func (s Store) DeleteExpiredAuthSessions(ctx context.Context) error {
	return s.Exec(ctx, authSessionDeleteQuery(s.Dialect.GOQU(), goqu.C("expires_at").Lt(time.Now())))
}

func (s Store) DeleteAuthSessionsByUserID(ctx context.Context, userID uint64) error {
	return s.Exec(ctx, authSessionDeleteQuery(s.Dialect.GOQU(), goqu.C("rel_user").Eq(userID)))
}
