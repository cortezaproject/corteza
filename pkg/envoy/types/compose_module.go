package types

import (
	compTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/util"
	"gopkg.in/yaml.v3"
)

type (
	ComposeModuleField struct {
		compTypes.ModuleField `yaml:",inline"`
	}
	ComposeModuleFieldSet []*ComposeModuleField

	ComposeModule struct {
		compTypes.Module `yaml:",inline"`
		Fields           ComposeModuleFieldSet `yaml:"fields"`
		Rbac             `yaml:",inline"`
	}
	ComposeModuleSet []*ComposeModule
)

func (mm *ComposeModuleSet) UnmarshalYAML(n *yaml.Node) error {
	cms := ComposeModuleSet{}

	err := util.YamlIterator(n, func(n, m *yaml.Node) error {
		handle := ""
		if n != nil {
			handle = n.Value
		}

		mod := &ComposeModule{}
		err := m.Decode(mod)
		if err != nil {
			return err
		}

		if mod.Handle == "" {
			mod.Handle = handle
		}
		if mod.Name == "" {
			mod.Name = handle
		}
		cms = append(cms, mod)
		return nil
	})

	if err != nil {
		return err
	}

	*mm = cms
	return nil
}

func (ff *ComposeModuleFieldSet) UnmarshalYAML(n *yaml.Node) error {
	ffs := ComposeModuleFieldSet{}

	err := util.YamlIterator(n, func(n, m *yaml.Node) error {
		name := ""
		if n != nil {
			name = n.Value
		}
		f := &ComposeModuleField{}

		err := m.Decode(f)
		if err != nil {
			return err
		}
		if f.Name == "" {
			f.Name = name
		}
		if f.Label == "" {
			f.Label = name
		}
		ffs = append(ffs, f)
		return nil
	})

	if err != nil {
		return err
	}

	*ff = ffs
	return nil
}
