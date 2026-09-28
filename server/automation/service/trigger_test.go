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

func TestValidateTriggerInterval(t *testing.T) {
	var (
		req = require.New(t)

		trigger = func(event string, vv ...string) *types.Trigger {
			return &types.Trigger{
				EventType:   event,
				Constraints: types.TriggerConstraintSet{{Values: vv}},
			}
		}
	)

	req.NoError(validateTriggerInterval(trigger("onInterval", "0 6 * * *")))
	req.NoError(validateTriggerInterval(trigger("onInterval", "")))
	req.NoError(validateTriggerInterval(trigger("onTimestamp", "0 0 6 * * *")))
	req.Error(validateTriggerInterval(trigger("onInterval", "0 0 6 * * *")))
	req.Error(validateTriggerInterval(trigger("onInterval", "0 6 * * *", "every day")))
}
