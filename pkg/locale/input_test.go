package locale

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SanitizeMessage(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{"simple", "abc", "abc"},
		{"accents", "čšž", "čšž"},
		{"html", "<b>čšž</b>", "čšž"},
		{"broken html 1", "<b>čšž</b", "čšž"},
		{"broken html 2", "b>čšž</b", "b>čšž"},
		{"broken html 3", "<b fff=\"čšž</b", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.out, SanitizeMessage(tt.in))
		})
	}
}
