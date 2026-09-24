package agentic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
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

	return fmt.Errorf("%w (%s)", err, strings.Join(valueIssueLines(dd), "; "))
}

// valueIssueLines renders one record value issue per line, field first.
//
// Identical lines are collapsed. Duplicate detection reports once per matching
// record, so a value held by three of them arrives as the same sentence three
// times — which says nothing the first one did not, and the record IDs it would
// take to tell them apart are not in the message.
func valueIssueLines(dd *cmpTypes.RecordValueErrorSet) []string {
	if dd.IsValid() {
		return nil
	}

	var (
		out  = make([]string, 0, dd.Len())
		seen = make(map[string]struct{}, dd.Len())
	)

	for _, e := range dd.Set {
		line := fillValuePlaceholders(e.Message, e.Meta)
		if field, ok := e.Meta["field"].(string); ok && field != "" {
			line = fmt.Sprintf("%s: %s", field, line)
		}

		if _, dup := seen[line]; dup {
			continue
		}
		seen[line] = struct{}{}
		out = append(out, line)
	}
	return out
}

// valuePlaceholder matches the slots a record value message is written with.
var valuePlaceholder = regexp.MustCompile(`{{\s*(\w+)\s*}}`)

// fillValuePlaceholders puts each issue's own meta into its message.
//
// The messages are templates — `The value "{{value}}" already exists in another
// record` — and the value lives in the issue's meta beside them. The webapp
// fills them (sections/compose/lib/record-errors.js) and nothing else did, so a
// caller here was told that a value it was never shown is taken.
func fillValuePlaceholders(msg string, meta map[string]any) string {
	if len(meta) == 0 || !strings.Contains(msg, "{{") {
		return msg
	}

	return valuePlaceholder.ReplaceAllStringFunc(msg, func(whole string) string {
		v, ok := meta[valuePlaceholder.FindStringSubmatch(whole)[1]]
		if !ok || v == nil || v == "" {
			return whole
		}
		return fmt.Sprintf("%v", v)
	})
}

// duplicateWarningNote states the duplicate rules a write matched without being
// refused by them.
//
// A rule with strict off is meant to warn, and the service hands its matches
// back beside the record rather than in an error. Dropping them made such a
// rule do nothing at all that a caller could see: the record was created, the
// response was a plain success, and the duplicate it was warning about was
// never mentioned.
func duplicateWarningNote(dd *cmpTypes.RecordValueErrorSet) string {
	lines := valueIssueLines(dd)
	if len(lines) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"Saved, and duplicate detection matched an existing record: %s. The rule is not strict, so the write went through — check whether this record should exist before relying on it.",
		strings.Join(lines, "; "),
	)
}

type (
	toolRegistrar interface {
		RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc, opts ...hmcp.RegisterOption)
		RegisterUIResource(ui hmcp.UIResource)
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
	ns, err := lookupNamespaceArg(ctx, args)
	if err != nil {
		return 0, 0, err
	}

	mod, err := resolveModuleArg(ctx, ns, args["module"])
	if err != nil {
		return 0, 0, err
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
		extra := map[string]any{}
		for k, v := range recordLinks(ctx, rec) {
			extra[k] = v
		}
		mod, err := cmpService.DefaultModule.FindByID(ctx, nsID, modID)
		if err != nil {
			mod = nil
		} else if refs := refLabels(ctx, mod, cmpTypes.RecordSet{rec}); refs != nil {
			extra["refs"] = refs
		}
		res, err := toolkit.JSONResultWithAny(rec, extra)
		return withRecordView(ctx, res, mod, cmpTypes.RecordSet{rec}), err
	}

	// A caller holding several IDs — the group keys of a report, say — asks for
	// them all at once or not at all, and the query language has no IN over
	// recordID to write that with. Both bracketed and parenthesised attempts
	// come back as parser noise, so the batch is a parameter instead.
	if ids := parseIDList(toolkit.Str(args, "recordIDs")); len(ids) > 0 {
		return h.lookupByIDs(ctx, nsID, modID, ids)
	}

	// One module load for the four things that need it: the Select check, the
	// reference paths in the filter, the sort columns and the reference labels.
	mod, _ := cmpService.DefaultModule.FindByID(ctx, nsID, modID)

	query := toolkit.Str(args, "filter")
	if err = checkSelectValues(mod, "record list", query); err != nil {
		return nil, err
	}
	if query, err = resolveRefPaths(ctx, mod, query); err != nil {
		return nil, err
	}

	f := cmpTypes.RecordFilter{
		NamespaceID: nsID,
		ModuleID:    modID,
		Query:       query,
	}

	if err = applyRecordSort(&f, mod, toolkit.Str(args, "sort")); err != nil {
		return nil, err
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
		return nil, fieldNameError(ctx, "record list", nsID, modID, err)
	}

	// Pass the cursor itself, not its String(): String is a human-readable debug
	// rendering ("<id: 123, [FWD]>") while parseCursor expects base64 of the
	// cursor JSON, which is what MarshalJSON emits. json.Marshal renders a nil
	// pointer as null, so the last page is safe.
	res := map[string]any{
		"records":        set,
		"nextPageCursor": out.NextPage,
	}
	if refs := refLabels(ctx, mod, set); refs != nil {
		res["refs"] = refs
	}

	result, err := toolkit.JSONResult(res)
	return withRecordView(ctx, result, mod, set), err
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
	dimension := toolkit.Str(args, "dimension")

	var mod *cmpTypes.Module
	if cmpService.DefaultModule != nil {
		mod, _ = cmpService.DefaultModule.FindByID(ctx, nsID, modID)
	}

	if err = validateMetrics(metrics, mod); err != nil {
		return nil, err
	}

	query := toolkit.Str(args, "filter")
	if err = checkSelectValues(mod, "record report", query); err != nil {
		return nil, err
	}
	if query, err = resolveRefPaths(ctx, mod, query); err != nil {
		return nil, err
	}

	out, err := cmpService.DefaultRecord.Report(
		ctx,
		nsID,
		modID,
		metrics,
		dimension,
		query,
		filter.StateExcluded,
	)
	if err != nil {
		return nil, reportError(ctx, nsID, modID, err)
	}

	out = sortAndLimitRows(out, parseSortSpec(toolkit.Str(args, "sort")), reportLimit(args))

	res := map[string]any{"rows": out}

	// A bare number carries no unit, and a caller that has to guess one guesses
	// wrong — a total off a field prefixed "$ " came back reported in euros.
	// The prefix is on the field, so it travels with the answer rather than
	// costing a second lookup nobody remembers to make.
	if mod != nil {
		if units := metricUnits(mod, metrics); len(units) > 0 {
			res["units"] = units
		}
		if refs := dimensionRefs(ctx, mod, dimension, out); refs != nil {
			res["refs"] = refs
		}
	}

	return toolkit.JSONResult(res)
}

// reportError names the fields that do exist when the caller asked for one that
// does not.
//
// The DAL answers "unknown attribute card.rarity" — true, and not enough to act
// on. A model given only that guesses again, or falls back to reading every
// record and adding the numbers up by hand, which is the thing this tool exists
// to stop. Field names are cheap to list and turn the retry into the right one.
func reportError(ctx context.Context, nsID, modID uint64, err error) error {
	names := moduleFieldNames(ctx, nsID, modID)
	if names == "" || !strings.Contains(err.Error(), "unknown attribute") {
		return toolkit.Errf("record report", err)
	}

	return fmt.Errorf(
		"record report failed: %w. Metrics and dimension name fields on this module only — a dotted path through a Record reference is not one. Available: %s. To group by a field of a referenced record, aggregate that module instead",
		err, names,
	)
}

// fieldNameError names the module's fields when a query used one that is not
// there.
//
// The usual cause is a dotted path through a reference — `card.name = 'Bolt'`
// — which reads so naturally that a model writes it twice before giving up.
// The store's own answer, "unknown attribute", says which name failed and
// nothing about what would have worked.
// recordSortColumns are the columns a record has beside its module's fields.
// 'recordID' is taken as a synonym for 'ID' — the rest of this surface calls it
// recordID and the store does not.
var recordSortColumns = []string{
	"ID", "moduleID", "namespaceID", "revision", "ownedBy",
	"createdAt", "createdBy", "updatedAt", "updatedBy", "deletedAt", "deletedBy",
}

// applyRecordSort puts the caller's ordering on the filter.
//
// Ordering is the store's job: without it a caller after the ten newest rows
// reads every page and orders them itself, which the result ceiling stops it
// doing on any module big enough to want ordering.
//
// A bare column is checked against the module before the query runs, because
// the store answers an unknown one with the same "unknown attribute" a filter
// produces — sending a caller whose sort was wrong to look at a filter that is
// fine. An expression carrying a modifier is left to the store, which is the
// only thing that knows what the modifier accepts.
func applyRecordSort(f *cmpTypes.RecordFilter, mod *cmpTypes.Module, expr string) error {
	if strings.TrimSpace(expr) == "" {
		return nil
	}

	if err := f.Sort.Set(expr); err != nil {
		return fmt.Errorf("invalid sort %q: %w", expr, err)
	}

	known := make(map[string]struct{}, len(recordSortColumns))
	for _, c := range recordSortColumns {
		known[c] = struct{}{}
	}
	if mod != nil {
		for _, fld := range mod.Fields {
			known[fld.Name] = struct{}{}
		}
	}

	for _, e := range f.Sort {
		if strings.EqualFold(e.Column, "recordID") {
			e.SetColumns("ID")
		}

		if mod == nil || e.Modifier() != "" {
			continue
		}

		if _, ok := known[e.Column]; !ok {
			return fmt.Errorf(
				"cannot sort by %q: this module has no such field. Fields: %s. Record columns: %s",
				e.Column, moduleFieldNamesOf(mod), strings.Join(recordSortColumns, ", "))
		}
	}

	return nil
}

func fieldNameError(ctx context.Context, subject string, nsID, modID uint64, err error) error {
	names := moduleFieldNames(ctx, nsID, modID)
	if names == "" || !strings.Contains(err.Error(), "unknown attribute") {
		return toolkit.Errf(subject, err)
	}

	return fmt.Errorf(
		"%s failed: %w. A filter names fields on this module only — a dotted path through a Record reference is not one. Available: %s. To filter by a referenced record's field, look that record up first and filter on the reference field by its ID",
		subject, err, names,
	)
}

func moduleFieldNames(ctx context.Context, nsID, modID uint64) string {
	if cmpService.DefaultModule == nil {
		return ""
	}

	mod, err := cmpService.DefaultModule.FindByID(ctx, nsID, modID)
	if err != nil || mod == nil {
		return ""
	}

	return moduleFieldNamesOf(mod)
}

func moduleFieldNamesOf(mod *cmpTypes.Module) string {
	if mod == nil {
		return ""
	}

	names := make([]string, 0, len(mod.Fields))
	for _, f := range mod.Fields {
		names = append(names, f.Name)
	}
	return strings.Join(names, ", ")
}

// metricFieldRef pulls what is being aggregated out of an expression:
// "SUM(line_value) AS total" -> "line_value".
var metricFieldRef = regexp.MustCompile(`\(([^()]*)\)`)

// metricIdent picks the field names out of that, since what is aggregated is
// often an expression rather than a bare field: SUM(quantity * unit_value).
var metricIdent = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]*`)

// metricUnits maps each aggregated field to how the module says it should read.
// Only fields carrying a prefix or suffix are reported — a plain number has no
// unit to state, and listing it would only invite one to be invented.
//
// Every field the expression touches is reported, not just a lone one: reading
// the whole parenthesised body as a field name meant a total over
// "quantity_wanted * current_price" carried no unit at all, and a wishlist
// priced in dollars came back quoted in euros.
func metricUnits(mod *cmpTypes.Module, metrics string) map[string]any {
	units := make(map[string]any)

	for _, m := range metricFieldRef.FindAllStringSubmatch(metrics, -1) {
		for _, name := range metricIdent.FindAllString(m[1], -1) {
			addFieldUnit(units, mod, name)
		}
	}

	return units
}

func addFieldUnit(units map[string]any, mod *cmpTypes.Module, name string) {
	if _, seen := units[name]; seen {
		return
	}

	f := mod.Fields.FindByName(name)
	if f == nil {
		return
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

func (h *recordHandler) create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	values, _, err := parseValues(args["values"])
	if err != nil {
		return nil, err
	}

	rec := &cmpTypes.Record{NamespaceID: nsID, ModuleID: modID, Values: values}

	rec, dd, err := cmpService.DefaultRecord.Create(ctx, rec)
	if err != nil {
		return nil, toolkit.Errf("record creation", withValueIssues(err, dd))
	}
	return toolkit.JSONResultWith(rec, withNote(recordLinks(ctx, rec), duplicateWarningNote(dd)))
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

	values, named, err := parseValues(args["values"])
	if err != nil {
		return nil, err
	}

	if !toolkit.Bool(args, "replace") {
		if values, err = mergeOntoRecord(ctx, nsID, modID, recID, values, named); err != nil {
			return nil, err
		}
	}

	rec := &cmpTypes.Record{ID: recID, NamespaceID: nsID, ModuleID: modID, Values: values}

	rec, dd, err := cmpService.DefaultRecord.Update(ctx, rec)
	if err != nil {
		return nil, toolkit.Errf("record update", withValueIssues(err, dd))
	}
	return toolkit.JSONResultWith(rec, withNote(recordLinks(ctx, rec), duplicateWarningNote(dd)))
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
// parseValues turns the values argument into a value set, and reports which
// fields it named — which is not the same thing: a field sent as an empty list
// contributes no values, and only the name says it was meant to be cleared.
func parseValues(raw any) (cmpTypes.RecordValueSet, []string, error) {
	var m map[string]any

	switch v := raw.(type) {
	case nil:
		return nil, nil, fmt.Errorf("values is required")
	case string:
		if v == "" {
			return nil, nil, fmt.Errorf("values is required")
		}
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, nil, fmt.Errorf("invalid values JSON: %w", err)
		}
	case map[string]any:
		m = v
	default:
		return nil, nil, fmt.Errorf("invalid values: expected JSON string or object")
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
					return nil, nil, err
				}
				out = append(out, &cmpTypes.RecordValue{Name: name, Value: str, Place: uint(i)})
			}
		default:
			str, err := recordValueString(name, val)
			if err != nil {
				return nil, nil, err
			}
			out = append(out, &cmpTypes.RecordValue{Name: name, Value: str})
		}
	}

	return out, names, nil
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

// aggregateCall matches an expression that begins with one of the aggregate
// functions the report understands.
var aggregateCall = regexp.MustCompile(`(?i)^(SUM|AVG|MIN|MAX|COUNT|COUNTD)\s*\(`)

// unknownCall matches a metric that calls a function this report does not have.
var unknownCall = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*\(`)

// asAlias strips the "AS name" a metric may end with.
var asAlias = regexp.MustCompile(`(?i)\s+AS\s+[A-Za-z_][A-Za-z0-9_]*\s*$`)

// validateMetrics refuses a metric that is a bare field name.
//
// The aggregation reaches the database with such a field ungrouped, and what
// comes back is the driver's own complaint — `pq: column
// "compose_record.values" must appear in the GROUP BY clause` — which names
// neither the metric at fault nor anything the caller wrote. It is a plain
// mistake with a plain correction, so make it here.
func validateMetrics(metrics string, mod *cmpTypes.Module) error {
	for _, m := range splitMetrics(metrics) {
		expr := strings.TrimSpace(asAlias.ReplaceAllString(m, ""))
		if expr == "" || aggregateCall.MatchString(expr) {
			continue
		}

		// A bare field listed as a metric is usually one meant to label the
		// groups, not to be summed — the same mistake as writing a column into
		// a GROUP BY query. Say so when the field cannot be aggregated at all.
		if f := fieldOf(mod, expr); f != nil && !aggregatableKind(f.Kind) {
			return fmt.Errorf(
				"%q is a %s field and cannot be aggregated: pass it as 'dimension' to group by it, and keep 'metrics' to numeric aggregates",
				expr, f.Kind,
			)
		}

		// An expression that already calls something is a function this report
		// does not have, not a field waiting to be wrapped — suggesting
		// SUM(FIRST(price)) helps nobody.
		if fn := unknownCall.FindStringSubmatch(expr); fn != nil {
			return fmt.Errorf(
				"%q is not an aggregate function here; metrics use SUM, AVG, MIN, MAX or COUNT only. The record count is returned as 'count' without asking, and ordering is done with 'sort' rather than by picking a row",
				fn[1],
			)
		}

		return fmt.Errorf(
			"metric %q is not an aggregate: wrap it in a function, e.g. %q. Every metric must be SUM, AVG, MIN, MAX or COUNT of a field; the record count is returned as 'count' without asking",
			expr, "SUM("+expr+") AS "+expr,
		)
	}
	return nil
}

func fieldOf(mod *cmpTypes.Module, name string) *cmpTypes.ModuleField {
	if mod == nil {
		return nil
	}
	return mod.Fields.FindByName(name)
}

// aggregatableKind is the set SUM and AVG mean anything over.
func aggregatableKind(kind string) bool {
	return kind == "Number"
}

// splitMetrics splits on the commas between metrics, leaving the ones inside a
// function call alone.
func splitMetrics(metrics string) []string {
	var (
		out   []string
		depth int
		cur   strings.Builder
	)

	for _, r := range metrics {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				if seg := strings.TrimSpace(cur.String()); seg != "" {
					out = append(out, seg)
				}
				cur.Reset()
				continue
			}
		}
		cur.WriteRune(r)
	}

	if s := strings.TrimSpace(cur.String()); s != "" {
		out = append(out, s)
	}
	return out
}

// parseIDList reads a comma-separated list of record IDs, ignoring whatever
// punctuation a caller wrapped it in.
func parseIDList(raw string) map[uint64]struct{} {
	raw = strings.Trim(strings.TrimSpace(raw), "[]()")
	if raw == "" {
		return nil
	}

	out := map[uint64]struct{}{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.Trim(strings.TrimSpace(part), `"'`)
		if id, err := strconv.ParseUint(part, 10, 64); err == nil && id > 0 {
			out[id] = struct{}{}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// lookupByIDs answers a batch the way the list path answers a page, so a caller
// reads one shape either way.
func (h *recordHandler) lookupByIDs(ctx context.Context, nsID, modID uint64, ids map[uint64]struct{}) (*mcp.CallToolResult, error) {
	set := findRecordsByID(ctx, nsID, modID, ids)

	res := map[string]any{"records": set}
	if len(set) < len(ids) {
		res["note"] = fmt.Sprintf("%d of %d requested records were found; the rest do not exist in this module or are deleted", len(set), len(ids))
	}

	mod, err := cmpService.DefaultModule.FindByID(ctx, nsID, modID)
	if err != nil {
		mod = nil
	} else if refs := refLabels(ctx, mod, set); refs != nil {
		res["refs"] = refs
	}

	out, err := toolkit.JSONResult(res)
	return withRecordView(ctx, out, mod, set), err
}

// sortSpec is one "column DESC" ordering over report rows.
type sortSpec struct {
	key  string
	desc bool
}

func parseSortSpec(raw string) *sortSpec {
	fields := strings.Fields(strings.TrimSpace(raw))
	if len(fields) == 0 {
		return nil
	}

	s := &sortSpec{key: strings.Trim(fields[0], `"'`)}
	if len(fields) > 1 && strings.EqualFold(fields[1], "DESC") {
		s.desc = true
	}
	return s
}

// sortAndLimitRows orders a report's groups and keeps the first few.
//
// The aggregation itself orders by the dimension and returns every group, so
// "the five that moved most" meant handing a model every group and asking it to
// rank them — which it does approximately, and which grows until the result
// exceeds the size ceiling. The ranking happens here instead. It is applied to
// the computed groups rather than pushed into the query: the database still
// aggregates everything, and only the answer is trimmed.
func sortAndLimitRows(rows any, spec *sortSpec, limit int) any {
	if spec == nil && limit <= 0 {
		return rows
	}

	enc, err := json.Marshal(rows)
	if err != nil {
		return rows
	}
	var decoded []map[string]any
	if err = json.Unmarshal(enc, &decoded); err != nil {
		return rows
	}

	if spec != nil {
		sort.SliceStable(decoded, func(i, j int) bool {
			less := rowLess(decoded[i][spec.key], decoded[j][spec.key])
			if spec.desc {
				return rowLess(decoded[j][spec.key], decoded[i][spec.key])
			}
			return less
		})
	}

	if limit > 0 && len(decoded) > limit {
		decoded = decoded[:limit]
	}
	return decoded
}

// rowLess orders two group values, numbers numerically and everything else as
// text. A missing value sorts first so it never displaces a real one from the
// top of a descending ranking.
func rowLess(a, b any) bool {
	af, aok := a.(float64)
	bf, bok := b.(float64)
	switch {
	case aok && bok:
		return af < bf
	case aok != bok:
		return bok
	case a == nil || b == nil:
		return a == nil && b != nil
	default:
		return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
	}
}

// reportLimit reads the limit the caller actually gave.
//
// toolkit.Page defaults to a page size, which is right for a list of records
// and wrong here: a report with no limit must return every group, and silently
// keeping the first fifty would drop data with nothing to show for it.
func reportLimit(args map[string]any) int {
	switch v := args["limit"].(type) {
	case string:
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n > 0 {
			return n
		}
	case float64:
		if v > 0 {
			return int(v)
		}
	}
	return 0
}

// mergeOntoRecord folds the values a caller sent onto the ones the record
// already holds.
//
// Compose stores a record's values as one set and an update replaces it whole,
// so naming three fields deletes every other one. Read conversationally —
// "put them in the binder" — that is never what was meant, and it happened: an
// update naming quantity and price silently dropped the holding's condition,
// language, foil flag and storage. The named fields are replaced entire, so a
// multi-value field still takes the whole list; everything unnamed survives.
func mergeOntoRecord(ctx context.Context, nsID, modID, recID uint64, values cmpTypes.RecordValueSet, named []string) (cmpTypes.RecordValueSet, error) {
	cur, _, err := cmpService.DefaultRecord.FindByID(ctx, nsID, modID, recID)
	if err != nil {
		return nil, toolkit.Errf("record lookup for update", err)
	}

	replacing := make(map[string]struct{}, len(named))
	for _, n := range named {
		replacing[n] = struct{}{}
	}

	out := make(cmpTypes.RecordValueSet, 0, len(cur.Values)+len(values))
	for _, v := range cur.Values {
		if _, ok := replacing[v.Name]; ok {
			continue
		}
		out = append(out, &cmpTypes.RecordValue{Name: v.Name, Value: v.Value, Ref: v.Ref, Place: v.Place})
	}

	return append(out, values...), nil
}
