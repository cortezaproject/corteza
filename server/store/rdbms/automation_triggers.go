package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

func (s Store) convertAutomationTriggerFilter(f types.TriggerFilter) (query squirrel.SelectBuilder, err error) {
	query = s.automationTriggersSelectBuilder()

	query = filter.StateCondition(query, "atmt.deleted_at", f.Deleted)
	query = filter.StateConditionNegBool(query, "atmt.enabled", f.Disabled)

	if len(f.TriggerID) > 0 {
		query = query.Where(squirrel.Eq{"atmt.id": f.TriggerID})
	}

	if len(f.WorkflowID) > 0 {
		query = query.Where(squirrel.Eq{"atmt.rel_workflow": f.WorkflowID})
	}

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"atmt.id": f.LabeledIDs})
	}

	if len(f.EventType) > 0 {
		query = query.Where(squirrel.Eq{"atmt.event_type": f.EventType})
	}

	if len(f.ResourceType) > 0 {
		query = query.Where(squirrel.Eq{"atmt.resource_type": f.ResourceType})
	}

	return
}
