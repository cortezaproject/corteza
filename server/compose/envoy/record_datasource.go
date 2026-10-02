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
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/store"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/modern-go/reflect2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type (
	// RecordDatasource provides a mechanism for you to access large
	// record datasets optimally
	RecordDatasource struct {
		Mapping  envoyx.DatasourceMapping
		Provider envoyx.Provider

		multivalues    map[string]bool
		datetimeFields map[string]bool

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
		resolveRefs  bool
		includeRefID bool
		relMods      map[string]refModWrap
		dal          dal.FullService

		// Access control for field-level permissions
		ac             recordValueAccessController
		mod            *types.Module
		readableFields map[string]bool

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

	// recordValueAccessController interface for checking field-level read permissions
	recordValueAccessController interface {
		CanReadRecordValueOnModuleField(context.Context, *types.ModuleField) bool
	}
)

const (
	bufferPullChunkSize = int(100)
)

func mkIteratorProvider(ctx context.Context, ac recordValueAccessController, s store.Storer, dl dal.FullService, iter dal.Iterator, mod *types.Module, resolveRefs bool, includeRefID bool) (out *iteratorProvider, err error) {
	out = &iteratorProvider{
		iter:         iter,
		dal:          dl,
		resolveRefs:  resolveRefs,
		includeRefID: includeRefID,

		ac:  ac,
		mod: mod,
	}

	// Pre-compute readable fields based on access control permissions
	// This avoids checking permissions on every record iteration
	if !reflect2.IsNil(ac) {
		out.readableFields = make(map[string]bool)
		for _, f := range mod.Fields {
			// Ensure NamespaceID is set on the field (it may not be loaded from store)
			// RBAC checks require the full resource path including namespace
			if f.NamespaceID == 0 {
				f.NamespaceID = mod.NamespaceID
			}
			out.readableFields[f.Name] = ac.CanReadRecordValueOnModuleField(ctx, f)
		}
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

// ResolvedFields returns field names that will have a companion "X ID" column when includeRefID is enabled.
func (ip *iteratorProvider) ResolvedFields() []string {
	if !ip.resolveRefs || !ip.includeRefID {
		return nil
	}

	out := make([]string, 0, len(ip.userFields)+len(ip.relMods))
	for f := range ip.userFields {
		out = append(out, f)
	}

	for f := range ip.relMods {
		out = append(out, f)
	}

	return out
}

// ResolvedFields returns field names that will have a companion "X value" column.
func (rd *RecordDatasource) ResolvedFields() []string {
	if ip, ok := rd.Provider.(*iteratorProvider); ok {
		return ip.ResolvedFields()
	}
	return nil
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

	// Apply field-level access control filtering
	ip.filterUnreadableFields(rowCache)

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
		ip.buffIndex = 0

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

		if len(ip.rows) == 0 {
			return false, nil
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

	// Apply field-level access control filtering
	ip.filterUnreadableFields(rowCache)

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

			// Build a map from recordID to resolved label for ordered assignment
			labelByID := make(map[string]string, len(relRecords))
			for _, rec := range relRecords {
				if rec.Values == nil {
					continue
				}
				v := rec.Values.Get(resLab, 0)
				if v == nil {
					continue
				}
				labelByID[strconv.FormatUint(rec.ID, 10)] = v.Value
			}

			for j, origID := range value.Values {
				label, ok := labelByID[origID]
				if !ok {
					continue
				}
				if ip.includeRefID {
					cacheRecord.SetValue(fmt.Sprintf("%s ID", refField), uint(j), origID)
				}
				cacheRecord.SetValue(refField, uint(j), label)
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

func mkRecordRefWrap(ctx context.Context, s store.Storer, f *types.ModuleField) (wrap refModWrap, err error) {
	relModID := f.Options.UInt64("moduleID")
	if relModID == 0 {
		return
	}

	var relMod *types.Module
	relMod, err = s.LookupComposeModuleByID(ctx, relModID)
	if err != nil {
		// Referenced module may have been deleted or made inaccessible since the
		// field was configured. Skip ref resolution for this field instead of
		// failing the entire export — the column is still emitted, just unresolved.
		if errors.IsNotFound(err) {
			logger.Default().Warn(
				"skipping record ref resolution: referenced module not found",
				zap.String("field", f.Name),
				zap.Uint64("moduleID", relModID),
			)
			err = nil
		}
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
		// Two-level label resolution: the labelField on the related module must
		// itself be a Record field pointing to another module. If that chain is
		// broken (labelField renamed/removed, target module missing), skip the
		// nested resolution instead of aborting the export.
		labelFieldName := f.Options.String("labelField")
		nestedRef := wrap.modLvl1.Fields.FindByName(labelFieldName)
		if nestedRef == nil {
			logger.Default().Warn(
				"skipping nested record ref resolution: labelField not found on related module",
				zap.String("field", f.Name),
				zap.String("labelField", labelFieldName),
				zap.Uint64("relatedModuleID", relMod.ID),
			)
			return
		}
		nestedModID := nestedRef.Options.UInt64("moduleID")
		if nestedModID == 0 {
			logger.Default().Warn(
				"skipping nested record ref resolution: labelField has no moduleID option",
				zap.String("field", f.Name),
				zap.String("labelField", labelFieldName),
			)
			return
		}

		wrap.modLvl2, err = s.LookupComposeModuleByID(ctx, nestedModID)
		if err != nil {
			if errors.IsNotFound(err) {
				logger.Default().Warn(
					"skipping nested record ref resolution: nested module not found",
					zap.String("field", f.Name),
					zap.Uint64("nestedModuleID", nestedModID),
				)
				err = nil
			}
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

	// Collect unique user IDs across the current chunk
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
		UserID:    ids,
		Paging:    filter.Paging{Limit: 0},
		Deleted:   filter.StateInclusive,
		Suspended: filter.StateInclusive,
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
				if val == "0" || val == "" {
					row.SetValue(fieldName, uint(j), "")
					continue
				}
				if label, ok := labels[val]; ok {
					if ip.includeRefID {
						row.SetValue(fmt.Sprintf("%s ID", fieldName), uint(j), val)
					}
					row.SetValue(fieldName, uint(j), label)
				}
			}
		}
		ip.rows[i] = row
	}
	return
}

func userLabel(u *systemTypes.User) string {
	if u.Name != "" {
		return u.Name
	}
	if u.Email != "" {
		return u.Email
	}
	if u.Handle != "" {
		return u.Handle
	}
	return strconv.FormatUint(u.ID, 10)
}

// filterUnreadableFields removes values the given user does not have access to
func (ip *iteratorProvider) filterUnreadableFields(rowCache datasource.RawRecord) {
	if reflect2.IsNil(ip.ac) {
		return
	}

	for fieldName := range rowCache {
		// Skip system fields (ID, createdAt, etc.) - these are always readable
		// System fields are lowercase and don't have associated module field definitions
		if _, hasField := ip.readableFields[fieldName]; !hasField {
			// Field not in module definition, likely a system field - keep it
			continue
		}

		// Check if user can read this field
		if !ip.readableFields[fieldName] {
			delete(rowCache, fieldName)
		}
	}
}
