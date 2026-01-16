package envoy

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/envoyx/datasource"
	"github.com/spf13/cast"
)

type (
	CsvEncoder struct{}
)

func (e CsvEncoder) Encode(ctx context.Context, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tt envoyx.Traverser) (err error) {
	w, err := e.getWriter(p)
	if err != nil {
		return
	}

	cw := csv.NewWriter(w)

	switch rt {
	case ComposeRecordDatasourceAuxType:
		_, err = e.encodeRecordDatasources(ctx, cw, p, nodes, tt)
		if err != nil {
			return
		}
	}

	cw.Flush()
	return
}

func (e CsvEncoder) encodeRecordDatasources(ctx context.Context, writer *csv.Writer, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out any, err error) {
	for _, n := range nodes {
		_, err = e.encodeRecordDatasource(ctx, writer, p, n, tt)
		if err != nil {
			return
		}
	}

	return
}

func (e CsvEncoder) encodeRecordDatasource(ctx context.Context, writer *csv.Writer, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (_ any, err error) {
	rds := node.Datasource.(*RecordDatasource)
	resolved := map[string]string{}

	header := make([]string, 0, 4)

	hasID := false
	for _, m := range p.FieldMapping {
		header = append(header, m.Field)

		hasID = hasID || strings.ToLower(m.Field) == "id"
	}

	if !hasID {
		header = append([]string{"ID"}, header...)
	}

	mvDelimiter := ";"
	wrapBrackets := false

	if p.Params["multiValueDelimiter"] == nil {
		wrapBrackets = false
		mvDelimiter = ";"
	}

	switch cast.ToString(p.Params["multiValueDelimiter"]) {
	case ",":
		mvDelimiter = ","
		wrapBrackets = false
	case ";":
		mvDelimiter = ";"
		wrapBrackets = false
	case "|":
		mvDelimiter = "|"
		wrapBrackets = false
	case "[,]":
		mvDelimiter = ","
		wrapBrackets = true
	case "[;]":
		mvDelimiter = ";"
		wrapBrackets = true
	case "[|]":
		mvDelimiter = "|"
		wrapBrackets = true
	}

	getValue := func(cache datasource.RawRecord, h string) string {
		v := cache[h]
		if len(v.Values) == 0 {
			return ""
		}

		// Encode as mf if the OG field is multi value OR we're on the resolved record column
		src, ok := resolved[h]
		if !rds.multivalues[h] && !(ok && rds.multivalues[src]) {
			return v.Values[0]
		}

		auxv := make([]string, 0, len(v.Values))
		for _, vv := range v.Values {
			if strings.Contains(vv, mvDelimiter) {
				auxv = append(auxv, fmt.Sprintf("\"%s\"", vv))
			} else {
				auxv = append(auxv, vv)
			}
		}

		out := strings.Join(auxv, mvDelimiter)
		if wrapBrackets {
			out = fmt.Sprintf("[%s]", out)
		}

		return out
	}

	row := make([]string, 0, 4)
	var (
		more     bool
		hWritten bool
	)
	for {
		cache := make(datasource.RawRecord)
		_, more, err = rds.Next(ctx, cache)
		if err != nil || !more {
			return
		}

		if len(header) == 0 {
			for k := range cache {
				header = append(header, k)
			}
		}

		if !hWritten {
			hWritten = true

			// Splice in resolved ref values
			for i, h := range header {
				if _, ok := cache[fmt.Sprintf("%s value", h)]; ok {
					header = append(header, "")
					copy(header[i+1:], header[i:])
					header[i] = fmt.Sprintf("%s value", h)

					resolved[header[i]] = h
				}
			}

			err = writer.Write(header)
			if err != nil {
				return
			}
		}

		for _, h := range header {
			row = append(row, getValue(cache, h))
		}

		err = writer.Write(row)
		if err != nil {
			return
		}

		row = nil
	}
}

func (e CsvEncoder) getWriter(p envoyx.EncodeParams) (out io.Writer, err error) {
	aux, ok := p.Params[paramsKeyWriter]
	if ok {
		out, ok = aux.(io.Writer)
		if ok {
			return
		}
	}

	err = fmt.Errorf("csv encoder expects a writer conforming to io.Writer interface")
	return
}
