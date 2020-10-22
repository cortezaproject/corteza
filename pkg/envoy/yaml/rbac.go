package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"gopkg.in/yaml.v3"
)

type (
	rbacRules struct {
		// mapping unresolved roles to defined RBAC rules
		rules map[string]rbac.RuleSet
	}
)

func decodeRbacOperations(rr *rbacRules, access rbac.Access, role string, res rbac.Resource) func(*yaml.Node) error {
	if _, set := rr.rules[role]; !set {
		rr.rules[role] = rbac.RuleSet{}
	}

	return func(v *yaml.Node) error {
		rule := &rbac.Rule{Resource: res, Operation: rbac.Operation(v.Value), Access: access}
		rr.rules[role] = append(rr.rules[role], rule)
		return nil
	}
}

func decodeResourceAccessControl(res rbac.Resource, n *yaml.Node) (*rbacRules, error) {
	var (
		rr = &rbacRules{rules: make(map[string]rbac.RuleSet)}
	)

	return rr, eachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "allow":
			return rr.decodeResourceAccessControl(rbac.Allow, res, v)
		case "deny":
			return rr.decodeResourceAccessControl(rbac.Deny, res, v)
		}

		return nil
	})
}

func decodeGlobalAccessControl(n *yaml.Node) (*rbacRules, error) {
	var (
		rr = &rbacRules{rules: make(map[string]rbac.RuleSet)}
	)

	return rr, eachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "allow":
			return rr.decodeGlobalAccessControl(rbac.Allow, v)
		case "deny":
			return rr.decodeGlobalAccessControl(rbac.Deny, v)
		}

		return nil
	})
}

// Decodes RBAC rules defined as access/role/ops...
func (wrap *rbacRules) decodeResourceAccessControl(access rbac.Access, res rbac.Resource, v *yaml.Node) error {
	return eachMap(v, func(k, v *yaml.Node) error {
		return eachSeq(v, decodeRbacOperations(wrap, access, k.Value, res))
	})
}

// Decodes RBAC rules defined as access/role/resource/ops...
func (wrap *rbacRules) decodeGlobalAccessControl(access rbac.Access, v *yaml.Node) error {
	return eachMap(v, func(k, v *yaml.Node) error {
		var role = k.Value

		return eachMap(v, func(k, v *yaml.Node) error {
			var res = rbac.Resource(k.Value)
			if !res.IsValid() {
				return nodeErr(k, "RBAC resource %s invalid", k.Value)
			}

			return eachSeq(v, decodeRbacOperations(wrap, access, role, res))
		})
	})
}
