package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func TestModuleImport_CastSet(t *testing.T) {
	stdFieldsAsTest := func(t *testing.T, module *Module) {
		// map definition's sort order is unreliable
		req := require.New(t)

		req.Len(module.set, 1)

		tc := module.set.FindByHandle("mod1")
		req.Equal(tc.Name, "Module with fields")
		req.Len(tc.Fields, 2)
		req.Equal(tc.Fields.FindByName("f1"), &types.ModuleField{
			Name:  "f1",
			Kind:  "Number",
			Label: "f1",
		})

		req.Equal(tc.Fields.FindByName("f2"), &types.ModuleField{
			Name:     "f2",
			Kind:     "String",
			Label:    "F2",
			Place:    1,
			Required: true,
			Options: map[string]interface{}{
				"multiLine":         true,
				"useRichTextEditor": true,
				"multiDelimiter":    "\n",
			},
		})
	}

	impFixTester(t, "module_fields_as_slice", stdFieldsAsTest)
	impFixTester(t, "module_fields_as_map", stdFieldsAsTest)
}
