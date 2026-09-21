package service

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	httpClient "github.com/crusttech/human/server/pkg/http"
	"github.com/crusttech/human/server/pkg/mail"
	"github.com/stretchr/testify/require"
)

// Attachments are downloaded from URLs sent by the users,
// download must follow the network restrictions
func TestNotification_attachmentNetworkRestrictions(t *testing.T) {
	var (
		svc = notification{}
		srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("internal-secret"))
		}))

		render = func() string {
			msg := mail.New()
			msg.SetHeader("From", "from@example.tld")
			msg.SetHeader("To", "to@example.tld")
			msg.SetBody("text/plain", "body")

			require.NoError(t, svc.procEmailAttachments(context.Background(), msg, srv.URL+"/secret.txt"))

			buf := &bytes.Buffer{}
			_, err := msg.WriteTo(buf)
			require.NoError(t, err)
			return buf.String()
		}
	)

	defer srv.Close()

	t.Run("allowed address", func(t *testing.T) {
		httpClient.SetupRestricted(0, false, httpClient.RestrictLinkLocal)
		require.Contains(t, render(), "secret.txt")
	})

	t.Run("restricted address", func(t *testing.T) {
		httpClient.SetupRestricted(0, false, httpClient.RestrictPrivate)
		defer httpClient.SetupRestricted(0, false, httpClient.RestrictLinkLocal)

		require.NotContains(t, render(), "secret.txt")
	})
}
