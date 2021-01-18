package yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	. "github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	// Document defines the supported yaml structure
	Document struct {
		compose      *compose
		messaging    *messaging
		roles        roleSet
		users        userSet
		applications applicationSet
		settings     *settings
		rbac         rbacRuleSet
	}
)

func (doc *Document) UnmarshalYAML(n *yaml.Node) (err error) {
	if err = n.Decode(&doc.compose); err != nil {
		return
	}

	if err = n.Decode(&doc.messaging); err != nil {
		return
	}

	if doc.rbac, err = decodeRbac(n); err != nil {
		return
	}

	return EachMap(n, func(k, v *yaml.Node) error {
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
func (doc *Document) Decode(ctx context.Context) ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, 100)

	mm := make([]envoy.Marshaller, 0, 20)
	if doc.compose != nil {
		mm = append(mm, doc.compose)
	}
	if doc.messaging != nil {
		mm = append(mm, doc.messaging)
	}
	if doc.roles != nil {
		mm = append(mm, doc.roles)
	}
	if doc.users != nil {
		mm = append(mm, doc.users)
	}
	if doc.applications != nil {
		mm = append(mm, doc.applications)
	}
	if doc.settings != nil {
		mm = append(mm, doc.settings)
	}
	if doc.rbac != nil {
		mm = append(mm, doc.rbac)
	}

	for _, m := range mm {
		if tmp, err := m.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}
