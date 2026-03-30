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
	"github.com/mark3labs/mcp-go/server"
)

type (
	toolRegistrar interface {
		RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc)
	}

	taqHandler struct {
		reg toolRegistrar
	}
)

func TAQHandler(reg toolRegistrar) *taqHandler {
	h := &taqHandler{reg: reg}
	h.register()
	return h
}

func (h *taqHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_lookup",
			mcp.WithDescription("List all TAQs or look up a specific one by ID or handle. Omit 'taq' to list all."),
			mcp.WithString("taq", mcp.Description("TAQ ID or handle. Omit to list all.")),
			mcp.WithString("query", mcp.Description("Search query to filter TAQs")),
		),
		"Lookup TAQ",
		h.lookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_exec",
			mcp.WithDescription("Execute a TAQ by ID or handle and wait for the result"),
			mcp.WithString("taq", mcp.Required(), mcp.Description("TAQ ID or handle")),
			mcp.WithString("entryPoint", mcp.Description("Entry point (trigger handle) to invoke")),
			mcp.WithString("input", mcp.Description("JSON object of input variables")),
		),
		"Execute TAQ",
		h.exec,
	)
	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_executions",
			mcp.WithDescription("List executions for a TAQ"),
			mcp.WithString("taq", mcp.Required(), mcp.Description("TAQ ID or handle")),
		),
		"List TAQ executions",
		h.executions,
	)
	h.reg.RegisterTool(
		mcp.NewTool("automation_taq_execution_trace",
			mcp.WithDescription("Get the execution trace for a specific TAQ execution"),
			mcp.WithString("taq", mcp.Required(), mcp.Description("TAQ ID or handle")),
			mcp.WithString("executionID", mcp.Required(), mcp.Description("Execution ID")),
		),
		"Get TAQ execution trace",
		h.executionTrace,
	)
}

func (h *taqHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]interface{})

	taqRef, _ := args["taq"].(string)
	query, _ := args["query"].(string)

	if taqRef == "" {
		set, _, err := autoService.DefaultNgAutomation.Search(ctx, autoTypes.NgAutomationFilter{Query: query})
		if err != nil {
			return nil, fmt.Errorf("TAQ list failed: %w", err)
		}

		if len(set) == 0 && query != "" {
			slugQuery := strings.ReplaceAll(strings.ToLower(query), " ", "_")
			if slugQuery != query {
				set, _, err = autoService.DefaultNgAutomation.Search(ctx, autoTypes.NgAutomationFilter{Query: slugQuery})
				if err != nil {
					return nil, fmt.Errorf("TAQ list failed: %w", err)
				}
			}
		}

		out, err := json.Marshal(set)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal TAQs: %w", err)
		}
		return mcp.NewToolResultText(string(out)), nil
	}

	taq, err := h.resolve(ctx, taqRef)
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(taq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TAQ: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *taqHandler) exec(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	taq, err := h.resolve(ctx, args["taq"].(string))
	if err != nil {
		return nil, err
	}

	params := autoTypes.NgAutomationExecParams{}

	if ep, ok := args["entryPoint"].(string); ok && ep != "" {
		params.EntryPoint = ep
	}

	if inputStr, ok := args["input"].(string); ok && inputStr != "" {
		var inputMap map[string]interface{}
		if err := json.Unmarshal([]byte(inputStr), &inputMap); err != nil {
			return nil, fmt.Errorf("invalid input JSON: %w", err)
		}
		vars, err := expr.NewVars(inputMap)
		if err != nil {
			return nil, fmt.Errorf("failed to build input vars: %w", err)
		}
		params.Input = vars
	}

	result, err := autoService.DefaultNgAutomation.ExecAndWait(ctx, taq.ID, params)
	if err != nil {
		return nil, fmt.Errorf("TAQ execution failed: %w", err)
	}
	out, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *taqHandler) executions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	taq, err := h.resolve(ctx, args["taq"].(string))
	if err != nil {
		return nil, err
	}

	results, err := autoService.DefaultNgAutomation.GetExecutions(ctx, taq.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions: %w", err)
	}
	out, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal executions: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *taqHandler) executionTrace(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	taq, err := h.resolve(ctx, args["taq"].(string))
	if err != nil {
		return nil, err
	}

	execID, err := strconv.ParseUint(args["executionID"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid executionID: %w", err)
	}

	trace, err := autoService.DefaultNgAutomation.GetExecutionTrace(ctx, taq.ID, execID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution trace: %w", err)
	}
	out, err := json.Marshal(trace)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trace: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *taqHandler) resolve(ctx context.Context, refStr string) (*autoTypes.NgAutomation, error) {
	if refStr == "" {
		return nil, fmt.Errorf("taq identifier required")
	}

	if id, err := strconv.ParseUint(refStr, 10, 64); err == nil {
		return autoService.DefaultNgAutomation.LookupByID(ctx, id)
	}

	set, _, err := autoService.DefaultNgAutomation.Search(ctx, autoTypes.NgAutomationFilter{Handle: refStr})
	if err != nil {
		return nil, fmt.Errorf("TAQ lookup failed: %w", err)
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("TAQ %q not found", refStr)
	}
	return set[0], nil
}
