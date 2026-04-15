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
		RegisterHiddenTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc)
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
	h.reg.RegisterHiddenTool(
		mcp.NewTool("automation_taq_lookup",
			mcp.WithDescription("List all TAQs or look up a specific one by ID or handle. Omit 'taq' to list all."),
			mcp.WithString("taq", mcp.Description("TAQ ID as string (to prevent precision loss) or handle. Omit to list all.")),
			mcp.WithString("query", mcp.Description("Search query to filter TAQs")),
		),
		"Lookup TAQ",
		h.lookup,
	)
	h.reg.RegisterHiddenTool(
		mcp.NewTool("automation_taq_executions",
			mcp.WithDescription("List executions for a TAQ"),
			mcp.WithString("taq", mcp.Required(), mcp.Description("TAQ ID as string (to prevent precision loss) or handle")),
		),
		"List TAQ executions",
		h.executions,
	)
	h.reg.RegisterHiddenTool(
		mcp.NewTool("automation_taq_execution_trace",
			mcp.WithDescription("Get the execution trace for a specific TAQ execution"),
			mcp.WithString("taq", mcp.Required(), mcp.Description("TAQ ID as string (to prevent precision loss) or handle")),
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

	taqRef, _ := args["taq"].(string)
	taq, err := h.resolve(ctx, taqRef)
	if err != nil {
		return nil, err
	}

	params := autoTypes.NgAutomationExecParams{}

	if ep, ok := args["entryPoint"].(string); ok && ep != "" {
		params.EntryPoint = ep
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

	taqRef, _ := args["taq"].(string)
	taq, err := h.resolve(ctx, taqRef)
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

	taqRef, _ := args["taq"].(string)
	taq, err := h.resolve(ctx, taqRef)
	if err != nil {
		return nil, err
	}

	execIDStr, _ := args["executionID"].(string)
	execID, err := strconv.ParseUint(execIDStr, 10, 64)
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
