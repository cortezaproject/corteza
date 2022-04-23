package crs

import (
	"context"
	"database/sql"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/google/uuid"
)

type (
	iterator struct {
		ms   *crs
		rows *sql.Rows
		err  error
		next *types.Record

		// @todo should filter also be here?

		// buffer for scanned rows
		scanBuf []any
	}
)

func (i *iterator) Next(ctx context.Context) bool {
	if i.err == nil && i.rows == nil {
		i.rows, i.err = i.fetch(ctx)
	}

	if i.err != nil {
		return false
	}

	return i.rows.Next()
}

func (i *iterator) fetch(ctx context.Context) (*sql.Rows, error) {
	if i.scanBuf == nil {
		// we're going to init scan buffer only once
		// and rely on the sql.Rows.Scan function to
		// fill it up with fresh values!
		i.scanBuf = newScanBuffer(i.ms)
	}

	q := i.ms.selectSql()

	// todo filter?
	// todo cursor constraints?
	sql, args, err := q.ToSQL()
	if err != nil {
		return nil, err
	}

	return i.ms.conn.QueryContext(ctx, sql, args...)
}

func (i *iterator) Scan(r *types.Record) (err error) {
	if err = i.rows.Scan(i.scanBuf...); err != nil {
		return
	}

	return decode(i.ms.columns, i.scanBuf, r)
}

func (i *iterator) Err() error {
	return i.err
}

// Close iterator and cleanup
func (i *iterator) Close() error {
	return i.rows.Close()
}

// prepare slice of (typed) placeholders we're going to
// push to be used by rows.Scan
func newScanBuffer(d *crs) []any {
	var (
		vv = make([]any, len(d.columns))
	)

	for i, c := range d.columns {
		switch c.columnType.(type) {
		case *data.TypeID, *data.TypeRef:
			vv[i] = new(uint64)
		case *data.TypeTimestamp, *data.TypeDate:
			vv[i] = new(sql.NullTime)
		case *data.TypeTime:
			vv[i] = new(sql.NullString)
		case *data.TypeNumber:
			vv[i] = new(sql.NullString)
		case *data.TypeText:
			vv[i] = new(sql.NullString)
		case *data.TypeBoolean:
			vv[i] = new(sql.NullBool)
		case *data.TypeEnum:
			vv[i] = new(sql.NullString)
		case *data.TypeJSON, *data.TypeGeometry, *data.TypeBlob:
			vv[i] = new(sql.RawBytes)
		case *data.TypeUUID:
			vv[i] = new(uuid.UUID)
		}

		i++
	}

	return vv
}

func CursorBackward(r *types.Record, sort filter.SortExprSet) *filter.PagingCursor { return nil }
func CursorForward(r *types.Record, sort filter.SortExprSet) *filter.PagingCursor  { return nil }
func Drain(*iterator) ([]*types.Record, error)                                     { return nil, nil }
