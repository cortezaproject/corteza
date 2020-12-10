package rdbms

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

func (s Store) convertComposeRecordValueFilter(_ *types.Module, f types.RecordValueFilter) (query squirrel.SelectBuilder, err error) {
	// Always filter by record IDs
	query = s.composeRecordValuesSelectBuilder().Where(squirrel.Eq{"crv.record_id": f.RecordID})
	query = filter.StateCondition(query, "crv.deleted_at", f.Deleted)

	return query, nil
}

func (s Store) ComposeRecordValueRefLookup(ctx context.Context, m *types.Module, field string, ref uint64) (uint64, error) {
	q := s.composeRecordValuesSelectBuilder().
		Join(s.composeRecordTable("crd"), "crv.record_id = crd.id").
		Where(squirrel.Eq{
			"crv.name":       field,
			"crv.ref":        ref,
			"crv.deleted_at": nil,
			"crd.module_id":  m.ID,
			"crd.deleted_at": nil,
		}).
		Column("record_id").
		Limit(1)

	row, err := s.QueryRow(ctx, q)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	var recordID uint64
	if err = row.Scan(&recordID); err != nil {
		return 0, err
	}

	return recordID, nil
}

// PartialComposeRecordValueUpdate updates specific record values across multiple records
func (s Store) PartialComposeRecordValueUpdate(ctx context.Context, m *types.Module, vv ...*types.RecordValue) (err error) {
	{
		// handle standard record-value storage
		for _, v := range vv {
			if err = s.execUpsertComposeRecordValues(ctx, s.internalComposeRecordValueEncoder(v)); err != nil {
				return
			}
		}
	}

	return nil
}
