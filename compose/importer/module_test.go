package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func TestModuleImport_CastSet(t *testing.T) {
	stdFieldsAsTest := func(t *testing.T, module *Module, testOrder bool) {
		// map definition's sort order is unreliable
		req := require.New(t)

		req.Len(module.set, 1)

		tc := module.set.FindByHandle("mod1")
		req.Equal(tc.Name, "Module with fields")
		req.Len(tc.Fields, 2)
		req.NotNil(tc.Fields.FindByName("f1"))
		req.NotNil(tc.Fields.FindByName("f2"))

		f1 := tc.Fields.FindByName("f1")
		f2 := tc.Fields.FindByName("f2")

		if testOrder {
			req.Equal(0, f1.Place)
			req.Equal(1, f2.Place)
		} else {
			// Check if field's place property has been set
			req.Equal(1, f1.Place+f2.Place)
		}

		// due to "module_fields_as_map" & unstable order of map[..]..
		// we need to reset order here
		f1.Place = 0
		f2.Place = 0

		req.Equal(f1, &types.ModuleField{
			Name:  "f1",
			Kind:  "Number",
			Label: "f1",
		})

		req.Equal(f2, &types.ModuleField{
			Name:     "f2",
			Kind:     "String",
			Label:    "F2",
			Required: true,
			Options: map[string]interface{}{
				"multiLine":         true,
				"useRichTextEditor": true,
				"multiDelimiter":    "\n",
			},
		})
	}

	impFixTester(t, "module_fields_as_slice", func(t *testing.T, module *Module) {
		stdFieldsAsTest(t, module, true)

	})

	impFixTester(t, "module_fields_as_map", func(t *testing.T, module *Module) {
		stdFieldsAsTest(t, module, false)
	})
}
