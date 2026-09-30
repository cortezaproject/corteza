package automation

import (
	"context"
	"testing"

	pkgHttp "github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/stretchr/testify/require"
)

// Workflows must not be able to reach cloud metadata endpoints
func TestHttpRequestRestrictedNetworks(t *testing.T) {
	for _, url := range []string{
		"http://169.254.169.254/latest/meta-data/iam/security-credentials/",
		"http://[fd00:ec2::254]/latest/meta-data/",
		"http://100.100.100.200/latest/meta-data/",
		"http://[::ffff:169.254.169.254]/",
	} {
		t.Run(url, func(t *testing.T) {
			_, err := httpRequestHandler{}.send(context.Background(), &httpRequestSendArgs{Url: url, Method: "GET"})
			require.ErrorIs(t, err, pkgHttp.ErrRestrictedNetwork)
		})
	}
}
