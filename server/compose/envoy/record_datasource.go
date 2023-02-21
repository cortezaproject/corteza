package envoy

import (
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
)

func (rd *RecordDatasource) SetProvider(s envoyx.Provider) bool {
	if rd.mapping.SourceIdent != s.Ident() {
		return false
	}

	rd.provider = s
	return true
}

func (rd *RecordDatasource) Next(out map[string]string) (ident string, more bool, err error) {
	if rd.rowCache == nil {
		rd.rowCache = make(map[string]string)
	}

	more, err = rd.provider.Next(rd.rowCache)
	if err != nil || !more {
		return
	}

	rd.applyMapping(rd.rowCache, out)

	ident = out[rd.mapping.KeyField]

	return
}

func (rd *RecordDatasource) Reset() (err error) {
	return rd.provider.Reset()
}

func (rd *RecordDatasource) applyMapping(in, out map[string]string) {
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
