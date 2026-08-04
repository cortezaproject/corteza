package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	autoService "github.com/crusttech/human/server/automation/service"
	autoTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/filter"
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type (
	// toolRegistrar is shared by every handler in this package; it lives here
	// rather than in a file of its own because taq is the package's first
	// handler.
	toolRegistrar interface {
		RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc, opts ...hmcp.RegisterOption)
	}

	taqHandler struct {
		reg toolRegistrar
	}
)

// Declarations for these handlers are in taq_tools.go, in the same order.
func TAQHandler(reg toolRegistrar) *taqHandler {
	h := &taqHandler{reg: reg}
	h.register()
	return h
}

// taqItem is the compact projection returned by list mode. A single-item
// lookup returns the service type unchanged; a list must not, because every
// TAQ carries its full trigger, step and path graph and one graph per row
// would flood the caller's context.
type taqItem struct {
	AutomationID string `json:"automationID"`
	Handle       string `json:"handle"`
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
}

func (h *taqHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	// Reference param is never Required, so an empty ref means "list".
	if taqRef := toolkit.Str(args, "taq"); taqRef != "" {
		taq, err := h.resolve(ctx, taqRef)
		if err != nil {
			return nil, err
		}
		return toolkit.JSONResult(taq)
	}

	query := toolkit.Str(args, "query")

	f := autoTypes.NgAutomationFilter{Query: query}
	if toolkit.Bool(args, "includeDisabled") {
		f.Disabled = filter.StateInclusive
	}

	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := autoService.DefaultNgAutomation.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("TAQ list", err)
	}

	// The query only matches the handle, and handles are slugs, so a
	// human-shaped query ("my taq") finds nothing where the slug would.
	if len(set) == 0 && query != "" {
		if slugQuery := strings.ReplaceAll(strings.ToLower(query), " ", "_"); slugQuery != query {
			f.Query = slugQuery
			if set, out, err = autoService.DefaultNgAutomation.Search(ctx, f); err != nil {
				return nil, toolkit.Errf("TAQ list", err)
			}
		}
	}

	items := make([]taqItem, 0, len(set))
	for _, taq := range set {
		item := taqItem{
			AutomationID: strconv.FormatUint(taq.ID, 10),
			Handle:       taq.Handle,
			Enabled:      taq.Enabled,
		}
		if taq.Meta != nil {
			item.Name = taq.Meta.Short
		}
		items = append(items, item)
	}

	// The cursor value, not its String(): String is a debug rendering, while
	// parseCursor expects the base64 form MarshalJSON emits. Nil marshals to
	// null, so the last page is safe.
	return toolkit.JSONResult(map[string]any{
		"taqs":           items,
		"nextPageCursor": out.NextPage,
	})
}

// taqWriteResult is what create and update return.
//
// The service stores a TAQ whose graph does not validate — at HTTP success,
// with the defects reported in Issues — and skips registering it with the
// runtime. Echoing Issues verbatim is therefore the entire feedback loop, and a
// result that only said "created" would be actively misleading. Runnable states
// out loud what Issues implies, because the consequence is not obvious from the
// issue list: any issue at all, of any severity, means exec cannot find the
// executable.
//
// Triggers and steps come back trimmed rather than whole: the full graph is
// what automation_taq_lookup is for, and echoing it here would double the size
// of every write for nothing. The IDs are the part the caller cannot get any
// other way — a trigger ID minted by the service on create is needed to author
// a path to it later.
type taqWriteResult struct {
	AutomationID string                         `json:"automationID"`
	Handle       string                         `json:"handle"`
	Name         string                         `json:"name"`
	Enabled      bool                           `json:"enabled"`
	Runnable     bool                           `json:"runnable"`
	Triggers     []taqTriggerItem               `json:"triggers"`
	Steps        []taqStepItem                  `json:"steps"`
	Issues       autoTypes.NgAutomationIssueSet `json:"issues,omitempty"`
	DroppedSteps []string                       `json:"droppedSteps,omitempty"`
	Note         string                         `json:"note,omitempty"`
}

type taqTriggerItem struct {
	TriggerID    string `json:"triggerID"`
	Handle       string `json:"handle"`
	ResourceType string `json:"resourceType"`
	EventType    string `json:"eventType"`
	Enabled      bool   `json:"enabled"`
}

type taqStepItem struct {
	StepID string `json:"stepID"`
	Handle string `json:"handle"`
	Kind   string `json:"kind"`
	Ref    string `json:"ref,omitempty"`
}

func (h *taqHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	taqHandle, err := toolkit.ReqStr(args, "handle")
	if err != nil {
		return nil, err
	}

	name, err := toolkit.ReqStr(args, "name")
	if err != nil {
		return nil, err
	}

	// Meta is where the name lives — Meta.Short, not a Name field — and the
	// service rejects a nil Meta or an empty Short outright.
	newTaq := &autoTypes.NgAutomation{
		Handle: taqHandle,
		Meta: &autoTypes.NgAutomationMeta{
			Short:       name,
			Description: toolkit.Str(args, "description"),
		},
		// Enabled gates both trigger registration and exec, so a TAQ created
		// without it could not even be tested. Default to live and say so in
		// the description rather than hand back something inert.
		Enabled: true,
	}
	if v, present := toolkit.OptBool(args, "enabled"); present {
		newTaq.Enabled = v
	}

	if _, err = toolkit.JSONArg(args, "triggers", "triggers", &newTaq.Triggers); err != nil {
		return nil, err
	}
	if _, err = toolkit.JSONArg(args, "steps", "steps", &newTaq.Steps); err != nil {
		return nil, err
	}
	if _, err = toolkit.JSONArg(args, "paths", "paths", &newTaq.Paths); err != nil {
		return nil, err
	}

	if err = prepSteps(newTaq.Steps); err != nil {
		return nil, err
	}
	submitted := stepIDs(newTaq.Steps)

	taq, err := autoService.DefaultNgAutomation.Create(ctx, newTaq)
	if err != nil {
		return nil, toolkit.Errf("TAQ create", err)
	}

	return toolkit.JSONResult(taqWriteResultOf(taq, submitted))
}

func (h *taqHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	taqRef, err := toolkit.ReqRef(args, "taq")
	if err != nil {
		return nil, err
	}

	current, err := h.resolve(ctx, taqRef)
	if err != nil {
		return nil, err
	}

	// Read-modify-write, because the service's update is only half sparse:
	// handleUpdate nil-guards Meta, Scope, Triggers, Steps, Paths and Labels but
	// assigns Handle, Enabled, RunAs and OwnedBy straight across from whatever
	// it is given. A payload carrying only the fields the caller named would
	// clear the owner and the run-as user and be rejected for an invalid empty
	// handle. UpdatedAt is echoed for the same reason it is loaded: it is what
	// isStale compares against.
	upd := &autoTypes.NgAutomation{
		ID:        current.ID,
		Handle:    current.Handle,
		Enabled:   current.Enabled,
		RunAs:     current.RunAs,
		OwnedBy:   current.OwnedBy,
		Labels:    current.Labels,
		Scope:     current.Scope,
		UpdatedAt: current.UpdatedAt,
	}

	// A copy of the stored meta, not the pointer: Update rejects a nil Meta or
	// an empty Short outright, so the current values have to be carried over,
	// and writing through the loaded pointer would mutate the object
	// handleUpdate then compares against to decide whether anything changed.
	upd.Meta = &autoTypes.NgAutomationMeta{}
	if current.Meta != nil {
		meta := *current.Meta
		upd.Meta = &meta
	}

	if v, present := toolkit.OptStr(args, "handle"); present {
		upd.Handle = v
	}
	if v, present := toolkit.OptStr(args, "name"); present {
		if v == "" {
			return nil, fmt.Errorf("name cannot be cleared: a TAQ must have a name, so omit 'name' to keep the current one")
		}
		upd.Meta.Short = v
	}
	if v, present := toolkit.OptStr(args, "description"); present {
		upd.Meta.Description = v
	}
	if v, present := toolkit.OptBool(args, "enabled"); present {
		upd.Enabled = v
	}

	// Absent leaves a collection alone, [] clears it (CONVENTIONS.md §8.2).
	// That is the service's own rule — handleUpdate branches on a nil slice —
	// so the two only agree if an omitted param stays nil here.
	if _, err = toolkit.JSONArg(args, "triggers", "triggers", &upd.Triggers); err != nil {
		return nil, err
	}
	if _, err = toolkit.JSONArg(args, "steps", "steps", &upd.Steps); err != nil {
		return nil, err
	}
	if _, err = toolkit.JSONArg(args, "paths", "paths", &upd.Paths); err != nil {
		return nil, err
	}

	if err = prepSteps(upd.Steps); err != nil {
		return nil, err
	}
	submitted := stepIDs(upd.Steps)

	taq, err := autoService.DefaultNgAutomation.Update(ctx, upd)
	if err != nil {
		return nil, toolkit.Errf("TAQ update", err)
	}

	return toolkit.JSONResult(taqWriteResultOf(taq, submitted))
}

// undelete takes an ID rather than the ID-or-handle 'taq' reference every other
// tool here accepts, because resolve cannot reach a deleted TAQ by handle:
// NgAutomationFilter defaults Deleted to StateExcluded, so the handle search
// would report "not found" for exactly the TAQs this tool exists to restore.
// The ID goes straight to the service, which loads deleted rows.
func (h *taqHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	taqID, err := toolkit.ReqID(args, "taqID")
	if err != nil {
		return nil, err
	}

	if err := autoService.DefaultNgAutomation.DeleteByID(ctx, taqID); err != nil {
		return nil, toolkit.Errf("TAQ delete", err)
	}
	return toolkit.TextResult("TAQ %d deleted", taqID), nil
}

func (h *taqHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	taqID, err := toolkit.ReqID(args, "taqID")
	if err != nil {
		return nil, err
	}

	if err = autoService.DefaultNgAutomation.UndeleteByID(ctx, taqID); err != nil {
		return nil, toolkit.Errf("TAQ undelete", err)
	}
	return toolkit.TextResult("TAQ %d restored", taqID), nil
}

func (h *taqHandler) exec(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	taqRef, err := toolkit.ReqStr(args, "taq")
	if err != nil {
		return nil, err
	}

	taq, err := h.resolve(ctx, taqRef)
	if err != nil {
		return nil, err
	}

	params := autoTypes.NgAutomationExecParams{
		EventType:    "onAgentic",
		ResourceType: "automation:trigger:agentic",
		EntryPoint:   toolkit.Str(args, "entryPoint"),
	}

	if inputMap, err := parseInput(args["input"]); err != nil {
		return nil, err
	} else if inputMap != nil {
		inputMap, err = h.validateAndFixInput(ctx, taq, params.EntryPoint, inputMap)
		if err != nil {
			return nil, err
		}

		vars, err := expr.NewVars(inputMap)
		if err != nil {
			return nil, toolkit.Errf("input vars build", err)
		}
		params.Input = vars
	}

	result, err := autoService.DefaultNgAutomation.ExecAndWait(ctx, taq.ID, params)
	if err != nil {
		return nil, toolkit.Errf("TAQ execution", err)
	}
	return toolkit.JSONResult(result)
}

func (h *taqHandler) executions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	taqRef, err := toolkit.ReqStr(args, "taq")
	if err != nil {
		return nil, err
	}

	taq, err := h.resolve(ctx, taqRef)
	if err != nil {
		return nil, err
	}

	results, err := autoService.DefaultNgAutomation.GetExecutions(ctx, taq.ID)
	if err != nil {
		return nil, toolkit.Errf("TAQ execution list", err)
	}
	return toolkit.JSONResult(results)
}

func (h *taqHandler) executionTrace(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	taqRef, err := toolkit.ReqStr(args, "taq")
	if err != nil {
		return nil, err
	}

	taq, err := h.resolve(ctx, taqRef)
	if err != nil {
		return nil, err
	}

	execID, err := toolkit.ReqID(args, "executionID")
	if err != nil {
		return nil, err
	}

	trace, err := autoService.DefaultNgAutomation.GetExecutionTrace(ctx, taq.ID, execID, 0)
	if err != nil {
		return nil, toolkit.Errf("TAQ execution trace", err)
	}
	return toolkit.JSONResult(trace)
}

func parseInput(raw interface{}) (map[string]interface{}, error) {
	if raw == nil {
		return nil, nil
	}
	switch v := raw.(type) {
	case string:
		if v == "" {
			return nil, nil
		}
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, fmt.Errorf("invalid input JSON: %w", err)
		}
		return m, nil
	case map[string]interface{}:
		return v, nil
	default:
		return nil, fmt.Errorf("invalid input: expected JSON string or object")
	}
}

// prepSteps checks a submitted step set before anything is written.
//
// Kind is a plain string on the generated type where Workflow has a typed enum,
// so nothing downstream stops a workflow kind being stored: it converts to
// nothing, the step drops out of the executable graph, and what the caller gets
// back is an issue about a step they thought was fine. Rejecting here turns
// that into a refused call naming the six kinds that do exist, which matters
// because a caller who has just authored a workflow will reach for
// "expressions".
//
// Results are dropped because the converter overwrites them from the function
// definition on the very object it then persists — anything a caller sends is
// silently replaced, so the tool does not offer the field and does not carry it.
func prepSteps(steps autoTypes.NgAutomationStepSet) error {
	for i, s := range steps {
		if s == nil {
			return fmt.Errorf("step %d is null", i)
		}

		if !autoTypes.IsValidNgAutomationStepKind(s.Kind) {
			return fmt.Errorf(
				"step %d has unsupported kind %q: a TAQ step kind is one of %s. TAQ and workflow step kinds are not interchangeable — \"expressions\" in particular exists only on a workflow",
				i, s.Kind, strings.Join(autoTypes.NgAutomationStepKinds, ", "),
			)
		}

		s.Results = nil
	}

	return nil
}

func stepIDs(steps autoTypes.NgAutomationStepSet) []string {
	out := make([]string, 0, len(steps))
	for _, s := range steps {
		if s != nil {
			out = append(out, strconv.FormatUint(s.ID, 10))
		}
	}
	return out
}

// taqWriteResultOf projects what the service returned, and diffs the steps that
// came back against the ones that went in.
//
// The diff is a second net rather than the primary one. A step that fails
// conversion is now reported as a step.invalid issue, which is the real signal;
// the service also echoes the submitted set today, so in practice this finds
// nothing. It stays because a step silently vanishing from a TAQ was invisible
// for a long time and cost real hours, and the check is three lines.
func taqWriteResultOf(taq *autoTypes.NgAutomation, submitted []string) taqWriteResult {
	out := taqWriteResult{
		AutomationID: strconv.FormatUint(taq.ID, 10),
		Handle:       taq.Handle,
		Enabled:      taq.Enabled,
		Issues:       taq.Issues,
		Runnable:     taq.Enabled && len(taq.Issues) == 0,
		Triggers:     make([]taqTriggerItem, 0, len(taq.Triggers)),
		Steps:        make([]taqStepItem, 0, len(taq.Steps)),
	}

	if taq.Meta != nil {
		out.Name = taq.Meta.Short
	}

	for _, t := range taq.Triggers {
		if t == nil {
			continue
		}
		out.Triggers = append(out.Triggers, taqTriggerItem{
			TriggerID:    strconv.FormatUint(t.ID, 10),
			Handle:       t.Handle,
			ResourceType: t.ResourceType,
			EventType:    t.EventType,
			Enabled:      t.Enabled,
		})
	}

	stored := make(map[string]bool, len(taq.Steps))
	for _, s := range taq.Steps {
		if s == nil {
			continue
		}
		id := strconv.FormatUint(s.ID, 10)
		stored[id] = true
		out.Steps = append(out.Steps, taqStepItem{
			StepID: id,
			Handle: s.Handle,
			Kind:   s.Kind,
			Ref:    s.Ref,
		})
	}

	for _, id := range submitted {
		if !stored[id] {
			out.DroppedSteps = append(out.DroppedSteps, id)
		}
	}

	var notes []string
	if len(out.DroppedSteps) > 0 {
		notes = append(notes, fmt.Sprintf(
			"steps %s were submitted but are not in the stored TAQ.",
			strings.Join(out.DroppedSteps, ", "),
		))
	}

	switch {
	case len(taq.Issues) > 0:
		notes = append(notes, "Stored, but NOT registered with the automation runtime: any issue at all, "+
			"of any severity, blocks registration, so automation_taq_exec will answer \"manager: executable "+
			"not found\" until every issue above is resolved. Fix them with automation_taq_update.")
	case !taq.Enabled:
		notes = append(notes, "Stored and valid, but disabled: its triggers are not listening and "+
			"automation_taq_exec refuses a disabled TAQ. Set enabled to true to make it live.")
	default:
		notes = append(notes, "Registered with the automation runtime. Storing is not proof of working: run it "+
			"with automation_taq_exec, then read automation_taq_execution_trace, because exec reports status "+
			"only. A \"completed\" status with no step frames in the trace — nothing, or only the trigger "+
			"frame — means no step ran, and a frame with populated \"args\" proves the arguments bound but "+
			"not that the step had its effect: verify any visible effect separately.")
	}

	out.Note = strings.Join(notes, " ")

	return out
}

// resolve accepts a TAQ ID or a handle, because both are stable identifiers a
// caller may hold. There is no FindByAny on this service.
func (h *taqHandler) resolve(ctx context.Context, refStr string) (*autoTypes.NgAutomation, error) {
	if refStr == "" {
		return nil, fmt.Errorf("taq identifier required")
	}

	if id, err := strconv.ParseUint(refStr, 10, 64); err == nil {
		return autoService.DefaultNgAutomation.FindByID(ctx, id)
	}

	// Disabled TAQs are addressable by handle, matching the ID path above:
	// NgAutomationFilter defaults Disabled to StateExcluded, so without this a
	// disabled TAQ resolved by ID but reported "not found" by handle — and
	// "not found" is exactly the wrong answer for someone working out why their
	// TAQ is not firing.
	set, _, err := autoService.DefaultNgAutomation.Search(ctx, autoTypes.NgAutomationFilter{
		Handle:   refStr,
		Disabled: filter.StateInclusive,
	})
	if err != nil {
		return nil, toolkit.Errf("TAQ lookup", err)
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("TAQ %q not found", refStr)
	}
	return set[0], nil
}

func (h *taqHandler) validateAndFixInput(ctx context.Context, taq *autoTypes.NgAutomation, entryPoint string, inputMap map[string]interface{}) (map[string]interface{}, error) {
	if inputMap == nil {
		return nil, nil
	}

	var trigger *autoTypes.NgAutomationTrigger
	if entryPoint != "" {
		for _, t := range taq.Triggers {
			if t.Handle == entryPoint {
				trigger = t
				break
			}
		}
	} else if len(taq.Triggers) > 0 {
		trigger = taq.Triggers[0]
	}

	if trigger == nil {
		return inputMap, nil
	}

	// Programmatically refactor casing mapping
	correctedMap := fixInputCasing(inputMap, trigger.InputSchema)

	// Validate required schema properties
	var missing []string
	for _, schemaParam := range trigger.InputSchema {
		if schemaParam.Required {
			if val, exists := correctedMap[schemaParam.Name]; !exists || val == nil {
				missing = append(missing, schemaParam.Name)
			}
		}
	}

	if len(missing) > 0 {
		var provided []string
		for k := range correctedMap {
			provided = append(provided, k)
		}
		return nil, fmt.Errorf("Schema validation failed. Missing required parameters: [%s]. (Provided keys: [%s]). Please re-evaluate the schema and provide the missing parameters.", strings.Join(missing, ", "), strings.Join(provided, ", "))
	}

	return correctedMap, nil
}

func fixInputCasing(inputMap map[string]interface{}, schema autoTypes.NgAutomationTriggerSchema) map[string]interface{} {
	correctedMap := make(map[string]interface{}, len(inputMap))
	for k, v := range inputMap {
		matched := false
		for _, schemaParam := range schema {
			if strings.EqualFold(k, schemaParam.Name) {
				correctedMap[schemaParam.Name] = v
				matched = true
				break
			}
		}
		if !matched {
			correctedMap[k] = v
		}
	}
	return correctedMap
}
