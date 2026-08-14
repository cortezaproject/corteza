package types

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var updateFieldKinds = flag.Bool("update-field-kinds", false, "regenerate testdata/module_field_kinds.json")

// TestModuleFieldKindsSnapshot pins the canonical kind list to a committed
// snapshot, which is the cross-language contract:
// lib/js/src/compose/types/module-field/kind-contract.test.ts checks it against
// the classes the webapp actually registers, so a kind named here that no
// editor can render fails CI there, and a kind the webapp gains without being
// named here fails there too.
//
// It exists because nothing else could catch the drift. An unknown kind is not
// rejected anywhere — the DAL falls through to TypeText and the webapp falls
// back to the String editor — so a field of a kind that does not exist looks
// like a working text field, which is how the agentic tool documentation came
// to advertise Currency and Duration for months.
//
// Regenerate with:
//
//	go test ./compose/types/ -run ModuleFieldKindsSnapshot -update-field-kinds
func TestModuleFieldKindsSnapshot(t *testing.T) {
	got, err := json.MarshalIndent(ModuleFieldKinds, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')

	path := filepath.Join("testdata", "module_field_kinds.json")

	if *updateFieldKinds {
		if err = os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("snapshot updated: %s — re-run the lib/js kind-contract test", path)
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read snapshot (regenerate with -update-field-kinds): %v", err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("ModuleFieldKinds changed but the snapshot was not regenerated;\n"+
			"run: go test ./compose/types/ -run ModuleFieldKindsSnapshot -update-field-kinds\n"+
			"then run the lib/js contract test\ngot:\n%s", got)
	}
}

func TestIsValidModuleFieldKind(t *testing.T) {
	for _, kind := range ModuleFieldKinds {
		if !IsValidModuleFieldKind(kind) {
			t.Errorf("%q is in the canonical list but not reported valid", kind)
		}
	}

	// The two the agentic tool used to advertise, and a kind nobody has ever
	// claimed — all three were accepted and stored as written.
	for _, kind := range []string{"Currency", "Duration", "NotAKindAtAll", ""} {
		if IsValidModuleFieldKind(kind) {
			t.Errorf("%q is not a field kind but was reported valid", kind)
		}
	}
}
