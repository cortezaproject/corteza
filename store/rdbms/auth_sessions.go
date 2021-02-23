package rdbms

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertAuthSessionFilter(f types.AuthSessionFilter) (query squirrel.SelectBuilder, err error) {
	query = s.authSessionsSelectBuilder()

	if f.UserID > 0 {
		query = query.Where(squirrel.Eq{"ses.rel_user": f.UserID})
	}

	return
}
func (s Store) DeleteAuthSessionsByUserID(ctx context.Context, userID uint64) error {
	return s.execDeleteAuthSessions(ctx, squirrel.Eq{"ses.rel_user": userID})
}

func (s Store) DeleteExpiredAuthSessions(ctx context.Context) error {
	return s.execDeleteAuthSessions(ctx, squirrel.Lt{"ses.expires_at": "NOW()"})
}
