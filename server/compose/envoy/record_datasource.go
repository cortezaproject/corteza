package envoy

import (
	"context"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/spf13/cast"
)

type (
	// RecordDatasource provides a mechanism for you to access large
	// record datasets optimally
	RecordDatasource struct {
		mapping  datasourceMapping
		provider envoyx.Provider

		// Reusable buffer for reading records
		rowCache map[string]string

		// Index to map from ref to ID
		// @todo we might need to flush these to the disc in case a huge dataset is passed in
		refToID map[string]uint64
	}

	// iteratorProvider is a wrapper around the dal.Iterator to conform to the
	// envoy.Provider interface
	iteratorProvider struct {
		iter dal.Iterator

		rowCache auxRecord
	}

	auxRecord map[string]string
)

func (rd *RecordDatasource) SetProvider(s envoyx.Provider) bool {
	if rd.mapping.SourceIdent != s.Ident() {
		return false
	}

	rd.provider = s
	return true
}

func (rd *RecordDatasource) Next(ctx context.Context, out map[string]string) (ident []string, more bool, err error) {
	if rd.rowCache == nil {
		rd.rowCache = make(map[string]string)
	}

	more, err = rd.provider.Next(ctx, rd.rowCache)
	if err != nil || !more {
		return
	}

	rd.applyMapping(rd.rowCache, out)

	for _, k := range rd.mapping.KeyField {
		ident = append(ident, rd.rowCache[k])
	}

	return
}

func (rd *RecordDatasource) Reset(ctx context.Context) (err error) {
	return rd.provider.Reset(ctx)
}

func (rd *RecordDatasource) applyMapping(in, out map[string]string) {
	if len(rd.mapping.Mapping.m) == 0 {
		if !rd.mapping.Defaultable {
			return
		}

		for k, v := range in {
			out[k] = v
		}
		return
	}

	if rd.mapping.Defaultable {
		rd.applyMappingWithDefaults(in, out)
	} else {
		rd.applyMappingWoDefaults(in, out)
	}
}

func (rd *RecordDatasource) applyMappingWithDefaults(in, out map[string]string) {
	maps := make(map[string]mapEntry)
	for k, v := range rd.mapping.Mapping.m {
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

func (rd *RecordDatasource) applyMappingWoDefaults(in, out map[string]string) {
	for _, m := range rd.mapping.Mapping.m {
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

func (ar auxRecord) SetValue(name string, pos uint, value any) (err error) {
	ar[name] = cast.ToString(value)
	return
}

func (ip *iteratorProvider) Next(ctx context.Context, out map[string]string) (more bool, err error) {
	if ip.rowCache == nil {
		ip.rowCache = make(auxRecord)
	}

	if !ip.iter.Next(ctx) {
		return false, ip.iter.Err()
	}

	err = ip.iter.Scan(ip.rowCache)
	if err != nil {
		return
	}

	for k, v := range ip.rowCache {
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
