// Package schemas_test writes the JSON Schema reference a configurator's
// assistant reads: the shape of a module, a page, a page layout, each block
// kind's options, a chart and a TAQ, reflected from the Go types the server
// stores, with the closed value sets (field kinds, block kinds, step refs,
// trigger events) filled from the catalogues. HUMAN_UPDATE_DOCS=1 rewrites
// the files; otherwise the test holds them to the code.
package schemas_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	autoTypes "github.com/crusttech/human/server/automation/types"
	cmpAgentic "github.com/crusttech/human/server/compose/agentic"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
)

const schemasDir = "docs/public/schemas"

func repoRoot(t *testing.T) string {
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}

// reflector maps the Go types that do not describe themselves: IDs are
// strings on the wire, variables and options are free objects, labels are
// strings.
func reflector() *jsonschema.Reflector {
	idSchema := func(desc string) *jsonschema.Schema {
		return &jsonschema.Schema{Type: "string", Pattern: "^[0-9]+$", Description: desc}
	}
	return &jsonschema.Reflector{
		// Pages hold child pages and TAQ steps hold expressions, so nested
		// types stay as $defs references: inlining them never ends.
		ExpandedStruct:            true,
		AllowAdditionalProperties: false,
		Mapper: func(t reflect.Type) *jsonschema.Schema {
			switch t {
			case reflect.TypeOf(uint64(0)):
				return idSchema("64-bit integer as a string")
			case reflect.TypeOf(time.Time{}):
				return &jsonschema.Schema{Type: "string", Format: "date-time"}
			case reflect.TypeOf(cmpTypes.ModuleFieldOptions{}):
				return &jsonschema.Schema{Type: "object", Description: "Per-kind options. Select: options [{value, text}], selectType. Record: moduleID, labelField, queryFields, selectType. Number: precision, format. DateTime: onlyDate, onlyTime. File: mode, allowImages, allowDocuments. User: presetWithAuthenticated, selectType. Every kind: multiDelimiter, when isMulti."}
			case reflect.TypeOf(expr.Vars{}):
				return &jsonschema.Schema{Type: "object", Description: "Typed variables: {\"name\": {\"@type\": \"String\", \"@value\": \"…\"}}"}
			}
			return nil
		},
	}
}

type doc struct {
	file   string
	schema *jsonschema.Schema
}

func describe(s *jsonschema.Schema, title, desc string) *jsonschema.Schema {
	s.Title, s.Description = title, desc
	return s
}

// at walks properties and array items by name, following $ref into the
// root's $defs; "[]" steps into items.
func at(t *testing.T, root *jsonschema.Schema, path ...string) *jsonschema.Schema {
	resolve := func(s *jsonschema.Schema) *jsonschema.Schema {
		if s.Ref == "" {
			return s
		}
		name := strings.TrimPrefix(s.Ref, "#/$defs/")
		def, ok := root.Definitions[name]
		require.Truef(t, ok, "unresolved %s at %v", s.Ref, path)
		return def
	}
	cur := resolve(root)
	for _, p := range path {
		if p == "[]" {
			require.NotNilf(t, cur.Items, "no items at %v", path)
			cur = resolve(cur.Items)
			continue
		}
		require.NotNilf(t, cur.Properties, "no properties at %v", path)
		next, ok := cur.Properties.Get(p)
		require.Truef(t, ok, "no property %q at %v", p, path)
		cur = resolve(next)
	}
	return cur
}

func enum(values []string) []any {
	out := make([]any, 0, len(values))
	for _, v := range values {
		out = append(out, v)
	}
	return out
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func buildDocs(t *testing.T) []doc {
	r := reflector()
	var docs []doc

	// Block options: one schema per kind, required from the block's own list.
	blockKinds := sortedKeys(cmpTypes.PageBlockOptionSchemas)
	required := cmpAgentic.RequiredBlockOptions()
	blocks := &jsonschema.Schema{Type: "object", Properties: jsonschema.NewProperties()}
	for _, kind := range blockKinds {
		s := r.ReflectFromType(reflect.TypeOf(cmpTypes.PageBlockOptionSchemas[kind]))
		s.Required = required[kind]
		s.Title = kind
		blocks.Properties.Set(kind, s)
	}
	docs = append(docs, doc{"page-block-options.schema.json", describe(blocks, "Page block options",
		"The options object of each page block kind, keyed by kind. required lists the options a block cannot render without; compose_page_block_schema returns the same per kind at run time.")})

	module := r.Reflect(&cmpTypes.Module{})
	at(t, module, "fields", "[]", "kind").Enum = enum(cmpTypes.ModuleFieldKinds)
	docs = append(docs, doc{"module.schema.json", describe(module, "Module",
		"A compose module: the data structure records of one kind follow. fields[].kind is the closed list of field kinds; fields[].options depends on the kind.")})

	page := r.Reflect(&cmpTypes.Page{})
	at(t, page, "blocks", "[]", "kind").Enum = enum(blockKinds)
	at(t, page, "blocks", "[]", "options").Description = "See page-block-options.schema.json under the block's kind."
	docs = append(docs, doc{"page.schema.json", describe(page, "Page",
		"A compose page: what its blocks are. Where they sit on the grid is the page layout's (page-layout.schema.json). A record page sets moduleID.")})

	layout := r.Reflect(&cmpTypes.PageLayout{})
	docs = append(docs, doc{"page-layout.schema.json", describe(layout, "Page layout",
		"One arrangement of a page's blocks on the 48-column grid, with the visibility condition that selects it.")})

	chart := r.Reflect(&cmpTypes.Chart{})
	docs = append(docs, doc{"chart.schema.json", describe(chart, "Chart",
		"A compose chart. config.reports[] names the module, filter, dimensions and metrics; dimension and metric entries are free objects whose keys the chart_building skill lists.")})

	// Step refs and trigger events come from handlers that register at boot,
	// which this unit test does not run; the generated reference pages under
	// docs/reference/taq hold them and are tested where the app is booted.
	taq := r.Reflect(&autoTypes.NgAutomation{})
	at(t, taq, "steps", "[]", "ref").Description = "A function or construct ref from the TAQ steps reference (docs/reference/taq/steps)."
	at(t, taq, "triggers", "[]", "resourceType").Description = "A resource type from the TAQ triggers reference (docs/reference/taq/triggers)."
	at(t, taq, "triggers", "[]", "eventType").Description = "An event type the resource type offers, from the TAQ triggers reference."
	docs = append(docs, doc{"taq.schema.json", describe(taq, "TAQ (Trigger Action Query)",
		"A TAQ definition: triggers that start it, steps (ref names the function or construct), and paths that connect them. Step arguments and results are {scope, expr} references; the taq_authoring skill and the TAQ reference pages carry the rules.")})

	return docs
}

// TestSchemasAreCurrent holds docs/public/schemas/*.schema.json to the types.
func TestSchemasAreCurrent(t *testing.T) {
	dir := filepath.Join(repoRoot(t), filepath.FromSlash(schemasDir))
	for _, d := range buildDocs(t) {
		got, err := json.MarshalIndent(d.schema, "", "  ")
		require.NoError(t, err)
		got = append(got, '\n')
		p := filepath.Join(dir, d.file)

		if os.Getenv("HUMAN_UPDATE_DOCS") == "1" {
			require.NoError(t, os.MkdirAll(dir, 0o755))
			require.NoError(t, os.WriteFile(p, got, 0o644))
			continue
		}
		want, err := os.ReadFile(p)
		require.NoErrorf(t, err, "%s is missing; run `make docs-reference` in server/", d.file)
		require.Equalf(t, string(want), string(got), "%s is out of date with the types; run `make docs-reference` in server/ and commit the result", d.file)
	}
}
