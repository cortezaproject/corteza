package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	autoService "github.com/cortezaproject/corteza/server/automation/service"
	autoTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/mark3labs/mcp-go/mcp"
)

type (
	workflowHandler struct {
		reg toolRegistrar
	}
)

func WorkflowHandler(reg toolRegistrar) *workflowHandler {
	h := &workflowHandler{reg: reg}
	h.register()
	return h
}

func (h *workflowHandler) register() {
	h.reg.RegisterHiddenTool(
		mcp.NewTool("automation_workflow_lookup",
			mcp.WithDescription("List all workflows or look up a specific one by ID or handle. Omit 'workflow' to list all."),
			mcp.WithString("workflow", mcp.Description("Workflow ID or handle. Omit to list all.")),
			mcp.WithString("query", mcp.Description("Search query to filter workflows by handle")),
		),
		"Lookup workflow",
		h.lookup,
	)
	h.reg.RegisterHiddenTool(
		mcp.NewTool("automation_workflow_exec",
			mcp.WithDescription("Execute a workflow by ID or handle and wait for the result"),
			mcp.WithString("workflow", mcp.Required(), mcp.Description("Workflow ID or handle")),
			mcp.WithString("input", mcp.Description("JSON object of input variables")),
		),
		"Execute workflow",
		h.exec,
	)
}

func (h *workflowHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]interface{})
	wfRef, _ := args["workflow"].(string)
	query, _ := args["query"].(string)

	if wfRef == "" {
		set, _, err := autoService.DefaultWorkflow.Search(ctx, autoTypes.WorkflowFilter{Query: query})
		if err != nil {
			return nil, fmt.Errorf("workflow list failed: %w", err)
		}

		if len(set) == 0 && query != "" {
			slugQuery := strings.ReplaceAll(strings.ToLower(query), " ", "_")
			if slugQuery != query {
				set, _, err = autoService.DefaultWorkflow.Search(ctx, autoTypes.WorkflowFilter{Query: slugQuery})
				if err != nil {
					return nil, fmt.Errorf("workflow list failed: %w", err)
				}
			}
		}

		out, err := json.Marshal(set)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal workflows: %w", err)
		}
		return mcp.NewToolResultText(string(out)), nil
	}

	wf, err := h.resolve(ctx, wfRef)
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(wf)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *workflowHandler) exec(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	wf, err := h.resolve(ctx, args["workflow"])
	if err != nil {
		return nil, err
	}

	params := autoTypes.WorkflowExecParams{}

	if inputMap, err := parseInput(args["input"]); err != nil {
		return nil, err
	} else if inputMap != nil {
		vars, err := expr.NewVars(inputMap)
		if err != nil {
			return nil, fmt.Errorf("failed to build input vars: %w", err)
		}
		params.Input = vars
	}

	result, _, _, err := autoService.DefaultWorkflow.Exec(ctx, wf.ID, params)
	if err != nil {
		return nil, fmt.Errorf("workflow execution failed: %w", err)
	}

	out, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *workflowHandler) resolve(ctx context.Context, ref interface{}) (*autoTypes.Workflow, error) {
	refStr, ok := ref.(string)
	if !ok || refStr == "" {
		return nil, fmt.Errorf("workflow identifier required")
	}

	if id, err := strconv.ParseUint(refStr, 10, 64); err == nil {
		return autoService.DefaultWorkflow.LookupByID(ctx, id)
	}

	set, _, err := autoService.DefaultWorkflow.Search(ctx, autoTypes.WorkflowFilter{Handle: refStr})
	if err != nil {
		return nil, fmt.Errorf("workflow lookup failed: %w", err)
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("workflow %q not found", refStr)
	}
	return set[0], nil
}
