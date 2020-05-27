package repository

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func TestRecordReportBuilder2(t *testing.T) {
	builder := NewRecordReportBuilder(&types.Module{
		ID: 1000,
		Fields: types.ModuleFieldSet{
			&types.ModuleField{Name: "single1"},
			&types.ModuleField{Name: "multi1", Multi: true},
			&types.ModuleField{Name: "ref1", Kind: "Record"},
			&types.ModuleField{Name: "multiRef1", Kind: "Record", Multi: true},
		}},
	)

	expected := "SELECT (COUNT(*)) AS count, (CAST(max(rv_single1.value) AS DECIMAL(14,2))) AS metric_0, " +
		"(QUARTER(rv_ref1.value)) AS dimension_0 " +
		"FROM compose_record AS r " +
		"LEFT JOIN compose_record_value AS rv_single1 ON (rv_single1.record_id = r.id AND rv_single1.name = ? AND rv_single1.deleted_at IS NULL) " +
		"LEFT JOIN compose_record_value AS rv_ref1 ON (rv_ref1.record_id = r.id AND rv_ref1.name = ? AND rv_ref1.deleted_at IS NULL) " +
		"WHERE r.deleted_at IS NULL AND r.module_id = ? AND (rv_ref1.value = 2) " +
		"GROUP BY dimension_0 " +
		"ORDER BY dimension_0"

	sql, _, err := builder.Build("max(single1)", "QUARTER(ref1)", "ref1 = 2")
	require.NoError(t, err)
	require.Equal(t, expected, sql)
}
