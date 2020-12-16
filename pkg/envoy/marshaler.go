package envoy

import (
	"fmt"
	"reflect"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	Marshaller interface {
		MarshalEnvoy() ([]resource.Interface, error)
	}
)

// MarshalMerge takes one or more nodes and Marshals and merges all
func CollectNodes(ii ...interface{}) (nn []resource.Interface, err error) {
	for _, i := range ii {
		// i == nil is not enough in go
		if i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil()) {
			continue
		}

		switch c := i.(type) {
		// case NodeSet:
		// 	nn = append(nn, c...)
		case resource.Interface:
			nn = append(nn, c)

		case Marshaller:
			if tmp, err := c.MarshalEnvoy(); err != nil {
				return nil, err
			} else {
				nn = append(nn, tmp...)
			}
		default:
			// @todo resoruceset?
			return nil, fmt.Errorf("failed to merge %T; expecting Resource, ResourceSet? or Marshaller interface", i)
		}
	}

	return nn, nil
}
