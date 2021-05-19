package saml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	defaultSAMLEmailPayload     = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
	defaultSamlTestEmailPayload = "urn:oasis:names:tc:SAML:attribute:subject-id"
)

func Test_guessIdentifier(t *testing.T) {
	var (
		sp = SamlSPService{
			IDPUserMeta: &IdpIdentityPayload{
				Name:       "name_field",
				Handle:     "handle_field",
				Identifier: "identifier_field",
			},
		}

		tcc = []struct {
			name    string
			payload map[string][]string
			expect  string
		}{
			{
				name:    "default email",
				payload: map[string][]string{"email": {"test@example.com"}},
				expect:  "test@example.com",
			},
			{
				name:    "samltest.id email payload",
				payload: map[string][]string{defaultSamlTestEmailPayload: {"test@example.com"}},
				expect:  "test@example.com",
			},
			{
				name:    "default SAML emailAddress payload",
				payload: map[string][]string{defaultSAMLEmailPayload: {"test@example.com"}},
				expect:  "test@example.com",
			},
			{
				name:    "emailAddress payload",
				payload: map[string][]string{"emailAddress": {"test@example.com"}},
				expect:  "test@example.com",
			},
			{
				name:    "missing email payload",
				payload: map[string][]string{"non-existing-identifier": {"test@example.com"}},
				expect:  "",
			},
			{
				name:    "default payload as set in settings",
				payload: map[string][]string{"identifier_field": {"test@example.com"}},
				expect:  "test@example.com",
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			req.Equal(tc.expect, sp.GuessIdentifier(tc.payload))
		})
	}
}
