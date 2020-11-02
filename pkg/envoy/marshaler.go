package envoy

import (
	"fmt"

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
		switch c := i.(type) {
		// case NodeSet:
		// 	nn = append(nn, c...)
		// case Node:
		// 	nn = append(nn, c)

		case Marshaller:
			if tmp, err := c.MarshalEnvoy(); err != nil {
				println(err)
				return nil, err
			} else {
				tmp = append(nn, tmp...)
			}
		default:
			return nil, fmt.Errorf("failed to merge %T; expecting Node, NodeSet or Marshaller interface", i)
		}
	}

	return nn, nil
}
