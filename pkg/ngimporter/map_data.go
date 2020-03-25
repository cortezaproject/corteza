package ngimporter

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/ngimporter/types"
)

type (
	// MapBuffer helps us keep track of new SourceNodes, defined by the map operation.
	MapBuffer struct {
		buffer    *bytes.Buffer
		name      string
		row       []string
		header    []string
		hasHeader bool
		writer    *csv.Writer

		// field: masterID: [value]
		joins map[string]map[string][]string
	}
)

// maps data from the original ImportSource node into new ImportSource nodes
// based on the provided DataMap.
// Algorithm outline:
//   * parse data map
//   * for each record in the original import source, based on the map, create new
//     MapBuffers
//   * create new import sources based on map buffers
func mapData(is types.ImportSource) ([]types.ImportSource, error) {
	var rr []types.ImportSource
	if is.DataMap == nil {
		rr = append(rr, is)
		return rr, nil
	}

	// unpack the map
	// @todo provide a better structure!!
	var dataMap []map[string]interface{}
	src, _ := ioutil.ReadAll(is.DataMap)
	err := json.Unmarshal(src, &dataMap)
	if err != nil {
		return nil, err
	}

	// get header fields
	r := csv.NewReader(is.Source)
	header, err := r.Read()
	if err == io.EOF {
		return rr, nil
	}

	if err != nil {
		return nil, err
	}

	// maps { header field: field index } for a nicer lookup
	hMap := make(map[string]int)
	for i, h := range header {
		hMap[h] = i
	}

	bufs := make(map[string]*MapBuffer)

	// data mapping
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// on next row, currently acquired headers are marked as final
		for _, b := range bufs {
			b.hasHeader = true
		}

		// find applicable maps, that can be used for the given row.
		// the system allows composition, so all applicable maps are used.
		for _, strmp := range dataMap {
			if ok, err := checkWhere(strmp["where"], record, hMap); ok && err == nil {
				maps, ok := strmp["map"].([]interface{})
				if !ok {
					return nil, errors.New("dataMap.invalidMap " + is.Name)
				}

				// handle current record and it's values
				for _, mp := range maps {
					mm, ok := mp.(map[string]interface{})
					if !ok {
						return nil, errors.New("dataMap.map.invalidEntry " + is.Name)
					}

					from, ok := mm["from"].(string)
					if !ok {
						return nil, errors.New("dataMap.map.entry.invalidFrom " + is.Name)
					}

					to, ok := mm["to"].(string)
					if !ok {
						return nil, errors.New("dataMap.map.invalidTo " + is.Name)
					}

					vv := strings.Split(to, ".")
					nm := vv[0]
					nmF := vv[1]

					if bufs[nm] == nil {
						var bb bytes.Buffer
						ww := csv.NewWriter(&bb)
						defer ww.Flush()
						bufs[nm] = &MapBuffer{
							buffer:    &bb,
							writer:    ww,
							name:      nm,
							hasHeader: false,
						}
					}

					val := record[hMap[from]]

					// handle data join
					if strings.Contains(from, ".") {
						// construct a `alias.joinOnID` value, so we can perform a simple map lookup
						pts := strings.Split(from, ".")
						baseFieldAlias := pts[0]
						originalOn := is.AliasMap[baseFieldAlias]
						joinField := pts[1]

						oo := []string{}
						for _, ff := range originalOn {
							oo = append(oo, record[hMap[ff]])
						}
						val = baseFieldAlias + "." + strings.Join(oo[:], ".")

						// modify header field to specify what joined node field to use
						nmF += ":" + joinField
					}

					bufs[nm].row = append(bufs[nm].row, val)
					if !bufs[nm].hasHeader {
						bufs[nm].header = append(bufs[nm].header, nmF)
					}
				}
			} else if err != nil {
				return nil, err
			}
		}

		// write csv rows
		for _, v := range bufs {
			if len(v.row) > 0 {
				v.writer.Write(v.row)
				v.row = []string{}
			}
		}
	}

	// construct output import source nodes
	for _, v := range bufs {
		rr = append(rr, types.ImportSource{
			Name:     v.name,
			Source:   v.buffer,
			Header:   &v.header,
			FieldMap: is.FieldMap,
			AliasMap: is.AliasMap,
			ValueMap: is.ValueMap,
		})
	}

	return rr, nil
}

// checks if the given condition passes for the given row
func checkWhere(where interface{}, row []string, hMap map[string]int) (bool, error) {
	if where == nil {
		return true, nil
	}

	wh, ok := where.(string)
	if !ok {
		return true, nil
	}

	ev, err := types.ExprLang.NewEvaluable(wh)
	if err != nil {
		return false, err
	}

	// prep payload
	pr := map[string]string{}
	for k, v := range hMap {
		pr[k] = row[v]
	}

	rr, err := ev.EvalBool(context.Background(), pr)
	if err != nil {
		return false, err
	}

	return rr, nil
}
