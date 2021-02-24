package yaml

import (
	"reflect"

	"github.com/cortezaproject/corteza-server/pkg/y7s"
	. "github.com/cortezaproject/corteza-server/pkg/y7s"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

type (
	compose struct {
		Namespaces composeNamespaceSet
		Modules    composeModuleSet
		Records    composeRecordSet
		Pages      composePageSet
		Charts     composeChartSet

		EncoderConfig *EncoderConfig `yaml:"-"`
	}
)

func (c *compose) MarshalYAML() (interface{}, error) {
	nsDefined := c.Namespaces != nil && len(c.Namespaces) > 0
	cn, _ := makeMap()
	var err error

	addNsRef := func(n *yaml.Node, ref string) (*yaml.Node, error) {
		if !nsDefined {
			n, err = addMap(n,
				"namespace", ref,
			)
			if err != nil {
				return nil, err
			}
		}

		return n, nil
	}

	if nsDefined {
		c.Namespaces.ConfigureEncoder(c.EncoderConfig)

		cn, err = encodeResource(cn, "namespaces", c.Namespaces, c.EncoderConfig.MappedOutput, "slug")
		if err != nil {
			return nil, err
		}
	}

	if len(c.Modules) > 0 {
		cn, err = addNsRef(cn, c.Modules[0].refNamespace)
		if err != nil {
			return nil, err
		}

		c.Modules.ConfigureEncoder(c.EncoderConfig)

		cn, err = encodeResource(cn, "modules", c.Modules, c.EncoderConfig.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if len(c.Records) > 0 {
		cn, err = addNsRef(cn, c.Records[0].refNamespace)
		if err != nil {
			return nil, err
		}

		c.Records.ConfigureEncoder(c.EncoderConfig)

		// Records don't have this
		cn, err = encodeResource(cn, "records", c.Records, false, "")
		if err != nil {
			return nil, err
		}
	}

	if len(c.Pages) > 0 {
		cn, err = addNsRef(cn, c.Pages[0].refNamespace)
		if err != nil {
			return nil, err
		}

		c.Pages.ConfigureEncoder(c.EncoderConfig)

		// @todo A bit of a complication with pages and handles...
		//       Will probably just leave it as so, but might change it later.
		cn, err = encodeResource(cn, "pages", c.Pages, false, "handle")
		if err != nil {
			return nil, err
		}
	}

	if len(c.Charts) > 0 {
		cn, err = addNsRef(cn, c.Charts[0].refNamespace)
		if err != nil {
			return nil, err
		}

		c.Charts.ConfigureEncoder(c.EncoderConfig)

		cn, err = encodeResource(cn, "charts", c.Charts, c.EncoderConfig.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	return cn, nil
}

func (c *compose) UnmarshalYAML(n *yaml.Node) error {
	var (
		nsRef string
		err   error
	)

	// 1st pass: handle doc-level references
	err = y7s.EachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "namespace":
			if def := FindKeyNode(n, "namespaces"); def != nil {
				return y7s.NodeErr(def, "cannot combine namespace reference and namespaces definition")
			}

			if err := decodeRef(v, "namespace", &nsRef); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 2nd pass: handle definitions
	return y7s.EachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "namespaces":
			err = v.Decode(&c.Namespaces)
			if err != nil {
				return err
			}
			// If we're defining a fresh namespace and we're not nesting
			if nsRef == "" {
				nsRef = c.Namespaces[0].res.Slug
			}

		case "modules":
			if err = v.Decode(&c.Modules); err != nil {
				return err
			}

			return c.Modules.setNamespaceRef(nsRef)

		case "pages":
			if err = v.Decode(&c.Pages); err != nil {
				return err
			}

			return c.Pages.setNamespaceRef(nsRef)

		case "charts":
			if err = v.Decode(&c.Charts); err != nil {
				return err
			}

			return c.Charts.setNamespaceRef(nsRef)

		case "records":
			if err = v.Decode(&c.Records); err != nil {
				return err
			}

			return c.Records.setNamespaceRef(nsRef)

		}

		return nil
	})
}

func (c compose) MarshalEnvoy() ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, 100)

	rf := reflect.ValueOf(c)
	for i := 0; i < rf.NumField(); i++ {
		if mr, ok := rf.Field(i).Interface().(EnvoyMarshler); ok {
			tmp, err := mr.MarshalEnvoy()
			if err != nil {
				return nil, err
			}
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}
