package types

import (
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/spf13/cast"
)

// resourceRefsExt extends the generated Module.ResourceRefs() with field
// delegation.
func (m Module) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	for _, f := range m.Fields {
		out = append(out, f.ResourceRefs()...)
	}

	return out
}

// resourceRefsExt extends the generated ModuleField.ResourceRefs() with the
// JSON-option module ref.
func (f ModuleField) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	return resourceref.Append(out, resourceref.MakeIdent(
		resourceref.KindComposeModule,
		cast.ToString(f.Options["moduleID"]),
		resourceref.ReasonModuleFieldRef,
	))
}

// resourceRefsExt extends the generated Chart.ResourceRefs() with the report
// module refs.
func (c Chart) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	for _, r := range c.Config.Reports {
		if r == nil {
			continue
		}

		out = resourceref.Append(out, resourceref.Make(
			resourceref.KindComposeModule,
			r.ModuleID,
			resourceref.ReasonChartModule,
		))
	}

	return out
}

// resourceRefsExt extends the generated Page.ResourceRefs() with page-block
// refs.
func (p Page) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	return append(out, p.Blocks.ResourceRefs()...)
}

// ResourceRefs returns configuration-level references found in page blocks
func (bb PageBlocks) ResourceRefs() (out []resourceref.Ref) {
	for _, b := range bb {
		out = append(out, b.resourceRefs()...)
	}

	return
}

func (b PageBlock) resourceRefs() (out []resourceref.Ref) {
	switch b.Kind {
	case "RecordList", "Comment", "RecordOrganizer":
		out = resourceref.Append(out, resourceref.MakeIdent(
			resourceref.KindComposeModule,
			blockOptString(b.Options, "module", "moduleID"),
			resourceref.ReasonPageModule,
		))

	case "Chart":
		out = resourceref.Append(out, resourceref.MakeIdent(
			resourceref.KindComposeChart,
			blockOptString(b.Options, "chart", "chartID"),
			resourceref.ReasonPageChart,
		))

	case "Calendar":
		ff, _ := b.Options["feeds"].([]interface{})
		for _, f := range ff {
			feed, _ := f.(map[string]interface{})
			opt, _ := feed["options"].(map[string]interface{})

			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindComposeModule,
				blockOptString(opt, "module", "moduleID"),
				resourceref.ReasonPageModule,
			))
		}

	case "Metric":
		mm, _ := b.Options["metrics"].([]interface{})
		for _, m := range mm {
			mops, _ := m.(map[string]interface{})

			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindComposeModule,
				blockOptString(mops, "module", "moduleID"),
				resourceref.ReasonPageModule,
			))
		}

	case "Progress":
		for _, k := range []string{"minValue", "maxValue", "value"} {
			opt, _ := b.Options[k].(map[string]interface{})
			if opt == nil {
				continue
			}

			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindComposeModule,
				blockOptString(opt, "module", "moduleID"),
				resourceref.ReasonPageModule,
			))
		}

	case "Automation":
		bb, _ := b.Options["buttons"].([]interface{})
		for _, btn := range bb {
			button, _ := btn.(map[string]interface{})

			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindAutomationWorkflow,
				blockOptString(button, "workflow", "workflowID"),
				resourceref.ReasonPageWorkflow,
			))
		}
	}

	return
}

// blockOptString returns the first present option value as a string
func blockOptString(opt map[string]interface{}, kk ...string) string {
	for _, k := range kk {
		if v, has := opt[k]; has {
			return cast.ToString(v)
		}
	}

	return ""
}
