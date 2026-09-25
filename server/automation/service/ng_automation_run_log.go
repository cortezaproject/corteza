package service

import (
	"context"
	"strconv"

	"github.com/crusttech/human/server/pkg/actionlog"
	execTypes "github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/id"
)

const (
	// NgAutomationRunAction is the action-log entry every finished TAQ run leaves behind
	NgAutomationRunAction = "run"

	ngAutomationRunResource = "automation:ng-automation"
)

// recordNgAutomationRun writes one action-log entry per terminal TAQ execution.
//
// The execution ledger is in-memory, so this entry is the durable record of a
// run; the admin dashboard counts runs and failures from it.
func recordNgAutomationRun(ctx context.Context, ex execTypes.Execution) {
	if DefaultActionlog == nil {
		return
	}

	a := &actionlog.Action{
		RequestOrigin: actionlog.RequestOrigin_Automation,
		Resource:      ngAutomationRunResource,
		Action:        NgAutomationRunAction,
		Severity:      actionlog.Notice,
		Description:   "TAQ run " + string(ex.Status),
		Meta: actionlog.Meta{
			"automationID": plainID(ex.ExecutableID),
			"executionID":  plainID(ex.ID),
			"revision":     ex.Revision,
			"status":       string(ex.Status),
			"eventType":    ex.EventType,
			"resourceType": ex.ResourceType,
		},
	}

	if ex.EndedAt != nil {
		a.Meta["durationMs"] = ex.EndedAt.Sub(ex.CreatedAt).Milliseconds()
	}

	if ex.Status == execTypes.StatusFailed {
		a.Severity = actionlog.Error
		if ex.Error != nil {
			a.Error = ex.Error.Error()
		}
	}

	DefaultActionlog.Record(ctx, a)
}

// plainID renders an ID without the quotes id.ID.String() adds
func plainID(v id.ID) string {
	if v.IsZero() {
		return ""
	}

	if n := v.Num(); n != 0 {
		return strconv.FormatUint(n, 10)
	}

	return v.Str()
}
