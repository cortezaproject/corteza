package yaml

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
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
		rbac         rbacRuleSet
	}
)

func (doc *Document) UnmarshalYAML(n *yaml.Node) (err error) {
	if err = n.Decode(&doc.compose); err != nil {
		return
	}

	if doc.rbac, err = decodeRbac(n); err != nil {
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
func (doc *Document) Decode(ctx context.Context, l loader) ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, 100)

	if doc.compose != nil {
		if tmp, err := doc.compose.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}
