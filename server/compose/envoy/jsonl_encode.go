package envoy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/envoyx/datasource"
	"github.com/cortezaproject/corteza/server/pkg/j7s"
)

type (
	JsonlEncoder struct{}
)

func (e JsonlEncoder) Encode(ctx context.Context, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tt envoyx.Traverser) (err error) {
	w, err := e.getWriter(p)
	if err != nil {
		return
	}

	jw := json.NewEncoder(w)

	switch rt {
	case ComposeRecordDatasourceAuxType:
		_, err = e.encodeRecordDatasources(ctx, jw, p, nodes, tt)
		if err != nil {
			return
		}
	}

	return
}

func (e JsonlEncoder) encodeRecordDatasources(ctx context.Context, writer *json.Encoder, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out any, err error) {
	for _, n := range nodes {
		_, err = e.encodeRecordDatasource(ctx, writer, p, n, tt)
		if err != nil {
			return
		}
	}

	return
}

func (e JsonlEncoder) encodeRecordDatasource(ctx context.Context, writer *json.Encoder, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (_ any, err error) {
	rds := node.Datasource.(*RecordDatasource)
	header := make([]string, 0, 4)

	hasID := false
	for _, m := range p.FieldMapping {
		header = append(header, m.Field)

		hasID = hasID || strings.ToLower(m.Field) == "id"
	}

	if !hasID {
		header = append([]string{"ID"}, header...)
	}

	var more bool
	for {
		out := make(datasource.RawRecord)
		_, more, err = rds.Next(ctx, out)
		if err != nil || !more {
			break
		}

		if len(header) == 0 {
			for k := range out {
				header = append(header, k)
			}
		}

		row, _ := j7s.MakeMap()
		for _, h := range header {
			v := out[h]
			if len(v.Values) == 0 {
				continue
			}

			if !rds.multivalues[h] {
				row, err = j7s.AddMap(row, h, v.Values[0])
				if err != nil {
					return
				}
			} else {
				auxv, _ := j7s.MakeSeq()
				for _, vv := range v.Values {
					auxv, err = j7s.AddSeq(auxv, vv)
					if err != nil {
						return
					}
				}

				row, err = j7s.AddMap(row, h, auxv)
				if err != nil {
					return
				}
			}
		}

		err = writer.Encode(row)
		if err != nil {
			return
		}
	}

	return
}

func (e JsonlEncoder) getWriter(p envoyx.EncodeParams) (out io.Writer, err error) {
	aux, ok := p.Params[paramsKeyWriter]
	if ok {
		out, ok = aux.(io.Writer)
		if ok {
			return
		}
	}

	err = fmt.Errorf("jsonl encoder expects a writer conforming to io.Writer interface")
	return
}
