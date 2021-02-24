package yaml

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/y7s"
	. "github.com/cortezaproject/corteza-server/pkg/y7s"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"gopkg.in/yaml.v3"
)

func decodeEnvoyConfig(n *yaml.Node) (*resource.EnvoyConfig, error) {
	ec := &resource.EnvoyConfig{}

	var ecNode *yaml.Node
	for i, k := range n.Content {
		if k.Value == "(envoy)" {
			ecNode = n.Content[i+1]
			break
		}
	}

	if ecNode == nil {
		return nil, nil
	}

	return ec, y7s.EachMap(ecNode, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "skipIf", "skip":
			return y7s.DecodeScalar(v, "decode skip if", &ec.SkipIf)
		case "onExisting", "mergeAlg":
			return decodeMergeAlg(v, "decode merge alg", &ec.OnExisting)
		}

		return nil
	})
}

func decodeMergeAlg(n *yaml.Node, refType string, val *resource.MergeAlg) error {
	if n == nil {
		return nil
	}

	if !IsKind(n, yaml.ScalarNode) {
		return y7s.NodeErr(n, "%s reference must be scalar", refType)
	}

	switch strings.ToLower(n.Value) {
	case "skip",
		"s":
		*val = resource.Skip
	case "replace",
		"r":
		*val = resource.Replace
	case "mergeleft",
		"left",
		"ml":
		*val = resource.MergeLeft
	case "mergeright",
		"right",
		"mr":
		*val = resource.MergeRight
	default:
		return y7s.NodeErr(n, "%s unknown algorithm", refType)
	}

	return nil
}
