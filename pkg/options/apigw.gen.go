package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/apigw.yaml

import (
	"time"
)

type (
	ApigwOpt struct {
		Enabled              bool          `env:"APIGW_ENABLED"`
		Debug                bool          `env:"APIGW_DEBUG"`
		LogEnabled           bool          `env:"APIGW_LOG_ENABLED"`
		LogRequestBody       bool          `env:"APIGW_LOG_REQUEST_BODY"`
		ProxyEnableDebugLog  bool          `env:"APIGW_PROXY_ENABLE_DEBUG_LOG"`
		ProxyFollowRedirects bool          `env:"APIGW_PROXY_FOLLOW_REDIRECTS"`
		ProxyOutboundTimeout time.Duration `env:"APIGW_PROXY_OUTBOUND_TIMEOUT"`
	}
)

// Apigw initializes and returns a ApigwOpt with default values
func Apigw() (o *ApigwOpt) {
	o = &ApigwOpt{
		Enabled:              true,
		Debug:                false,
		LogEnabled:           false,
		LogRequestBody:       false,
		ProxyEnableDebugLog:  false,
		ProxyFollowRedirects: true,
		ProxyOutboundTimeout: time.Second * 30,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Apigw) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
