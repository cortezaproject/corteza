package handle

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		name   string
		handle string
		want   bool
	}{
		{"empty", "", true},
		{"alphanum", "a1", true},
		{"alpha", "abc", true},
		{"num", "123", false},
		{"valid1", "Valid1Handle", true},
		{"valid2", "Valid-Handle", true},
		{"valid3", "Valid_Handle", true},
		{"weirdo", "a$&_.!", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.handle); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
