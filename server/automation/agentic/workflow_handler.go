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
	"github.com/crusttech/human/server/system/agentic/toolkit"
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

	set, f, err := h.search(ctx, query, paging)
	if err != nil {
		return nil, err
	}

	// Handles are slugs, so a model searching for what it sees in the UI ("My
	// Workflow") matches nothing. Retry once with the slugified query — but not
	// mid-pagination, where switching the query would silently change what the
	// caller is paging through.
	if len(set) == 0 && query != "" && page.Cursor == "" {
		if slugQuery := strings.ReplaceAll(strings.ToLower(query), " ", "_"); slugQuery != query {
			if set, f, err = h.search(ctx, slugQuery, paging); err != nil {
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
func (h *workflowHandler) search(ctx context.Context, query string, paging filter.Paging) (autoTypes.WorkflowSet, autoTypes.WorkflowFilter, error) {
	set, f, err := autoService.DefaultWorkflow.Search(ctx, autoTypes.WorkflowFilter{
		Query:  query,
		Paging: paging,
	})
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

	set, _, err := autoService.DefaultWorkflow.Search(ctx, autoTypes.WorkflowFilter{Handle: ref})
	if err != nil {
		return nil, toolkit.Errf("workflow lookup", err)
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("workflow %q not found", ref)
	}
	return set[0], nil
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
