package agentic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// withValueIssues folds the per-field detail the record service returns beside
// its error into the message. The set is a separate return value, so discarding
// it leaves the caller with "invalid record value input" — a message naming
// neither the field nor the rule that refused the write, and unactionable for
// anything trying to correct itself.
func withValueIssues(err error, dd *cmpTypes.RecordValueErrorSet) error {
	if dd.IsValid() {
		// Duplicate detection hands the set back as a return value, but the
		// validation path attaches it to the error instead
		// (RecordErrValueInput().Wrap(rve)), so the chain is the other place
		// it lives — and that path is the one that reports a required field.
		var wrapped *cmpTypes.RecordValueErrorSet
		if errors.As(err, &wrapped) {
			dd = wrapped
		}
	}

	if dd.IsValid() {
		return err
	}

	issues := make([]string, 0, dd.Len())
	for _, e := range dd.Set {
		if field, ok := e.Meta["field"].(string); ok && field != "" {
			issues = append(issues, fmt.Sprintf("%s: %s", field, e.Message))
			continue
		}
		issues = append(issues, e.Message)
	}

	return fmt.Errorf("%w (%s)", err, strings.Join(issues, "; "))
}

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
		return toolkit.JSONResultWith(rec, recordLinks(ctx, rec))
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

	// Pass the cursor itself, not its String(): String is a human-readable debug
	// rendering ("<id: 123, [FWD]>") while parseCursor expects base64 of the
	// cursor JSON, which is what MarshalJSON emits. json.Marshal renders a nil
	// pointer as null, so the last page is safe.
	return toolkit.JSONResult(map[string]any{
		"records":        set,
		"nextPageCursor": out.NextPage,
	})
}

// report aggregates server-side.
//
// It exists because the alternative is arithmetic done by a language model over
// rows it was handed, which is wrong often enough to be untrustworthy and wrong
// differently on each run — and which silently ignores every record past the
// first page.
func (h *recordHandler) report(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	metrics := toolkit.Str(args, "metrics")

	out, err := cmpService.DefaultRecord.Report(
		ctx,
		nsID,
		modID,
		metrics,
		toolkit.Str(args, "dimension"),
		toolkit.Str(args, "filter"),
		filter.StateExcluded,
	)
	if err != nil {
		return nil, toolkit.Errf("record report", err)
	}

	res := map[string]any{"rows": out}

	// A bare number carries no unit, and a caller that has to guess one guesses
	// wrong — a total off a field prefixed "$ " came back reported in euros.
	// The prefix is on the field, so it travels with the answer rather than
	// costing a second lookup nobody remembers to make.
	if mod, err := cmpService.DefaultModule.FindByID(ctx, nsID, modID); err == nil {
		if units := metricUnits(mod, metrics); len(units) > 0 {
			res["units"] = units
		}
	}

	return toolkit.JSONResult(res)
}

// metricFieldRef pulls the field name out of an aggregate expression:
// "SUM(line_value) AS total" -> "line_value".
var metricFieldRef = regexp.MustCompile(`\(([^()]*)\)`)

// metricUnits maps each aggregated field to how the module says it should read.
// Only fields carrying a prefix or suffix are reported — a plain number has no
// unit to state, and listing it would only invite one to be invented.
func metricUnits(mod *cmpTypes.Module, metrics string) map[string]any {
	units := make(map[string]any)

	for _, m := range metricFieldRef.FindAllStringSubmatch(metrics, -1) {
		name := strings.TrimSpace(m[1])
		f := mod.Fields.FindByName(name)
		if f == nil {
			continue
		}

		u := make(map[string]any)
		if p := f.Options.String("prefix"); p != "" {
			u["prefix"] = p
		}
		if sfx := f.Options.String("suffix"); sfx != "" {
			u["suffix"] = sfx
		}
		if len(u) > 0 {
			u["label"] = f.Label
			units[name] = u
		}
	}

	return units
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

	values, err := parseValues(args["values"])
	if err != nil {
		return nil, err
	}

	rec := &cmpTypes.Record{NamespaceID: nsID, ModuleID: modID, Values: values}

	rec, dd, err := cmpService.DefaultRecord.Create(ctx, rec)
	if err != nil {
		return nil, toolkit.Errf("record creation", withValueIssues(err, dd))
	}
	return toolkit.JSONResultWith(rec, recordLinks(ctx, rec))
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

	values, err := parseValues(args["values"])
	if err != nil {
		return nil, err
	}

	rec := &cmpTypes.Record{ID: recID, NamespaceID: nsID, ModuleID: modID, Values: values}

	rec, dd, err := cmpService.DefaultRecord.Update(ctx, rec)
	if err != nil {
		return nil, toolkit.Errf("record update", withValueIssues(err, dd))
	}
	return toolkit.JSONResultWith(rec, recordLinks(ctx, rec))
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

func (h *recordHandler) undelete(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// One record per call, even though the service takes a variadic set: a bulk
	// undelete swallows per-record errors and reports success, which would tell
	// the caller nothing about what actually came back.
	if err = cmpService.DefaultRecord.UndeleteByID(ctx, nsID, modID, recID); err != nil {
		return nil, toolkit.Errf("record undelete", err)
	}
	return toolkit.TextResult("record %d restored", recID), nil
}

// parseValues turns the `values` argument into the record values the store
// actually holds: a LIST of {name, value, place}, not a map of name to value.
//
// A multi-value field is several values sharing a name, told apart by place, so
// a map cannot express one — and the tool's own parameter is a JSON object,
// which cannot carry a key twice. An array is therefore how a caller says
// "these several", and each element becomes its own value.
//
// Everything that is not a scalar, an array or an object is refused rather than
// stringified. The predecessor ran every value through fmt.Sprintf("%v"), which
// turned ["a","b"] into the literal `[a b]` and {"coordinates":[46,14]} into
// `map[coordinates:[46 14]]`; whether that surfaced as an error or as stored
// garbage was decided by how strict the field kind's validator happened to be,
// so String, Number, DateTime and Geometry took it silently.
func parseValues(raw any) (cmpTypes.RecordValueSet, error) {
	var m map[string]any

	switch v := raw.(type) {
	case nil:
		return nil, fmt.Errorf("values is required")
	case string:
		if v == "" {
			return nil, fmt.Errorf("values is required")
		}
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, fmt.Errorf("invalid values JSON: %w", err)
		}
	case map[string]any:
		m = v
	default:
		return nil, fmt.Errorf("invalid values: expected JSON string or object")
	}

	// Map iteration is unordered and these become rows; fix an order so the
	// same payload always produces the same record.
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make(cmpTypes.RecordValueSet, 0, len(m))
	for _, name := range names {
		switch val := m[name].(type) {
		case []any:
			for i, item := range val {
				str, err := recordValueString(name, item)
				if err != nil {
					return nil, err
				}
				out = append(out, &cmpTypes.RecordValue{Name: name, Value: str, Place: uint(i)})
			}
		default:
			str, err := recordValueString(name, val)
			if err != nil {
				return nil, err
			}
			out = append(out, &cmpTypes.RecordValue{Name: name, Value: str})
		}
	}

	return out, nil
}

// recordValueString renders one value the way the store keeps it: scalars as
// themselves, and an object as its JSON, which is the form a Geometry value is
// held in. A nested array has no meaning — a field is one value or a list of
// them, never a list of lists — so it is refused rather than guessed at.
func recordValueString(name string, val any) (string, error) {
	switch v := val.(type) {
	case nil:
		return "", nil
	case string:
		return v, nil
	case bool, float64, float32, int, int64, uint64, json.Number:
		return fmt.Sprintf("%v", v), nil
	case map[string]any:
		b, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("value of %q cannot be encoded: %w", name, err)
		}
		return string(b), nil
	default:
		return "", fmt.Errorf(
			"value of %q is a %T, which is not a field value; send a string, number, "+
				"boolean, an object, or an array of those for a multi-value field", name, val)
	}
}
