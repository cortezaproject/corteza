package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertCredentialsFilter(f types.CredentialsFilter) (query squirrel.SelectBuilder, err error) {
	query = s.QueryCredentials()

	if f.Kind != "" {
		query = query.Where(squirrel.Eq{"crd.kind": f.Kind})
	}

	return
}
