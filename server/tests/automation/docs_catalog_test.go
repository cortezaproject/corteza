package automation

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/crusttech/human/server/automation/service"
	"github.com/crusttech/human/server/automation/types"
	"gopkg.in/yaml.v3"
)

// Connection-defined constructs exist per instance, so the reference leaves them out.
const (
	connFunctionPrefix = "conn_"
	connTriggerPrefix  = "corteza::system:connection-webhook/"
)

type (
	docInput struct {
		label, types, desc string
		required           bool
	}

	docOutput struct {
		name, types, desc string
	}

	docEntry struct {
		group, label, anchor, meta, desc, note string
		weight                                 int
		inputs                                 []docInput
		outputs                                []docOutput
	}
)

// TestDocsCatalog holds docs/reference/taq/*.gen.md to the construct catalog the
// builder reads; HUMAN_UPDATE_DOCS=1 rewrites the pages.
func TestDocsCatalog(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..", "..")
	dir := filepath.Join(root, "docs", "reference", "taq")

	branches := loadBranchStrings(t, filepath.Join(root, "locale", "en", "human-webapp", "builder.yaml"))

	pages := map[string]string{
		"triggers.gen.md": renderTriggersPage(t, triggerEntries(service.ConstructLibrary().Triggers())),
		"steps.gen.md":    renderStepsPage(t, append(branches, stepEntries(service.ConstructLibrary().Functions())...)),
	}

	names := []string{"steps.gen.md", "triggers.gen.md"}
	if os.Getenv("HUMAN_UPDATE_DOCS") == "1" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, n := range names {
			if err := os.WriteFile(filepath.Join(dir, n), []byte(pages[n]), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return
	}

	for _, n := range names {
		have, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil || !bytes.Equal(have, []byte(pages[n])) {
			t.Errorf("docs/reference/taq/%s is stale against the construct catalog; run make docs-reference", n)
		}
	}
}

func triggerEntries(tt []types.ConstructTrigger) (out []docEntry) {
	for _, tr := range tt {
		if strings.HasPrefix(tr.ResourceType, connTriggerPrefix) {
			continue
		}

		e := docEntry{
			group: firstOr(tr.Groups, "General"),
			label: tr.EventType,
			meta:  fmt.Sprintf("`%s` · `%s`", tr.ResourceType, tr.EventType),
		}
		if tr.Meta != nil {
			e.label = orStr(tr.Meta.Short, tr.EventType)
			e.desc = tr.Meta.Description
		}
		e.anchor = slug(e.label)

		ctypes := map[string][]string{}
		for _, c := range tr.Constraints {
			ctypes[c.Name] = c.Types
		}
		for _, in := range segmentInputs(tr.Segments) {
			e.inputs = append(e.inputs, docInput{
				label:    orStr(in.Label, in.Argument),
				types:    typeList(ctypes[in.Argument]),
				required: in.Required,
				desc:     withOptions(in.Description, in),
			})
		}

		for _, p := range tr.Properties {
			e.outputs = append(e.outputs, docOutput{
				name:  p.Name,
				types: typeList([]string{p.Type}),
				desc:  orStr(p.Meta.Description, p.Meta.Short),
			})
		}

		if tr.EventType == "onAgentic" {
			e.note = "The values the agent passes are declared per TAQ, in the trigger's input schema."
		}

		out = append(out, e)
	}
	return
}

func stepEntries(ff []types.ConstructFunction) (out []docEntry) {
	for _, fn := range ff {
		if fn.Kind == "gateway" || strings.HasPrefix(fn.Ref, connFunctionPrefix) {
			continue
		}

		e := docEntry{
			group:  firstOr(fn.Groups, "System"),
			label:  fn.Ref,
			anchor: fn.Ref,
		}
		if fn.Meta != nil {
			e.label = orStr(fn.Meta.Short, fn.Ref)
			e.desc = fn.Meta.Description
			e.weight = fn.Meta.Weight
		}

		meta := []string{"`" + fn.Ref + "`"}
		if fn.Kind == "iterator" {
			meta = append(meta, "Loop")
		}
		if _, ok := fn.Labels["corredor"]; ok {
			meta = append(meta, "Needs Corredor (`CORREDOR_ENABLED`)")
		}
		e.meta = strings.Join(meta, " · ")

		params := map[string]*types.Param{}
		for _, p := range fn.Parameters {
			params[p.ArgumentName] = p
		}
		for _, in := range segmentInputs(fn.Segments) {
			di := docInput{label: orStr(in.Label, in.Argument), required: in.Required, desc: in.Description}
			if p := params[in.Argument]; p != nil {
				di.types = typeList(p.Types)
				di.required = di.required || p.Required
				if di.desc == "" && p.Meta != nil {
					di.desc = p.Meta.Description
				}
			} else {
				di.types = typeList(nil)
			}
			di.desc = withOptions(di.desc, in)
			e.inputs = append(e.inputs, di)
		}

		for _, r := range fn.Results {
			o := docOutput{name: r.ArgumentName, types: typeList(r.Types)}
			if r.Meta != nil {
				o.desc = orStr(r.Meta.Description, r.Meta.Label)
			}
			e.outputs = append(e.outputs, o)
		}

		out = append(out, e)
	}
	return
}

// loadBranchStrings reads the two branch nodes the step picker adds on its own
// from the webapp's strings.
func loadBranchStrings(t *testing.T, path string) []docEntry {
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var doc struct {
		NodePicker struct {
			Categories struct {
				Branches string `yaml:"branches"`
			} `yaml:"categories"`
			Nodes struct {
				Branches map[string]struct {
					Label       string `yaml:"label"`
					Description string `yaml:"description"`
				} `yaml:"branches"`
			} `yaml:"nodes"`
		} `yaml:"nodePicker"`
	}
	if err = yaml.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}

	var out []docEntry
	for i, k := range []string{"exclusive", "inclusive"} {
		b, ok := doc.NodePicker.Nodes.Branches[k]
		if !ok || b.Label == "" {
			t.Fatalf("builder.yaml has no nodePicker.nodes.branches.%s", k)
		}
		out = append(out, docEntry{
			group:  doc.NodePicker.Categories.Branches,
			label:  b.Label,
			anchor: slug(b.Label),
			desc:   b.Description,
			weight: i,
			note:   "Each path out of the branch carries its own condition, an [expression](/reference/expressions/) that decides whether the path is taken.",
		})
	}
	return out
}

// withOptions names the choices a select-like input offers, after its
// placeholder when it has no description.
func withOptions(desc string, in types.SectionElementInput) string {
	if len(in.Options) == 0 {
		return desc
	}
	desc = orStr(desc, in.Placeholder)
	oo := make([]string, len(in.Options))
	for i, o := range in.Options {
		oo[i] = orStr(o.Label, o.Value)
	}
	return strings.TrimSpace(sentence(desc) + " Options: " + strings.Join(oo, ", ") + ".")
}

func segmentInputs(ss []types.ConstructSegment) (out []types.SectionElementInput) {
	for _, s := range ss {
		for _, sec := range s.Sections {
			for _, el := range sec.Elements {
				if el.Input.Argument != "" {
					out = append(out, el.Input)
				}
			}
		}
	}
	return
}

func renderTriggersPage(t *testing.T, ee []docEntry) string {
	var b strings.Builder
	b.WriteString(`---
title: TAQ triggers
description: Every trigger a TAQ can start from, with the inputs it takes and the values it provides.
outline: [2, 3]
---

<!-- This file is auto-generated from the TAQ construct catalog (server/tests/automation/docs_catalog_test.go). -->

# TAQ triggers

A trigger decides when a [TAQ](/platform/automation) runs. These are the
triggers the builder offers on every server, grouped as in the trigger picker.
**Inputs** are what you set on the trigger; **Provides** lists the values it
hands to the steps that follow. The steps are in [TAQ steps](./steps).

Triggers that come from a connection, such as its webhooks, depend on the
connections configured on your instance and are not listed here.
`)
	renderEntries(t, &b, ee, "Provides", "Value")
	return b.String()
}

func renderStepsPage(t *testing.T, ee []docEntry) string {
	var b strings.Builder
	b.WriteString(`---
title: TAQ steps
description: Every step a TAQ can run, with the inputs it takes and the values it returns.
outline: [2, 3]
---

<!-- This file is auto-generated from the TAQ construct catalog (server/tests/automation/docs_catalog_test.go). -->

# TAQ steps

Steps do the work of a [TAQ](/platform/automation); branches and loops decide
which steps run and how often. These are the steps the builder offers on every
server, grouped as in the step picker. **Inputs** are what you set on the step;
**Outputs** are the values later steps can use. What starts a TAQ is in
[TAQ triggers](./triggers).

Steps marked **Needs Corredor** are listed on every server but can only be
added where Corredor is turned on. Operations that come from a connection
depend on the connections configured on your instance and are not listed here.
`)
	renderEntries(t, &b, ee, "Outputs", "Output")
	return b.String()
}

func renderEntries(t *testing.T, b *strings.Builder, ee []docEntry, outTitle, outCol string) {
	sort.SliceStable(ee, func(i, j int) bool {
		if ee[i].group != ee[j].group {
			return ee[i].group < ee[j].group
		}
		if ee[i].weight != ee[j].weight {
			return ee[i].weight < ee[j].weight
		}
		if ee[i].label != ee[j].label {
			return ee[i].label < ee[j].label
		}
		return ee[i].anchor < ee[j].anchor
	})

	seen := map[string]bool{}
	group := ""
	for _, e := range ee {
		if e.group != group {
			group = e.group
			if seen[slug(group)] {
				t.Errorf("anchor %q is used twice", slug(group))
			}
			seen[slug(group)] = true
			fmt.Fprintf(b, "\n## %s\n", text(group))
		}
		if seen[e.anchor] {
			t.Errorf("anchor %q is used twice", e.anchor)
		}
		seen[e.anchor] = true

		fmt.Fprintf(b, "\n### %s {#%s}\n", text(e.label), e.anchor)
		if e.meta != "" {
			fmt.Fprintf(b, "\n<div class=\"env-meta\">\n\n%s\n\n</div>\n", e.meta)
		}
		if e.desc != "" {
			fmt.Fprintf(b, "\n%s\n", sentence(text(e.desc)))
		}
		if e.note != "" {
			fmt.Fprintf(b, "\n%s\n", e.note)
		}

		if len(e.inputs) > 0 {
			rows := [][]string{{"Input", "Type", "Required", "Description"}}
			for _, in := range e.inputs {
				req := ""
				if in.required {
					req = "Yes"
				}
				rows = append(rows, []string{cell(in.label), cell(in.types), req, cell(in.desc)})
			}
			writeTable(b, "Inputs", rows)
		}

		if len(e.outputs) > 0 {
			rows := [][]string{{outCol, "Type", "Description"}}
			for _, o := range e.outputs {
				rows = append(rows, []string{cell("`" + o.name + "`"), cell(o.types), cell(o.desc)})
			}
			writeTable(b, outTitle, rows)
		}
	}
}

// writeTable renders rows[0] as the header and leaves out every column past
// the second that no row fills.
func writeTable(b *strings.Builder, title string, rows [][]string) {
	for c := len(rows[0]) - 1; c >= 2; c-- {
		empty := true
		for _, r := range rows[1:] {
			empty = empty && r[c] == ""
		}
		if empty {
			for i := range rows {
				rows[i] = append(rows[i][:c], rows[i][c+1:]...)
			}
		}
	}

	fmt.Fprintf(b, "\n**%s**\n\n", title)
	for i, r := range rows {
		b.WriteString("|")
		for _, c := range r {
			if c == "" {
				b.WriteString(" |")
			} else {
				b.WriteString(" " + c + " |")
			}
		}
		b.WriteString("\n")
		if i == 0 {
			b.WriteString("|" + strings.Repeat(" --- |", len(r)) + "\n")
		}
	}
}

var (
	reSpace   = regexp.MustCompile(`\s+`)
	reNonSlug = regexp.MustCompile(`[^a-z0-9]+`)
	reTagLike = regexp.MustCompile(`<([A-Za-z/!])`)
)

// text flattens whitespace and escapes what VitePress would read as Vue or HTML
// outside code spans.
func text(s string) string {
	parts := strings.Split(reSpace.ReplaceAllString(strings.TrimSpace(s), " "), "`")
	for i := 0; i < len(parts); i += 2 {
		p := parts[i]
		p = strings.ReplaceAll(p, "{{", "&#123;&#123;")
		p = reTagLike.ReplaceAllString(p, "&lt;$1")
		parts[i] = p
	}
	return strings.Join(parts, "`")
}

func cell(s string) string {
	return strings.ReplaceAll(text(s), "|", `\|`)
}

func sentence(s string) string {
	if s == "" || strings.HasSuffix(s, ".") || strings.HasSuffix(s, "`") {
		return s
	}
	return s + "."
}

func slug(s string) string {
	return strings.Trim(reNonSlug.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func typeList(tt []string) string {
	var out []string
	for _, t := range tt {
		if t != "" {
			out = append(out, "`"+t+"`")
		}
	}
	if len(out) == 0 {
		return "Any"
	}
	return strings.Join(out, ", ")
}

func firstOr(ss []string, def string) string {
	if len(ss) > 0 && ss[0] != "" {
		return ss[0]
	}
	return def
}

func orStr(s, def string) string {
	if s != "" {
		return s
	}
	return def
}
