package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertCredentialsFilter(f types.CredentialsFilter) (query squirrel.SelectBuilder, err error) {
	query = s.credentialsSelectBuilder()

	query = filter.StateCondition(query, "crd.deleted_at", f.Deleted)

	if f.Kind != "" {
		query = query.Where(squirrel.Eq{"crd.kind": f.Kind})
	}

	if f.Credentials != "" {
		query = query.Where(squirrel.Eq{"crd.credentials": f.Credentials})
	}

	if f.OwnerID > 0 {
		query = query.Where(squirrel.Eq{"crd.rel_owner": f.OwnerID})
	}

	return
}
