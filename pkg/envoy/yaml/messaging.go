package yaml

import (
	"reflect"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	messaging struct {
		Channels messagingChannelSet

		EncoderConfig *EncoderConfig `yaml:"-"`
	}
)

func (c messaging) MarshalEnvoy() ([]resource.Interface, error) {
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

func (c *messaging) MarshalYAML() (interface{}, error) {
	cn, _ := makeMap()
	var err error

	if len(c.Channels) > 0 {
		c.Channels.ConfigureEncoder(c.EncoderConfig)

		cn, err = encodeResource(cn, "channels", c.Channels, c.EncoderConfig.MappedOutput, "name")
		if err != nil {
			return nil, err
		}
	}

	return cn, nil
}
