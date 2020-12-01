package yaml

import (
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

	return ec, eachMap(ecNode, func(k, v *yaml.Node) (err error) {
		switch k.Value {
		case "skipIf", "skip":
			return decodeScalar(v, "decode skip if", &ec.SkipIf)
		}

		return nil
	})
}
