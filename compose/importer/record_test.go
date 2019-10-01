package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func TestRecordImport(t *testing.T) {

	// Pre fill with module that imported record is referring to
	imp.namespaces.Setup(ns)
	imp.GetModuleImporter(ns.Slug).set = types.ModuleSet{
		{
			NamespaceID: ns.ID,
			Handle:      "TestModule",
			Fields: types.ModuleFieldSet{
				{Name: "field1"},
				{Name: "field2"},
			},
		},
	}

	impFixTester(t, "records", func(t *testing.T, record *Record) {
		req := require.New(t)

		req.Len(record.set["TestModule"], 2)

		r0v := record.set["TestModule"][0].Values
		r1v := record.set["TestModule"][1].Values

		req.Equal(r0v.FilterByName("field1"), types.RecordValueSet{
			{Name: "field1", Value: "val1.1"},
		})
		req.Equal(r0v.FilterByName("field2"), types.RecordValueSet{
			{Name: "field2", Value: "val1.2"},
		})

		req.Equal(r1v.FilterByName("field1"), types.RecordValueSet{
			{Name: "field1", Value: "val2.1"},
		})

		req.Equal(r1v.FilterByName("field2"), types.RecordValueSet{
			{Name: "field2", Value: "val2.2"},
		})

	})
}
