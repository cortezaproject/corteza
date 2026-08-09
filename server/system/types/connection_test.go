package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectionTemplate_UnmarshalNestedAuth(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		wantValue  string
		wantHeader string
		wantPHName string
	}{
		{
			name:      "bare string",
			json:      `"https://api.example.com"`,
			wantValue: "https://api.example.com",
		},
		{
			name:      "flat template",
			json:      `{"value":"{{apiKey}}","placeholders":[{"name":"apiKey","type":"String","required":true}]}`,
			wantValue: "{{apiKey}}", wantPHName: "apiKey",
		},
		{
			// Notion/Stripe/Mailchimp catalog shape
			name: "nested token with Authorization header",
			json: `{
				"token": {"value":"{{apiToken}}","placeholders":[{"name":"apiToken","type":"String","required":true}]},
				"headerName": {"value":"Authorization"}
			}`,
			wantValue: "{{apiToken}}", wantHeader: "Authorization", wantPHName: "apiToken",
		},
		{
			// Datadog catalog shape: custom, non-Authorization header
			name: "nested token with custom header",
			json: `{
				"token": {"value":"{{apiKey}}","placeholders":[{"name":"apiKey","type":"String","required":true}]},
				"headerName": {"value":"DD-API-KEY"}
			}`,
			wantValue: "{{apiKey}}", wantHeader: "DD-API-KEY", wantPHName: "apiKey",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ct ConnectionTemplate
			require.NoError(t, json.Unmarshal([]byte(tc.json), &ct))
			require.Equal(t, tc.wantValue, ct.Value)
			require.Equal(t, tc.wantHeader, ct.HeaderName)
			if tc.wantPHName != "" {
				require.Len(t, ct.Placeholders, 1)
				require.Equal(t, tc.wantPHName, ct.Placeholders[0].Name)
			}
		})
	}
}

// The header name must survive a marshal/unmarshal round-trip (persistence),
// without re-emitting the nested catalog shape.
func TestConnectionTemplate_HeaderNameRoundTrip(t *testing.T) {
	src := ConnectionTemplate{
		Value:      "{{apiKey}}",
		HeaderName: "DD-API-KEY",
	}

	b, err := json.Marshal(src)
	require.NoError(t, err)

	var got ConnectionTemplate
	require.NoError(t, json.Unmarshal(b, &got))
	require.Equal(t, "{{apiKey}}", got.Value)
	require.Equal(t, "DD-API-KEY", got.HeaderName)
}
