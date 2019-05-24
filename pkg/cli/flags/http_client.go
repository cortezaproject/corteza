package flags

import (
	"time"

	"github.com/spf13/cobra"
)

type (
	HttpClientOpt struct {
		ClientTSLInsecure bool
		HttpClientTimeout time.Duration
	}
)

func HttpClient(cmd *cobra.Command) (o *HttpClientOpt) {
	o = &HttpClientOpt{}

	BindBool(cmd, &o.ClientTSLInsecure,
		"http-client-tsl-insecure", false,
		"Skip insecure TSL verification on outbound HTTP requests (allow invalid/self-signed certificates")

	BindDuration(cmd, &o.HttpClientTimeout,
		"http-client-timeout", 30*time.Second,
		"Default HTTP client timeout")

	return
}
