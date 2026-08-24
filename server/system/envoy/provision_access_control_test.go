package envoy

import (
	"os"
	"testing"

	automationService "github.com/crusttech/human/server/automation/service"
	composeService "github.com/crusttech/human/server/compose/service"
	federationService "github.com/crusttech/human/server/federation/service"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/pkg/y7s"
	systemService "github.com/crusttech/human/server/system/service"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// provisionRuleKey identifies one provisioned rule the way the store does.
type provisionRuleKey struct{ role, resource, operation string }

// decodeProvisionedRules reads an access-control config the way provisioning
// reads it, rather than as plain YAML: role handles arrive as references and
// the allow/deny split is structural, so a hand-rolled walk would be testing a
// different file than the one that gets imported.
func decodeProvisionedRules(t *testing.T, path string) map[provisionRuleKey]rbac.Access {
	t.Helper()

	f, err := os.ReadFile(path)
	require.NoError(t, err)

	var doc yaml.Node
	require.NoError(t, yaml.Unmarshal(f, &doc))

	got := map[provisionRuleKey]rbac.Access{}

	require.NoError(t, y7s.EachMap(doc.Content[0], func(k, v *yaml.Node) error {
		var (
			key string
			acc rbac.Access
		)

		if err := y7s.DecodeScalar(k, "key", &key); err != nil {
			return err
		}

		switch key {
		case "allow":
			acc = rbac.Allow
		case "deny":
			acc = rbac.Deny
		default:
			return nil
		}

		nn, err := unmarshalRBACNode(v, acc)
		if err != nil {
			return err
		}

		for _, n := range nn {
			r, ok := n.Resource.(*rbac.Rule)
			require.True(t, ok)

			role := ""
			for _, ref := range n.References {
				if len(ref.Identifiers.Slice) > 0 {
					role = ref.Identifiers.Slice[0]
				}
			}

			got[provisionRuleKey{role, r.Resource, r.Operation}] = r.Access
		}

		return nil
	}))

	return got
}

// operationCatalog is every (wildcard resource, operation) pair the four
// components will actually evaluate. List() is generated from the same source
// as the Can* functions, and its "any" is the wildcard form provisioning
// writes, so the two are directly comparable.
func operationCatalog() map[string]map[string]bool {
	out := map[string]map[string]bool{}

	for _, def := range [][]map[string]string{
		systemService.AccessControl(nil).List(),
		composeService.AccessControl(nil).List(),
		automationService.AccessControl(nil).List(),
		federationService.AccessControl(nil).List(),
	} {
		for _, op := range def {
			if out[op["any"]] == nil {
				out[op["any"]] = map[string]bool{}
			}
			out[op["any"]][op["op"]] = true
		}
	}

	return out
}

// A provisioned rule that names an operation no component evaluates, or writes
// a resource path with the wrong number of segments, imports without complaint
// and then never matches anything. Nothing at run time reports it: the role
// simply behaves as if the grant were not there. Read both sides and compare.
func TestProvisionedRulesAreEvaluable(t *testing.T) {
	var (
		catalog = operationCatalog()

		files = []string{
			"../../provision/000_base/system_access_control.yaml",
			"../../provision/000_base/compose_access_control.yaml",
			"../../provision/200_federation/2000_access_control.yaml",
			"../../provision/300_automation/2000_access_control.yaml",
			"../../provision/400_discovery/2000_access_control.yaml",
		}
	)

	require.NotEmpty(t, catalog)

	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			for k := range decodeProvisionedRules(t, f) {
				ops, ok := catalog[k.resource]
				require.True(t, ok,
					"%s grants %q on %q, which is not a resource any component evaluates",
					k.role, k.operation, k.resource)
				require.True(t, ops[k.operation],
					"%s grants %q on %q, which is not an operation of that resource",
					k.role, k.operation, k.resource)
			}
		})
	}
}

// The compose baseline is the whole of what a logged-in user may reach without
// a role, plus what an older install has had taken back off it.
func TestProvisionComposeAccessControl(t *testing.T) {
	got := decodeProvisionedRules(t, "../../provision/000_base/compose_access_control.yaml")

	// rbac.Deny is the zero value of rbac.Access, so a rule that is simply
	// absent reads as a denial. Every assertion below goes through this.
	access := func(k provisionRuleKey) string {
		v, ok := got[k]
		if !ok {
			return "<no rule>"
		}
		return v.String()
	}

	// An earlier baseline let every logged-in user read all compose data and
	// export a namespace whole. Provisioning only ever adds rules, so the deny
	// is what takes those back on an install that already carries them.
	for _, revoked := range []provisionRuleKey{
		{"authenticated", "corteza::compose/", "namespaces.search"},
		{"authenticated", "corteza::compose:namespace/*", "read"},
		{"authenticated", "corteza::compose:namespace/*", "export"},
		{"authenticated", "corteza::compose:module/*/*", "read"},
		{"authenticated", "corteza::compose:module/*/*", "records.search"},
		{"authenticated", "corteza::compose:module-field/*/*/*", "record.value.read"},
		{"authenticated", "corteza::compose:record/*/*/*", "read"},
		{"authenticated", "corteza::compose:chart/*/*", "read"},
		{"authenticated", "corteza::compose:page/*/*", "read"},
		{"authenticated", "corteza::compose:page-layout/*/*/*", "read"},
	} {
		require.Equal(t, "deny", access(revoked), "expected a deny for %v", revoked)
	}

	// Nothing is granted to everyone here; compose access is a role's to give.
	for k, v := range got {
		if k.role == "authenticated" {
			require.Equal(t, rbac.Deny, v, "compose baseline grew an allow: %v", k)
		}
	}

	// Screens that were reachable but answered "not allowed" to the roles that
	// own them: page layouts, record undelete and record ownership.
	for _, role := range []string{"admin", "low-code-admin"} {
		for _, held := range []provisionRuleKey{
			{role, "corteza::compose:page/*/*", "page-layout.create"},
			{role, "corteza::compose:page/*/*", "page-layouts.search"},
			{role, "corteza::compose:page-layout/*/*/*", "read"},
			{role, "corteza::compose:record/*/*/*", "undelete"},
			{role, "corteza::compose:record/*/*/*", "owner.manage"},
			{role, "corteza::compose:record/*/*/*", "revisions.search"},
			{role, "corteza::compose:module/*/*", "owned-record.create"},
			{role, "corteza::compose/", "resource-translations.manage"},
		} {
			require.Equal(t, "allow", access(held), "expected an allow for %v", held)
		}
	}

	// An agent may create a chart, so it holds the chart it created.
	for _, held := range []provisionRuleKey{
		{"agent-creator", "corteza::compose:namespace/*", "chart.create"},
		{"agent-creator", "corteza::compose:chart/*/*", "read"},
		{"agent-creator", "corteza::compose:chart/*/*", "update"},
		{"agent-creator", "corteza::compose:chart/*/*", "delete"},
	} {
		require.Equal(t, "allow", access(held), "expected an allow for %v", held)
	}
}
