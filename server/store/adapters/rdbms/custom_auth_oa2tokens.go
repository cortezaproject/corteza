package rdbms

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
)

func (s Store) DeleteExpiredAuthOA2Tokens(ctx context.Context) error {
	return s.Exec(ctx, authOa2tokenDeleteQuery(s.Dialect.GOQU(), goqu.C("expires_at").Lt(time.Now())))
}

func (s Store) DeleteAuthOA2TokenByCode(ctx context.Context, code string) error {
	return s.Exec(ctx, authOa2tokenDeleteQuery(s.Dialect.GOQU(), goqu.C("code").Eq(code)))
}

func (s Store) DeleteAuthOA2TokenByAccess(ctx context.Context, access string) error {
	return s.Exec(ctx, authOa2tokenDeleteQuery(s.Dialect.GOQU(), goqu.C("access").Eq(access)))
}

func (s Store) DeleteAuthOA2TokenByRefresh(ctx context.Context, refresh string) error {
	return s.Exec(ctx, authOa2tokenDeleteQuery(s.Dialect.GOQU(), goqu.C("refresh").Eq(refresh)))
}

func (s Store) DeleteAuthOA2TokenByUserID(ctx context.Context, userID uint64) error {
	return s.Exec(ctx, authOa2tokenDeleteQuery(s.Dialect.GOQU(), goqu.C("rel_user").Eq(userID)))
}
