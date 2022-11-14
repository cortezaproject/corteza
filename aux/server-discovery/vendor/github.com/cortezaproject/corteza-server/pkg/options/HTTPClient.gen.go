package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/HTTPClient.yaml

import (
	"time"
)

type (
	HTTPClientOpt struct {
		ClientTSLInsecure bool          `env:"HTTP_CLIENT_TLS_INSECURE"`
		HttpClientTimeout time.Duration `env:"HTTP_CLIENT_TIMEOUT"`
	}
)

// HTTPClient initializes and returns a HTTPClientOpt with default values
func HTTPClient() (o *HTTPClientOpt) {
	o = &HTTPClientOpt{
		ClientTSLInsecure: false,
		HttpClientTimeout: 30 * time.Second,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *HTTPClient) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
