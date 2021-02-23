package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertAuthClientFilter(f types.AuthClientFilter) (query squirrel.SelectBuilder, err error) {
	query = s.authClientsSelectBuilder()

	query = filter.StateCondition(query, "ac.deleted_at", f.Deleted)

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"ac.id": f.LabeledIDs})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"ac.handle": f.Handle})
	}

	return
}
