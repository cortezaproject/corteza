package util

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func yamlNodeErr(n *yaml.Node, format string, aa ...interface{}) error {
	format += " (%d:%d)"
	aa = append(aa, n.Line, n.Column)
	return fmt.Errorf(format, aa...)
}

// YamlIterator helps iterate over mapping and sequence nodes fairly trivially
func YamlIterator(n *yaml.Node, fn func(*yaml.Node, *yaml.Node) error) error {
	if n.Kind == yaml.MappingNode {
		for i := 0; i < len(n.Content); i += 2 {
			if err := fn(n.Content[i], n.Content[i+1]); err != nil {
				return err
			}
		}

		return nil
	}

	if n.Kind == yaml.SequenceNode {
		var placeholder *yaml.Node
		for i := 0; i < len(n.Content); i++ {
			if err := fn(placeholder, n.Content[i]); err != nil {
				return err
			}
		}

		return nil
	}

	return yamlNodeErr(n, "expecting mapping or sequence node")
}
