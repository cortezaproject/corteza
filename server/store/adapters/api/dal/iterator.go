package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers"
)

type (
	iterator struct {
		// source model; how data we are reading from is shaped
		src *model

		// destination model; how data we are reading into is shaped
		// this is used to create scan buffer
		// when not doing plain selection from one table final results might
		// require a different list of columns and scanning needs to be adjusted
		dst *model

		method   string
		endpoint string
		query    map[string]string

		// payload payloadBits

		i    int
		rows [][]byte

		// last error
		err error

		sorting filter.SortExprSet
		cursor  *filter.PagingCursor
		limit   uint
	}
)

func (i *iterator) Next(ctx context.Context) bool {
	if i.err == nil && i.rows == nil {
		// @todo improve buffering
		i.i--
		i.rows, i.err = i.fetch(ctx)
	}

	if i.err != nil {
		return false
	}

	i.i++
	if i.i >= len(i.rows) {
		return false
	}

	return true
}

// More fetches more records from the point of last record
func (i *iterator) More(max uint, last dal.ValueGetter) (err error) {
	// @todo can we support this?
	return
}

func (i *iterator) Preload(_ context.Context, max uint, cur *filter.PagingCursor) (err error) {
	// @todo do we need to support this?
	return
}

func (i *iterator) Sorting() filter.SortExprSet {
	return i.sorting
}

func (i *iterator) fetch(ctx context.Context) (rows [][]byte, err error) {
	// Proc and prep payload stuff
	xr := drivers.XRequest{}
	xr, err = i.prepSort(xr, i.sorting)
	if err != nil {
		return
	}

	xr, err = i.prepLimit(xr, i.limit)
	if err != nil {
		return
	}

	// Run the request
	data, err := i.src.conn.Run(ctx, i.method, i.dst.dialect.EncrichEndpoint(i.endpoint, xr), nil, nil)
	if err != nil {
		return
	}

	// Unpack the response
	// @todo process meta
	rows, meta, err := i.unpackResponse(data)
	_ = meta
	if err != nil {
		return
	}

	return rows, nil
}

func (i *iterator) unpackResponse(data []byte) (rows [][]byte, meta []byte, err error) {
	// @todo should we push this down to the dialect?
	rows, meta, _ = UnpackRowsWithPaths(
		data,
		strings.Split(i.dst.dialect.SearchDataPath(), "."),
		strings.Split(i.dst.dialect.SearchMetaPath(), "."),
	)

	return
}

func UnpackRowsWithPaths(
	b []byte,
	dataPath []string,
	metaPath []string,
) (rows [][]byte, meta []byte, err error) {

	// 1. Root array case
	if len(b) > 0 && b[0] == '[' {
		var rawRows []json.RawMessage
		if err = json.Unmarshal(b, &rawRows); err != nil {
			return nil, nil, err
		}

		for _, r := range rawRows {
			rows = append(rows, []byte(r))
		}
		return rows, nil, nil
	}

	// 2. Root object
	var root map[string]json.RawMessage
	if err = json.Unmarshal(b, &root); err != nil {
		return nil, nil, err
	}

	// 3. Extract data array
	dataRaw, ok := getRawAtPath(root, dataPath)
	if !ok {
		return nil, nil, fmt.Errorf("data path not found: %v", dataPath)
	}

	var rawRows []json.RawMessage
	if err = json.Unmarshal(dataRaw, &rawRows); err != nil {
		return nil, nil, fmt.Errorf("data is not an array")
	}

	for _, r := range rawRows {
		rows = append(rows, []byte(r))
	}

	// 4. Extract meta (optional)
	if metaPath != nil {
		meta, _ = getRawAtPath(root, metaPath)
	}

	return rows, meta, nil
}

func getRawAtPath(
	root map[string]json.RawMessage,
	path []string,
) (json.RawMessage, bool) {

	current, ok := root[path[0]]
	if !ok {
		return nil, false
	}

	for _, p := range path[1:] {
		var m map[string]json.RawMessage
		if err := json.Unmarshal(current, &m); err != nil {
			return nil, false
		}
		current, ok = m[p]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

func (i *iterator) prepSort(req drivers.XRequest, sort filter.SortExprSet) (out drivers.XRequest, err error) {
	out = req

	for _, s := range sort {
		// assuming all columns were pre-validated!
		tmp, _ := i.src.table.AttributeExpression(s.Column)

		out, err = i.dst.dialect.AddSort(out, tmp, s.Descending)
		if err != nil {
			return
		}
	}

	return
}

func (i *iterator) prepLimit(req drivers.XRequest, limit uint) (out drivers.XRequest, err error) {
	out = req

	out, err = i.dst.dialect.AddLimit(out, limit)
	if err != nil {
		return
	}

	return
}

func (i *iterator) Scan(r dal.ValueSetter) (err error) {
	row := i.rows[i.i]
	err = i.dst.table.Decode(row, r)
	if err != nil {
		return
	}

	return
}

func (i *iterator) Err() error {
	return i.err
}

// Close iterator and cleanup
func (i *iterator) Close() error {
	// @todo do we need to close things?
	return nil
}

func (i *iterator) BackCursor(r dal.ValueGetter) (cur *filter.PagingCursor, err error) {
	return
}

func (i *iterator) ForwardCursor(r dal.ValueGetter) (*filter.PagingCursor, error) {
	return i.collectCursorValues(r)
}

func (i *iterator) collectCursorValues(r dal.ValueGetter) (_ *filter.PagingCursor, err error) {
	return
}
