package y7s

// y7s (YAML Utils)

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

func NodeErr(n *yaml.Node, format string, aa ...interface{}) error {
	format += " (%d:%d)"
	aa = append(aa, n.Line, n.Column)
	return fmt.Errorf(format, aa...)
}
