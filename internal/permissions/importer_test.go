package permissions

import (
	"testing"
)

func TestPermissionRulesImport_CastResourcesSet(t *testing.T) {
	t.Skip()
	// tests := []struct {
	// 	name     string
	// 	resource string
	// 	yaml     string
	// 	rules    map[string]RuleSet
	// }{
	// 	{name: "empty", yaml: ``},
	// 	{name: "empty map", yaml: `{}`},
	// 	{name: "empty slice", yaml: `[]`},
	// 	{
	// 		name: "one role, one resource, two ops",
	// 		yaml: `admins: { resource: [ read, write ] }`,
	// 		rules: map[string]RuleSet{
	// 			"admins": {
	// 				AllowRule(0, "resource", "read"),
	// 				AllowRule(0, "resource", "write"),
	// 			},
	// 		},
	// 	},
	// 	{
	// 		name: "op as string",
	// 		yaml: `admins: { resource: read }`,
	// 		rules: map[string]RuleSet{
	// 			"admins": {
	// 				AllowRule(0, "resource", "read"),
	// 			},
	// 		},
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		imp := &Importer{}
	//
	// 		aux, err := importer.ParseYAML([]byte(tt.yaml))
	// 		require.NoError(t, err)
	// 		require.NoError(t, imp.CastResourcesSet("allow", aux))
	// 		require.Equal(t, tt.rules, imp.rules)
	// 	})
	// }
}
