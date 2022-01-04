package httprate

import (
	"net"
	"net/http"
	"strings"
	"time"
)

func Limit(requestLimit int, windowLength time.Duration, options ...Option) func(next http.Handler) http.Handler {
	return NewRateLimiter(requestLimit, windowLength, options...).Handler
}

type KeyFunc func(r *http.Request) (string, error)
type Option func(rl *rateLimiter)

func LimitAll(requestLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return Limit(requestLimit, windowLength)
}

func LimitByIP(requestLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return Limit(requestLimit, windowLength, WithKeyFuncs(KeyByIP))
}

func KeyByIP(r *http.Request) (string, error) {
	var ip string

	if tcip := r.Header.Get("True-Client-IP"); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else {
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
	}

	return ip, nil
}

func KeyByEndpoint(r *http.Request) (string, error) {
	return r.URL.Path, nil
}

func WithKeyFuncs(keyFuncs ...KeyFunc) Option {
	return func(rl *rateLimiter) {
		if len(keyFuncs) > 0 {
			rl.keyFn = composedKeyFunc(keyFuncs...)
		}
	}
}

func WithLimitHandler(h http.HandlerFunc) Option {
	return func(rl *rateLimiter) {
		rl.onRequestLimit = h
	}
}

func WithLimitCounter(c LimitCounter) Option {
	return func(rl *rateLimiter) {
		rl.limitCounter = c
	}
}

func composedKeyFunc(keyFuncs ...KeyFunc) KeyFunc {
	return func(r *http.Request) (string, error) {
		var key strings.Builder
		for i := 0; i < len(keyFuncs); i++ {
			k, err := keyFuncs[i](r)
			if err != nil {
				return "", err
			}
			key.WriteString(k)
		}
		return key.String(), nil
	}
}
