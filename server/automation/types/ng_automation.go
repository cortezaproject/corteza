package types

import (
	"strings"

	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	NgAutomationFilter struct {
		AutomationID []string `json:"automationID"`
		TenantID     uint64   `json:"tenantID,string,omitempty"`
		ProjectID    uint64   `json:"projectID,string,omitempty"`

		Handle string `json:"handle"`

		Query string `json:"query"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*NgAutomation) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	NgAutomationTriggerSchema []NgAutomationTriggerParam
)

const (
	NgAutomationSeverityError   = "error"
	NgAutomationSeverityWarning = "warning"
	NgAutomationSeverityInfo    = "info"
)

// Step kinds a TAQ's runtime understands.
//
// NgAutomationStep.Kind is a plain string on the generated type, unlike
// Workflow's typed WorkflowStepKind enum — so nothing stops a workflow kind
// name being stored on a TAQ step. It is not interchangeable: every workflow
// kind other than function, iterator, error and termination is unknown here,
// and "expressions" in particular is the one most likely to be carried across
// by hand.
//
// The authoritative switch is stepConv in
// server/automation/service/ng_automation_converter.go; these constants name
// the same set so in-tree callers can stop spelling it out, and
// IsValidNgAutomationStepKind gives a caller a cheap check before a write.
const (
	NgAutomationStepKindFunction         = "function"
	NgAutomationStepKindIterator         = "iterator"
	NgAutomationStepKindTermination      = "termination"
	NgAutomationStepKindGatewayExclusive = "gatewayExclusive"
	NgAutomationStepKindGatewayInclusive = "gatewayInclusive"
	NgAutomationStepKindError            = "error"
)

// NgAutomationStepKinds is the full set, in the order a caller is most likely
// to need them.
var NgAutomationStepKinds = []string{
	NgAutomationStepKindFunction,
	NgAutomationStepKindIterator,
	NgAutomationStepKindGatewayExclusive,
	NgAutomationStepKindGatewayInclusive,
	NgAutomationStepKindTermination,
	NgAutomationStepKindError,
}

// IsValidNgAutomationStepKind reports whether the runtime can execute a step of
// this kind. A step of any other kind is dropped from the executable graph.
func IsValidNgAutomationStepKind(kind string) bool {
	for _, k := range NgAutomationStepKinds {
		if k == kind {
			return true
		}
	}
	return false
}

const (
	IssueCodeScopeUnknown        = "scope.unknown"
	IssueCodeTypeMismatch        = "type.mismatch"
	IssueCodeStepDuplicateID     = "step.duplicateID"
	IssueCodeStepEmptyID         = "step.emptyID"
	IssueCodeStepIDCollision     = "step.idCollision"
	IssueCodeStepInvalid         = "step.invalid"
	IssueCodeTriggerDuplicateID  = "trigger.duplicateID"
	IssueCodeTriggerEmptyID      = "trigger.emptyID"
	IssueCodeTriggerMultiPaths   = "trigger.multiplePaths"
	IssueCodeFunctionUnknown     = "function.unknown"
	IssueCodePathEmpty           = "path.empty"
	IssueCodePathSelfLoop        = "path.selfLoop"
	IssueCodePathUnknownChild    = "path.unknownChild"
	IssueCodePathUnknownParent   = "path.unknownParent"
	IssueCodeGatewayTooFewPaths  = "gateway.tooFewPaths"
	IssueCodeGatewayNoElse       = "gateway.noElse"
	IssueCodeGatewayMultiElse    = "gateway.multipleElse"
	IssueCodeGraphNoEntry        = "graph.noEntry"
	IssueCodeGraphCycle          = "graph.cycle"
	IssueCodeGraphAmbiguousEntry = "graph.ambiguousEntry"
	IssueCodeAutomationNil       = "automation.nil"
	IssueCodeInternal            = "internal"
	IssueCodeRunAsLoadFailed     = "runAs.loadFailed"
	IssueCodeRunAsInvalid        = "runAs.invalid"
)

const (
	IssueDetailMissingReference = "missingReference"
	IssueDetailInvalidType      = "invalidType"
	IssueDetailDuplicateID      = "duplicateID"
	IssueDetailCycle            = "cycle"
	IssueDetailResourceRef      = "resourceRef"
	IssueDetailGatewayPaths     = "gatewayPaths"
	IssueDetailEmptyField       = "emptyField"
)

func (set NgAutomationIssueSet) Error() string {
	out := make([]string, 0, 4)
	for _, s := range set {
		out = append(out, s.Message)
	}

	return strings.Join(out, ", ")
}
