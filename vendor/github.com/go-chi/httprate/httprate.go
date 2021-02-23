package httprate

import (
	"net"
	"net/http"
	"strings"
	"time"
)

func Limit(requestLimit int, windowLength time.Duration, keyFuncs ...KeyFunc) func(next http.Handler) http.Handler {
	return NewRateLimiter(requestLimit, windowLength, nil, keyFuncs...).Handler
}

type KeyFunc func(r *http.Request) (string, error)

func LimitAll(requestLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return Limit(requestLimit, windowLength)
}

func LimitByIP(requestLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return Limit(requestLimit, windowLength, KeyByIP)
}

func KeyByIP(r *http.Request) (string, error) {
	var ip string

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
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
