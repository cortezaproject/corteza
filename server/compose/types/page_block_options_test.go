package types

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var updateBlockSchemas = flag.Bool("update-block-schemas", false, "regenerate testdata/page_block_option_schemas.json")

// TestPageBlockOptionSchemasSnapshot pins the serialized form of
// PageBlockOptionSchemas (what the compose_page_block_schema MCP tool serves)
// to a committed JSON snapshot. The snapshot is the cross-language contract:
// lib/js/src/compose/types/page-block/schema-contract.test.ts verifies every
// key in it against the actual webapp page-block classes, so a schema key the
// webapp does not read (the module-vs-moduleID class of bug) fails CI on the
// FE side, and any Go-side change fails here until the snapshot is
// regenerated with:
//
//	go test ./compose/types/ -run PageBlockOptionSchemasSnapshot -update-block-schemas
func TestPageBlockOptionSchemasSnapshot(t *testing.T) {
	got, err := json.MarshalIndent(PageBlockOptionSchemas, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')

	path := filepath.Join("testdata", "page_block_option_schemas.json")

	if *updateBlockSchemas {
		if err = os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("snapshot updated: %s — re-run the lib/js schema-contract test", path)
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read snapshot (regenerate with -update-block-schemas): %v", err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("PageBlockOptionSchemas changed but the snapshot was not regenerated;\n"+
			"run: go test ./compose/types/ -run PageBlockOptionSchemasSnapshot -update-block-schemas\n"+
			"then run the lib/js contract test: cd lib/js && npx vitest run src/compose/types/page-block/schema-contract.test.ts\ngot:\n%s", got)
	}
}
