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
		buffer *bytes.Buffer
		name   string
		row    []string
		header []string
		writer *csv.Writer
	}
)

// this function splits the stream of the given migrateable node.
// See readme for more info
func splitStream(m types.Migrateable) ([]types.Migrateable, error) {
	var rr []types.Migrateable
	rr = append(rr, m)
	if m.Map == nil {
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

	// @fix this hack will not always work.
	// replace with a set or something similar
	i := -1
	for {
		i++

		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// find first applicable map, that can be used for the given row.
		// default maps should not inclide a where field
		for _, strmp := range streamMap {
			if checkWhere(strmp["where"], record, hMap) {
				maps, ok := strmp["map"].([]interface{})
				if !ok {
					return nil, errors.New("streamMap.invalidMap")
				}

				// populate splitted streams
				for _, mp := range maps {
					mm, ok := mp.(map[string]interface{})
					if !ok {
						return nil, errors.New("streamMap.map.invalidEntry")
					}

					from, ok := mm["from"].(string)
					if !ok {
						return nil, errors.New("streamMap.map.entry.invalidFrom")
					}

					to, ok := mm["to"].(string)
					if !ok {
						return nil, errors.New("streamMap.map.invalidTo")
					}

					vv := strings.Split(to, ".")
					nm := vv[0]
					nmF := vv[1]

					if bufs[nm] == nil {
						var bb bytes.Buffer
						ww := csv.NewWriter(&bb)
						defer ww.Flush()
						bufs[nm] = &SplitBuffer{
							buffer: &bb,
							writer: ww,
							name:   nm,
						}
					}

					bufs[nm].row = append(bufs[nm].row, record[hMap[from]])
					if i == 0 {
						bufs[nm].header = append(bufs[nm].header, nmF)
					}
				}
			}

			// write csv rows
			for _, v := range bufs {
				v.writer.Write(v.row)
				var nn []string
				v.row = nn
			}
		}
	}

	// make migrateable nodes from the generated streams
	for _, v := range bufs {
		rr = append(rr, types.Migrateable{
			Name:   v.name,
			Source: v.buffer,
			Header: &v.header,
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
