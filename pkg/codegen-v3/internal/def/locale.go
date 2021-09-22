package def

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/tpl"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	locale struct {
		ResourceType string `yaml:"resourceType"`
		Resource     *resourceTranslation
		Extended     bool
		SkipSvc      bool `yaml:"skipSvc"`
		Keys         localeKeys
	}

	resourceTranslation struct {
		References []*resourceTranslationRef
	}

	resourceTranslationRef struct {
		Field        string
		ResourceType string
		Resource     string
		Component    string

		custom bool
	}

	localeKeys []*localeKey

	localeKey struct {
		Name          string
		Path          string
		Custom        bool
		CustomHandler string `yaml:"customHandler"`
		Description   string
		Field         string
	}
)

func (r *locale) proc(component, resource string) error {
	if r == nil {
		return nil
	}

	if r.ResourceType == "" {
		if strings.ToLower(resource) == "component" {
			r.ResourceType = component
		} else {
			r.ResourceType = fmt.Sprintf("%s:%s", component, resource)
		}
	}

	if r.Resource == nil {
		r.Resource = &resourceTranslation{References: []*resourceTranslationRef{{Field: "ID"}}}
	}

	// check types of each referenced component
	// and self-references (field==ID) with own resource type
	for _, rc := range r.Resource.References {
		if !rc.custom {
			if rc.Field == "ID" {
				rc.ResourceType = r.ResourceType
				rc.Component = component
				rc.Resource = resource
			} else {
				rc.ResourceType = fmt.Sprintf("%s:%s", component, rc.Field)
				rc.Component = component
				rc.Resource = rc.Field
				rc.Field = rc.Field + "ID"
			}
		}
	}

	// assure missing key field paths
	for _, k := range r.Keys {
		k.Custom = k.Custom || k.CustomHandler != ""

		// Guess the key path
		if k.Field == "" {
			k.Field = tpl.Export(k.Path)
		}

		// Guess the name
		if k.Name == "" {
			k.Name = k.Path
		}
	}

	return nil
}

func (op *localeKey) UnmarshalYAML(n *yaml.Node) error {
	type auxType localeKey

	switch {
	case y7s.IsKind(n, yaml.ScalarNode):
		return y7s.DecodeScalar(n, "locale key", &op.Path)

	case y7s.IsKind(n, yaml.MappingNode):
		var aux = (*auxType)(op)
		return n.Decode(aux)
	}

	return y7s.NodeErr(n, "unsupported formatting")
}

func (op *resourceTranslationRef) UnmarshalYAML(n *yaml.Node) error {
	if y7s.IsKind(n, yaml.ScalarNode) {
		op.Field = n.Value
		if n.Value != "ID" {
			op.ResourceType = n.Value
		}

		return nil
	}

	type auxType resourceTranslationRef
	var aux = (*auxType)(op)
	aux.custom = true
	return n.Decode(aux)
}
