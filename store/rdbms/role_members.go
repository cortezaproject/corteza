package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertRoleMemberFilter(f types.RoleMemberFilter) (query squirrel.SelectBuilder, err error) {
	query = s.roleMembersSelectBuilder()

	if f.RoleID > 0 {
		query = query.Where(squirrel.Eq{"rm.rel_role": f.RoleID})
	}

	if f.UserID > 0 {
		query = query.Where(squirrel.Eq{"rm.rel_user": f.UserID})
	}

	return
}
