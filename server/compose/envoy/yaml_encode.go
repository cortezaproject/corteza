package envoy

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/modern-go/reflect2"
	"gopkg.in/yaml.v3"
)

func (e YamlEncoder) encode(ctx context.Context, base *yaml.Node, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	return
}

func (e YamlEncoder) encodeChartConfigC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, chart *types.Chart, cfg types.ChartConfig) (_ any, err error) {

	reports, _ := y7s.MakeSeq()

	for i, r := range cfg.Reports {
		modRef, ok := n.References[fmt.Sprintf("Config.Reports.%d.ModuleID", i)]
		if !ok {
			continue
		}

		r, err := y7s.MakeMap(
			"filter", r.Filter,
			"module", safeParentIdentifier(tt, n, modRef),
			"metrics", r.Metrics,
			"dimensions", r.Dimensions,
			"yAxis", r.YAxis,
		)
		if err != nil {
			return nil, err
		}

		reports, err = y7s.AddSeq(reports, r)
		if err != nil {
			return nil, err
		}
	}

	return y7s.MakeMap(
		"reports", reports,
		"colorScheme", cfg.ColorScheme,
	)
}

func (e YamlEncoder) encodeModuleFieldsC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, mod *types.Module, fields types.ModuleFieldSet) (_ any, err error) {

	fn := tt.ChildrenForResourceType(n, types.ModuleFieldResourceType)

	out, err := e.encodeModuleFields(ctx, p, fn, tt)
	return out, err
}

func (e YamlEncoder) encodeModuleFieldOptionsC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, f *types.ModuleField, opt types.ModuleFieldOptions) (_ any, err error) {
	if opt == nil {
		opt = make(types.ModuleFieldOptions)
	}
	switch f.Kind {
	case "Record":
		modRef := n.References["Options.ModuleID"]
		opt["module"] = safeParentIdentifier(tt, n, modRef)
		delete(opt, "moduleID")

	case "User":
		aux := make([]string, 0, 2)
		for i := range opt.Strings("roles") {
			roleRef := n.References[fmt.Sprintf("Options.RoleID.%d", i)]
			aux = append(aux, safeParentIdentifier(tt, n, roleRef))
		}
		opt["roles"] = aux
		delete(opt, "role")
		delete(opt, "roleID")
	}

	nopt, _ := y7s.MakeMap()
	for k, v := range opt {
		nopt, err = y7s.AddMap(nopt, k, v)
		if err != nil {
			return nil, err
		}
	}

	return nopt, nil
}

func (e YamlEncoder) encodePageBlocksC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, pg *types.Page, bb types.PageBlocks) (_ any, err error) {
	var aux any
	out, _ := y7s.MakeSeq()

	for i, b := range pg.Blocks {
		aux, err = e.encodePageBlockC(ctx, p, tt, n, i, b.Kind, b.Options)
		if err != nil {
			return
		}

		out, err = y7s.AddSeq(out, aux)
		if err != nil {
			return
		}
	}

	return out, nil
}

func (e YamlEncoder) encodeNamespaceBlocksC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, ns *types.Namespace, bb types.GlobalBlocks) (_ any, err error) {
	var aux any
	out, _ := y7s.MakeSeq()

	for i, gB := range ns.Blocks {
		aux, err = e.encodePageBlockC(ctx, p, tt, n, i, gB.Kind, gB.Options)
		if err != nil {
			return
		}

		out, err = y7s.AddSeq(out, aux)
		if err != nil {
			return
		}
	}

	return out, nil
}

func (e YamlEncoder) encodePageBlockC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, index int, kind string, options map[string]interface{}) (_ any, err error) {
	switch kind {
	case "RecordList":
		options = e.cleanupPageblockRecordList(options)

		modRef := n.References[fmt.Sprintf("Blocks.%d.Options.ModuleID", index)]
		options["module"] = safeParentIdentifier(tt, n, modRef)
		delete(options, "moduleID")
		break

	case "RecordOrganizer":
		modRef := n.References[fmt.Sprintf("Blocks.%d.Options.ModuleID", index)]
		options["module"] = safeParentIdentifier(tt, n, modRef)
		delete(options, "moduleID")
		break

	case "Chart":
		chrRef := n.References[fmt.Sprintf("Blocks.%d.Options.ChartID", index)]
		options["chart"] = safeParentIdentifier(tt, n, chrRef)
		delete(options, "chartID")
		break

	case "Calendar":
		ff, _ := options["feeds"].([]interface{})
		for i, f := range ff {
			feed, _ := f.(map[string]interface{})
			fOpts, _ := (feed["options"]).(map[string]interface{})

			modRef := n.References[fmt.Sprintf("Blocks.%d.Options.feeds.%d.ModuleID", index, i)]
			fOpts["module"] = safeParentIdentifier(tt, n, modRef)
			delete(fOpts, "moduleID")
		}
		break

	case "Automation":
		bb, _ := options["buttons"].([]interface{})
		for i, b := range bb {
			button, _ := b.(map[string]interface{})
			if _, has := button["workflowID"]; !has {
				continue
			}

			wfRef := n.References[fmt.Sprintf("Blocks.%d.Options.buttons.%d.WorkflowID", index, i)]
			button["workflow"] = safeParentIdentifier(tt, n, wfRef)
			delete(button, "workflowID")
			i++
		}
		break

	case "Metric":
		mm, _ := options["metrics"].([]interface{})
		for i, m := range mm {
			modRef := n.References[fmt.Sprintf("Blocks.%d.Options.metrics.%d.ModuleID", index, i)]

			mops, _ := m.(map[string]interface{})
			mops["module"] = safeParentIdentifier(tt, n, modRef)
			delete(mops, "moduleID")
		}
		break

	case "Comment":
		modRef := n.References[fmt.Sprintf("Blocks.%d.Options.ModuleID", index)]
		options["module"] = safeParentIdentifier(tt, n, modRef)
		delete(options, "moduleID")
		break

	case "Progress":
		err = e.encodeProgressPageblockVal("minValue", index, n, tt, options)
		if err != nil {
			return
		}

		err = e.encodeProgressPageblockVal("maxValue", index, n, tt, options)
		if err != nil {
			return
		}

		err = e.encodeProgressPageblockVal("value", index, n, tt, options)
		if err != nil {
			return
		}
		break
	}

	return
}

func (e YamlEncoder) encodeProgressPageblockVal(k string, index int, n *envoyx.Node, tt envoyx.Traverser, options map[string]interface{}) (err error) {
	if reflect2.IsNil(options[k]) {
		return
	}

	modRef := n.References[fmt.Sprintf("Blocks.%d.Options.%s.ModuleID", index, k)]
	opt := options[k].(map[string]any)
	opt["moduleID"] = safeParentIdentifier(tt, n, modRef)
	delete(opt, "moduleID")

	return
}

func (e YamlEncoder) cleanupPageblockRecordList(options map[string]interface{}) (out map[string]interface{}) {
	out = options
	rawFF, has := out["fields"]
	if !has {
		return
	}

	ff, ok := rawFF.([]interface{})
	if !ok {
		return
	}

	retFF := make([]interface{}, 0, len(ff))
	for _, rawF := range ff {
		switch c := rawF.(type) {
		case string:
			retFF = append(retFF, map[string]interface{}{"name": c})
		case map[string]interface{}, map[string]string:
			retFF = append(retFF, c)
		default:
			retFF = append(retFF, rawF)
		}
	}

	out["fields"] = retFF

	return
}
