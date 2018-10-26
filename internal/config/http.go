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

var http *HTTP

func (c *HTTP) Validate() error {
	if c == nil {
		return nil
	}
	if c.Addr == "" {
		return errors.New("No HTTP Addr is set, can't listen for HTTP")
	}
	return nil
}

func (*HTTP) Init(prefix ...string) *HTTP {
	if http != nil {
		return http
	}

	p := func(s string) string {
		return prefix[0] + "-" + s
	}

	http = new(HTTP)
	flag.StringVar(&http.Addr, p("http-addr"), ":3000", "Listen address for HTTP server")
	flag.BoolVar(&http.Logging, p("http-log"), true, "Enable/disable HTTP request log")
	flag.BoolVar(&http.Pretty, p("http-pretty-json"), false, "Prettify returned JSON output")
	flag.BoolVar(&http.Tracing, p("http-error-tracing"), false, "Return error stack frame")
	flag.BoolVar(&http.Metrics, p("http-metrics"), false, "Provide metrics export for prometheus")
	return http
}
