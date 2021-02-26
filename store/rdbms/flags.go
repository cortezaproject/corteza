package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/flag/types"
)

func (s Store) convertFlagFilter(f types.FlagFilter) (query squirrel.SelectBuilder, err error) {
	query = s.flagsSelectBuilder()

	query = query.Where(squirrel.Eq{"flg.kind": f.Kind})

	if len(f.ResourceID) > 0 {
		query = query.Where(squirrel.Eq{"flg.rel_resource": f.ResourceID})
	}

	if len(f.OwnedBy) > 0 {
		query = query.Where(squirrel.Eq{"flg.owned_by": f.OwnedBy})
	}

	if len(f.Name) > 0 {
		query = query.Where(squirrel.Eq{"flg.name": f.Name})
	}

	return
}
