package types

import (
	"fmt"
	"net/http"

	h "github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/cortezaproject/corteza/server/pkg/options"
)

type (
	Scp map[string]interface{}
)

func (s Scp) Keys() (kk []string) {
	for i := range s {
		kk = append(kk, i)
	}

	return
}

func (s Scp) Request() *h.Request {
	if _, ok := s["request"]; ok {
		return s["request"].(*h.Request)
	}

	return nil
}

func (s Scp) Writer() http.ResponseWriter {
	if _, ok := s["writer"]; ok {
		return s["writer"].(http.ResponseWriter)
	}

	return nil
}

// Opts returns gateway options from the scope
//
// Routes store a Config (built from settings), older callers may store the
// options struct itself. Never returns nil; callers dereference the result
// directly.
func (s Scp) Opts() *options.ApigwOpt {
	if ss, ok := s["opts"]; ok {
		switch sss := ss.(type) {
		case *options.ApigwOpt:
			return withOptDefaults(sss)
		case options.ApigwOpt:
			return withOptDefaults(&sss)
		case *Config:
			return withOptDefaults(sss.toOpt())
		case Config:
			return withOptDefaults(sss.toOpt())
		}
	}

	return withOptDefaults(&options.ApigwOpt{})
}

// toOpt maps the gateway config onto the options struct
//
// ProxyEnableDebugLog has no counterpart in settings and stays off.
func (c Config) toOpt() *options.ApigwOpt {
	return &options.ApigwOpt{
		Enabled:              c.Enabled,
		ProfilerEnabled:      c.Profiler.Enabled,
		ProfilerGlobal:       c.Profiler.Global,
		ProxyFollowRedirects: c.Proxy.FollowRedirects,
		ProxyOutboundTimeout: c.Proxy.OutboundTimeout,
	}
}

// withOptDefaults fills in values that would otherwise break the caller
//
// An unset outbound timeout yields an already-expired context, so fall back
// to the same default the options carry.
func withOptDefaults(o *options.ApigwOpt) *options.ApigwOpt {
	if o.ProxyOutboundTimeout <= 0 {
		o.ProxyOutboundTimeout = options.Apigw().ProxyOutboundTimeout
	}

	return o
}

func (s Scp) Set(k string, v interface{}) {
	s[k] = v
}

func (s Scp) Get(k string) (v interface{}, err error) {
	var ok bool

	if v, ok = s[k]; !ok {
		err = fmt.Errorf("could not get key on index: %s", k)
		return
	}

	return
}

func (s *Scp) Dict() map[string]interface{} {
	return *s
}

func (s *Scp) Filter(fn func(k string, v interface{}) bool) *Scp {
	ss := Scp{}

	for k, v := range *s {
		if !fn(k, v) {
			continue
		}

		ss[k] = v
	}

	return &ss
}

func (s Scp) Has(k string) (has bool) {
	_, has = s[k]
	return
}
