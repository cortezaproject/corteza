package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (s Store) convertWorkflowFilter(f types.WorkflowFilter) (query squirrel.SelectBuilder, err error) {
	query = s.workflowsSelectBuilder()

	query = filter.StateCondition(query, "usr.deleted_at", f.Deleted)

	if len(f.WorkflowID) > 0 {
		query = query.Where(squirrel.Eq{"usr.id": f.WorkflowID})
	}

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"usr.id": f.LabeledIDs})
	}

	if f.Query != "" {
		qs := f.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"usr.handle": qs},
		})
	}

	return
}
