package yaml

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
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
				},
				refRole: roleRef,
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
	ref := &resource.Ref{
		ResourceType: resI.ResourceType(),
		Identifiers:  resI.Identifiers(),
	}

	var path []*resource.Ref
	if resRbac, ok := resI.(resource.RBACInterface); ok {
		path = resRbac.RBACPath()
	}

	rtr := make(rbacRuleSet, 0, len(rr))
	for _, r := range rr {
		_ = r.SetResource(ref.ResourceType)
		r.refRes = ref
		r.refPathRes = path
		rtr = append(rtr, r)
	}

	return rtr
}

func (rr rbacRuleSet) MarshalEnvoy() ([]resource.Interface, error) {
	var nn = make([]resource.Interface, 0, len(rr))

	for _, r := range rr {
		nn = append(nn, resource.NewRbacRule(r.res, r.refRole, r.refRes, r.refPathRes...))
	}
	return nn, nil
}

func (r *rbacRule) SetResource(res string) error {
	res, ref, pp, err := resource.ParseRule(res)
	if err != nil {
		return err
	}

	if res != "" {
		r.res.Resource = res
	}

	return r.bindRefs(ref, pp, err)
}
