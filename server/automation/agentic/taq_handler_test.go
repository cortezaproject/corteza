package agentic

import (
	"testing"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "JSON string",
			input: `{"key":"value","count":3}`,
			want:  map[string]interface{}{"key": "value", "count": float64(3)},
		},
		{
			name:  "map object",
			input: map[string]interface{}{"key": "value"},
			want:  map[string]interface{}{"key": "value"},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("parseInput() len = %d, want %d", len(got), len(tt.want))
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("parseInput()[%q] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}
