package saml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TemplateProvider(t *testing.T) {
	var (
		tcc = []struct {
			name   string
			url    string
			expect templateProvider
		}{
			{
				name: "Default SAML provider",
				url:  "http://example.tld",
				expect: templateProvider{
					Label:  "Default SAML provider",
					Handle: "saml/init",
					Icon:   "key",
				},
			},
			{
				name: "",
				url:  "http://example.tld",
				expect: templateProvider{
					Label:  "http://example.tld",
					Handle: "saml/init",
					Icon:   "key",
				},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			req.Equal(tc.expect, TemplateProvider(tc.url, tc.name))
		})
	}
}
