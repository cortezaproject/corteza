package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	cmpService "github.com/cortezaproject/corteza/server/compose/service"
	cmpTypes "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/mark3labs/mcp-go/mcp"
)

type recordHandler struct {
	reg *Registry
}

func RecordHandler(reg *Registry) *recordHandler {
	h := &recordHandler{
		reg: reg,
	}
	h.register()
	return h
}
func (h *recordHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_lookup",
			mcp.WithDescription("Look up a record by ID"),
			mcp.WithString("namespaceID",mcp.Required(),mcp.Description("Namespace ID")),
			mcp.WithString("moduleID",mcp.Required(),mcp.Description("Module ID")),
			mcp.WithString("recordID",mcp.Required(),mcp.Description("Record ID")),
		),
		h.lookup,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_create",
			mcp.WithDescription("Create a new record"),
			mcp.WithString("namespaceID",mcp.Required(),mcp.Description("Namespace ID")),
			mcp.WithString("moduleID",mcp.Required(),mcp.Description("Module ID")),
			mcp.WithString("values",mcp.Required(),mcp.Description("JSON object of field name-value pairs")),
		),
		h.create,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_update",
			mcp.WithDescription("Update an existing record"),
			mcp.WithString("namespaceID",mcp.Required(),mcp.Description("Namespace ID")),
			mcp.WithString("moduleID",mcp.Required(),mcp.Description("Module ID")),
			mcp.WithString("recordID",mcp.Required(),mcp.Description("Record ID")),
			mcp.WithString("values",mcp.Required(),mcp.Description("JSON object of field name-value pairs")),
		),
		h.update,
	)
	h.reg.RegisterTool(
		mcp.NewTool("compose_record_delete",
			mcp.WithDescription("Delete a record by ID"),
			mcp.WithString("namespaceID",mcp.Required(),mcp.Description("Namespace ID")),
			mcp.WithString("moduleID",mcp.Required(),mcp.Description("Module ID")),
			mcp.WithString("recordID",mcp.Required(),mcp.Description("Record ID")),
		),
		h.del,
	)
}

func (h *recordHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}
	nsID, err := strconv.ParseUint(args["namespaceID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid namespaceID : %w", err)
	}

	modID, err := strconv.ParseUint(args["moduleID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid moduleID : %w", err)
	}

	recID, err := strconv.ParseUint(args["recordID"].(string),10,64)
	if err != nil {
		return nil,fmt.Errorf("invalid recordID : %w", err)
	}

	rec, _, err := cmpService.DefaultRecord.FindByID(ctx, nsID, modID, recID)
	if err != nil {
		return nil, fmt.Errorf("record lookup failed: %w", err)
	}
	out, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *recordHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	nsID, err := strconv.ParseUint(args["namespaceID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid namespaceID: %w", err)
	}

	modID, err := strconv.ParseUint(args["moduleID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid moduleID: %w", err)
	}

	recID, err := strconv.ParseUint(args["recordID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid recordID: %w", err)
	}

	var valuesMap map[string]string
	if err := json.Unmarshal([]byte(args["values"].(string)), &valuesMap); err != nil {
		return nil, fmt.Errorf("invalid values JSON: %w", err)
	}

	rec := &cmpTypes.Record{
		ID:          recID,
		NamespaceID: nsID,
		ModuleID:    modID,
	}

	for name, value := range valuesMap {
		rec.Values = append(rec.Values, &cmpTypes.RecordValue{
			Name:  name,
			Value: value,
		})
	}

	rec, _, err = cmpService.DefaultRecord.Update(ctx, rec)
	if err != nil {
		return nil, fmt.Errorf("record update failed: %w", err)
	}

	out, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (h *recordHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	nsID, err := strconv.ParseUint(args["namespaceID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid namespaceID: %w", err)
	}

	modID, err := strconv.ParseUint(args["moduleID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid moduleID: %w", err)
	}

	recID, err := strconv.ParseUint(args["recordID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid recordID: %w", err)
	}

	err = cmpService.DefaultRecord.DeleteByID(ctx, nsID, modID, recID)
	if err != nil {
		return nil, fmt.Errorf("record delete failed: %w", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("record %d deleted", recID)), nil
}

func (h *recordHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	nsID, err := strconv.ParseUint(args["namespaceID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid namespaceID: %w", err)
	}

	modID, err := strconv.ParseUint(args["moduleID"].(string),10,64)
	if err != nil {
		return nil, fmt.Errorf("invalid moduleID: %w", err)
	}

	var valuesMap map[string]string
	if err := json.Unmarshal([]byte(args["values"].(string)), &valuesMap); err != nil {
		return nil, fmt.Errorf("invalid values JSON: %w", err)
	}

	rec := &cmpTypes.Record{
		NamespaceID: nsID,
		ModuleID:    modID,
	}

	for name, value := range valuesMap {
		rec.Values = append(rec.Values, &cmpTypes.RecordValue{
			Name:  name,
			Value: value,
		})
	}

	rec, _, err = cmpService.DefaultRecord.Create(ctx, rec)
	if err != nil {
		return nil, fmt.Errorf("record creation failed: %w", err)
	}

	out, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record: %w", err)
	}
	return mcp.NewToolResultText(string(out)), nil
}