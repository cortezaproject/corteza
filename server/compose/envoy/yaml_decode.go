package envoy

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	mapEntry struct {
		Column string
		Field  string
		Skip   bool
	}

	fieldMapping struct {
		m map[string]mapEntry
	}

	datasourceMapping struct {
		SourceIdent string `yaml:"source"`
		KeyField    string `yaml:"key"`
		References  map[string]string
		Scope       map[string]string
		Mapping     fieldMapping
	}
)

const (
	ComposeRecordDatasourceAuxType = "corteza::compose:record-datasource"
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

func unmarshalModuleFieldOptionsNode(r *types.ModuleField, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	refs = make(map[string]envoyx.Ref)

	err = y7s.EachMap(n, func(k, v *yaml.Node) error {
		switch strings.ToLower(k.Value) {
		case "module", "module_id", "moduleid":
			var aux any
			err = y7s.DecodeScalar(v, "moduleID", &aux)
			if err != nil {
				return err
			}
			refs["Options.ModuleID"] = envoyx.Ref{
				ResourceType: types.ModuleResourceType,
				Identifiers:  envoyx.MakeIdentifiers(aux),
			}
		default:
			return nil
		}
		return nil
	})

	return
}

func (d *auxYamlDoc) unmarshalYAML(k string, n *yaml.Node) (out envoyx.NodeSet, err error) {
	return
}

func (d *auxYamlDoc) postProcessNestedModuleNodes(nn envoyx.NodeSet) (out envoyx.NodeSet, err error) {
	// Get all references from all module fields
	refs := make(map[string]envoyx.Ref)
	for _, n := range nn {
		if n.ResourceType != types.ModuleFieldResourceType {
			continue
		}

		if len(n.References) == 0 {
			continue
		}

		r, ok := n.References["Options.ModuleID"]
		if !ok {
			continue
		}

		refs[fmt.Sprintf("%s.module", n.Identifiers.FriendlyIdentifier())] = r
		refs[fmt.Sprintf("%s.datasource", n.Identifiers.FriendlyIdentifier())] = envoyx.Ref{
			ResourceType: ComposeRecordDatasourceAuxType,
			Identifiers:  r.Identifiers,
			Scope:        r.Scope,
			Optional:     true,
		}
	}

	// Update datasources with references

	for _, n := range nn {
		if n.Datasource == nil {
			continue
		}

		n.References = envoyx.MergeRefs(n.References, refs)
	}
	out = nn
	return
}

func (d *auxYamlDoc) unmarshalSourceExtendedNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r datasourceMapping

	// @todo we're omitting errors because there will be a bunch due to invalid
	//       resource field types. This might be a bit unstable as other errors may
	//       also get ignored.
	//
	//       A potential fix would be to firstly unmarshal into an any, check errors
	//       and then unmarshal into the resource while omitting errors.
	n.Decode(&r)

	// @todo for now we only support record datasources; extend when needed
	auxN := &envoyx.Node{
		Datasource: &RecordDatasource{
			mapping: r,
		},

		ResourceType: ComposeRecordDatasourceAuxType,
		Identifiers:  envoyx.MakeIdentifiers(r.SourceIdent),
	}
	auxN.References, auxN.Scope = d.procMappingRefs(r.References)
	out = append(out, auxN)

	return
}

// UnmarshalYAML is used to get the yaml parsed into a series of nodes so
// we can easily pass it down
func (d *fieldMapping) UnmarshalYAML(n *yaml.Node) (err error) {
	d.m = make(map[string]mapEntry)
	if y7s.IsSeq(n) {
		err = y7s.EachSeq(n, func(n *yaml.Node) error {
			a, err := d.unmarshalMappingNode(n)
			d.m[a.Column] = a
			return err
		})
	} else {
		err = y7s.EachMap(n, func(k, n *yaml.Node) error {
			a, err := d.unmarshalMappingNode(n)
			if a.Column == "" {
				err = y7s.DecodeScalar(k, "fieldMapping column", &a.Column)
				if err != nil {
					return err
				}
			}

			d.m[a.Column] = a
			return err
		})
	}

	return
}

func (d *fieldMapping) unmarshalMappingNode(n *yaml.Node) (out mapEntry, err error) {
	if y7s.IsKind(n, yaml.ScalarNode) {
		err = y7s.DecodeScalar(n, "Column", &out.Column)
		if err != nil {
			return
		}
		err = y7s.DecodeScalar(n, "Field", &out.Field)
		return
	}

	// @todo we're omitting errors because there will be a bunch due to invalid
	//       resource field types. This might be a bit unstable as other errors may
	//       also get ignored.
	//
	//       A potential fix would be to firstly unmarshal into an any, check errors
	//       and then unmarshal into the resource while omitting errors.
	n.Decode(&out)

	err = y7s.EachMap(n, func(k, v *yaml.Node) error {
		switch strings.ToLower(k.Value) {
		case "skip":
			if v.Value == "/" {
				out.Skip = true
			}
		}
		return nil
	})

	return
}

func (d *auxYamlDoc) procMappingRefs(in map[string]string) (out map[string]envoyx.Ref, scope envoyx.Scope) {
	out = make(map[string]envoyx.Ref)

	scope = envoyx.Scope{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  envoyx.MakeIdentifiers(in["namespace"]),
	}

	out["NamespaceID"] = envoyx.Ref{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  envoyx.MakeIdentifiers(in["namespace"]),
		Scope:        scope,
	}

	out["ModuleID"] = envoyx.Ref{
		ResourceType: types.ModuleResourceType,
		Identifiers:  envoyx.MakeIdentifiers(in["module"]),
		Scope:        scope,
	}

	return
}
