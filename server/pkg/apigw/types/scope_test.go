package types

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/stretchr/testify/require"
)

func TestScp_Opts(t *testing.T) {
	var (
		// what routes actually put in the scope
		cfg = func() (c Config) {
			c.Enabled = true
			c.Profiler.Enabled = true
			c.Profiler.Global = true
			c.Proxy.FollowRedirects = true
			c.Proxy.OutboundTimeout = time.Second * 5
			return
		}()

		noTimeout = func() (c Config) {
			c.Enabled = true
			return
		}()

		defaultTimeout = options.Apigw().ProxyOutboundTimeout
	)

	tests := []struct {
		name string
		set  func(s Scp)
		want *options.ApigwOpt
	}{
		{
			name: "config, as stored by the route",
			set:  func(s Scp) { s.Set("opts", cfg) },
			want: &options.ApigwOpt{
				Enabled:              true,
				ProfilerEnabled:      true,
				ProfilerGlobal:       true,
				ProxyFollowRedirects: true,
				ProxyOutboundTimeout: time.Second * 5,
			},
		},
		{
			name: "config pointer",
			set:  func(s Scp) { s.Set("opts", &cfg) },
			want: &options.ApigwOpt{
				Enabled:              true,
				ProfilerEnabled:      true,
				ProfilerGlobal:       true,
				ProxyFollowRedirects: true,
				ProxyOutboundTimeout: time.Second * 5,
			},
		},
		{
			// an unset timeout yields an already expired context
			name: "config without an outbound timeout",
			set:  func(s Scp) { s.Set("opts", noTimeout) },
			want: &options.ApigwOpt{Enabled: true, ProxyOutboundTimeout: defaultTimeout},
		},
		{
			name: "options value",
			set:  func(s Scp) { s.Set("opts", options.ApigwOpt{Debug: true, ProxyOutboundTimeout: time.Second}) },
			want: &options.ApigwOpt{Debug: true, ProxyOutboundTimeout: time.Second},
		},
		{
			name: "options pointer keeps the debug log flag",
			set: func(s Scp) {
				s.Set("opts", &options.ApigwOpt{ProxyEnableDebugLog: true, ProxyOutboundTimeout: time.Second})
			},
			want: &options.ApigwOpt{ProxyEnableDebugLog: true, ProxyOutboundTimeout: time.Second},
		},
		{
			name: "nothing in scope",
			set:  func(s Scp) {},
			want: &options.ApigwOpt{ProxyOutboundTimeout: defaultTimeout},
		},
		{
			name: "unexpected type in scope",
			set:  func(s Scp) { s.Set("opts", "nonsense") },
			want: &options.ApigwOpt{ProxyOutboundTimeout: defaultTimeout},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				req = require.New(t)
				s   = Scp{}
			)

			tt.set(s)

			opts := s.Opts()

			// callers dereference this without checking
			req.NotNil(opts)
			req.Equal(tt.want, opts)
		})
	}
}
