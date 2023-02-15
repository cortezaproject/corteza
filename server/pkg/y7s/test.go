package y7s

// y7s (YAML Utils)

import (
	"gopkg.in/yaml.v3"
)

func IsSeq(n *yaml.Node) bool {
	return n.Kind == yaml.SequenceNode
}

func IsMapping(n *yaml.Node) bool {
	return n.Kind == yaml.MappingNode
}

func IsKind(n *yaml.Node, tt ...yaml.Kind) bool {
	if n != nil {
		for _, t := range tt {
			if t == n.Kind {
				return true
			}
		}
	}

	return false
}
