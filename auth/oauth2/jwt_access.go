package oauth2

import (
	"context"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/spf13/cast"
)

// JWTAccessGenerate generate the jwt access token
type (
	JWTAccessGenerate struct {
		tm auth.TokenGenerator
	}
)

// NewJWTAccessGenerate create to generate the jwt access token instance
//
// @todo move this to pkg/auth (??) so it can be re-used
func NewJWTAccessGenerate(tg auth.TokenGenerator) *JWTAccessGenerate {
	return &JWTAccessGenerate{tg}
}

// Token based on the UUID generated token
func (a *JWTAccessGenerate) Token(_ context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (_ string, refresh string, err error) {
	var (
		user     auth.Identifiable
		rawToken []byte
	)

	{
		// extract user ID and roles from a space-delimited list of IDs stored in userID
		userIdWithRoles := strings.Split(data.TokenInfo.GetUserID(), " ")
		if len(userIdWithRoles) == 1 {
			user = auth.Authenticated(cast.ToUint64(userIdWithRoles[0]))
		} else {
			user = auth.Authenticated(
				cast.ToUint64(userIdWithRoles[0]),
				payload.ParseUint64s(userIdWithRoles)...,
			)
		}
	}

	rawToken, err = a.tm.Encode(user, cast.ToUint64(data.Client.GetID()), data.TokenInfo.GetScope())
	if err != nil {
		return
	}

	if isGenRefresh {
		refresh = string(rand.Bytes(48))
	}

	return string(rawToken), refresh, nil
}
