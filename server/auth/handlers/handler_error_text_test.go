package handlers

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/cortezaproject/corteza/server/auth/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/stretchr/testify/require"
)

// localeWith answers the given translations and the key for everything else,
// the way the locale service does
type localeWith map[string]string

func (l localeWith) NS(_ context.Context, ns string) func(key string, rr ...string) string {
	return func(key string, rr ...string) string {
		if msg, has := l[ns+":"+key]; has {
			return msg
		}

		return key
	}
}

func Test_errorText(t *testing.T) {
	var (
		req = require.New(t)
		r   = httptest.NewRequest("GET", "/auth/login", nil)
		err = service.AuthErrPasswordNotSecure()
	)

	// translated when the language has the key the error carries
	authReq := &request.AuthReq{Request: r, Locale: localeWith{"system:auth.errors.passwordNotSecure": "Geslo ni dovolj varno"}}
	req.Equal("Geslo ni dovolj varno", errorText(authReq, err))

	// original message when there is no translation
	authReq = &request.AuthReq{Request: r, Locale: localeWith{}}
	req.Equal(err.Error(), errorText(authReq, err))

	// plain errors are passed through
	req.Equal("boom", errorText(authReq, errors.New("boom")))
	req.Equal("", errorText(authReq, nil))
}
