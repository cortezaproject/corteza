package envoy

import (
	"fmt"
)

type (
	Marshaller interface {
		MarshalEnvoy() ([]Node, error)
	}
)

// MarshalMerge takes one or more nodes and Marshals and merges all
func CollectNodes(ii ...interface{}) (nn []Node, err error) {
	for _, i := range ii {
		switch c := i.(type) {
		case NodeSet:
			nn = append(nn, c...)
		case Node:
			nn = append(nn, c)

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
