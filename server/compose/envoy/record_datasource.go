package envoy

import (
	"context"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/envoyx/datasource"
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
	}
)

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
