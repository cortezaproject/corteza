package yaml

import (
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	automationWorkflow struct {
		res      *types.Workflow
		triggers automationTriggerSet
		steps    automationWorkflowStepSet
		paths    automationWorkflowPathSet

		ts *resource.Timestamps
		us *resource.Userstamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		rbac rbacRuleSet
	}
	automationWorkflowSet []*automationWorkflow

	automationTrigger struct {
		res *types.Trigger

		ts *resource.Timestamps
		us *resource.Userstamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig
	}
	automationTriggerSet []*automationTrigger

	automationWorkflowStep struct {
		res *types.WorkflowStep

		ts *resource.Timestamps
		us *resource.Userstamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig
	}
	automationWorkflowStepSet []*automationWorkflowStep

	automationWorkflowPath struct {
		res *types.WorkflowPath

		ts *resource.Timestamps
		us *resource.Userstamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig
	}
	automationWorkflowPathSet []*automationWorkflowPath
)

func (nn automationWorkflowSet) ConfigureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.encoderConfig = cfg
	}
}
