package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"strings"
)

func (s Store) convertFederationNodeFilter(f types.NodeFilter) (query squirrel.SelectBuilder, err error) {
	query = s.federationNodesSelectBuilder()

	query = filter.StateCondition(query, "fdn.deleted_at", f.Deleted)

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"LOWER(fdn.name)": q},
			squirrel.Like{"LOWER(fdn.base_url)": q},
		})
	}

	if f.Status != "" {
		query = query.Where(squirrel.Eq{"fdn.status": strings.ToLower(f.Status)})
	}

	return
}
