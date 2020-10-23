package yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"gopkg.in/yaml.v3"
)

type (
	// Document defines the supported yaml structure
	Document struct {
		compose      *compose
		roles        roleSet
		users        userSet
		applications applicationSet
		settings     settings
		*rbacRules
	}
)

func (doc *Document) UnmarshalYAML(n *yaml.Node) (err error) {
	if err = n.Decode(&doc.compose); err != nil {
		return
	}

	if doc.rbacRules, err = decodeGlobalAccessControl(n); err != nil {
		return
	}

	return eachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "roles":
			return v.Decode(&doc.roles)

		case "users":
			return v.Decode(&doc.users)

		case "applications":
			return v.Decode(&doc.applications)

		case "settings":
			return v.Decode(&doc.settings)

		}

		return nil
	})
}

//
func (doc *Document) Decode(ctx context.Context, l loader) ([]envoy.Node, error) {
	nn := make([]envoy.Node, 0, 100)

	if doc.compose != nil {
		if tmp, err := doc.compose.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil

	//// In case of namespaces...
	//if doc.Namespaces != nil {
	//	nodes, err := doc.Namespaces.Decode(ctx, l)
	//	if err != nil {
	//		return nil, err
	//	}
	//	nn = append(nn, nodes...)
	//}
	//
	//ns := &types.Namespace{}
	//if doc.Namespace != "" {
	//	// In case of a namespace to provide dependencies
	//	ns.Slug = doc.Namespace
	//	ns.Name = doc.Namespace
	//} else if len(nn) > 0 {
	//	// Try to fall back to a namespace node
	//	ns = ((nn[0]).(*envoy.ComposeNamespaceNode)).Ns
	//} else {
	//	// No good; we can't link with a namespace.
	//	// @note This should be checked when converting Compose resources only.
	//	//			 Some resources don't belong to a namespace.
	//	return nil, fmt.Errorf("cannot resolve namespace")
	//}
	//
	//// In case of modules...
	//if doc.Modules != nil {
	//	nodes, err := y.convertModules(doc.Modules, ns)
	//	if err != nil {
	//		return nil, err
	//	}
	//	nn = append(nn, nodes...)
	//}
	//
	//if doc.Records != nil {
	//	for modRef, rr := range doc.Records {
	//		// We can define a basic module representation as it will be updated later
	//		// during validation/runtime
	//		mod := &types.module{}
	//		mod.Handle = modRef
	//		mod.Name = modRef
	//
	//		nodes, err := y.convertRecords(rr, mod)
	//		if err != nil {
	//			return nil, err
	//		}
	//		nn = append(nn, nodes...)
	//	}
	//}

	return nn, nil
}
