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
		{"safe html", "<b>čšž</b>", "<b>čšž</b>"},
		{"unsafe html", `<a href="javascript:document.location='https://planetcrust.com/'">XSS</A>`, "XSS"},

		{"regular link", `<a href="https://planetcrust.com/">home</a>`, `<a href="https://planetcrust.com/" rel="nofollow">home</a>`},
		{"link with target", `<a href="https://planetcrust.com/" target="_blank">home</a>`, `<a href="https://planetcrust.com/" target="_blank" rel="nofollow noopener">home</a>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.out, SanitizeMessage(tt.in))
		})
	}
}
