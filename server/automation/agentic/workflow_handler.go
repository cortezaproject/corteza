package agentic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	autoService "github.com/crusttech/human/server/automation/service"
	autoTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

type (
	workflowHandler struct {
		reg toolRegistrar
	}

	// workflowSummary is the list-mode projection.
	//
	// A workflow carries its whole step and path graph, so a list of raw
	// workflows would flood the caller's context with definitions it did not ask
	// for. The summary carries only what is needed to pick one out of a list;
	// the caller then fetches that one workflow in full.
	workflowSummary struct {
		WorkflowID string `json:"workflowID"`
		Handle     string `json:"handle"`
		Name       string `json:"name,omitempty"`
		Enabled    bool   `json:"enabled"`
	}

	// workflowWritten is what create and update return.
	//
	// The workflow already carries its issues, but they are omitempty and buried
	// under the graph the caller just sent, and the service stores a broken
	// workflow rather than refusing it — so a tool that returned the raw type
	// would report success on a workflow that can never run. Issues are lifted
	// out, always present, and accompanied by a warning a model cannot read past.
	workflowWritten struct {
		Workflow *autoTypes.Workflow        `json:"workflow"`
		Issues   autoTypes.WorkflowIssueSet `json:"issues"`
		Warning  string                     `json:"warning,omitempty"`
	}
)

// Declarations for these handlers are in workflow_tools.go, in the same order.
func WorkflowHandler(reg toolRegistrar) *workflowHandler {
	h := &workflowHandler{reg: reg}
	h.register()
	return h
}

func (h *workflowHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	if err = h.refuseNumericRef(args, "workflow"); err != nil {
		return nil, err
	}

	if toolkit.Str(args, "workflow") != "" {
		// Single-item mode returns the raw service type: the caller named this
		// workflow, so its graph is the answer rather than noise.
		wf, err := h.resolve(ctx, args)
		if err != nil {
			return nil, err
		}
		return toolkit.JSONResult(wf)
	}

	page := toolkit.Page(args)
	paging, err := filter.NewPaging(page.Limit, page.Cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	query := toolkit.Str(args, "query")

	set, f, err := h.search(ctx, query, paging, toolkit.Bool(args, "includeDisabled"))
	if err != nil {
		return nil, err
	}

	// Handles are slugs, so a model searching for what it sees in the UI ("My
	// Workflow") matches nothing. Retry once with the slugified query — but not
	// mid-pagination, where switching the query would silently change what the
	// caller is paging through.
	if len(set) == 0 && query != "" && page.Cursor == "" {
		if slugQuery := strings.ReplaceAll(strings.ToLower(query), " ", "_"); slugQuery != query {
			if set, f, err = h.search(ctx, slugQuery, paging, toolkit.Bool(args, "includeDisabled")); err != nil {
				return nil, err
			}
		}
	}

	items := make([]workflowSummary, len(set))
	for i, wf := range set {
		items[i] = workflowSummary{
			WorkflowID: strconv.FormatUint(wf.ID, 10),
			Handle:     wf.Handle,
			Enabled:    wf.Enabled,
		}
		if wf.Meta != nil {
			items[i].Name = wf.Meta.Name
		}
	}

	return toolkit.JSONResult(map[string]any{
		"workflows":      items,
		"nextPageCursor": f.NextPage,
	})
}

func (h *workflowHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	hdl, err := toolkit.ReqStr(args, "handle")
	if err != nil {
		return nil, err
	}

	name, err := toolkit.ReqStr(args, "name")
	if err != nil {
		return nil, err
	}

	wf := &autoTypes.Workflow{
		Handle: hdl,
		Meta: &autoTypes.WorkflowMeta{
			Name:        name,
			Description: toolkit.Str(args, "description"),
			SubWorkflow: toolkit.Bool(args, "subWorkflow"),
		},
		Trace: toolkit.Bool(args, "trace"),

		// A workflow nobody asked to disable is meant to run. toolkit.Bool alone
		// would read an absent flag as false and quietly store a workflow that
		// exec reports as not found.
		Enabled: true,
	}

	if _, given := args["enabled"]; given {
		wf.Enabled = toolkit.Bool(args, "enabled")
	}

	if wf.RunAs, err = toolkit.ID(args, "runAs"); err != nil {
		return nil, err
	}

	if _, err = toolkit.JSONArg(args, "steps", "steps", &wf.Steps); err != nil {
		return nil, err
	}

	if _, err = toolkit.JSONArg(args, "paths", "paths", &wf.Paths); err != nil {
		return nil, err
	}

	if wf.Scope, err = workflowScope(args["scope"]); err != nil {
		return nil, err
	}

	entryNote := noteMultipleEntryPoints(wf.Steps, wf.Paths)

	res, err := autoService.DefaultWorkflow.Create(ctx, wf)
	if err != nil {
		return nil, toolkit.Errf("workflow create", err)
	}

	return workflowWriteResult(res, entryNote)
}

// update is read-modify-write, and has to be.
//
// The service compares handle, enabled, trace, keepSessions, runAs and ownedBy
// against the stored workflow field by field with no nil guard, so a struct
// carrying only the changed fields would clear every one of them — a partial
// update has to start from the stored workflow. Steps, paths, scope and meta
// are nil-guarded and so are genuinely optional, but the loaded workflow is
// also what carries 'updatedAt' back to the service's optimistic lock, and its
// Meta pointer is dereferenced unguarded on the way in.
func (h *workflowHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	existing, err := h.resolve(ctx, args)
	if err != nil {
		return nil, err
	}

	upd := existing.Clone()
	if upd.Meta == nil {
		upd.Meta = &autoTypes.WorkflowMeta{}
	}

	if raw, given := args["handle"]; given && raw != nil {
		if upd.Handle, err = toolkit.Ref(args, "handle"); err != nil {
			return nil, err
		}
	}

	if raw, given := args["name"]; given && raw != nil {
		name := toolkit.Str(args, "name")
		if name == "" {
			return nil, fmt.Errorf("name cannot be cleared: a workflow with no name is refused — pass a name, or omit it to leave it unchanged")
		}
		upd.Meta.Name = name
	}

	if raw, given := args["description"]; given && raw != nil {
		upd.Meta.Description = toolkit.Str(args, "description")
	}

	if _, given := args["enabled"]; given {
		upd.Enabled = toolkit.Bool(args, "enabled")
	}

	if _, given := args["trace"]; given {
		upd.Trace = toolkit.Bool(args, "trace")
	}

	if _, given := args["subWorkflow"]; given {
		upd.Meta.SubWorkflow = toolkit.Bool(args, "subWorkflow")
	}

	if raw, given := args["runAs"]; given && raw != nil {
		if upd.RunAs, err = toolkit.ID(args, "runAs"); err != nil {
			return nil, err
		}
	}

	// Collections replace (§8.2). The service reads a nil Steps or Paths as
	// "unchanged", so an explicitly empty array has to survive as an empty set
	// rather than as nothing — that is what makes [] destructive and why both
	// descriptions say so.
	if given, err := toolkit.JSONArg(args, "steps", "steps", &upd.Steps); err != nil {
		return nil, err
	} else if given && upd.Steps == nil {
		upd.Steps = autoTypes.WorkflowStepSet{}
	}

	if given, err := toolkit.JSONArg(args, "paths", "paths", &upd.Paths); err != nil {
		return nil, err
	} else if given && upd.Paths == nil {
		upd.Paths = autoTypes.WorkflowPathSet{}
	}

	if raw, given := args["scope"]; given && raw != nil {
		if upd.Scope, err = workflowScope(raw); err != nil {
			return nil, err
		}
		if upd.Scope == nil {
			if upd.Scope, err = expr.NewVars(map[string]interface{}{}); err != nil {
				return nil, toolkit.Errf("workflow scope build", err)
			}
		}
	}

	entryNote := noteMultipleEntryPoints(upd.Steps, upd.Paths)

	res, err := autoService.DefaultWorkflow.Update(ctx, upd)
	if err != nil {
		return nil, toolkit.Errf("workflow update", err)
	}

	return workflowWriteResult(res, entryNote)
}

// undelete takes an ID rather than the ID-or-handle 'workflow' reference the
// other tools accept: resolve falls back to a handle search, and WorkflowFilter
// defaults Deleted to StateExcluded, so a deleted workflow is unreachable by
// handle. The ID path loads deleted rows, so it goes straight to the service.
func (h *workflowHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// Accepts the same ID-or-handle reference every other tool in this family
	// takes. Undelete still takes an ID only, and has to: resolve searches, and
	// a deleted workflow is not in the search.
	wf, err := h.resolve(ctx, args)
	if err != nil {
		return nil, err
	}

	if err := autoService.DefaultWorkflow.DeleteByID(ctx, wf.ID); err != nil {
		return nil, toolkit.Errf("workflow delete", err)
	}
	return toolkit.TextResult("workflow %d deleted", wf.ID), nil
}

func (h *workflowHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	wfID, err := toolkit.ReqID(args, "workflowID")
	if err != nil {
		return nil, err
	}

	if err = autoService.DefaultWorkflow.UndeleteByID(ctx, wfID); err != nil {
		return nil, toolkit.Errf("workflow undelete", err)
	}
	return toolkit.TextResult("workflow %d restored", wfID), nil
}

func (h *workflowHandler) exec(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	wf, err := h.resolve(ctx, args)
	if err != nil {
		return nil, err
	}

	params := autoTypes.WorkflowExecParams{}

	if inputMap, err := parseInput(args["input"]); err != nil {
		return nil, err
	} else if inputMap != nil {
		vars, err := expr.NewVars(inputMap)
		if err != nil {
			return nil, toolkit.Errf("workflow input conversion", err)
		}
		params.Input = vars
	}

	result, _, _, err := autoService.DefaultWorkflow.Exec(ctx, wf.ID, params)
	if err != nil {
		return nil, toolkit.Errf("workflow execution", err)
	}

	return toolkit.JSONResult(result)
}

// search runs a workflow search with the house error wrapping.
// search lists workflows. includeDisabled is a parameter rather than a default
// because WorkflowFilter.Disabled defaults to StateExcluded everywhere in the
// product, and a lookup that quietly disagreed with the rest of it would be its
// own kind of surprise.
func (h *workflowHandler) search(ctx context.Context, query string, paging filter.Paging, includeDisabled bool) (autoTypes.WorkflowSet, autoTypes.WorkflowFilter, error) {
	ff := autoTypes.WorkflowFilter{
		Query:  query,
		Paging: paging,
	}

	if includeDisabled {
		ff.Disabled = filter.StateInclusive
	}

	set, f, err := autoService.DefaultWorkflow.Search(ctx, ff)
	if err != nil {
		return nil, f, toolkit.Errf("workflow list", err)
	}
	return set, f, nil
}

// resolve turns the 'workflow' argument into a workflow.
//
// The argument is a union of ID and handle, so it cannot go through toolkit.ID —
// but the ID half still crosses the boundary as a string, and refuseNumericRef
// enforces that.
func (h *workflowHandler) resolve(ctx context.Context, args map[string]any) (*autoTypes.Workflow, error) {
	if err := h.refuseNumericRef(args, "workflow"); err != nil {
		return nil, err
	}

	ref, err := toolkit.ReqStr(args, "workflow")
	if err != nil {
		return nil, err
	}

	if id, err := strconv.ParseUint(ref, 10, 64); err == nil {
		wf, err := autoService.DefaultWorkflow.FindByID(ctx, id)
		if err != nil {
			return nil, toolkit.Errf("workflow lookup", err)
		}
		return wf, nil
	}

	// Disabled is inclusive here, unlike in the listing: the caller named one
	// specific workflow, and leaving it out is what made disabling by handle a
	// one-way door — nothing holding only the handle could reach it to enable
	// it again.
	set, _, err := autoService.DefaultWorkflow.Search(ctx, autoTypes.WorkflowFilter{
		Handle:   ref,
		Disabled: filter.StateInclusive,
	})
	if err != nil {
		return nil, toolkit.Errf("workflow lookup", err)
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("workflow %q not found", ref)
	}
	return set[0], nil
}

// workflowWriteResult returns a stored workflow together with the issues the
// service reported while converting it.
//
// Issues do not block the write: automation/service/workflow.go stores the
// workflow and skips trigger registration when there are any, so a graph that
// will never run comes back as a success. Returning them verbatim, next to a
// warning, is the whole feedback loop for a caller that cannot read the server
// log.
func workflowWriteResult(wf *autoTypes.Workflow, notes ...string) (*mcp.CallToolResult, error) {
	out := workflowWritten{Workflow: wf, Issues: autoTypes.WorkflowIssueSet{}}

	for _, n := range notes {
		if n != "" {
			out.Warning = n
		}
	}

	if len(wf.Issues) > 0 {
		out.Issues = wf.Issues
		out.Warning = fmt.Sprintf(
			"stored with %d issue(s): this workflow will NOT run and its triggers are not registered. "+
				"Fix every issue listed here and call automation_workflow_update again.",
			len(wf.Issues),
		)
	}

	return toolkit.JSONResult(out)
}

// workflowScope turns the 'scope' argument into the variables every run starts
// with. Absent or empty yields nil, which both callers treat as "not given".
func workflowScope(raw any) (*expr.Vars, error) {
	scopeMap, err := parseInput(raw)
	if err != nil {
		return nil, err
	}
	if scopeMap == nil {
		return nil, nil
	}

	vars, err := expr.NewVars(scopeMap)
	if err != nil {
		return nil, toolkit.Errf("workflow scope build", err)
	}
	return vars, nil
}

// noteMultipleEntryPoints reports a graph with more than one step that nothing
// points at. It warns rather than refuses, because such a graph is legal.
//
// session.Start only resolves the entry structurally when the exec carries no
// StepID. An event-driven run always names its step, and a manual run now
// inherits the first trigger's step, so a multi-entry workflow with at least one
// trigger runs either way. What still fails is the case this tool produces by
// default: a workflow with no triggers at all, run by hand, where there is no
// step to inherit and session.Start finds several parentless candidates.
//
// The editor can build these graphs, so refusing to author one would make the
// tool weaker than the UI for no safety gain. Saying plainly what will and will
// not run is the useful thing.
//
// Visual steps are excluded because the converter drops them before building the
// graph, so they are never entry points however they are connected.
func noteMultipleEntryPoints(steps autoTypes.WorkflowStepSet, paths autoTypes.WorkflowPathSet) string {
	hasParent := make(map[uint64]bool, len(paths))
	for _, p := range paths {
		if p != nil {
			hasParent[p.ChildID] = true
		}
	}

	entries := make([]string, 0, 2)
	for _, s := range steps {
		if s == nil || s.Kind == autoTypes.WorkflowStepKindVisual {
			continue
		}
		if !hasParent[s.ID] {
			entries = append(entries, strconv.FormatUint(s.ID, 10))
		}
	}

	if len(entries) > 1 {
		return fmt.Sprintf(
			"this workflow has %d starting steps (stepIDs %s). That is legal, but running it with "+
				"automation_workflow_exec without a stepID will fail while it has no trigger, because there is "+
				"no single entry to resolve. Attach a trigger with automation_trigger_create, pass a stepID to "+
				"exec, or give every step but one an inbound path.",
			len(entries), strings.Join(entries, ", "),
		)
	}

	return ""
}

// refuseNumericRef rejects an ID-or-handle reference passed as a JSON number.
//
// Human IDs exceed JavaScript's safe integer range, so by the time a numeric
// argument reaches here it may already have been truncated; coercing it would
// launder a corrupted value into the service layer. This is toolkit.ID's rule
// applied to a param that toolkit.ID cannot parse, because a handle is not a
// number.
func (h *workflowHandler) refuseNumericRef(args map[string]any, key string) error {
	raw, ok := args[key]
	if !ok || raw == nil {
		return nil
	}
	if _, isStr := raw.(string); !isStr {
		return fmt.Errorf("%s must be a string to avoid precision loss, got %T", key, raw)
	}
	return nil
}
