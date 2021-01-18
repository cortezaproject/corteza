package yaml

import (
	"fmt"
	. "github.com/cortezaproject/corteza-server/pkg/y7s"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	composePage struct {
		res    *types.Page
		ts     *resource.Timestamps
		config *resource.EnvoyConfig

		children composePageSet

		refNamespace string
		refModule    string
		refParent    string

		rbac rbacRuleSet
	}
	composePageSet []*composePage

	composePageBlock      = types.PageBlock
	composePageBlockStyle = types.PageBlockStyle
)

func (wset *composePageSet) UnmarshalYAML(n *yaml.Node) error {
	wx := make(map[uint64]int)

	return Each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &composePage{
				// Set this to something negative so we have an easier time determining
				// if we should fix the pages weight
				res: &types.Page{Weight: -1},
			}
		)

		if v == nil {
			return NodeErr(n, "malformed page definition")
		}

		if err = v.Decode(&wrap); err != nil {
			return
		}

		if err = decodeRef(k, "page", &wrap.res.Handle); err != nil {
			return
		}

		if wrap.res.Title == "" {
			// if name is not set, use handle
			wrap.res.Title = wrap.res.Handle
		}

		if wrap.res.Weight < 0 {
			wrap.res.Weight = wx[wrap.res.SelfID]
		}
		wx[wrap.res.SelfID]++

		*wset = append(*wset, wrap)
		return
	})
}

// @todo also do this for the pages property
func (wset composePageSet) setNamespaceRef(ref string) error {
	for _, res := range wset {
		if res.refNamespace != "" && ref != res.refNamespace {
			return fmt.Errorf("cannot override namespace reference %s with %s", res.refNamespace, ref)
		}

		res.refNamespace = ref
		if res.children != nil {
			res.children.setNamespaceRef(ref)
		}
	}

	return nil
}

func (wset composePageSet) MarshalEnvoy() ([]resource.Interface, error) {
	// namespace usually have bunch of sub-resources defined
	nn := make([]resource.Interface, 0, len(wset)*10)

	for _, res := range wset {
		if tmp, err := res.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}

func (wrap *composePage) UnmarshalYAML(n *yaml.Node) (err error) {
	if wrap.res == nil {
		wrap.rbac = make(rbacRuleSet, 0, 10)
		wrap.res = &types.Page{
			// Pages are visible by default
			Visible: true,
		}
	}

	if wrap.rbac, err = decodeRbac(n); err != nil {
		return
	}

	if wrap.config, err = decodeEnvoyConfig(n); err != nil {
		return
	}

	if wrap.ts, err = decodeTimestamps(n); err != nil {
		return
	}

	return EachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "title":
			return DecodeScalar(v, "page title", &wrap.res.Title)

		case "handle":
			return DecodeScalar(v, "page handle", &wrap.res.Handle)

		case "module":
			return decodeRef(v, "page module", &wrap.refModule)

		case "visible":
			return DecodeScalar(v, "page visible", &wrap.res.Visible)

		case "description":
			return DecodeScalar(v, "page description", &wrap.res.Description)

		case "blocks":
			var cpb = make([]composePageBlock, 0)
			if err = v.Decode(&cpb); err != nil {
				return err
			}

			wrap.res.Blocks = make([]types.PageBlock, len(cpb))
			for i, b := range cpb {
				wrap.res.Blocks[i] = b
			}

		case "children",
			"pages":
			return v.Decode(&wrap.children)

		}

		return nil
	})
}

func (wrap composePage) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewComposePage(wrap.res, wrap.refNamespace, wrap.refModule, wrap.refParent)
	rs.SetTimestamps(wrap.ts)
	rs.SetConfig(wrap.config)

	return envoy.CollectNodes(
		rs,
		wrap.children.bindParent(rs),
		wrap.rbac.bindResource(rs),
	)
}

func (rr composePageSet) bindParent(res resource.Interface) composePageSet {
	rtr := make(composePageSet, 0, len(rr))
	for _, r := range rr {
		idd := res.Identifiers().StringSlice()
		if len(idd) > 0 {
			r.refParent = idd[0]
		}
		rtr = append(rtr, r)
	}

	return rtr
}
