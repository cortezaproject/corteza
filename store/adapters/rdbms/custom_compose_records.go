package rdbms

import (
	"context"

	composeType "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/report"
)

func (s Store) ComposeRecordReport(ctx context.Context, mod *composeType.Module, metrics string, dimensions string, filters string) ([]map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (s Store) ComposeRecordDatasource(ctx context.Context, mod *composeType.Module, ld *report.LoadStepDefinition) (report.Datasource, error) {
	//TODO implement me
	panic("implement me")
}

func (s Store) PartialComposeRecordValueUpdate(ctx context.Context, mod *composeType.Module, values ...*composeType.RecordValue) error {
	//TODO implement me
	panic("implement me")
}

func (s Store) ComposeRecordValueRefLookup(ctx context.Context, m *composeType.Module, field string, ref uint64) (uint64, error) {
	//q := s.composeRecordValuesSelectBuilder().
	//	Join(s.composeRecordTable("crd"), "crv.record_id = crd.id").
	//	Where(squirrel.Eq{
	//		"crv.name":       field,
	//		"crv.ref":        ref,
	//		"crv.deleted_at": nil,
	//		"crd.module_id":  m.ID,
	//		"crd.deleted_at": nil,
	//	}).
	//	Column("record_id").
	//	Limit(1)
	//
	//row, err := s.QueryRow(ctx, q)
	//if errors.Is(err, sql.ErrNoRows) {
	//	return 0, nil
	//} else if err != nil {
	//	return 0, err
	//}
	//
	//var recordID uint64
	//if err = row.Scan(&recordID); err != nil {
	//	return 0, err
	//}
	//
	//return recordID, nil
	return 0, nil
}
