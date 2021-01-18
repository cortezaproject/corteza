package yaml

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	. "github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	rbacRule struct {
		res *rbac.Rule

		resRef *resource.Ref

		// To help us construct the resource
		resource    rbac.Resource
		refResource string
		// Related role
		refRole string
	}
	rbacRuleSet []*rbacRule
)

func decodeRbac(n *yaml.Node) (rbacRuleSet, error) {
	var (
		rr = make(rbacRuleSet, 0, 20)
	)

	return rr, EachMap(n, func(k, v *yaml.Node) (err error) {
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

func (rr rbacRuleSet) decodeRbac(a rbac.Access, rules *yaml.Node) (rbacRuleSet, error) {
	if rr == nil {
		rr = make(rbacRuleSet, 0, 10)
	}

	var err error

	parseOps := func(ops *yaml.Node, roleRef, res string) error {
		return EachSeq(ops, func(op *yaml.Node) error {
			rule := &rbacRule{
				res: &rbac.Rule{
					Access:    a,
					Operation: rbac.Operation(op.Value),
				},
				refRole: roleRef,
			}
			rule.SetResource(res)
			rr = append(rr, rule)
			return nil
		})
	}

	err = EachMap(rules, func(role, ops *yaml.Node) error {
		// If its a mapping node, keys represent resources
		if ops.Kind == yaml.MappingNode {
			err = EachMap(ops, func(res, ops *yaml.Node) error {
				return parseOps(ops, role.Value, res.Value)
			})
		} else {
			return parseOps(ops, role.Value, "")
		}

		return nil
	})

	return rr, err
}

func (rr rbacRuleSet) bindResource(resI resource.Interface) rbacRuleSet {
	ref := &resource.Ref{
		ResourceType: resI.ResourceType(),
		Identifiers:  resI.Identifiers(),
	}

	rtr := make(rbacRuleSet, 0, len(rr))
	for _, r := range rr {
		r.SetResource(ref.ResourceType)
		r.resRef = ref
		rtr = append(rtr, r)
	}

	return rtr
}

func (rr rbacRuleSet) setResource(res rbac.Resource) error {
	for _, r := range rr {
		if r.resource.String() != "" && res != r.resource {
			return fmt.Errorf("cannot override resource %s with %s", r.resource, res)
		}

		r.resource = res
	}

	return nil
}

func (rr rbacRuleSet) setResourceRef(ref string) error {
	for _, r := range rr {
		if r.refResource != "" && ref != r.refResource {
			return fmt.Errorf("cannot override resource reference %s with %s", r.refResource, ref)
		}

		r.refResource = ref
	}

	return nil
}

func (rr rbacRuleSet) MarshalEnvoy() ([]resource.Interface, error) {
	var nn = make([]resource.Interface, 0, len(rr))

	for _, r := range rr {
		nn = append(nn, resource.NewRbacRule(r.res, r.refRole, r.resRef))
	}
	return nn, nil
}

func (r *rbacRule) SetResource(res string) {
	if res == "" {
		return
	}

	sp := ":"

	res = strings.TrimSpace(res)
	res = strings.TrimRight(res, sp)
	rr := strings.Split(res, sp)

	// When len is 1; only top-level defined (system, compose, messaging)
	if len(rr) == 1 {
		r.res.Resource = rbac.Resource(res)
		return
	}

	// When len is 2; top-level and sub level defined (compose:namespace, system:user, ...)
	if len(rr) == 2 {
		r.res.Resource = rbac.Resource(res + sp)
		return
	}

	//When len is 3; both levels defined; resource ref also provided
	if len(rr) == 3 {
		res = strings.Join(rr[0:2], sp) + sp
		r.resRef = &resource.Ref{
			ResourceType: strings.Join(rr[0:2], sp) + sp,
			Identifiers:  resource.MakeIdentifiers(rr[2]),
		}
		r.res.Resource = rbac.Resource(res)
	}
}
