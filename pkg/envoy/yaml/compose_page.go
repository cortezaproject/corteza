package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"gopkg.in/yaml.v3"
)

type (
	composePage struct {
		res          *types.Page
		pages        composePageSet
		refNamespace string
		refModule    string
		*rbacRules
	}
	composePageSet []*composePage

	composePageBlock      = types.PageBlock
	composePageBlockStyle = types.PageBlockStyle
)

func (wset *composePageSet) UnmarshalYAML(n *yaml.Node) error {
	return each(n, func(k, v *yaml.Node) (err error) {
		var (
			wrap = &composePage{}
		)

		if v == nil {
			return nodeErr(n, "malformed page definition")
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

		*wset = append(*wset, wrap)
		return
	})
}

func (wset composePageSet) setNamespaceRef(ref string) error {
	for _, res := range wset {
		if res.refNamespace != "" && ref != res.refNamespace {
			return fmt.Errorf("cannot override namespace reference %s with %s", res.refNamespace, ref)
		}

		res.refNamespace = ref
	}

	return nil
}

func (wset composePageSet) MarshalEnvoy() ([]envoy.Node, error) {
	// namespace usually have bunch of sub-resources defined
	nn := make([]envoy.Node, 0, len(wset)*10)

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
		wrap.rbacRules = &rbacRules{}
		wrap.res = &types.Page{
			// Pages are visible by default
			Visible: true,
		}
	}

	if wrap.rbacRules, err = decodeResourceAccessControl(types.PageRBACResource, n); err != nil {
		return
	}

	return eachMap(n, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "title":
			return decodeScalar(v, "page title", &wrap.res.Title)

		case "handle":
			return decodeScalar(v, "page handle", &wrap.res.Handle)

		case "module":
			return decodeRef(v, "page module", &wrap.refModule)

		case "visible":
			return decodeScalar(v, "page visible", &wrap.res.Visible)

		case "description":
			return decodeScalar(v, "page description", &wrap.res.Description)

		case "blocks":
			var cpb = make([]composePageBlock, 0)
			if err = v.Decode(&cpb); err != nil {
				return err
			}

			wrap.res.Blocks = make([]types.PageBlock, len(cpb))
			for i, b := range cpb {
				wrap.res.Blocks[i] = b
			}

		case "pages":
			return v.Decode(&wrap.pages)

		}

		return nil
	})
}

func (wrap composePage) MarshalEnvoy() ([]envoy.Node, error) {
	nn := make([]envoy.Node, 0, 16)
	//nn = append(nn, &envoy.ComposePageNode{Page: wrap.res})

	return nn, nil
}
