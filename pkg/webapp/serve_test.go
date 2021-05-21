package webapp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_replaceBaseHrefPlaceholder(t *testing.T) {
	var (
		tcc = []struct {
			name string
			app  string
			base string
			body string
		}{
			{"unquoted condensed", "test-app", "", `... <base href=/> ... `},
			{"unquoted spacy", "test-app", "", `... <base href=/ > ... `},
			{"quoted spacy", "test-app", "", `... <base href="/" ... > `},
			{"preset", "test-app", "", `... <base href="/foo" > ... `},
			{"unquoted condensed", "test-app", "base", `... <base href=/> ... `},
			{"unquoted spacy", "test-app", "base", `... <base href=/ > ... `},
			{"quoted spacy", "test-app", "base", `... <base href="/" ... > `},
			{"preset", "test-app", "base", `... <base href="/foo" > ... `},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req   = require.New(t)
				fixed = string(replaceBaseHrefPlaceholder([]byte(tc.body), tc.app, tc.base))
			)

			req.NotContains(fixed, "Error!")
			req.Contains(fixed, "<base href=")
			req.Contains(fixed, tc.app)
		})

	}
}
