package envoy

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/envoyx/datasource"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
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

	// Get ref record fields and stuff
	refMods := make(map[string]refModWrap)

	for _, f := range mod.Fields {
		if f.Kind != "Record" {
			continue
		}

		relModID := f.Options.UInt64("moduleID")
		var relMod *types.Module
		relMod, err = s.LookupComposeModuleByID(ctx, relModID)
		if err != nil {
			return
		}

		relMod.Fields, _, err = s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{
			ModuleID: []uint64{relMod.ID},
		})

		wrap := refModWrap{
			modLvl1:   relMod,
			labelLvl1: f.Options.String("labelField"),
		}

		if f.Options.String("recordLabelField") != "" {
			nestedRef := wrap.modLvl1.Fields.FindByName(f.Options.String("labelField"))
			nestedModID := nestedRef.Options.UInt64("moduleID")

			relNestedMod, err := s.LookupComposeModuleByID(ctx, nestedModID)
			if err != nil {
				return nil, err
			}

			wrap.labelLvl2 = f.Options.String("recordLabelField")
			wrap.modLvl2 = relNestedMod
		}

		refMods[f.Name] = wrap
	}

	out.relMods = refMods
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
	}

	rowCache := ip.rows[ip.buffIndex]
	ip.buffIndex++

	for k, v := range rowCache {
		out[k] = v
	}

	return true, nil
}

func (ip *iteratorProvider) resolveReferences(ctx context.Context, ds dal.FullService) (err error) {
	// @todo I'll need to chunk these up for to reduce DB query count...

	for i, cacheRecord := range ip.rows {

		for refField, refWrap := range ip.relMods {
			value := cacheRecord[refField]
			if len(value.Values) == 0 {
				continue
			}

			aux := []string{}
			for _, v := range value.Values {
				aux = append(aux, fmt.Sprintf("recordID=%s", v))
			}

			var relRecords types.RecordSet

			resLab := refWrap.labelLvl1
			qq := fmt.Sprintf("(%s)", strings.Join(aux, " OR "))
			relRecords, _, err = dalutils.ComposeRecordsList(ctx, ds, refWrap.modLvl1, types.RecordFilter{
				Query: qq,
				Paging: filter.Paging{
					Limit: uint(len(aux)),
				},
			})

			if err != nil {
				return err
			}

			// lvl 2 nesting; current max lvl
			if refWrap.modLvl2 != nil {
				resLab = refWrap.labelLvl2

				// Iterate related records and collect lvl 2 identifiers
				aux := []string{}
				for _, rec := range relRecords {
					for _, v := range rec.Values.FilterByName(refWrap.labelLvl1) {
						aux = append(aux, fmt.Sprintf("recordID=%s", v.Value))
					}
				}

				qq := fmt.Sprintf("(%s)", strings.Join(aux, " OR "))
				relRecords, _, err = dalutils.ComposeRecordsList(ctx, ds, refWrap.modLvl2, types.RecordFilter{
					Query: qq,
					Paging: filter.Paging{
						Limit: uint(len(aux)),
					},
				})

				if err != nil {
					return err
				}
			}

			for i, rec := range relRecords {
				if rec.Values == nil {
					continue
				}

				v := rec.Values.Get(resLab, 0)
				if v == nil {
					continue
				}

				cacheRecord.SetValue(fmt.Sprintf("%s value", refField), uint(i), v.Value)
			}

			ip.rows[i] = cacheRecord
		}
	}

	return
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
