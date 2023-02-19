package envoy

import (
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"gopkg.in/yaml.v3"
)

func (d *auxYamlDoc) unmarshalYAML(k string, n *yaml.Node) (out envoyx.NodeSet, err error) {
	return
}
