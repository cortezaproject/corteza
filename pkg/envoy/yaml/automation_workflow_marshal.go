package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/envoy/util"
)

func automationWorkflowFromResource(r *resource.AutomationWorkflow, cfg *EncoderConfig) *automationWorkflow {
	tt := make(automationTriggerSet, len(r.Triggers))
	for i, t := range r.Triggers {
		tt[i] = &automationTrigger{
			res:           t.Res,
			encoderConfig: cfg,
		}
	}

	ss := make(automationWorkflowStepSet, len(r.Steps))
	for i, s := range r.Steps {
		ss[i] = &automationWorkflowStep{
			res:           s.Res,
			encoderConfig: cfg,
		}
	}

	pp := make(automationWorkflowPathSet, len(r.Paths))
	for i, p := range r.Paths {
		pp[i] = &automationWorkflowPath{
			res:           p.Res,
			encoderConfig: cfg,
		}
	}

	return &automationWorkflow{
		res:      r.Res,
		triggers: tt,
		steps:    ss,
		paths:    pp,

		encoderConfig: cfg,
	}
}

// Prepare prepares the automationWorkflow to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *automationWorkflow) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	wf, ok := state.Res.(*resource.AutomationWorkflow)
	if !ok {
		return encoderErrInvalidResource(resource.AUTOMATION_WORKFLOW_RESOURCE_TYPE, state.Res.ResourceType())
	}

	n.res = wf.Res
	n.us = wf.Userstamps()

	return nil
}

// Encode encodes the automationWorkflow to the document
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *automationWorkflow) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	if n.res.ID <= 0 {
		n.res.ID = util.NextID()
	}

	n.ts, err = resource.MakeCUDATimestamps(&n.res.CreatedAt, n.res.UpdatedAt, n.res.DeletedAt, nil).
		Model(n.encoderConfig.TimeLayout, n.encoderConfig.Timezone)
	if err != nil {
		return err
	}
	n.us, err = resolveUserstamps(state.ParentResources, n.us)
	if err != nil {
		return err
	}

	// @todo skip eval?

	doc.AddAutomationWorkflow(n)

	return err
}

func (wf *automationWorkflow) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"handle", wf.res.Handle,
		"meta", wf.res.Meta,
		"enabled", wf.res.Enabled,

		"trace", wf.res.Trace,
		"keepSessions", wf.res.KeepSessions,

		"scope", wf.res.Scope,
		"triggers", wf.triggers,
		"steps", wf.res.Steps,
		"paths", wf.res.Paths,

		// "issues", wf.res.Issues,
		"labels", wf.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	nn, err = mapTimestamps(nn, wf.ts)
	if err != nil {
		return nil, err
	}

	nn, err = mapUserstamps(nn, wf.us)
	if err != nil {
		return nil, err
	}

	return nn, nil
}

func (t *automationTrigger) MarshalYAML() (interface{}, error) {
	var err error

	nn, err := makeMap(
		"resourceType", t.res.ResourceType,
		"eventType", t.res.EventType,
		"constraints", t.res.Constraints,
		"enabled", t.res.Enabled,
		"workflowID", t.res.WorkflowID,

		"stepID", t.res.StepID,
		"input", t.res.Input,

		"meta", t.res.Meta,

		"labels", t.res.Labels,
	)
	if err != nil {
		return nil, err
	}

	nn, err = mapTimestamps(nn, t.ts)
	if err != nil {
		return nil, err
	}

	nn, err = mapUserstamps(nn, t.us)
	if err != nil {
		return nil, err
	}

	return nn, nil
}
