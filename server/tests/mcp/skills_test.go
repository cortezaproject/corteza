package mcp_test

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/system/agentic/skills"
	"github.com/stretchr/testify/require"
)

// The skill library is prose, and prose about a block's options is a second
// copy of a schema that lives in compose/types. It drifted: page_building
// documented a Metric tile as {"field", "reduce"} where the block reads
// {"metricField", "operation"}, a Comment as "commentField" where the block
// reads "contentField", and a Progress block's value as a number where it is an
// object. Every one of those saves without an error and renders an empty panel,
// so nothing surfaces until a person opens the page.
//
// These tests read the option shapes out of the skills and check them against
// PageBlockOptionSchemas — the same source compose_page_block_schema serves. A
// skill may say less than the schema; it may not say something different.

// blockOptionsInline matches the reference lines in page_building:
//
//	**RecordList**: `{"moduleID": "…"}`
var blockOptionsInline = regexp.MustCompile("(?m)^\\*\\*([A-Za-z]+)\\*\\*: `(\\{.*\\})`")

// blockArrayFence matches a fenced json block holding an array of page blocks,
// which is how page_record and page_record_list show a whole `blocks` argument.
var blockArrayFence = regexp.MustCompile("(?s)```json\\n(\\[.*?\\])\\n```")

// documentedBlock is one option shape a skill states, with enough context to
// name it in a failure.
type documentedBlock struct {
	skill   string
	kind    string
	where   string
	options map[string]any
}

func TestSkillsDocumentRealBlockOptions(t *testing.T) {
	found := documentedBlocks(t)
	require.NotEmpty(t, found,
		"no block option shapes were found in the skill library; the extractor and the "+
			"skills have diverged and this test is asserting nothing")

	for _, b := range found {
		schema, ok := cmpTypes.PageBlockOptionSchemas[b.kind]
		if !ok {
			t.Errorf("skill %q documents block kind %q, which no block schema defines", b.skill, b.kind)
			continue
		}

		checkOptions(t, b, b.kind, b.options, schemaMap(t, schema))
	}
}

// TestEveryBlockKindTheSkillsNameIsRegistered is the other half: a skill naming
// a kind the server does not have sends an agent to build a block that cannot
// exist. Checked by the same extraction, so it also fails when a heading is
// renamed out of the shape the extractor reads.
func TestEveryBlockKindTheSkillsNameIsRegistered(t *testing.T) {
	kinds := map[string]bool{}
	for _, b := range documentedBlocks(t) {
		kinds[b.kind] = true
	}

	for _, kind := range sortedStrings(kinds) {
		_, ok := cmpTypes.PageBlockOptionSchemas[kind]
		require.Truef(t, ok, "the skill library documents a %q block; the server has no such kind", kind)
	}
}

// documentedBlocks pulls every option shape the library states out of its
// markdown.
//
// A shape that is not valid JSON fails here rather than being skipped: the
// whole point is that these are checkable, and an example an agent cannot paste
// is not one worth documenting.
func documentedBlocks(t *testing.T) []documentedBlock {
	t.Helper()

	lib, err := skills.LoadLibrary()
	require.NoError(t, err)

	var out []documentedBlock

	for _, s := range lib.All() {
		for _, m := range blockOptionsInline.FindAllStringSubmatch(s.Body, -1) {
			kind, raw := m[1], m[2]

			var opts map[string]any
			require.NoErrorf(t, json.Unmarshal([]byte(raw), &opts),
				"skill %q documents %s options that are not valid JSON: %s", s.Name, kind, raw)

			out = append(out, documentedBlock{
				skill: s.Name, kind: kind, where: "**" + kind + "**", options: opts,
			})
		}

		for _, m := range blockArrayFence.FindAllStringSubmatch(s.Body, -1) {
			var blocks []struct {
				Kind    string         `json:"kind"`
				Options map[string]any `json:"options"`
			}
			require.NoErrorf(t, json.Unmarshal([]byte(m[1]), &blocks),
				"skill %q has a fenced json block array that is not valid JSON: %s", s.Name, m[1])

			for _, b := range blocks {
				if b.Kind == "" {
					continue
				}
				out = append(out, documentedBlock{
					skill: s.Name, kind: b.Kind, where: "the " + b.Kind + " example", options: b.Options,
				})
			}
		}
	}

	return out
}

// checkOptions walks a documented shape against the schema's shape, so a key
// documented at the wrong depth is caught as well as one that does not exist.
// The schema is authoritative about names only — values here are placeholders.
func checkOptions(t *testing.T, b documentedBlock, path string, got, want map[string]any) {
	t.Helper()

	for _, key := range sortedKeysOf(got) {
		here := path + "." + key

		w, ok := want[key]
		if !ok {
			t.Errorf("skill %q, %s: option %q is not in the %s block schema (%s). "+
				"A block saves with it and renders empty, so this reads as working. "+
				"Run compose_page_block_schema for the real names.",
				b.skill, b.where, here, b.kind, nearest(key, want))
			continue
		}

		switch v := got[key].(type) {
		case map[string]any:
			if nested, ok := w.(map[string]any); ok {
				checkOptions(t, b, here, v, nested)
			}
		case []any:
			nested, ok := elementSchema(w)
			if !ok {
				continue
			}
			for _, item := range v {
				if m, ok := item.(map[string]any); ok {
					checkOptions(t, b, here+"[]", m, nested)
				}
			}
		}
	}
}

// elementSchema is the object shape an array option holds, taken from the
// skeleton's first element. An empty array in the skeleton says nothing about
// its items, so it stops the walk rather than failing it.
func elementSchema(want any) (map[string]any, bool) {
	arr, ok := want.([]any)
	if !ok || len(arr) == 0 {
		return nil, false
	}
	m, ok := arr[0].(map[string]any)
	return m, ok
}

// nearest suggests the schema key a wrong one was probably meant to be, so the
// failure carries the fix. Substring either way is enough: the mistakes are
// "field" for "metricField" and "commentField" for "contentField".
func nearest(key string, want map[string]any) string {
	var hits []string
	lower := strings.ToLower(key)
	for _, k := range sortedKeysOf(want) {
		lk := strings.ToLower(k)
		if strings.Contains(lk, lower) || strings.Contains(lower, lk) {
			hits = append(hits, k)
		}
	}
	if len(hits) == 0 {
		return "the schema has none of that name"
	}
	return "did you mean " + strings.Join(hits, " or ") + "?"
}

func schemaMap(t *testing.T, schema any) map[string]any {
	t.Helper()

	raw, err := json.Marshal(schema)
	require.NoError(t, err)

	var out map[string]any
	require.NoError(t, json.Unmarshal(raw, &out))
	return out
}

func sortedKeysOf(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sortedStrings(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
