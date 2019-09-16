package permissions

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
)

type (
	PermissionRulesImport struct {
		whitelist interface {
			Check(rule *Rule) bool
		}

		// Rules per role
		rules map[string]RuleSet
	}
)

// CastSet - resolves permission rules for specific resource:
//   <role>: [<operation>, ...]
func (imp *PermissionRulesImport) CastSet(resource, accessStr string, roles interface{}) (err error) {
	if !deinterfacer.IsMap(roles) {
		return errors.New("expecting map of roles")
	}
	return deinterfacer.Each(roles, func(_ int, roleHandle string, oo interface{}) error {
		return imp.appendPermissionRule(roleHandle, accessStr, resource, oo)
	})
}

// CastResourcesSet - resolves permission rules:
//   <role> { <resource>: [<operation>, ...] }
func (imp *PermissionRulesImport) CastResourcesSet(accessStr string, roles interface{}) (err error) {
	// if !IsIterable(roles) {
	// 	return errors.New("expecting map of roles")
	// }

	return deinterfacer.Each(roles, func(_ int, roleHandle string, perResource interface{}) error {
		if !deinterfacer.IsMap(perResource) {
			return errors.New("expecting map of resources")
		}

		// Each resource
		return deinterfacer.Each(perResource, func(_ int, resource string, oo interface{}) error {
			return imp.appendPermissionRule(roleHandle, accessStr, resource, oo)
		})
	})
}

func (imp *PermissionRulesImport) appendPermissionRule(roleHandle, accessStr, res string, oo interface{}) (err error) {
	var access Access

	if err = access.UnmarshalJSON([]byte(accessStr)); err != nil {
		return
	}

	if imp.rules == nil {
		imp.rules = map[string]RuleSet{}
	}

	if imp.rules[roleHandle] == nil {
		imp.rules[roleHandle] = RuleSet{}
	}

	operations := deinterfacer.ToStrings(oo)
	if operations == nil {
		return errors.New("could not resolve permission rule operations")
	}

	sort.Strings(operations)

	for _, op := range operations {
		rule := &Rule{
			Access:    access,
			Resource:  Resource(res),
			Operation: Operation(op),
		}

		if imp.whitelist != nil && !imp.whitelist.Check(rule) {
			return errors.Errorf("invalid rule: %q on %q", res, op)
		}

		imp.rules[roleHandle] = append(imp.rules[roleHandle], rule)
	}

	return nil
}

// UpdateResources iterates over all rules and replaces resource (foo:bar => foo:42)
func (imp *PermissionRulesImport) UpdateResources(rwHandle, rwID Resource) {
	for _, rules := range imp.rules {
		for _, rule := range rules {
			if rule.Resource == rwHandle {
				rule.Resource = rwID
			}
		}
	}
}
