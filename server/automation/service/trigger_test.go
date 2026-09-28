package service

import (
	"github.com/crusttech/human/server/automation/types"
	sysEvent "github.com/crusttech/human/server/system/service/event"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateWorkflowTriggersEmpty(t *testing.T) {
	var (
		req = require.New(t)

		issues = validateWorkflowTriggers(
			&types.Workflow{},
		)
	)

	req.Empty(issues)
}
func TestValidateWorkflowTriggersRunAs(t *testing.T) {
	var (
		req = require.New(t)
		soi = sysEvent.SystemOnInterval()

		issues = validateWorkflowTriggers(
			&types.Workflow{},
			&types.Trigger{
				Enabled:      true,
				ResourceType: soi.ResourceType(),
				EventType:    soi.EventType(),
			},
		)
	)

	req.Len(issues, 1)
	req.Contains(issues[0].String(), "requires run-as to be set")
}

func TestValidateWorkflowTriggersSubWorkflow(t *testing.T) {
	var (
		req = require.New(t)

		issues = validateWorkflowTriggers(
			&types.Workflow{Meta: &types.WorkflowMeta{SubWorkflow: true}},
			&types.Trigger{Enabled: true},
		)
	)

	req.Len(issues, 1)
	req.Contains(issues[0].String(), "marked as sub-workflow")
}

func TestValidateTrigger(t *testing.T) {
	var (
		req = require.New(t)

		trigger = func(event string, enabled bool, cc ...*types.TriggerConstraint) *types.Trigger {
			return &types.Trigger{EventType: event, Enabled: enabled, Constraints: cc}
		}

		value = func(op string, vv ...string) *types.TriggerConstraint {
			return &types.TriggerConstraint{Name: "module", Op: op, Values: vv}
		}
	)

	req.NoError(validateTrigger(trigger("onInterval", true, value("", "0 6 * * *"))))
	req.NoError(validateTrigger(trigger("onInterval", true, value("", "@daily"))))
	req.NoError(validateTrigger(trigger("onInterval", false, value("", ""))))
	req.NoError(validateTrigger(trigger("onInterval", false)))
	req.NoError(validateTrigger(trigger("onTimestamp", true, value("", "2026-10-01T06:00:00Z"))))
	req.NoError(validateTrigger(trigger("afterUpdate", true, value("=", "ticket"))))
	req.NoError(validateTrigger(trigger("afterUpdate", true, value("~", "^tick"))))
	req.NoError(validateTrigger(trigger("afterUpdate", true)))

	req.Error(validateTrigger(trigger("onInterval", true)))
	req.Error(validateTrigger(trigger("onInterval", true, value("", ""))))
	req.Error(validateTrigger(trigger("onTimestamp", true)))
	req.Error(validateTrigger(trigger("onInterval", true, value("", "0 0 6 * * *"))))
	req.Error(validateTrigger(trigger("onInterval", true, value("", "0 6 * * *", "every day"))))
	req.Error(validateTrigger(trigger("onInterval", false, value("", "every day"))))
	req.Error(validateTrigger(trigger("onTimestamp", true, value("", "tomorrow at six"))))
	req.Error(validateTrigger(trigger("afterUpdate", true, value("between", "a"))))
	req.Error(validateTrigger(trigger("afterUpdate", true, value("~", "("))))
}
