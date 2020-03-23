package migrate

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/migrate/types"
)

type (
	SplitBuffer struct {
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

// this function splits the stream of the given migrateable node.
// See readme for more info
func splitStream(m types.Migrateable) ([]types.Migrateable, error) {
	var rr []types.Migrateable
	if m.Map == nil {
		rr = append(rr, m)
		return rr, nil
	}

	// unpack the map
	// @todo provide a better structure!!
	var streamMap []map[string]interface{}
	src, _ := ioutil.ReadAll(m.Map)
	err := json.Unmarshal(src, &streamMap)
	if err != nil {
		return nil, err
	}

	// get header fields
	r := csv.NewReader(m.Source)
	header, err := r.Read()
	if err == io.EOF {
		return rr, nil
	}

	if err != nil {
		return nil, err
	}

	// maps header field -> field index for a nicer lookup
	hMap := make(map[string]int)
	for i, h := range header {
		hMap[h] = i
	}

	bufs := make(map[string]*SplitBuffer)

	// splitting magic

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// on next row, old stream's headers are finished
		for _, b := range bufs {
			b.hasHeader = true
		}

		// find first applicable map, that can be used for the given row.
		// default maps should not inclide a where field
		for _, strmp := range streamMap {
			if checkWhere(strmp["where"], record, hMap) {
				maps, ok := strmp["map"].([]interface{})
				if !ok {
					return nil, errors.New("streamMap.invalidMap " + m.Name)
				}

				// populate splitted streams
				for _, mp := range maps {
					mm, ok := mp.(map[string]interface{})
					if !ok {
						return nil, errors.New("streamMap.map.invalidEntry " + m.Name)
					}

					from, ok := mm["from"].(string)
					if !ok {
						return nil, errors.New("streamMap.map.entry.invalidFrom " + m.Name)
					}

					to, ok := mm["to"].(string)
					if !ok {
						return nil, errors.New("streamMap.map.invalidTo " + m.Name)
					}

					vv := strings.Split(to, ".")
					nm := vv[0]
					nmF := vv[1]

					if bufs[nm] == nil {
						var bb bytes.Buffer
						ww := csv.NewWriter(&bb)
						defer ww.Flush()
						bufs[nm] = &SplitBuffer{
							buffer:    &bb,
							writer:    ww,
							name:      nm,
							hasHeader: false,
						}
					}

					val := record[hMap[from]]

					// handle joins
					if strings.Contains(from, ".") {
						// construct a `alias.joinOnID` value, so we can perform a simple map lookup
						pts := strings.Split(from, ".")
						baseFieldAlias := pts[0]
						originalOn := m.AliasMap[baseFieldAlias]
						joinField := pts[1]
						val = baseFieldAlias + "." + record[hMap[originalOn]]

						// modify header field to specify what joined node field to use
						nmF += ":" + joinField
					}

					bufs[nm].row = append(bufs[nm].row, val)
					if !bufs[nm].hasHeader {
						bufs[nm].header = append(bufs[nm].header, nmF)
					}
				}

				// write csv rows
				for _, v := range bufs {
					v.writer.Write(v.row)
					var nn []string
					v.row = nn
				}
				break
			}
		}
	}

	// make migrateable nodes from the generated streams
	for _, v := range bufs {
		rr = append(rr, types.Migrateable{
			Name:     v.name,
			Source:   v.buffer,
			Header:   &v.header,
			FieldMap: m.FieldMap,
			AliasMap: m.AliasMap,
			ValueMap: m.ValueMap,
		})
	}

	return rr, nil
}

// quick and dirty function to check the map's where condition.
// improve with our QL package
func checkWhere(where interface{}, row []string, hMap map[string]int) bool {
	if where == nil {
		return true
	}

	ww, ok := where.(string)
	if !ok {
		return true
	}

	pts := strings.Split(ww, "=")
	org := pts[0]
	val := pts[1]
	return row[hMap[org]] == val
}
