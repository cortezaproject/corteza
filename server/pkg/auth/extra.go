package auth

import (
	"context"
	"net/http"
)

type (
	ExtraReqInfo struct {
		RemoteAddr string
		UserAgent  string
	}
)

func ExtraReqInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ExtraReqInfo{}, ExtraReqInfo{
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
		})))
	})
}

func GetExtraReqInfoFromContext(ctx context.Context) ExtraReqInfo {
	eti := ctx.Value(ExtraReqInfo{})
	if eti != nil {
		return eti.(ExtraReqInfo)
	} else {
		return ExtraReqInfo{}
	}
}
