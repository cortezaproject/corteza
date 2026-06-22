package types

import (
	"fmt"

	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	WorkflowFilter struct {
		WorkflowID []string `json:"workflowID"`

		Handle string `json:"handle"`

		Query string `json:"query"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		// include sub-workflows
		SubWorkflow filter.State `json:"subWorkflow"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Workflow) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	WorkflowMeta struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`

		// list as one of the sub-workflows, when set to true
		// there should be no enabled triggers on this workflow
		SubWorkflow bool `json:"subWorkflow,omitempty"`

		// Named input/output contract, read by the "Run Workflow" step
		Input  []WorkflowIODef `json:"input,omitempty"`
		Output []WorkflowIODef `json:"output,omitempty"`
	}

	// WorkflowIODef declares a single named input or output field of a workflow
	WorkflowIODef struct {
		Name     string   `json:"name"`
		Label    string   `json:"label,omitempty"`
		Types    []string `json:"types,omitempty"`
		Required bool     `json:"required,omitempty"`
	}

	WorkflowIssue struct {
		// url encoded location of the error:
		Culprit     map[string]int `json:"culprit"`
		Description string         `json:"description"`
	}

	WorkflowExecParams struct {
		// When executed as a sub-workflow
		CallerWorkflowID uint64

		// When executed as a sub-workflow
		CallerSessionID uint64

		// When executed as a sub-workflow
		CallerStepID uint64

		// Start with this specific step
		StepID uint64

		EventType    string
		ResourceType string

		// Enable execution tracing
		Trace bool

		// Do not wait for workflow to be finished
		Async bool

		// Wait for workflow to be executed even if it's deferred
		Wait bool

		Input *expr.Vars
	}
)

// CheckDeferred returns true if any of the steps is deferred.
//
// Workflow is considered deferred when delay or prompt step types are used.
// Deferred workflows cannot short-circuit triggers or prevent creation/update on before triggers
//
// @todo add flag on workflow to explicitly mark workflow as deferred even when there are no delay or prompt steps
func (r Workflow) CheckDeferred() bool {
	return r.Steps.HasDeferred()
}

// Executable returns true if workflow is valid and enabled
func (r Workflow) Executable() bool {
	return r.DeletedAt == nil && r.Enabled
}

func (r Workflow) Dict() map[string]interface{} {
	return map[string]interface{}{
		"ID":         r.ID,
		"workflowID": r.ID,
		"labels":     r.Labels,
		"ownedBy":    r.OwnedBy,
		"createdAt":  r.CreatedAt,
		"createdBy":  r.CreatedBy,
		"updatedAt":  r.UpdatedAt,
		"updatedBy":  r.UpdatedBy,
		"deletedAt":  r.DeletedAt,
		"deletedBy":  r.DeletedBy,
	}
}

func (issue *WorkflowIssue) String() string {
	return fmt.Sprintf("%s [%v]", issue.Description, issue.Culprit)
}

func (set WorkflowIssueSet) Error() string {
	switch len(set) {
	case 0:
		return fmt.Sprintf("no workflow issue found")
	case 1:
		return fmt.Sprintf("1 workflow issue found")
	default:
		return fmt.Sprintf("%d workflow issues found", len(set))
	}
}

func (set WorkflowIssueSet) Append(err error, culprit map[string]int) WorkflowIssueSet {
	if culprit == nil {
		culprit = make(map[string]int)
	}

	return append(set, &WorkflowIssue{
		Culprit:     culprit,
		Description: err.Error(),
	})
}

// Distinct returns set of issues without duplicates
func (set WorkflowIssueSet) Distinct() (out WorkflowIssueSet) {
	idx := make(map[string]bool)
	out = make([]*WorkflowIssue, 0, len(set))

	for i := range set {
		if idx[set[i].String()] {
			continue
		}

		out = append(out, set[i])
		idx[set[i].String()] = true
	}

	return
}

func (set WorkflowIssueSet) SetCulprit(name string, pos int) WorkflowIssueSet {
	for i := range set {
		set[i].Culprit[name] = pos
	}

	return set
}

