package options

import (
	"time"
)

type (
	HttpClientOpt struct {
		ClientTSLInsecure bool
		HttpClientTimeout time.Duration
	}
)

func HttpClient(pfix string) (o *HttpClientOpt) {
	o = &HttpClientOpt{
		ClientTSLInsecure: EnvBool(pfix, "HTTP_CLIENT_TSL_INSECURE", false),
		HttpClientTimeout: EnvDuration(pfix, "HTTP_CLIENT_TIMEOUT", 30*time.Second),
	}

	return
}
