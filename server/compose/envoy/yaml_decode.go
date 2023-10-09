package envoy

import (
	"fmt"
	"strings"

	automationTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/modern-go/reflect2"
	"gopkg.in/yaml.v3"
)

const (
	ComposeRecordDatasourceAuxType = "corteza::compose:record-datasource"
)

func (d *auxYamlDoc) unmarshalYAML(k string, n *yaml.Node) (out envoyx.NodeSet, err error) { return }

func (d *auxYamlDoc) unmarshalChartConfigNode(r *types.Chart, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
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

				auxRefs, auxIdents, err = d.unmarshalChartConfigReportNode(r, c, i)
				refs = envoyx.MergeRefs(refs, auxRefs)
				idents = idents.Merge(auxIdents)
				return err
			})
			if err != nil {
				return err
			}
		} else {
			refs, idents, err = d.unmarshalChartConfigReportNode(r, v, 0)
			return err
		}
		return nil
	})

	return
}

func (d *auxYamlDoc) unmarshalChartConfigReportNode(r *types.Chart, n *yaml.Node, index int) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
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

func (d *auxYamlDoc) unmarshalPageBlocksNode(r *types.Page, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	refs = map[string]envoyx.Ref{}

	for index, b := range r.Blocks {
		switch b.Kind {
		case "RecordList":
			refs = envoyx.MergeRefs(refs, getPageBlockRecordListRefs(b, index))

		case "Automation":
			refs = envoyx.MergeRefs(refs, getPageBlockAutomationRefs(b, index))

		case "RecordOrganizer":
			refs = envoyx.MergeRefs(refs, getPageBlockRecordOrganizerRefs(b, index))

		case "Chart":
			refs = envoyx.MergeRefs(refs, getPageBlockChartRefs(b, index))

		case "Calendar":
			refs = envoyx.MergeRefs(refs, getPageBlockCalendarRefs(b, index))

		case "Metric":
			refs = envoyx.MergeRefs(refs, getPageBlockMetricRefs(b, index))

		case "Comment":
			refs = envoyx.MergeRefs(refs, getPageBlockCommentRefs(b, index))

		case "Progress":
			refs = envoyx.MergeRefs(refs, getPageBlockProgressRefs(b, index))
		}
	}

	return
}

func getPageBlockRecordListRefs(b types.PageBlock, index int) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref)

	id := optString(b.Options, "module", "moduleID")
	if id == "" || id == "0" {
		return
	}

	refs[fmt.Sprintf("Blocks.%d.Options.ModuleID", index)] = envoyx.Ref{
		ResourceType: types.ModuleResourceType,
		Identifiers:  envoyx.MakeIdentifiers(id),
	}

	return
}

func getPageBlockChartRefs(b types.PageBlock, index int) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref)

	id := optString(b.Options, "chart", "chartID")
	if id == "" || id == "0" {
		return
	}

	refs[fmt.Sprintf("Blocks.%d.Options.ChartID", index)] = envoyx.Ref{
		ResourceType: types.ChartResourceType,
		Identifiers:  envoyx.MakeIdentifiers(id),
	}

	return
}

func getPageBlockCalendarRefs(b types.PageBlock, index int) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref)

	ff, _ := b.Options["feeds"].([]interface{})
	for j, f := range ff {
		feed, _ := f.(map[string]interface{})
		opt, _ := (feed["options"]).(map[string]interface{})

		id := optString(opt, "module", "moduleID")
		if id == "" || id == "0" {
			return
		}

		refs[fmt.Sprintf("Blocks.%d.Options.feeds.%d.ModuleID", index, j)] = envoyx.Ref{
			ResourceType: types.ModuleResourceType,
			Identifiers:  envoyx.MakeIdentifiers(id),
		}
	}

	return
}

func getPageBlockMetricRefs(b types.PageBlock, index int) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref)

	mm, _ := b.Options["metrics"].([]interface{})
	for j, m := range mm {
		mops, _ := m.(map[string]interface{})

		id := optString(mops, "module", "moduleID")
		if id == "" || id == "0" {
			return
		}

		refs[fmt.Sprintf("Blocks.%d.Options.metrics.%d.ModuleID", index, j)] = envoyx.Ref{
			ResourceType: types.ModuleResourceType,
			Identifiers:  envoyx.MakeIdentifiers(id),
		}
	}

	return
}

func getPageBlockCommentRefs(b types.PageBlock, index int) (refs map[string]envoyx.Ref) {
	// Same difference
	return getPageBlockRecordListRefs(b, index)
}

func getPageBlockProgressRefs(b types.PageBlock, index int) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref)
	var aux *envoyx.Ref

	aux = getPageBlockProgressValueRefs(b.Options["minValue"])
	if aux != nil {
		refs[fmt.Sprintf("Blocks.%d.Options.minValue.ModuleID", index)] = *aux
	}

	aux = getPageBlockProgressValueRefs(b.Options["maxValue"])
	if aux != nil {
		refs[fmt.Sprintf("Blocks.%d.Options.maxValue.ModuleID", index)] = *aux
	}

	aux = getPageBlockProgressValueRefs(b.Options["value"])
	if aux != nil {
		refs[fmt.Sprintf("Blocks.%d.Options.value.ModuleID", index)] = *aux
	}

	return
}

func getPageBlockProgressValueRefs(val any) (ref *envoyx.Ref) {
	if reflect2.IsNil(val) {
		return
	}

	opt, ok := val.(map[string]interface{})
	if !ok {
		return
	}

	id := optString(opt, "module", "moduleID")
	if id == "" || id == "0" {
		return
	}

	return &envoyx.Ref{
		ResourceType: types.ModuleResourceType,
		Identifiers:  envoyx.MakeIdentifiers(id),
	}
}

func getPageBlockRecordOrganizerRefs(b types.PageBlock, index int) (refs map[string]envoyx.Ref) {
	// Same difference
	return getPageBlockRecordListRefs(b, index)
}

func getPageBlockAutomationRefs(b types.PageBlock, index int) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref)

	bb, _ := b.Options["buttons"].([]interface{})
	for buttonIx, b := range bb {
		button, _ := b.(map[string]interface{})
		id := optString(button, "workflow", "workflowID")
		if id == "" || id == "0" {
			return
		}

		refs[fmt.Sprintf("Blocks.%d.Options.buttons.%d.WorkflowID", index, buttonIx)] = envoyx.Ref{
			ResourceType: automationTypes.WorkflowResourceType,
			Identifiers:  envoyx.MakeIdentifiers(id),
		}
	}

	return
}

func optString(opt map[string]interface{}, kk ...string) string {
	for _, k := range kk {
		if vr, has := opt[k]; has {
			v, _ := vr.(string)
			return v
		}
	}
	return ""
}

func (d *auxYamlDoc) unmarshalModuleFieldOptionsNode(r *types.ModuleField, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
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

func (d *auxYamlDoc) unmarshalModuleFieldDefaultValueNode(r *types.ModuleField, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	var rvs = types.RecordValueSet{}

	switch n.Kind {
	case yaml.ScalarNode:
		rvs = rvs.Set(&types.RecordValue{Value: n.Value})

	case yaml.SequenceNode:
		_ = y7s.EachSeq(n, func(v *yaml.Node) error {
			rvs = rvs.Set(&types.RecordValue{Value: n.Value, Place: uint(len(rvs))})
			return nil
		})
	}

	r.DefaultValue = rvs
	return
}

func (d *auxYamlDoc) unmarshalModuleFieldExpressionsNode(r *types.ModuleField, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {

	err = y7s.EachMap(n, func(k *yaml.Node, v *yaml.Node) error {
		switch k.Value {
		case "sanitizer":
			var aux string
			err = y7s.DecodeScalar(v, "sanitizer", &aux)
			if err != nil {
				return err
			}
			r.Expressions.Sanitizers = append(r.Expressions.Sanitizers, aux)

		case "sanitizers":
			return y7s.EachSeq(v, func(san *yaml.Node) error {
				r.Expressions.Sanitizers = append(r.Expressions.Sanitizers, san.Value)
				return nil
			})

		case "validator":
			var aux types.ModuleFieldValidator
			err = v.Decode(&aux)
			if err != nil {
				return err
			}
			r.Expressions.Validators = append(r.Expressions.Validators, aux)

		case "validators":
			return y7s.Each(v, func(k *yaml.Node, v *yaml.Node) error {
				var aux types.ModuleFieldValidator
				if y7s.IsKind(v, yaml.MappingNode) {
					err = v.Decode(&aux)
					if err != nil {
						return err
					}
				} else {
					aux.Test = k.Value
					aux.Error = v.Value
				}

				r.Expressions.Validators = append(r.Expressions.Validators, aux)
				return nil
			})

		case "formatter":
			r.Expressions.Formatters = append(r.Expressions.Formatters, v.Value)
			return nil

		case "formatters":
			return y7s.EachSeq(v, func(san *yaml.Node) error {
				r.Expressions.Formatters = append(r.Expressions.Formatters, san.Value)
				return nil
			})
		}
		return nil
	})

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

func (d *auxYamlDoc) unmarshalPagesExtendedNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	return d.unmarshalPageNode(dctx, n, meta...)
}

func (d *auxYamlDoc) unmarshalSourceExtendedNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r envoyx.DatasourceMapping

	// @todo we're omitting errors because there will be a bunch due to invalid
	//       resource field types. This might be a bit unstable as other errors may
	//       also get ignored.
	//
	//       A potential fix would be to firstly unmarshal into an any, check errors
	//       and then unmarshal into the resource while omitting errors.
	n.Decode(&r)

	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		var auxNodeValue any
		_ = auxNodeValue

		switch strings.ToLower(k.Value) {
		case "origin", "from":
			err = y7s.DecodeScalar(n, "origin", &r.SourceIdent)
			if err != nil {
				return err
			}

		case "key", "index", "pk":
			if !y7s.IsKind(n, yaml.SequenceNode) {
				r.KeyField = []string{n.Value}
			} else {
				r.KeyField = make([]string, 0, 3)
				y7s.EachSeq(n, func(n *yaml.Node) error {
					r.KeyField = append(r.KeyField, n.Value)
					return nil
				})
			}

		case "map":
			return n.Decode(&r.Mapping)
		}

		return nil
	})

	// @todo for now we only support record datasources; extend when needed
	auxN := &envoyx.Node{
		Datasource: &RecordDatasource{
			Mapping: r,
		},

		ResourceType: ComposeRecordDatasourceAuxType,
		Identifiers:  dctx.parentIdent,
	}
	auxN.References, auxN.Scope = d.procMappingRefs(r.References)
	out = append(out, auxN)

	return
}

func (d *auxYamlDoc) procMappingRefs(in map[string]string) (out map[string]envoyx.Ref, scope envoyx.Scope) {
	out = make(map[string]envoyx.Ref)

	nsII := envoyx.MakeIdentifiers(in["namespace"])
	if len(nsII.Slice) > 0 {
		scope = envoyx.Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  nsII,
		}
	}

	if len(nsII.Slice) > 0 {
		out["NamespaceID"] = envoyx.Ref{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  nsII,
			Scope:        scope,
		}
	}

	modII := envoyx.MakeIdentifiers(in["module"])
	if len(modII.Slice) > 0 {
		out["ModuleID"] = envoyx.Ref{
			ResourceType: types.ModuleResourceType,
			Identifiers:  envoyx.MakeIdentifiers(in["module"]),
			Scope:        scope,
		}
	}

	return
}
