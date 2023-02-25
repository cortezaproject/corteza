package envoy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
)

func (e YamlEncoder) encodeChartConfigC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, chart *types.Chart, cfg types.ChartConfig) (_ any, err error) {

	reports, _ := y7s.MakeSeq()

	for i, r := range cfg.Reports {
		modRef, ok := n.References[fmt.Sprintf("Config.Reports.%d.ModuleID", i)]
		if !ok {
			continue
		}

		mNode := tt.ParentForRef(n, modRef)
		if mNode == nil {
			err = fmt.Errorf("module for ref not found")
			return
		}

		r, err := y7s.MakeMap(
			"filter", r.Filter,
			"module", mNode.Identifiers.FriendlyIdentifier(),
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

	return reports, nil
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
		mNode := tt.ParentForRef(n, n.References["Options.ModuleID"])
		if mNode == nil {
			err = fmt.Errorf("module for ref not found")
			return
		}
		opt["module"] = mNode.Identifiers.FriendlyIdentifier()
		delete(opt, "moduleID")

	case "User":
		aux := make([]string, 0, 2)
		for i := range opt.Strings("roles") {
			rNode := tt.ParentForRef(n, n.References[fmt.Sprintf("Options.RoleID.%d", i)])
			if rNode == nil {
				err = fmt.Errorf("role for ref not found")
				return
			}

			aux = append(aux, rNode.Identifiers.FriendlyIdentifier())
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
	out, _ := y7s.MakeSeq()

	var aux any
	for i, b := range pg.Blocks {
		aux, err = e.encodePageBlockC(ctx, p, tt, n, pg, i, b)
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

func (e YamlEncoder) encodePageBlockC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, pg *types.Page, index int, b types.PageBlock) (_ any, err error) {

	switch b.Kind {
	case "RecordList":
		b = e.cleanupPageblockRecordList(b)

		node := tt.ParentForRef(n, n.References[fmt.Sprintf("Blocks.%d.Options.ModuleID", index)])
		if node == nil {
			err = fmt.Errorf("module for ref not found")
			return
		}
		b.Options["module"] = node.Identifiers.FriendlyIdentifier()
		delete(b.Options, "moduleID")
		break

	case "RecordOrganizer":
		node := tt.ParentForRef(n, n.References[fmt.Sprintf("Blocks.%d.Options.ModuleID", index)])
		if node == nil {
			err = fmt.Errorf("module for ref not found")
			return
		}
		b.Options["module"] = node.Identifiers.FriendlyIdentifier()
		delete(b.Options, "moduleID")
		break

	case "Chart":
		node := tt.ParentForRef(n, n.References[fmt.Sprintf("Blocks.%d.Options.ChartID", index)])
		if node == nil {
			err = fmt.Errorf("chart for ref not found")
			return
		}
		b.Options["chart"] = node.Identifiers.FriendlyIdentifier()
		delete(b.Options, "chartID")
		break

	case "Calendar":
		ff, _ := b.Options["feeds"].([]interface{})
		for i, f := range ff {
			feed, _ := f.(map[string]interface{})
			fOpts, _ := (feed["options"]).(map[string]interface{})

			node := tt.ParentForRef(n, n.References[fmt.Sprintf("Blocks.%d.Options.feeds.%d.ModuleID", index, i)])
			if node == nil {
				err = fmt.Errorf("module for ref not found")
				return
			}
			fOpts["module"] = node.Identifiers.FriendlyIdentifier()
			delete(fOpts, "moduleID")
		}
		break

	case "Automation":
		bb, _ := b.Options["buttons"].([]interface{})
		for i, b := range bb {
			button, _ := b.(map[string]interface{})
			if _, has := button["workflowID"]; !has {
				continue
			}

			node := tt.ParentForRef(n, n.References[fmt.Sprintf("Blocks.%d.Options.buttons.%d.WorkflowID", index, i)])
			if node == nil {
				err = fmt.Errorf("chart for ref not found")
				return
			}

			button["workflow"] = node.Identifiers.FriendlyIdentifier()
			delete(button, "workflowID")
			i++
		}
		break

	case "Metric":
		mm, _ := b.Options["metrics"].([]interface{})
		for i, m := range mm {
			node := tt.ParentForRef(n, n.References[fmt.Sprintf("Blocks.%d.Options.metrics.%d.ModuleID", index, i)])
			if node == nil {
				err = fmt.Errorf("chart for ref not found")
				return
			}

			mops, _ := m.(map[string]interface{})
			mops["module"] = node.Identifiers.FriendlyIdentifier()
			delete(mops, "moduleID")
		}
		break

	case "Comment":
		node := tt.ParentForRef(n, n.References[fmt.Sprintf("Blocks.%d.Options.ModuleID", index)])
		if node == nil {
			err = fmt.Errorf("module for ref not found")
			return
		}
		b.Options["module"] = node.Identifiers.FriendlyIdentifier()
		delete(b.Options, "moduleID")
		break
	}

	return
}

func (e YamlEncoder) cleanupPageblockRecordList(b types.PageBlock) (_ types.PageBlock) {
	rawFF, has := b.Options["fields"]
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

	b.Options["fields"] = retFF

	return b
}
