package importer

import (
	"gopkg.in/yaml.v2"
)

// Handles importing of YAML structures for compose
func ParseYAML(in []byte) (aux interface{}, err error) {
	return aux, yaml.Unmarshal(in, &aux)
}
