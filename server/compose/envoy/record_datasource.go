package envoy

import (
	"context"

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

func (rd *RecordDatasource) Next(ctx context.Context, out map[string]string) (ident string, more bool, err error) {
	if rd.rowCache == nil {
		rd.rowCache = make(map[string]string)
	}

	more, err = rd.provider.Next(ctx, rd.rowCache)
	if err != nil || !more {
		return
	}

	rd.applyMapping(rd.rowCache, out)

	ident = out[rd.mapping.KeyField]

	return
}

func (rd *RecordDatasource) Reset(ctx context.Context) (err error) {
	return rd.provider.Reset(ctx)
}

func (rd *RecordDatasource) applyMapping(in, out map[string]string) {
	if len(rd.mapping.Mapping.m) == 0 {
		for k, v := range in {
			out[k] = v
		}
		return
	}

	for _, m := range rd.mapping.Mapping.m {
		if m.Skip {
			continue
		}

		// @todo expand when needed (expressions and such)
		out[m.Field] = in[m.Column]
	}
}

func (rd *RecordDatasource) ResolveRef(ref any) (out uint64, err error) {
	r, err := cast.ToStringE(ref)
	if err != nil {
		return
	}

	out = rd.refToID[r]
	return
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
