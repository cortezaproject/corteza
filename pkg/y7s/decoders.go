package y7s

// y7s (YAML Utils)

import (
	"gopkg.in/yaml.v3"
)

// findKeyNode returns key node from mapping
// key value checked in lower case
func FindKeyNode(n *yaml.Node, key string) *yaml.Node {
	// compare it with lowercase values
	for i := 0; i < len(n.Content); i += 2 {
		if key == n.Content[i].Value {
			return n.Content[i+1]
		}
	}

	return nil
}

// Checks validity of ref node and sets the value to given arg ptr
func DecodeScalar(n *yaml.Node, name string, val interface{}) error {
	if !IsKind(n, yaml.ScalarNode) {
		return NodeErr(n, "expecting scalar value for %s", name)
	}

	return n.Decode(val)
}
