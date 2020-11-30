package yaml

import (
	"reflect"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	messaging struct {
		Channels messagingChannelSet `yaml:"channels"`
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
