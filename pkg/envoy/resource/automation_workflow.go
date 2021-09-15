package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/automation/types"
)

type (
	AutomationWorkflow struct {
		*base
		Res *types.Workflow

		Triggers []*AutomationTrigger
		Steps    []*AutomationWorkflowStep
		Paths    []*AutomationWorkflowPath
	}

	AutomationTrigger struct {
		*base
		Res *types.Trigger
	}

	AutomationWorkflowStep struct {
		*base
		Res *types.WorkflowStep
	}

	AutomationWorkflowPath struct {
		*base
		Res *types.WorkflowPath

		ParentStep Identifiers
		ChildStep  Identifiers
	}
)

func NewAutomationWorkflow(res *types.Workflow) *AutomationWorkflow {
	r := &AutomationWorkflow{
		base: &base{},
	}
	r.SetResourceType(types.WorkflowResourceType)
	r.Res = res

	r.AddIdentifier(identifiers(res.Handle, res.Meta.Name, res.ID)...)

	// Initial stamps
	r.SetTimestamps(MakeTimestampsCUDA(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, nil))
	us := MakeUserstampsCUDO(res.CreatedBy, res.UpdatedBy, res.DeletedBy, res.OwnedBy)
	us.RunAs = MakeUserstampFromID(res.RunAs)
	r.SetUserstamps(us)

	return r
}

func (r *AutomationWorkflow) AddAutomationTrigger(res *types.Trigger) *AutomationTrigger {
	t := &AutomationTrigger{
		base: &base{},
	}

	t.Res = res

	// Initial stamps
	t.SetTimestamps(MakeTimestampsCUDA(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, nil))
	t.SetUserstamps(MakeUserstampsCUDO(res.CreatedBy, res.UpdatedBy, res.DeletedBy, res.OwnedBy))

	if r.Triggers == nil {
		r.Triggers = make([]*AutomationTrigger, 0, 2)
	}
	r.Triggers = append(r.Triggers, t)

	return t
}

func (r *AutomationWorkflow) ResourceTranslationParts() (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref()
	path = nil
	resource = fmt.Sprintf(types.WorkflowResourceTranslationTpl(), types.WorkflowResourceTranslationType, firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Handle))

	return
}

func (r *AutomationWorkflow) AddAutomationWorkflowStep(res *types.WorkflowStep) *AutomationWorkflowStep {
	s := &AutomationWorkflowStep{
		base: &base{},
	}

	s.Res = res

	if r.Steps == nil {
		r.Steps = make([]*AutomationWorkflowStep, 0, 100)
	}
	r.Steps = append(r.Steps, s)

	return s
}

func (r *AutomationWorkflow) AddAutomationWorkflowPath(res *types.WorkflowPath) *AutomationWorkflowPath {
	p := &AutomationWorkflowPath{
		base: &base{},
	}

	p.Res = res

	if r.Paths == nil {
		r.Paths = make([]*AutomationWorkflowPath, 0, 100)
	}
	r.Paths = append(r.Paths, p)

	return p
}

func (r *AutomationWorkflow) SysID() uint64 {
	return r.Res.ID
}

// FindAutomationWorkflow looks for the workflow in the resource set
func FindAutomationWorkflow(rr InterfaceSet, ii Identifiers) (ns *types.Workflow) {
	var wfRes *AutomationWorkflow

	rr.Walk(func(r Interface) error {
		wr, ok := r.(*AutomationWorkflow)
		if !ok {
			return nil
		}

		if wr.Identifiers().HasAny(ii) {
			wfRes = wr
		}
		return nil
	})

	// Found it
	if wfRes != nil {
		return wfRes.Res
	}
	return nil
}

func AutomationWorkflowErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("automation workflow unresolved %v", ii.StringSlice())
}
