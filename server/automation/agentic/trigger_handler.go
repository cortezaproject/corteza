package agentic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	autoService "github.com/crusttech/human/server/automation/service"
	autoTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/agentic/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

type (
	triggerHandler struct {
		reg toolRegistrar
	}

	// triggerItem is the compact projection returned by list mode.
	//
	// Constraints and Input are the heavy fields and are deliberately absent:
	// a constraint set is the expressive part of a trigger but it is incidental
	// when picking one trigger out of a list, and one set per row would flood
	// the caller's context. A single-item lookup returns the service type
	// unchanged, so nothing is unreachable.
	triggerItem struct {
		TriggerID    string     `json:"triggerID"`
		WorkflowID   string     `json:"workflowID"`
		EventType    string     `json:"eventType"`
		ResourceType string     `json:"resourceType"`
		Enabled      bool       `json:"enabled"`
		DeletedAt    *time.Time `json:"deletedAt,omitempty"`
	}
)

// Declarations for these handlers are in trigger_tools.go, in the same order.
func TriggerHandler(reg toolRegistrar) *triggerHandler {
	h := &triggerHandler{reg: reg}
	h.register()
	return h
}

func (h *triggerHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// Reference param is never Required, so an absent triggerID means "list".
	triggerID, err := toolkit.ID(args, "triggerID")
	if err != nil {
		return nil, err
	}

	if triggerID != 0 {
		// Single-item mode returns the service type unchanged: the caller named
		// this trigger, so its constraints are the answer rather than noise.
		// The store lookup returns a trigger even when it is deleted, so this
		// is also the way to inspect one before restoring it.
		t, err := autoService.DefaultTrigger.FindByID(ctx, triggerID)
		if err != nil {
			return nil, toolkit.Errf("trigger lookup", err)
		}
		return toolkit.JSONResult(t)
	}

	f := autoTypes.TriggerFilter{
		EventType:    toolkit.Str(args, "eventType"),
		ResourceType: toolkit.Str(args, "resourceType"),
	}

	wfRef, err := toolkit.Ref(args, "workflow")
	if err != nil {
		return nil, err
	}
	if wfRef != "" {
		wf, err := h.resolveWorkflow(ctx, wfRef)
		if err != nil {
			return nil, err
		}
		f.WorkflowID = []string{strconv.FormatUint(wf.ID, 10)}
	}

	if toolkit.Bool(args, "includeDeleted") {
		f.Deleted = filter.StateInclusive
	}
	if toolkit.Bool(args, "includeDisabled") {
		f.Disabled = filter.StateInclusive
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := autoService.DefaultTrigger.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("trigger list", err)
	}

	items := make([]triggerItem, 0, len(set))
	for _, t := range set {
		items = append(items, triggerItem{
			TriggerID:    strconv.FormatUint(t.ID, 10),
			WorkflowID:   strconv.FormatUint(t.WorkflowID, 10),
			EventType:    t.EventType,
			ResourceType: t.ResourceType,
			Enabled:      t.Enabled,
			DeletedAt:    t.DeletedAt,
		})
	}

	// The cursor value, not its String(): String is a debug rendering, while
	// parseCursor expects the base64 form MarshalJSON emits. Nil marshals to
	// null, so the last page needs no special case.
	return toolkit.JSONResult(map[string]any{
		"triggers":       items,
		"nextPageCursor": out.NextPage,
	})
}

func (h *triggerHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	wfRef, err := toolkit.ReqRef(args, "workflow")
	if err != nil {
		return nil, err
	}

	wf, err := h.resolveWorkflow(ctx, wfRef)
	if err != nil {
		return nil, err
	}

	eventType, err := toolkit.ReqStr(args, "eventType")
	if err != nil {
		return nil, err
	}

	resourceType, err := toolkit.ReqStr(args, "resourceType")
	if err != nil {
		return nil, err
	}

	stepID, err := toolkit.ID(args, "stepID")
	if err != nil {
		return nil, err
	}

	// A trigger nobody asked to disable is meant to fire, so enabled defaults to
	// true; toolkit.Bool alone would read an absent flag as false and quietly
	// create a trigger that never runs.
	enabled := true
	if _, given := args["enabled"]; given {
		enabled = toolkit.Bool(args, "enabled")
	}

	t := &autoTypes.Trigger{
		WorkflowID:   wf.ID,
		StepID:       stepID,
		EventType:    eventType,
		ResourceType: resourceType,
		Enabled:      enabled,
	}

	if description := toolkit.Str(args, "description"); description != "" {
		t.Meta = &autoTypes.TriggerMeta{Description: description}
	}

	if t.Input, err = triggerInput(args["input"]); err != nil {
		return nil, err
	}

	res, err := autoService.DefaultTrigger.Create(ctx, t)
	if err != nil {
		return nil, toolkit.Errf("trigger create", err)
	}
	return toolkit.JSONResult(res)
}

// update is read-modify-write, and has to be.
//
// The service's Update assigns workflowID, enabled, eventType, resourceType and
// ownedBy from the passed trigger unconditionally, so handing it a struct that
// carries only the changed fields would blank everything else — a partial update
// has to start from the stored trigger. Loading it first also gives the stale
// check something real to compare against.
func (h *triggerHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	triggerID, err := toolkit.ReqID(args, "triggerID")
	if err != nil {
		return nil, err
	}

	existing, err := autoService.DefaultTrigger.FindByID(ctx, triggerID)
	if err != nil {
		return nil, toolkit.Errf("trigger lookup", err)
	}

	upd := existing.Clone()

	// workflow, eventType and resourceType are changeable but not clearable.
	// §8.2's present-and-empty rule clears a field; here clearing produces a
	// trigger that cannot fire, and in the workflow case one that can never be
	// loaded again — canManageTrigger resolves the workflow before every write,
	// so a trigger with workflowID 0 is permanently uneditable and undeletable.
	// An explicit empty value is refused instead of stored.
	if raw, given := args["workflow"]; given && raw != nil {
		wfRef, err := toolkit.Ref(args, "workflow")
		if err != nil {
			return nil, err
		}
		if wfRef == "" {
			return nil, fmt.Errorf("workflow cannot be cleared: a trigger with no workflow can no longer be edited, deleted or restored")
		}

		wf, err := h.resolveWorkflow(ctx, wfRef)
		if err != nil {
			return nil, err
		}
		upd.WorkflowID = wf.ID
	}

	if raw, given := args["eventType"]; given && raw != nil {
		eventType := toolkit.Str(args, "eventType")
		if eventType == "" {
			return nil, fmt.Errorf("eventType cannot be cleared: a trigger with no event type can never fire — pass a value from automation_event_type_lookup, or omit it to leave it unchanged")
		}
		upd.EventType = eventType
	}

	if raw, given := args["resourceType"]; given && raw != nil {
		resourceType := toolkit.Str(args, "resourceType")
		if resourceType == "" {
			return nil, fmt.Errorf("resourceType cannot be cleared: a trigger with no resource type can never fire — pass a value from automation_event_type_lookup, or omit it to leave it unchanged")
		}
		upd.ResourceType = resourceType
	}

	if raw, given := args["stepID"]; given && raw != nil {
		if upd.StepID, err = toolkit.ID(args, "stepID"); err != nil {
			return nil, err
		}
	}

	if _, given := args["enabled"]; given {
		upd.Enabled = toolkit.Bool(args, "enabled")
	}

	if raw, given := args["description"]; given && raw != nil {
		description := toolkit.Str(args, "description")
		if upd.Meta == nil {
			upd.Meta = &autoTypes.TriggerMeta{}
		}
		// Empty clears the note; the rest of Meta, the visual layout the
		// workflow editor keeps there, is carried over by Clone.
		upd.Meta.Description = description
	}

	if _, given := args["input"]; given {
		if upd.Input, err = triggerInput(args["input"]); err != nil {
			return nil, err
		}
		if upd.Input == nil {
			// The service reads a nil Input as "unchanged", so clearing has to
			// be expressed as an empty set rather than as nothing.
			if upd.Input, err = expr.NewVars(map[string]interface{}{}); err != nil {
				return nil, toolkit.Errf("trigger input build", err)
			}
		}
	}

	res, err := autoService.DefaultTrigger.Update(ctx, upd)
	if err != nil {
		return nil, toolkit.Errf("trigger update", err)
	}
	return toolkit.JSONResult(res)
}

func (h *triggerHandler) delete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	triggerID, err := toolkit.ReqID(args, "triggerID")
	if err != nil {
		return nil, err
	}

	if err = autoService.DefaultTrigger.DeleteByID(ctx, triggerID); err != nil {
		return nil, toolkit.Errf("trigger delete", err)
	}
	return toolkit.TextResult("trigger %d deleted", triggerID), nil
}

func (h *triggerHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	triggerID, err := toolkit.ReqID(args, "triggerID")
	if err != nil {
		return nil, err
	}

	if err = autoService.DefaultTrigger.UndeleteByID(ctx, triggerID); err != nil {
		return nil, toolkit.Errf("trigger undelete", err)
	}
	return toolkit.TextResult("trigger %d restored", triggerID), nil
}

// resolveWorkflow turns an ID-or-handle reference into a workflow.
//
// Triggers carry no handle of their own, so the workflow is the only thing a
// caller can name rather than number, and the trigger tools accept it for the
// same reason automation_workflow_lookup does: a handle is what a person reads
// off the workflow editor.
func (h *triggerHandler) resolveWorkflow(ctx context.Context, ref string) (*autoTypes.Workflow, error) {
	if ref == "" {
		return nil, fmt.Errorf("workflow identifier required")
	}

	if id, err := strconv.ParseUint(ref, 10, 64); err == nil {
		wf, err := autoService.DefaultWorkflow.FindByID(ctx, id)
		if err != nil {
			return nil, toolkit.Errf("workflow lookup", err)
		}
		return wf, nil
	}

	set, _, err := autoService.DefaultWorkflow.Search(ctx, autoTypes.WorkflowFilter{Handle: ref})
	if err != nil {
		return nil, toolkit.Errf("workflow lookup", err)
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("workflow %q not found", ref)
	}
	return set[0], nil
}

// triggerInput turns the 'input' argument into the fixed variables a trigger
// merges into its workflow's scope. Absent or empty yields nil, which every
// caller here treats as "not given".
func triggerInput(raw interface{}) (*expr.Vars, error) {
	inputMap, err := parseInput(raw)
	if err != nil {
		return nil, err
	}
	if inputMap == nil {
		return nil, nil
	}

	vars, err := expr.NewVars(inputMap)
	if err != nil {
		return nil, toolkit.Errf("trigger input build", err)
	}
	return vars, nil
}
