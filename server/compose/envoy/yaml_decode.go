package envoy

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

func unmarshalChartConfigNode(r *types.Chart, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	err = y7s.EachMap(n, func(k, v *yaml.Node) error {
		if k.Value != "reports" {
			return nil
		}

		if y7s.IsSeq(v) {
			var (
				auxRefs   = make(map[string]envoyx.Ref)
				auxIdents envoyx.Identifiers
				i         = -1
			)
			err = y7s.EachSeq(v, func(c *yaml.Node) error {
				i++

				auxRefs, auxIdents, err = unmarshalChartConfigReportNode(r, c, i)
				refs = envoyx.MergeRefs(refs, auxRefs)
				idents = idents.Merge(auxIdents)
				return err
			})
			if err != nil {
				return err
			}
		} else {
			refs, idents, err = unmarshalChartConfigReportNode(r, v, 0)
			return err
		}
		return nil
	})

	return
}

func unmarshalChartConfigReportNode(r *types.Chart, n *yaml.Node, index int) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	err = y7s.EachMap(n, func(k, v *yaml.Node) error {
		switch strings.ToLower(k.Value) {
		case "module", "mod", "moduleid", "module_id":
			var auxi any
			y7s.DecodeScalar(v, "moduleID", &auxi)
			refs = map[string]envoyx.Ref{
				fmt.Sprintf("Config.Reports.%d.ModuleID", index): {
					ResourceType: types.ModuleResourceType,
					Identifiers:  envoyx.MakeIdentifiers(auxi),
				},
			}
		}
		return nil
	})
	return
}

func (d *auxYamlDoc) unmarshalYAML(k string, n *yaml.Node) (out envoyx.NodeSet, err error) {
	return
}
