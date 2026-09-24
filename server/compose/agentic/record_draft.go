package agentic

import (
	"context"
	"sort"
	"strconv"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/crusttech/human/server/pkg/weburl"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// draftChoiceLimit caps how many records or users a reference field offers in
// the form. The form picks from what it was given; it cannot search.
const draftChoiceLimit = 100

// draftNote is what the model reads back: the draft is on screen and nothing
// has been written.
const draftNote = "Shown to the user as a form; nothing has been saved. The user reviews it and saves " +
	"it from the form, so do not create or update this record yourself unless they ask you to."

// draft puts a proposed record in front of the user without writing it. With
// recordID it proposes changes to that record: the fields named in values
// replace the record's own, and the rest are shown as they are.
func (h *recordHandler) draft(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	nsID, modID, err := h.resolveNsMod(ctx, args)
	if err != nil {
		return nil, err
	}

	mod, err := cmpService.DefaultModule.FindByID(ctx, nsID, modID)
	if err != nil {
		return nil, toolkit.Errf("record draft", err)
	}

	var (
		proposed cmpTypes.RecordValueSet
		named    []string
		current  cmpTypes.RecordValueSet
	)

	if raw, ok := args["values"]; ok && raw != "" && raw != nil {
		if proposed, named, err = parseValues(raw); err != nil {
			return nil, err
		}
	}

	recID, err := toolkit.ID(args, "recordID")
	if err != nil {
		return nil, err
	}
	if recID > 0 {
		rec, _, err := cmpService.DefaultRecord.FindByID(ctx, nsID, modID, recID)
		if err != nil {
			return nil, toolkit.Errf("record draft", err)
		}
		current = rec.Values
	}

	merged := mergeDraft(current, proposed, named)

	res := map[string]any{
		"draft":       draftValues(mod, merged),
		"namespaceID": strconv.FormatUint(nsID, 10),
		"moduleID":    strconv.FormatUint(modID, 10),
		"note":        draftNote,
	}
	if recID > 0 {
		res["recordID"] = strconv.FormatUint(recID, 10)
	}
	if refs := refLabels(ctx, mod, cmpTypes.RecordSet{{NamespaceID: nsID, ModuleID: modID, Values: merged}}); refs != nil {
		res["refs"] = refs
	}

	result, err := toolkit.JSONResult(res)
	if err != nil {
		return nil, err
	}

	view := recordFormViewOf(mod, draftChoices(ctx, mod))
	if recID > 0 {
		view.Link = weburl.ComposeRecord(nsURLPart(ctx, nsID), modID, recID)
	}

	return hmcp.WithViewData(result, view), nil
}

// mergeDraft replaces every field named in the proposal with the proposal's
// values and keeps the rest of current. A named field with no values is
// cleared, the way compose_record_update clears it.
func mergeDraft(current, proposed cmpTypes.RecordValueSet, named []string) cmpTypes.RecordValueSet {
	replaced := make(map[string]struct{}, len(named))
	for _, n := range named {
		replaced[n] = struct{}{}
	}

	out := make(cmpTypes.RecordValueSet, 0, len(current)+len(proposed))
	for _, v := range current {
		if _, ok := replaced[v.Name]; !ok {
			out = append(out, v)
		}
	}
	return append(out, proposed...)
}

// draftValues is the draft in the shape compose_record_create takes: one value
// per field, an array for a multi-value field. Values naming no field of mod
// are kept, so the form can say it does not know them.
func draftValues(mod *cmpTypes.Module, vv cmpTypes.RecordValueSet) map[string]any {
	out := make(map[string]any, len(vv))
	for _, v := range vv {
		if f := mod.Fields.FindByName(v.Name); f != nil && f.Multi {
			list, _ := out[v.Name].([]string)
			out[v.Name] = append(list, v.Value)
			continue
		}
		out[v.Name] = v.Value
	}
	return out
}

// draftChoices lists what each Record and User field can point at, labelled
// the way the lookup's refs are.
func draftChoices(ctx context.Context, mod *cmpTypes.Module) map[string]fieldChoices {
	out := map[string]fieldChoices{}
	var users *fieldChoices

	for _, f := range mod.Fields {
		switch f.Kind {
		case "Record":
			if c, ok := recordChoices(ctx, mod.NamespaceID, f); ok {
				out[f.Name] = c
			}
		case "User":
			if users == nil {
				c := userChoices(ctx)
				users = &c
			}
			out[f.Name] = *users
		}
	}

	return out
}

func recordChoices(ctx context.Context, nsID uint64, f *cmpTypes.ModuleField) (fieldChoices, bool) {
	refModID := f.Options.Uint64("moduleID")
	if refModID == 0 {
		return fieldChoices{}, false
	}

	set, out, err := cmpService.DefaultRecord.Search(ctx, cmpTypes.RecordFilter{
		NamespaceID: nsID,
		ModuleID:    refModID,
		Paging:      filter.Paging{Limit: draftChoiceLimit},
	})
	if err != nil {
		return fieldChoices{}, false
	}

	ids := make(map[uint64]struct{}, len(set))
	for _, r := range set {
		ids[r.ID] = struct{}{}
	}

	labels := map[string]string{}
	addRecordLabels(ctx, labels, nsID, refModID, &refTarget{
		ids:              ids,
		labelField:       f.Options.String("labelField"),
		recordLabelField: f.Options.String("recordLabelField"),
	})

	return choicesOf(set.IDs(), labels, out.NextPage != nil), true
}

func userChoices(ctx context.Context) fieldChoices {
	uu, out, err := sysService.DefaultUser.Find(ctx, sysTypes.UserFilter{
		Paging: filter.Paging{Limit: draftChoiceLimit},
	})
	if err != nil {
		return fieldChoices{}
	}

	ids := make([]uint64, 0, len(uu))
	labels := make(map[string]string, len(uu))
	for _, u := range uu {
		ids = append(ids, u.ID)
		labels[strconv.FormatUint(u.ID, 10)] = userLabel(u)
	}

	return choicesOf(ids, labels, out.NextPage != nil)
}

// choicesOf orders choices by label, falling back to the ID for one that has
// none.
func choicesOf(ids []uint64, labels map[string]string, more bool) fieldChoices {
	c := fieldChoices{Items: make([]fieldChoice, 0, len(ids)), More: more}
	for _, id := range ids {
		key := strconv.FormatUint(id, 10)
		label := labels[key]
		if label == "" {
			label = key
		}
		c.Items = append(c.Items, fieldChoice{ID: key, Label: label})
	}
	sort.SliceStable(c.Items, func(i, j int) bool { return c.Items[i].Label < c.Items[j].Label })
	return c
}
