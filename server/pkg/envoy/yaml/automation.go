package yaml

import (
	"reflect"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	automation struct {
		Workflows automationWorkflowSet

		EncoderConfig *EncoderConfig `yaml:"-"`
	}
)

func (c *automation) MarshalYAML() (interface{}, error) {
	cn, _ := makeMap()
	var err error

	if len(c.Workflows) > 0 {
		c.Workflows.configureEncoder(c.EncoderConfig)

		cn, err = encodeResource(cn, "workflows", c.Workflows, c.EncoderConfig.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	return cn, nil
}

func (a *automation) UnmarshalYAML(n *yaml.Node) error {
	var err error

	return y7s.EachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "workflows":
			err = v.Decode(&a.Workflows)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (a automation) MarshalEnvoy() ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, 100)
	rf := reflect.ValueOf(a)
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
