package y7s

// y7s (YAML Utils)

import (
	"gopkg.in/yaml.v3"
)

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
