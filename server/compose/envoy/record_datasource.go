package envoy

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/compose/dalutils"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/envoyx"
	"github.com/crusttech/human/server/pkg/envoyx/datasource"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
	systemTypes "github.com/crusttech/human/server/system/types"
	"github.com/spf13/cast"
)

type (
	// RecordDatasource provides a mechanism for you to access large
	// record datasets optimally
	RecordDatasource struct {
		Mapping  envoyx.DatasourceMapping
		Provider envoyx.Provider

		multivalues map[string]bool

		CheckExisting func(ctx context.Context, ref ...[]string) ([]uint64, error)

		currentIndex int

		// Index to map from ref to ID
		// @todo we might need to flush these to the disc in case a huge dataset is passed in
		refToID map[string]uint64
		// @todo might be worth putting both into one map; not sure how much space we'd save up
		existingIDs map[uint64]bool
	}

	// iteratorProvider is a wrapper around the dal.Iterator to conform to the
	// envoy.Provider interface
	iteratorProvider struct {
		iter dal.Iterator

		// Items related to ref resolution
		// @todo can be removed when reworked
		resolveRefs bool
		relMods     map[string]refModWrap
		dal         dal.FullService

		rows      []datasource.RawRecord
		buffIndex int
		done      bool

		// User ref resolution
		store      store.Storer
		userFields map[string]bool
	}

	refModWrap struct {
		modLvl1   *types.Module
		labelLvl1 string

		modLvl2   *types.Module
		labelLvl2 string
	}
)

const (
	bufferPullChunkSize = int(100)
)

func mkIteratorProvider(ctx context.Context, s store.Storer, dl dal.FullService, iter dal.Iterator, mod *types.Module, resolveRefs bool) (out *iteratorProvider, err error) {
	out = &iteratorProvider{
		iter:        iter,
		dal:         dl,
		resolveRefs: resolveRefs,
	}

	refMods := make(map[string]refModWrap)
	userFields := make(map[string]bool)

	for _, name := range []string{"createdBy", "updatedBy", "ownedBy", "deletedBy"} {
		userFields[name] = true
	}

	for _, f := range mod.Fields {
		switch f.Kind {
		case "Record":
			refMods[f.Name], err = mkRecordRefWrap(ctx, s, f)
			if err != nil {
				return
			}
		case "User":
			userFields[f.Name] = true
		}
	}

	out.relMods = refMods
	out.store = s
	out.userFields = userFields

	return
}

func (rd *RecordDatasource) SetProvider(s envoyx.Provider) bool {
	if rd.Mapping.SourceIdent != s.Ident() {
		return false
	}

	rd.Provider = s
	return true
}

func (rd *RecordDatasource) Next(ctx context.Context, out datasource.RawRecord) (ident []string, more bool, err error) {
	rowCache := make(datasource.RawRecord)

	more, err = rd.Provider.Next(ctx, rowCache)
	if err != nil || !more {
		return
	}

	rd.applyMapping(rowCache, out)

	if len(rd.Mapping.KeyField) == 0 {
		ident = append(ident, strconv.FormatInt(int64(rd.currentIndex), 10))
	} else {
		for _, k := range rd.Mapping.KeyField {
			ident = append(ident, strings.Join(rowCache[k].Values, ","))
		}
	}

	rd.currentIndex++

	return
}

func (rd *iteratorProvider) SetConfigs(map[string]any) error {
	return nil
}

func (rd *RecordDatasource) Reset(ctx context.Context) (err error) {
	rd.currentIndex = 0
	return rd.Provider.Reset(ctx)
}

func (rd *RecordDatasource) applyMapping(in, out datasource.RawRecord) {
	if len(rd.Mapping.Mapping.Map) == 0 {
		if !rd.Mapping.Defaultable {
			return
		}

		for k, v := range in {
			out[k] = v
		}
		return
	}

	if rd.Mapping.Defaultable {
		rd.applyMappingWithDefaults(in, out)
	} else {
		rd.applyMappingWoDefaults(in, out)
	}
}

func (rd *RecordDatasource) applyMappingWithDefaults(in, out datasource.RawRecord) {
	maps := make(map[string]envoyx.MapEntry)
	for k, v := range rd.Mapping.Mapping.Map {
		maps[k] = v
	}

	for k, v := range in {
		if m, ok := maps[k]; ok {
			if m.Skip {
				continue
			}
			out[m.Field] = v
		} else {
			out[k] = v
		}
	}
}

func (rd *RecordDatasource) applyMappingWoDefaults(in, out datasource.RawRecord) {
	for _, m := range rd.Mapping.Mapping.Map {
		if m.Skip {
			continue
		}

		out[m.Field] = in[m.Column]
	}
}

func (rd *RecordDatasource) ResolveRef(ref ...any) (out uint64, err error) {
	idents, err := cast.ToStringSliceE(ref)
	if err != nil {
		return
	}

	for i, ident := range idents {
		idents[i] = strings.Replace(ident, "-", "_", -1)
	}

	out = rd.refToID[strings.Join(idents, "-")]
	return
}

func (rd *RecordDatasource) ResolveRefS(ref ...string) (out uint64, err error) {
	aux := make([]any, len(ref))
	for i, r := range ref {
		aux[i] = r
	}

	return rd.ResolveRef(aux...)
}

// @todo this should be replaced by some smarter structure
func (rd *RecordDatasource) AddRef(id uint64, idents ...string) {
	for i, ident := range idents {
		idents[i] = strings.Replace(ident, "-", "_", -1)
	}

	rd.refToID[strings.Join(idents, "-")] = id
}

func (ip *iteratorProvider) Next(ctx context.Context, out datasource.RawRecord) (more bool, err error) {
	if ip.resolveRefs {
		return ip.nextResolved(ctx, out)
	}

	return ip.next(ctx, out)
}

func (ip *iteratorProvider) next(ctx context.Context, out datasource.RawRecord) (more bool, err error) {
	rowCache := make(datasource.RawRecord)

	if !ip.iter.Next(ctx) {
		return false, ip.iter.Err()
	}

	err = ip.iter.Scan(rowCache)
	if err != nil {
		return
	}

	for k, v := range rowCache {
		out[k] = v
	}

	return true, nil
}

func (ip *iteratorProvider) nextResolved(ctx context.Context, out datasource.RawRecord) (more bool, err error) {
	if ip.done && ip.buffIndex >= len(ip.rows) {
		return false, nil
	}

	if ip.buffIndex == len(ip.rows) {

		// pull chunk
		ip.rows = make([]datasource.RawRecord, 0)

		for i := 0; i < bufferPullChunkSize; i++ {
			rowCache := make(datasource.RawRecord)
			if !ip.iter.Next(ctx) {
				ip.done = true

				err := ip.iter.Err()
				if err != nil {
					return false, err
				}

				break
			}

			err = ip.iter.Scan(rowCache)
			if err != nil {
				return
			}

			ip.rows = append(ip.rows, rowCache)
		}

		// resolve stuff
		err = ip.resolveReferences(ctx, ip.dal)
		if err != nil {
			return
		}

		err = ip.resolveUsers(ctx, ip.store)
		if err != nil {
			return
		}
	}

	rowCache := ip.rows[ip.buffIndex]
	ip.buffIndex++

	for k, v := range rowCache {
		out[k] = v
	}

	return true, nil
}

func (ip *iteratorProvider) resolveReferences(ctx context.Context, ds dal.FullService) (err error) {
	for refField, refWrap := range ip.relMods {
		if refWrap.modLvl1 == nil {
			continue
		}

		var labels map[string]string
		labels, err = refLabels(ctx, ds, refWrap, refIDs(ip.rows, refField))
		if err != nil {
			return
		}

		for i, row := range ip.rows {
			for pos, v := range row[refField].Values {
				if label, ok := labels[v]; ok {
					row.SetValue(fmt.Sprintf("%s value", refField), uint(pos), label)
				}
			}

			ip.rows[i] = row
		}
	}

	return
}

// refIDs collects every record ID the rows reference through the field
func refIDs(rows []datasource.RawRecord, field string) (out []uint64) {
	seen := make(map[uint64]bool)
	for _, row := range rows {
		for _, v := range row[field].Values {
			if id := cast.ToUint64(v); id > 0 && !seen[id] {
				seen[id] = true
				out = append(out, id)
			}
		}
	}

	return
}

// refLabels maps each referenced record ID to its label, following
// recordLabelField one level down when the label field is itself a reference
func refLabels(ctx context.Context, ds dal.FullService, ref refModWrap, ids []uint64) (out map[string]string, err error) {
	lvl1, err := recordsByID(ctx, ds, ref.modLvl1, ids)
	if err != nil || len(lvl1) == 0 {
		return
	}

	out = make(map[string]string, len(lvl1))

	if ref.modLvl2 == nil {
		for _, rec := range lvl1 {
			if v := rec.Values.Get(ref.labelLvl1, 0); v != nil {
				out[strconv.FormatUint(rec.ID, 10)] = v.Value
			}
		}
		return
	}

	nestedIDs := make([]uint64, 0, len(lvl1))
	for _, rec := range lvl1 {
		if v := rec.Values.Get(ref.labelLvl1, 0); v != nil {
			nestedIDs = append(nestedIDs, cast.ToUint64(v.Value))
		}
	}

	lvl2, err := recordsByID(ctx, ds, ref.modLvl2, nestedIDs)
	if err != nil {
		return
	}

	nested := make(map[string]string, len(lvl2))
	for _, rec := range lvl2 {
		if v := rec.Values.Get(ref.labelLvl2, 0); v != nil {
			nested[strconv.FormatUint(rec.ID, 10)] = v.Value
		}
	}

	for _, rec := range lvl1 {
		if v := rec.Values.Get(ref.labelLvl1, 0); v != nil {
			if label, ok := nested[v.Value]; ok {
				out[strconv.FormatUint(rec.ID, 10)] = label
			}
		}
	}

	return
}

func recordsByID(ctx context.Context, ds dal.FullService, mod *types.Module, ids []uint64) (types.RecordSet, error) {
	if mod == nil || len(ids) == 0 {
		return nil, nil
	}

	rr, _, err := dalutils.ComposeRecordsList(ctx, ds, mod, types.RecordFilter{
		RecordID: ids,
		Paging:   filter.Paging{Limit: uint(len(ids))},
	})

	return rr, err
}

// @todo consider omitting these from the interface since they're not always needed
func (ip *iteratorProvider) Reset(ctx context.Context) (err error) {
	return
}

// @todo consider omitting these from the interface since they're not always needed
func (ip *iteratorProvider) Ident() (out string) {
	return
}

// @todo consider omitting these from the interface since they're not always needed
func (ip *iteratorProvider) SetIdent(string) {
}

func mkRecordRefWrap(ctx context.Context, s store.Storer, f *types.ModuleField) (wrap refModWrap, err error) {
	relModID := f.Options.UInt64("moduleID")
	if relModID == 0 {
		return
	}

	var relMod *types.Module
	relMod, err = s.LookupComposeModuleByID(ctx, relModID)
	if err != nil {
		return
	}

	relMod.Fields, _, err = s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{
		ModuleID: []uint64{relMod.ID},
	})
	if err != nil {
		return
	}

	wrap = refModWrap{
		modLvl1:   relMod,
		labelLvl1: f.Options.String("labelField"),
	}

	if f.Options.String("recordLabelField") != "" {
		nestedRef := wrap.modLvl1.Fields.FindByName(f.Options.String("labelField"))
		nestedModID := nestedRef.Options.UInt64("moduleID")

		wrap.modLvl2, err = s.LookupComposeModuleByID(ctx, nestedModID)
		if err != nil {
			return
		}
		wrap.labelLvl2 = f.Options.String("recordLabelField")
	}

	return
}

func (ip *iteratorProvider) resolveUsers(ctx context.Context, s store.Storer) (err error) {
	if len(ip.userFields) == 0 || s == nil {
		return
	}

	seen := make(map[string]bool)
	for _, row := range ip.rows {
		for fieldName := range ip.userFields {
			for _, val := range row[fieldName].Values {
				if val != "" && val != "0" {
					seen[val] = true
				}
			}
		}
	}
	if len(seen) == 0 {
		return
	}

	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}

	uu, _, err := store.SearchUsers(ctx, s, systemTypes.UserFilter{
		UserID: ids,
		Paging: filter.Paging{Limit: 0},
	})
	if err != nil {
		return
	}

	labels := make(map[string]string, len(uu))
	for _, u := range uu {
		labels[strconv.FormatUint(u.ID, 10)] = userLabel(u)
	}

	for i, row := range ip.rows {
		for fieldName := range ip.userFields {
			v := row[fieldName]
			if len(v.Values) == 0 {
				continue
			}
			for j, val := range v.Values {
				if label, ok := labels[val]; ok {
					row.SetValue(fmt.Sprintf("%s value", fieldName), uint(j), label)
				}
			}
		}
		ip.rows[i] = row
	}
	return
}

func userLabel(u *systemTypes.User) string {
	if u.Handle != "" {
		return u.Handle
	}
	if u.Email != "" {
		return u.Email
	}
	if u.Name != "" {
		return u.Name
	}
	return strconv.FormatUint(u.ID, 10)
}
