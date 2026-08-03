package agentic

import (
	"context"
	"encoding/json"
	"fmt"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/crusttech/human/server/system/agentic/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type (
	toolRegistrar interface {
		RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc, opts ...hmcp.RegisterOption)
	}

	recordHandler struct {
		reg toolRegistrar
	}
)

// Declarations for these handlers are in record_tools.go, in the same order.
func RecordHandler(reg toolRegistrar) *recordHandler {
	h := &recordHandler{reg: reg}
	h.register()
	return h
}

func (h *recordHandler) resolveNsMod(ctx context.Context, args map[string]any) (nsID, modID uint64, err error) {
	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
	if err != nil {
		return 0, 0, toolkit.Errf("namespace lookup", err)
	}

	mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, args["module"])
	if err != nil {
		return 0, 0, toolkit.Errf("module lookup", err)
	}

	return ns.ID, mod.ID, nil
}

func (h *recordHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	if recID, err := toolkit.ID(args, "recordID"); err != nil {
		return nil, err
	} else if recID > 0 {
		rec, _, err := cmpService.DefaultRecord.FindByID(ctx, nsID, modID, recID)
		if err != nil {
			return nil, toolkit.Errf("record lookup", err)
		}
		return toolkit.JSONResult(rec)
	}

	f := cmpTypes.RecordFilter{
		NamespaceID: nsID,
		ModuleID:    modID,
		Query:       toolkit.Str(args, "filter"),
	}

	// Records are returned in full rather than as a slim projection: unlike a
	// module's fields or a chart's config, a record's values are the thing the
	// caller asked for. Paging and the toolkit result ceiling do the work that
	// projection does elsewhere — and without paging this call drains an entire
	// module, because the DAL loops while the limit is zero.
	page := toolkit.Page(args)
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := cmpService.DefaultRecord.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("record list", err)
	}

	return toolkit.JSONResult(map[string]any{
		"records":        set,
		"nextPageCursor": out.NextPage.String(),
	})
}

func (h *recordHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	valuesMap, err := parseValues(args["values"])
	if err != nil {
		return nil, err
	}

	rec := &cmpTypes.Record{NamespaceID: nsID, ModuleID: modID}
	for name, value := range valuesMap {
		rec.Values = append(rec.Values, &cmpTypes.RecordValue{Name: name, Value: value})
	}

	rec, _, err = cmpService.DefaultRecord.Create(ctx, rec)
	if err != nil {
		return nil, toolkit.Errf("record creation", err)
	}
	return toolkit.JSONResult(rec)
}

func (h *recordHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	recID, err := toolkit.ReqID(args, "recordID")
	if err != nil {
		return nil, err
	}

	valuesMap, err := parseValues(args["values"])
	if err != nil {
		return nil, err
	}

	rec := &cmpTypes.Record{ID: recID, NamespaceID: nsID, ModuleID: modID}
	for name, value := range valuesMap {
		rec.Values = append(rec.Values, &cmpTypes.RecordValue{Name: name, Value: value})
	}

	rec, _, err = cmpService.DefaultRecord.Update(ctx, rec)
	if err != nil {
		return nil, toolkit.Errf("record update", err)
	}
	return toolkit.JSONResult(rec)
}

func (h *recordHandler) del(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	recID, err := toolkit.ReqID(args, "recordID")
	if err != nil {
		return nil, err
	}

	if err = cmpService.DefaultRecord.DeleteByID(ctx, nsID, modID, recID); err != nil {
		return nil, toolkit.Errf("record delete", err)
	}
	return toolkit.TextResult("record %d deleted", recID), nil
}

// parseValues accepts field values as either a JSON object or a JSON string
// holding one, because models produce both.
func parseValues(raw any) (map[string]string, error) {
	switch v := raw.(type) {
	case nil:
		return nil, fmt.Errorf("values is required")
	case string:
		if v == "" {
			return nil, fmt.Errorf("values is required")
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, fmt.Errorf("invalid values JSON: %w", err)
		}
		out := make(map[string]string, len(m))
		for key, val := range m {
			out[key] = fmt.Sprintf("%v", val)
		}
		return out, nil
	case map[string]any:
		out := make(map[string]string, len(v))
		for key, val := range v {
			out[key] = fmt.Sprintf("%v", val)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("invalid values: expected JSON string or object")
	}
}
