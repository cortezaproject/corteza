package types

import (
	"fmt"

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

		// RecordList selection buttons trigger workflows / ng-automations
		// (same Button shape as the Automation block)
		for _, btn := range blockOptSlice(b.Options, "selectionButtons") {
			out = append(out, buttonRefs(btn)...)
		}

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
		for _, btn := range blockOptSlice(b.Options, "buttons") {
			out = append(out, buttonRefs(btn)...)
		}

	case "AgentChat":
		out = resourceref.Append(out, resourceref.MakeIdent(
			resourceref.KindAgent,
			blockOptString(b.Options, "defaultAgentID"),
			resourceref.ReasonPageAgent,
		))
		for _, id := range cast.ToStringSlice(b.Options["allowedAgentIDs"]) {
			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindAgent, id, resourceref.ReasonPageAgent,
			))
		}

	case "Custom":
		out = resourceref.Append(out, resourceref.MakeIdent(
			resourceref.KindApplication,
			blockOptString(b.Options, "application", "applicationID"),
			resourceref.ReasonPageApplication,
		))

	case "ChatbotInbox":
		for _, id := range cast.ToStringSlice(b.Options["chatbotIDs"]) {
			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindChatbot, id, resourceref.ReasonPageChatbot,
			))
		}

	case "Geometry":
		for _, f := range blockOptSlice(b.Options, "feeds") {
			opt, _ := f["options"].(map[string]interface{})

			out = resourceref.Append(out, resourceref.MakeIdent(
				resourceref.KindComposeModule,
				blockOptString(opt, "module", "moduleID"),
				resourceref.ReasonPageModule,
			))
		}

	case "Navigation":
		for _, it := range blockOptSlice(b.Options, "navigationItems") {
			opt, _ := it["options"].(map[string]interface{})
			item, _ := opt["item"].(map[string]interface{})
			if item == nil {
				continue
			}

			out = resourceref.Append(out,
				resourceref.MakeIdent(resourceref.KindComposePage, blockOptString(item, "pageID"), resourceref.ReasonPageNavigation),
				resourceref.MakeIdent(resourceref.KindComposePageLayout, blockOptString(item, "pageLayoutID"), resourceref.ReasonPageNavigation),
				resourceref.MakeIdent(resourceref.KindComposeModule, blockOptString(item, "moduleID"), resourceref.ReasonPageModule),
			)
		}
	}

	return
}

// LocatedRefs returns the references of page blocks that an import can write
// back, keyed by where each sits in the page ("Blocks.<i>.Options.ModuleID",
// ...), which is the path Page.SetValue takes. Envoy keys a reference by that
// path so the ID it resolves on import lands in the option it came from.
func (bb PageBlocks) LocatedRefs() map[string]resourceref.Ref {
	out := make(map[string]resourceref.Ref)
	put := func(path string, ref resourceref.Ref) {
		if !ref.IsEmpty() {
			out[path] = ref
		}
	}
	module := func(opt map[string]interface{}) resourceref.Ref {
		return resourceref.MakeIdent(resourceref.KindComposeModule, blockOptString(opt, "module", "moduleID"), resourceref.ReasonPageModule)
	}

	for i, b := range bb {
		at := fmt.Sprintf("Blocks.%d.Options.", i)

		switch b.Kind {
		case "RecordList", "Comment", "RecordOrganizer":
			put(at+"ModuleID", module(b.Options))

		case "Chart":
			put(at+"ChartID", resourceref.MakeIdent(resourceref.KindComposeChart, blockOptString(b.Options, "chart", "chartID"), resourceref.ReasonPageChart))

		case "Calendar":
			ff, _ := b.Options["feeds"].([]interface{})
			for j, f := range ff {
				feed, _ := f.(map[string]interface{})
				opt, _ := feed["options"].(map[string]interface{})
				put(fmt.Sprintf("%sfeeds.%d.ModuleID", at, j), module(opt))
			}

		case "Metric":
			mm, _ := b.Options["metrics"].([]interface{})
			for j, m := range mm {
				opt, _ := m.(map[string]interface{})
				put(fmt.Sprintf("%smetrics.%d.ModuleID", at, j), module(opt))
			}

		case "Progress":
			for _, k := range []string{"minValue", "maxValue", "value"} {
				if opt, _ := b.Options[k].(map[string]interface{}); opt != nil {
					put(at+k+".ModuleID", module(opt))
				}
			}

		case "Automation":
			bb, _ := b.Options["buttons"].([]interface{})
			for j, raw := range bb {
				btn, _ := raw.(map[string]interface{})
				put(fmt.Sprintf("%sbuttons.%d.WorkflowID", at, j), resourceref.MakeIdent(resourceref.KindAutomationWorkflow, blockOptString(btn, "workflow", "workflowID"), resourceref.ReasonPageWorkflow))
			}

		case "Custom":
			put(at+"ApplicationID", resourceref.MakeIdent(resourceref.KindApplication, blockOptString(b.Options, "application", "applicationID"), resourceref.ReasonPageApplication))
		}
	}

	return out
}

// LocatedRefs returns the field's module reference keyed by the path
// ModuleField.SetValue takes.
func (f ModuleField) LocatedRefs() map[string]resourceref.Ref {
	out := make(map[string]resourceref.Ref)
	if ref := resourceref.MakeIdent(resourceref.KindComposeModule, cast.ToString(f.Options["moduleID"]), resourceref.ReasonModuleFieldRef); !ref.IsEmpty() {
		out["Options.ModuleID"] = ref
	}
	return out
}

// LocatedRefs returns the chart's report module references keyed by the path
// Chart.SetValue takes.
func (c Chart) LocatedRefs() map[string]resourceref.Ref {
	out := make(map[string]resourceref.Ref)
	for i, r := range c.Config.Reports {
		if r == nil {
			continue
		}
		if ref := resourceref.Make(resourceref.KindComposeModule, r.ModuleID, resourceref.ReasonChartModule); !ref.IsEmpty() {
			out[fmt.Sprintf("Config.Reports.%d.ModuleID", i)] = ref
		}
	}
	return out
}

// buttonRefs emits the workflow + ng-automation refs declared on a page-block
// button. Automation-block buttons and RecordList selection buttons share the
// Button shape (workflowID / automationID).
func buttonRefs(button map[string]interface{}) (out []resourceref.Ref) {
	out = resourceref.Append(out, resourceref.MakeIdent(
		resourceref.KindAutomationWorkflow,
		blockOptString(button, "workflow", "workflowID"),
		resourceref.ReasonPageWorkflow,
	))
	out = resourceref.Append(out, resourceref.MakeIdent(
		resourceref.KindNgAutomation,
		blockOptString(button, "automation", "automationID"),
		resourceref.ReasonPageAutomation,
	))

	return
}

// blockOptSlice returns a page-block option as a slice of string-keyed maps,
// skipping entries that are not objects.
func blockOptSlice(opt map[string]interface{}, key string) []map[string]interface{} {
	raw, _ := opt[key].([]interface{})
	out := make([]map[string]interface{}, 0, len(raw))
	for _, v := range raw {
		if m, ok := v.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}

	return out
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
