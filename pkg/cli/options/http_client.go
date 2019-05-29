package options

import (
	"time"
)

type (
	HttpClientOpt struct {
		ClientTSLInsecure bool          `env:"HTTP_CLIENT_TSL_INSECURE"`
		HttpClientTimeout time.Duration `env:"HTTP_CLIENT_TIMEOUT"`
	}
)

func HttpClient(pfix string) (o *HttpClientOpt) {
	o = &HttpClientOpt{
		ClientTSLInsecure: false,
		HttpClientTimeout: 30 * time.Second,
	}

	fill(o, pfix)

	return
}
