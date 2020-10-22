package yaml

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"strings"

	"gopkg.in/yaml.v3"
)

func nodeErr(n *yaml.Node, format string, aa ...interface{}) error {
	format += " (%d:%d)"
	aa = append(aa, n.Line, n.Column)
	return fmt.Errorf(format, aa...)
}

// iterator helps iterate over mapping and sequence nodes fairly trivially
func iterator(n *yaml.Node, fn func(*yaml.Node, *yaml.Node) error) error {
	if isKind(n, yaml.MappingNode) {
		for i := 0; i < len(n.Content); i += 2 {
			if err := fn(n.Content[i], n.Content[i+1]); err != nil {
				return err
			}
		}

		return nil
	}

	if isKind(n, yaml.SequenceNode) {
		var placeholder *yaml.Node
		for i := 0; i < len(n.Content); i++ {
			if err := fn(placeholder, n.Content[i]); err != nil {
				return err
			}
		}

		return nil
	}

	return nodeErr(n, "expecting mapping or sequence node")
}

func isKind(n *yaml.Node, tt ...yaml.Kind) bool {
	if n != nil {
		for _, t := range tt {
			if t == n.Kind {
				return true
			}
		}
	}

	return false
}

// returns true if one or more keys passed in as argument are present under mapping node
func anyKeysPresent(n *yaml.Node, kk ...string) bool {
	// make a map with lowercase strings
	h := make(map[string]bool)
	for i := 0; i < len(kk); i++ {
		h[strings.ToLower(kk[i])] = true
	}

	// compare it with lowercase values
	for i := 0; i < len(n.Content); i += 2 {
		if !h[n.Content[i+1].Value] {
			return false
		}
	}

	return true
}

// findKeyNode returns key node from mapping
// key value checked in lower case
func findKeyNode(n *yaml.Node, key string) *yaml.Node {
	// compare it with lowercase values
	for i := 0; i < len(n.Content); i += 2 {
		if key == n.Content[i].Value {
			return n.Content[i+1]
		}
	}

	return nil
}

// Checks validity of ref node and sets the value to given arg ptr
func decodeScalar(n *yaml.Node, name string, val interface{}) error {
	if !isKind(n, yaml.ScalarNode) {
		return nodeErr(n, "expecting scalar value for %s", name)
	}

	return n.Decode(val)
}

// Checks validity of ref node and sets the value to given arg ptr
func decodeRef(n *yaml.Node, refType string, ref *string) error {
	if !isKind(n, yaml.ScalarNode) {
		return nodeErr(n, "%s reference must be scalar", refType)
	}

	if !handle.IsValid(n.Value) {
		return nodeErr(n, "%s reference must be a valid handle", refType)
	}

	*ref = n.Value
	return nil
}
