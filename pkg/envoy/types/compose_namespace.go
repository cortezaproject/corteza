package types

import (
	compTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/util"
	"gopkg.in/yaml.v3"
)

type (
	ComposeNamespace struct {
		compTypes.Namespace `yaml:",inline"`
		Modules             ComposeModuleSet
		Rbac                `yaml:",inline"`
	}
	ComposeNamespaceSet []*ComposeNamespace
)

func (ss *ComposeNamespaceSet) UnmarshalYAML(n *yaml.Node) error {
	nss := ComposeNamespaceSet{}

	err := util.YamlIterator(n, func(n, m *yaml.Node) error {
		slug := ""
		if n != nil {
			slug = n.Value
		}

		ns := &ComposeNamespace{}
		err := m.Decode(ns)
		if err != nil {
			return err
		}
		if ns.Slug == "" {
			ns.Slug = slug
		}
		if ns.Name == "" {
			ns.Name = slug
		}
		nss = append(nss, ns)
		return nil
	})

	if err != nil {
		return err
	}

	*ss = nss
	return nil
}
