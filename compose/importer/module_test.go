package importer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
)

func TestModuleImport_CastSet(t *testing.T) {
	impFixTester(t, "module_full", func(t *testing.T, module *Module) {
		req := require.New(t)

		req.Len(module.set, 1)

		tc := module.set.FindByHandle("mod1")
		req.Equal(tc.Name, "Module with fields")
		req.Equal(tc.Fields, types.ModuleFieldSet{
			{
				Name:  "f1",
				Kind:  "Number",
				Label: "f1",
			},
			{

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
			},
		})
	})
}
