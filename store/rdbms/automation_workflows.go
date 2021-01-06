package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

func (s Store) convertAutomationWorkflowFilter(f types.WorkflowFilter) (query squirrel.SelectBuilder, err error) {
	query = s.automationWorkflowsSelectBuilder()

	query = filter.StateCondition(query, "atmwf.deleted_at", f.Deleted)
	query = filter.StateConditionNegBool(query, "atmwf.enabled", f.Disabled)

	if len(f.WorkflowID) > 0 {
		query = query.Where(squirrel.Eq{"atmwf.id": f.WorkflowID})
	}

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"atmwf.id": f.LabeledIDs})
	}

	if f.Query != "" {
		qs := f.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"atmwf.handle": qs},
		})
	}

	return
}
