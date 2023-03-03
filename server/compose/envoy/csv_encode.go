package envoy

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
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

	out := make(map[string]string)
	header := make([]string, 0, 4)

	row := make([]string, 0, 4)
	var more bool
	for {
		_, more, err = rds.Next(ctx, out)
		if err != nil || !more {
			return
		}

		if len(header) == 0 {
			for k := range out {
				header = append(header, k)
			}
			err = writer.Write(header)
			if err != nil {
				return
			}
		}

		for _, h := range header {
			row = append(row, out[h])
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
