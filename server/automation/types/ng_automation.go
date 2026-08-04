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
