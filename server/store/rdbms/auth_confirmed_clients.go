package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/system/types"
)

func (s Store) convertAuthConfirmedClientFilter(f types.AuthConfirmedClientFilter) (query squirrel.SelectBuilder, err error) {
	query = s.authConfirmedClientsSelectBuilder()
	query = query.Where(squirrel.Eq{"acc.rel_user": f.UserID})
	return
}
