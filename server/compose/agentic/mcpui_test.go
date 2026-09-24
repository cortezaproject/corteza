package agentic

import (
	"strings"
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordViewDescribesFieldsAndLinksRecords(t *testing.T) {
	mod := &cmpTypes.Module{
		ID:     42,
		Name:   "Tasks",
		Handle: "task",
		Fields: cmpTypes.ModuleFieldSet{
			{Name: "title", Label: "Title", Kind: "String"},
			{Name: "tags", Kind: "String", Multi: true},
			{Name: "priority", Kind: "Select", Options: cmpTypes.ModuleFieldOptions{
				"options":    []any{map[string]any{"value": "high", "text": "High"}},
				"selectType": "default",
			}},
		},
	}
	set := cmpTypes.RecordSet{{ID: 7}, {ID: 8}}

	v := recordViewOf(mod, "crm", set)

	assert.Equal(t, recordViewModule{Name: "Tasks", Handle: "task"}, v.Module)
	require.Len(t, v.Fields, 3)
	assert.Equal(t, recordViewField{Name: "title", Label: "Title", Kind: "String"}, v.Fields[0])
	assert.True(t, v.Fields[1].Multi)
	assert.Equal(t, []any{map[string]any{"value": "high", "text": "High"}}, v.Fields[2].Options,
		"a Select column carries its options and nothing else from the field's options")

	require.Len(t, v.Links, 2)
	assert.True(t, strings.HasSuffix(v.Links["7"], "/compose/namespace/crm/admin/modules/42/records/7"), v.Links["7"])
}

func TestRecordViewWithoutSlugHasNoLinks(t *testing.T) {
	v := recordViewOf(&cmpTypes.Module{ID: 42}, "", cmpTypes.RecordSet{{ID: 7}})
	assert.Nil(t, v.Links)
}

func TestReportViewNamesMetricsAndTheDimensionField(t *testing.T) {
	mod := &cmpTypes.Module{
		Name: "Deals",
		Fields: cmpTypes.ModuleFieldSet{
			{Name: "stage", Label: "Stage", Kind: "Select", Options: cmpTypes.ModuleFieldOptions{
				"options": []any{"lead", "won"},
			}},
			{Name: "closed", Kind: "DateTime"},
		},
	}

	v := reportViewOf(mod, "stage", "SUM(amount) AS total, AVG(quantity * price), COUNT(id)")

	require.NotNil(t, v.Dimension)
	assert.Equal(t, "Stage", v.Dimension.Label)
	assert.Equal(t, []any{"lead", "won"}, v.Dimension.Options)
	assert.Equal(t, []reportMetric{
		{Key: "total", Field: "amount"},
		{Key: "AVG(quantity * price)", Field: "quantity"},
		{Key: "COUNT(id)", Field: "id"},
	}, v.Metrics, "each metric keeps the key the rows carry it under, in the order asked")
}

func TestReportViewReadsTheFieldInsideADateBucket(t *testing.T) {
	mod := &cmpTypes.Module{Fields: cmpTypes.ModuleFieldSet{{Name: "closed", Kind: "DateTime"}}}

	v := reportViewOf(mod, "DATE(closed)", "")

	require.NotNil(t, v.Dimension)
	assert.Equal(t, "DateTime", v.Dimension.Kind)
	assert.Empty(t, v.Metrics)
}

func TestReportViewWithoutAKnownDimensionHasNone(t *testing.T) {
	assert.Nil(t, reportViewOf(&cmpTypes.Module{}, "", "").Dimension)
	assert.Nil(t, reportViewOf(&cmpTypes.Module{}, "missing", "").Dimension)
}
