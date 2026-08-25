package agentic

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
)

// refLabelCeiling caps how many distinct references one result resolves.
//
// The dictionary is meant to be small beside the records it explains; a page of
// records that each point somewhere different is the case where it stops being
// worth its size.
const refLabelCeiling = 500

// refTarget is one referenced module and everything this result points at in it.
type refTarget struct {
	ids        map[uint64]struct{}
	labelField string
}

// refLabels maps every record and user ID a result mentions to the label a
// person would see in its place.
//
// A Record value is stored as the target's bare ID, so a caller reading
// deck_card gets `card: "510764704641581057"` and has no way to know that is
// Lightning Bolt short of looking each one up. Resolving them here is the
// difference between one call and one call per distinct reference.
//
// The dictionary is keyed by ID rather than by field so that an ID appearing
// anywhere in the payload — including in a filter the caller wrote itself —
// resolves against it. Nothing here fails the call: an unresolvable reference
// is simply absent, which reads the same as it did before.
func refLabels(ctx context.Context, mod *cmpTypes.Module, rr cmpTypes.RecordSet) map[string]string {
	if mod == nil || len(rr) == 0 {
		return nil
	}

	var (
		byModule = map[uint64]*refTarget{}
		users    = map[uint64]struct{}{}
	)

	note := func(dst map[uint64]struct{}, raw string) {
		if id, err := strconv.ParseUint(raw, 10, 64); err == nil && id > 0 {
			dst[id] = struct{}{}
		}
	}

	for _, f := range mod.Fields {
		switch f.Kind {
		case "Record":
			refMod, err := strconv.ParseUint(f.Options.String("moduleID"), 10, 64)
			if err != nil || refMod == 0 {
				continue
			}
			t := byModule[refMod]
			if t == nil {
				t = &refTarget{ids: map[uint64]struct{}{}}
				byModule[refMod] = t
			}
			// Two fields may point at one module and name different label
			// fields; the dictionary holds one label per ID, so the first
			// choice wins and the rest read the same as the viewer fallback.
			if t.labelField == "" {
				t.labelField = f.Options.String("labelField")
			}
			for _, r := range rr {
				for _, v := range r.Values.FilterByName(f.Name) {
					note(t.ids, v.Value)
				}
			}
		case "User":
			for _, r := range rr {
				for _, v := range r.Values.FilterByName(f.Name) {
					note(users, v.Value)
				}
			}
		}
	}

	// Every record carries these three, and they are IDs in exactly the same way
	// a Record value is.
	for _, r := range rr {
		for _, id := range []uint64{r.OwnedBy, r.CreatedBy, r.UpdatedBy} {
			if id > 0 {
				users[id] = struct{}{}
			}
		}
	}

	out := map[string]string{}
	for refMod, t := range byModule {
		addRecordLabels(ctx, out, mod.NamespaceID, refMod, t)
	}
	addUserLabels(ctx, out, users)

	if len(out) == 0 {
		return nil
	}
	return out
}

// labelFieldOf picks the field whose value stands in for a record, mirroring
// what the webapp's viewers do (CFieldRecordViewer.vue labelFieldDef): the
// named labelField when it resolves, and the module's first field otherwise.
func labelFieldOf(mod *cmpTypes.Module, named string) *cmpTypes.ModuleField {
	if mod == nil || len(mod.Fields) == 0 {
		return nil
	}
	if named != "" {
		if f := mod.Fields.FindByName(named); f != nil {
			return f
		}
	}
	return mod.Fields[0]
}

func addRecordLabels(ctx context.Context, out map[string]string, nsID, modID uint64, t *refTarget) {
	if t == nil || len(t.ids) == 0 || len(out) >= refLabelCeiling {
		return
	}

	refMod, err := cmpService.DefaultModule.FindByID(ctx, nsID, modID)
	if err != nil || refMod == nil {
		return
	}

	lf := labelFieldOf(refMod, t.labelField)
	if lf == nil {
		return
	}

	for _, rec := range findRecordsByID(ctx, nsID, modID, t.ids) {
		if v := rec.Values.Get(lf.Name, 0); v != nil && v.Value != "" {
			out[strconv.FormatUint(rec.ID, 10)] = v.Value
		}
	}
}

// findRecordsByID loads a specific set of records in one search.
//
// The compose query language has no IN over recordID ("unsupported IN operator
// on a single value field"), so an equality chain is what a batch looks like.
func findRecordsByID(ctx context.Context, nsID, modID uint64, ids map[uint64]struct{}) cmpTypes.RecordSet {
	if len(ids) == 0 {
		return nil
	}

	var q strings.Builder
	for id := range ids {
		if q.Len() > 0 {
			q.WriteString(" OR ")
		}
		q.WriteString("recordID = ")
		q.WriteString(strconv.FormatUint(id, 10))
	}

	set, _, err := cmpService.DefaultRecord.Search(ctx, cmpTypes.RecordFilter{
		NamespaceID: nsID,
		ModuleID:    modID,
		Query:       q.String(),
		Paging:      filter.Paging{Limit: uint(len(ids))},
	})
	if err != nil {
		return nil
	}
	return set
}

func addUserLabels(ctx context.Context, out map[string]string, ids map[uint64]struct{}) {
	if len(ids) == 0 || len(out) >= refLabelCeiling {
		return
	}

	ss := make([]string, 0, len(ids))
	for id := range ids {
		ss = append(ss, strconv.FormatUint(id, 10))
	}

	uu, _, err := sysService.DefaultUser.Find(ctx, sysTypes.UserFilter{
		UserID:   ss,
		AllKinds: true,
		Paging:   filter.Paging{Limit: uint(len(ss))},
	})
	if err != nil {
		return
	}

	for _, u := range uu {
		if l := userLabel(u); l != "" {
			out[strconv.FormatUint(u.ID, 10)] = l
		}
	}
}

func userLabel(u *sysTypes.User) string {
	for _, s := range []string{u.Name, u.Username, u.Handle, u.Email} {
		if s != "" {
			return s
		}
	}
	return ""
}

// dimensionRefs resolves the group keys of a report grouped by a reference.
//
// A breakdown by a Record field comes back keyed by the target's bare ID —
// "dimension_0": "510764704641581057" — and reading five of those as names cost
// five more calls, one per row. The dictionary is the same one a record lookup
// returns, so a caller reads both the same way.
func dimensionRefs(ctx context.Context, mod *cmpTypes.Module, dimension string, rows any) map[string]string {
	if mod == nil || dimension == "" {
		return nil
	}

	f := mod.Fields.FindByName(dimension)
	if f == nil {
		return nil
	}

	ids := map[uint64]struct{}{}
	for _, raw := range dimensionValues(rows) {
		if id, err := strconv.ParseUint(raw, 10, 64); err == nil && id > 0 {
			ids[id] = struct{}{}
		}
	}
	if len(ids) == 0 {
		return nil
	}

	out := map[string]string{}
	switch f.Kind {
	case "Record":
		refMod, err := strconv.ParseUint(f.Options.String("moduleID"), 10, 64)
		if err != nil || refMod == 0 {
			return nil
		}
		addRecordLabels(ctx, out, mod.NamespaceID, refMod, &refTarget{
			ids:        ids,
			labelField: f.Options.String("labelField"),
		})
	case "User":
		addUserLabels(ctx, out, ids)
	default:
		return nil
	}

	if len(out) == 0 {
		return nil
	}
	return out
}

// dimensionValues pulls dimension_0 out of report rows.
//
// The report is typed as `any` by the service and its row type is unexported,
// so the rows are read back through their JSON form — the same form they are
// about to be returned in.
func dimensionValues(rows any) []string {
	enc, err := json.Marshal(rows)
	if err != nil {
		return nil
	}

	var decoded []map[string]any
	if err = json.Unmarshal(enc, &decoded); err != nil {
		return nil
	}

	out := make([]string, 0, len(decoded))
	for _, row := range decoded {
		if v, ok := row["dimension_0"]; ok {
			if s, ok := v.(string); ok && s != "" {
				out = append(out, s)
			}
		}
	}
	return out
}
