package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

func (s Store) convertAutomationSessionFilter(f types.SessionFilter) (query squirrel.SelectBuilder, err error) {
	query = s.automationSessionsSelectBuilder()

	query = filter.StateCondition(query, "atms.completed_at", f.Completed)

	if len(f.SessionID) > 0 {
		query = query.Where(squirrel.Eq{"atms.id": f.SessionID})
	}

	if len(f.Status) > 0 {
		query = query.Where(squirrel.Eq{"atms.status": f.Status})
	}

	if len(f.WorkflowID) > 0 {
		query = query.Where(squirrel.Eq{"atms.id": f.WorkflowID})
	}

	if len(f.EventType) > 0 {
		query = query.Where(squirrel.Eq{"atms.event_type": f.EventType})
	}

	if len(f.ResourceType) > 0 {
		query = query.Where(squirrel.Eq{"atms.resource_type": f.ResourceType})
	}

	return
}
