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
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/crusttech/human/server/system/agentic/toolkit"
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

// undelete takes an ID rather than the ID-or-handle 'taq' reference every other
// tool here accepts, because resolve cannot reach a deleted TAQ by handle:
// NgAutomationFilter defaults Deleted to StateExcluded, so the handle search
// would report "not found" for exactly the TAQs this tool exists to restore.
// The ID goes straight to the service, which loads deleted rows.
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
