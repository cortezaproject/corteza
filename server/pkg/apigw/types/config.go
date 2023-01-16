package types

import "time"

type (
	Config struct {
		Enabled bool

		Profiler struct {
			Enabled bool
			Global  bool
		}

		Proxy struct {
			FollowRedirects bool
			OutboundTimeout time.Duration
		}
	}
)
