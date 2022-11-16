package yaml

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/cortezaproject/corteza/server/system/types"
	"gopkg.in/yaml.v3"
)

func decodeRbac(n *yaml.Node) (rbacRuleSet, error) {
	var (
		rr = make(rbacRuleSet, 0, 20)
	)

	return rr, y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "allow":
			rr, err = rr.decodeRbac(rbac.Allow, v)
			if err != nil {
				return err
			}
		case "deny":
			rr, err = rr.decodeRbac(rbac.Deny, v)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (rr rbacRuleSet) decodeRbac(a rbac.Access, rules *yaml.Node) (oo rbacRuleSet, err error) {
	if rr == nil {
		oo = make(rbacRuleSet, 0, 10)
	} else {
		oo = rr
	}

	parseOps := func(ops *yaml.Node, roleRef, res string) error {
		return y7s.EachSeq(ops, func(op *yaml.Node) error {
			rule := &rbacRule{
				res: &rbac.Rule{
					Access:    a,
					Operation: op.Value,
					Resource:  res,
				},
				refRole: resource.MakeRef(types.RoleResourceType, resource.MakeIdentifiers(roleRef)),
			}

			if res != "" {
				if err = rule.SetResource(res); err != nil {
					return fmt.Errorf("failed to decode RBAC rule for role '%s': %w", roleRef, err)
				}
			}

			oo = append(oo, rule)
			return nil
		})
	}

	return oo, y7s.EachMap(rules, func(role, ops *yaml.Node) error {
		// If its a mapping node, keys represent resources
		if ops.Kind == yaml.MappingNode {
			err = y7s.EachMap(ops, func(res, ops *yaml.Node) error {
				return parseOps(ops, role.Value, res.Value)
			})
		} else {
			return parseOps(ops, role.Value, "")
		}

		return nil
	})
}

func (rr rbacRuleSet) bindResource(resI resource.Interface) rbacRuleSet {
	res, ref, pp := resI.(resource.RBACInterface).RBACParts()

	rtr := make(rbacRuleSet, 0, len(rr))
	for _, r := range rr {
		r.refRbacResource = res
		r.refRbacRes = ref
		r.refPathRes = pp
		rtr = append(rtr, r)
	}

	return rtr
}

func (rr rbacRuleSet) MarshalEnvoy() ([]resource.Interface, error) {
	var nn = make([]resource.Interface, 0, len(rr))

	for _, r := range rr {
		nn = append(nn, resource.NewRbacRule(r.res, r.refRole, r.refRbacRes, r.refRbacResource, r.refPathRes...))
	}
	return nn, nil
}

func (r *rbacRule) SetResource(res string) error {
	_, ref, pp, err := resource.ParseRule(res)
	if err != nil {
		return err
	}

	r.refRbacResource = res
	r.refRbacRes = ref
	r.refPathRes = pp
	return nil
}
