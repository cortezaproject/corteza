package y7s

// y7s (YAML Utils)

import (
	"gopkg.in/yaml.v3"
)

// eachKV iterates over map node
func EachMap(n *yaml.Node, fn func(*yaml.Node, *yaml.Node) error) error {
	if !IsKind(n, yaml.MappingNode) {
		// root node kind be mapping
		return NodeErr(n, "expecting mapping node")
	}

	for i := 0; i < len(n.Content); i += 2 {
		if err := fn(n.Content[i], n.Content[i+1]); err != nil {
			return err
		}
	}

	return nil
}

// eachKV iterates over sequence node
func EachSeq(n *yaml.Node, fn func(*yaml.Node) error) error {
	if !IsKind(n, yaml.SequenceNode) {
		// root node kind be mapping
		return NodeErr(n, "expecting sequence node")
	}

	for i := 0; i < len(n.Content); i++ {
		if err := fn(n.Content[i]); err != nil {
			return err
		}
	}

	return nil
}

// each helps iterate over mapping and sequence nodes fairly trivially
func Each(n *yaml.Node, fn func(*yaml.Node, *yaml.Node) error) error {
	if IsKind(n, yaml.MappingNode) {
		return EachMap(n, fn)
	}

	if IsKind(n, yaml.SequenceNode) {
		var placeholder *yaml.Node
		for i := 0; i < len(n.Content); i++ {
			if err := fn(placeholder, n.Content[i]); err != nil {
				return err
			}
		}

		return nil
	}

	return NodeErr(n, "expecting mapping or sequence node")
}
