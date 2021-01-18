package yaml

import (
	. "github.com/cortezaproject/corteza-server/pkg/y7s"
	"reflect"

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
	}
)

func (c *compose) UnmarshalYAML(n *yaml.Node) error {
	var (
		nsRef string
		err   error
	)

	// 1st pass: handle doc-level references
	err = EachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "namespace":
			if def := FindKeyNode(n, "namespaces"); def != nil {
				return NodeErr(def, "cannot combine namespace reference and namespaces definition")
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
	return EachMap(n, func(k, v *yaml.Node) error {
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
