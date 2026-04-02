package agentic

import (
	"testing"
)

func TestParseValues(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    map[string]string
		wantErr bool
	}{
		{
			name:  "JSON string",
			input: `{"participant":"Pat","time":"2026-04-02T17:00:00"}`,
			want:  map[string]string{"participant": "Pat", "time": "2026-04-02T17:00:00"},
		},
		{
			name:  "map object",
			input: map[string]interface{}{"participant": "Pat", "time": "2026-04-02T17:00:00"},
			want:  map[string]string{"participant": "Pat", "time": "2026-04-02T17:00:00"},
		},
		{
			name:  "map with non-string value",
			input: map[string]interface{}{"count": 5},
			want:  map[string]string{"count": "5"},
		},
		{
			name:    "invalid JSON string",
			input:   `not json`,
			wantErr: true,
		},
		{
			name:    "unsupported type",
			input:   42,
			wantErr: true,
		},
		{
			name:    "nil",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseValues(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseValues() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("parseValues() len = %d, want %d", len(got), len(tt.want))
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("parseValues()[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}
