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
