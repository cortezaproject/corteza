package types

import (
	"fmt"

	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/spf13/cast"
)

// ResourceRefs returns configuration-level references to other resources
//
// Field references are included when fields are loaded on the module.
func (m Module) ResourceRefs() (out []resourceref.Ref) {
	out = resourceref.Append(out, resourceref.Make(
		resourceref.KindDalConnection,
		m.Config.DAL.ConnectionID,
		resourceref.ReasonModuleConnection,
		"Config.DAL.ConnectionID",
	))

	for _, f := range m.Fields {
		for _, r := range f.ResourceRefs() {
			r.Path = fmt.Sprintf("Fields.%s.%s", f.Name, r.Path)
			out = append(out, r)
		}
	}

	return
}

// ResourceRefs returns configuration-level references to other resources
func (f ModuleField) ResourceRefs() (out []resourceref.Ref) {
	return resourceref.Append(out, resourceref.MakeIdent(
		resourceref.KindComposeModule,
		cast.ToString(f.Options["moduleID"]),
		resourceref.ReasonModuleFieldRef,
		"Options.ModuleID",
	))
}

// ResourceRefs returns configuration-level references to other resources
func (c Chart) ResourceRefs() (out []resourceref.Ref) {
	for i, r := range c.Config.Reports {
		if r == nil {
			continue
		}

		out = resourceref.Append(out, resourceref.Make(
			resourceref.KindComposeModule,
			r.ModuleID,
			resourceref.ReasonChartModule,
			fmt.Sprintf("Config.Reports.%d.ModuleID", i),
		))
	}

	return
}

// ResourceRefs returns configuration-level references to other resources
func (p Page) ResourceRefs() (out []resourceref.Ref) {
	out = resourceref.Append(out, resourceref.Make(
		resourceref.KindComposeModule,
		p.ModuleID,
		resourceref.ReasonPageModule,
		"ModuleID",
	))

	return append(out, p.Blocks.ResourceRefs()...)
}

// ResourceRefs returns configuration-level references found in page blocks
func (bb PageBlocks) ResourceRefs() (out []resourceref.Ref) {
	for i, b := range bb {
		out = append(out, b.resourceRefs(i)...)
	}

	return
}

func (b PageBlock) resourceRefs(index int) (out []resourceref.Ref) {
	switch b.Kind {
	case "RecordList", "Comment", "RecordOrganizer":
		out = resourceref.Append(out, resourceref.MakeIdent(
			resourceref.KindComposeModule,
			blockOptString(b.Options, "module", "moduleID"),
			resourceref.ReasonPageModule,
			fmt.Sprintf("Blocks.%d.Options.ModuleID", index),
		))

	case "Chart":
		out = resourceref.Append(out, resourceref.MakeIdent(
			resourceref.KindComposeChart,
			blockOptString(b.Options, "chart", "chartID"),
			resourceref.ReasonPageChart,
			fmt.Sprintf("Blocks.%d.Options.ChartID", index),
		))

	case "Calendar":
		ff, _ := b.Options["feeds"].([]interface{})
		for j, f := range ff {
			feed, _ := f.(map[string]interface{})
			opt, _ := feed["options"].(map[string]interface{})

			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindComposeModule,
				blockOptString(opt, "module", "moduleID"),
				resourceref.ReasonPageModule,
				fmt.Sprintf("Blocks.%d.Options.feeds.%d.ModuleID", index, j),
			))
		}

	case "Metric":
		mm, _ := b.Options["metrics"].([]interface{})
		for j, m := range mm {
			mops, _ := m.(map[string]interface{})

			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindComposeModule,
				blockOptString(mops, "module", "moduleID"),
				resourceref.ReasonPageModule,
				fmt.Sprintf("Blocks.%d.Options.metrics.%d.ModuleID", index, j),
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
				fmt.Sprintf("Blocks.%d.Options.%s.ModuleID", index, k),
			))
		}

	case "Automation":
		bb, _ := b.Options["buttons"].([]interface{})
		for j, btn := range bb {
			button, _ := btn.(map[string]interface{})

			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindAutomationWorkflow,
				blockOptString(button, "workflow", "workflowID"),
				resourceref.ReasonPageWorkflow,
				fmt.Sprintf("Blocks.%d.Options.buttons.%d.WorkflowID", index, j),
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
