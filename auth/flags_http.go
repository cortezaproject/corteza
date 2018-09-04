package auth

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	httpFlags struct {
		addr    string
		logging bool
		pretty  bool
		tracing bool
		metrics bool
	}
)

func (c *httpFlags) validate() error {
	if c == nil {
		return nil
	}
	if c.addr == "" {
		return errors.New("No HTTP Addr is set, can't listen for HTTP")
	}
	return nil
}

func (c *httpFlags) flags(prefix ...string) *httpFlags {
	p := func(s string) string {
		return prefix[0] + "-" + s
	}

	flag.StringVar(&c.addr, p("http-addr"), ":3000", "Listen address for HTTP server")
	flag.BoolVar(&c.logging, p("http-log"), true, "Enable/disable HTTP request log")
	flag.BoolVar(&c.pretty, p("http-pretty-json"), false, "Prettify returned JSON output")
	flag.BoolVar(&c.tracing, p("http-error-tracing"), false, "Return error stack frame")
	flag.BoolVar(&c.metrics, p("http-metrics"), false, "Provide metrics export for prometheus")

	return c
}
