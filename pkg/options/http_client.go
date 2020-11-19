package options

import (
	"time"
)

type (
	HTTPClientOpt struct {
		ClientTSLInsecure bool          `env:"HTTP_CLIENT_TSL_INSECURE"`
		HttpClientTimeout time.Duration `env:"HTTP_CLIENT_TIMEOUT"`
	}
)

func HttpClient(pfix string) (o *HTTPClientOpt) {
	o = &HTTPClientOpt{
		ClientTSLInsecure: false,
		HttpClientTimeout: 30 * time.Second,
	}

	fill(o)

	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
