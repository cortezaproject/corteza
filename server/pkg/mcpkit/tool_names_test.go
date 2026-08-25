package mcpkit

import (
	"strings"
	"testing"
)

func TestParseToolNames(t *testing.T) {
	want := []string{"compose_module_create", "compose_record_create"}

	accepted := map[string][]string{
		"json array":             {`["compose_module_create","compose_record_create"]`},
		"json array, spaced":     {`[ "compose_module_create", "compose_record_create" ]`},
		"space separated":        {"compose_module_create compose_record_create"},
		"comma separated":        {"compose_module_create,compose_record_create"},
		"comma and space":        {"compose_module_create, compose_record_create"},
		"newline separated":      {"compose_module_create\ncompose_record_create"},
		"quoted, not an array":   {`"compose_module_create" "compose_record_create"`},
		"surrounding whitespace": {"  compose_module_create compose_record_create  "},
	}

	for name, inputs := range accepted {
		t.Run(name, func(t *testing.T) {
			got, err := parseToolNames(inputs[0])
			if err != nil {
				t.Fatalf("parseToolNames(%q) errored: %v", inputs[0], err)
			}
			if len(got) != len(want) {
				t.Fatalf("got %v, want %v", got, want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("got %v, want %v", got, want)
				}
			}
		})
	}

	t.Run("single name", func(t *testing.T) {
		got, err := parseToolNames("compose_module_create")
		if err != nil || len(got) != 1 || got[0] != "compose_module_create" {
			t.Fatalf("got %v, err %v", got, err)
		}
	})

	t.Run("empty is refused with the shapes it takes", func(t *testing.T) {
		_, err := parseToolNames("   ")
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "system_role_create") {
			t.Errorf("error should show an example, got: %v", err)
		}
	})

	// A malformed array must not be silently split on its punctuation — that
	// would load tools named "[" and produce a confusing unknown-tool list.
	t.Run("broken json array says so", func(t *testing.T) {
		_, err := parseToolNames(`["compose_module_create",`)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "JSON array") {
			t.Errorf("error should name the shape it tried, got: %v", err)
		}
	})
}
