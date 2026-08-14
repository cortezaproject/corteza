package agentic

import (
	"testing"
)

// A record value is {name, value, place}, and a multi-value field is several of
// them sharing a name. These cases are the shapes an agent actually sends.
func TestParseValues(t *testing.T) {
	type kv struct {
		name  string
		value string
		place uint
	}

	tests := []struct {
		name    string
		input   any
		want    []kv
		wantErr bool
	}{
		{
			name:  "JSON string",
			input: `{"participant":"Pat","time":"2026-04-02T17:00:00"}`,
			want: []kv{
				{"participant", "Pat", 0},
				{"time", "2026-04-02T17:00:00", 0},
			},
		},
		{
			name:  "map object",
			input: map[string]any{"participant": "Pat"},
			want:  []kv{{"participant", "Pat", 0}},
		},
		{
			name:  "number and boolean keep their scalar rendering",
			input: map[string]any{"count": 5, "done": true},
			want:  []kv{{"count", "5", 0}, {"done", "true", 0}},
		},

		// The whole point. Before this, an array went through fmt.Sprintf("%v")
		// and a String field stored the literal `[red blue]`.
		{
			name:  "array becomes one value per element, placed in order",
			input: map[string]any{"tags": []any{"red", "blue", "green"}},
			want: []kv{
				{"tags", "red", 0},
				{"tags", "blue", 1},
				{"tags", "green", 2},
			},
		},
		{
			name:  "single-element array is still a placed value",
			input: `{"tags":["red"]}`,
			want:  []kv{{"tags", "red", 0}},
		},
		{
			name:  "empty array contributes nothing",
			input: map[string]any{"tags": []any{}},
			want:  nil,
		},

		// A Geometry value is held as JSON, so the object an agent naturally
		// sends has to be marshalled — %v rendered it map[coordinates:[46 14]].
		{
			name:  "object is stored as its JSON",
			input: `{"geo":{"coordinates":[46.05,14.5]}}`,
			want:  []kv{{"geo", `{"coordinates":[46.05,14.5]}`, 0}},
		},
		{
			name:  "array of objects places each one",
			input: `{"geo":[{"coordinates":[1,2]},{"coordinates":[3,4]}]}`,
			want: []kv{
				{"geo", `{"coordinates":[1,2]}`, 0},
				{"geo", `{"coordinates":[3,4]}`, 1},
			},
		},

		{
			name:  "names are ordered so the same payload gives the same record",
			input: map[string]any{"zebra": "z", "alpha": "a", "middle": "m"},
			want:  []kv{{"alpha", "a", 0}, {"middle", "m", 0}, {"zebra", "z", 0}},
		},
		{
			name:  "null clears the value",
			input: `{"note":null}`,
			want:  []kv{{"note", "", 0}},
		},

		// Refused rather than stringified: a value shape nobody anticipated
		// used to become Go's debug output and, on a permissive field kind, was
		// stored exactly as printed.
		{
			name:    "nested array is refused",
			input:   `{"tags":[["a","b"]]}`,
			wantErr: true,
		},
		{name: "invalid JSON string", input: `not json`, wantErr: true},
		{name: "unsupported type", input: 42, wantErr: true},
		{name: "nil", input: nil, wantErr: true},
		{name: "empty string", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseValues(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseValues() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if len(got) != len(tt.want) {
				t.Fatalf("parseValues() returned %d values, want %d: %+v", len(got), len(tt.want), got)
			}

			for i, w := range tt.want {
				if got[i].Name != w.name || got[i].Value != w.value || got[i].Place != w.place {
					t.Errorf("value %d = {%q, %q, %d}, want {%q, %q, %d}",
						i, got[i].Name, got[i].Value, got[i].Place, w.name, w.value, w.place)
				}
			}
		})
	}
}
