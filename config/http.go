package config

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	HTTP struct {
		Addr    string
		Logging bool
		Pretty  bool
		Tracing bool
		Metrics bool
	}
)

func (c *HTTP) Validate() error {
	if c == nil {
		return nil
	}
	if c.Addr == "" {
		return errors.New("No HTTP Addr is set, can't listen for HTTP")
	}
	return nil
}

func (c *HTTP) Init(prefix ...string) *HTTP {
	p := func(s string) string {
		return prefix[0] + "-" + s
	}

	flag.StringVar(&c.Addr, p("http-addr"), ":3000", "Listen address for HTTP server")
	flag.BoolVar(&c.Logging, p("http-log"), true, "Enable/disable HTTP request log")
	flag.BoolVar(&c.Pretty, p("http-pretty-json"), false, "Prettify returned JSON output")
	flag.BoolVar(&c.Tracing, p("http-error-tracing"), false, "Return error stack frame")
	flag.BoolVar(&c.Metrics, p("http-metrics"), false, "Provide metrics export for prometheus")

	return c
}
