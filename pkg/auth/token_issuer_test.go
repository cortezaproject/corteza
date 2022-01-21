package auth

import (
	"testing"

	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/require"
)

func TestIdentityDecoding(t *testing.T) {
	var (
		req = require.New(t)
		ii  = []Identifiable{
			&identity{id: 1, memberOf: []uint64{}},
			&identity{id: 2, memberOf: []uint64{2, 3, 4}},
		}

		tm, err = NewTokenIssuer(WithSecretSigner("test"))
	)

	req.NoError(err)

	for _, i := range ii {
		t.Run(i.String(), func(t *testing.T) {
			var (
				req    = require.New(t)
				token  jwt.Token
				signed []byte
			)

			signed, err = tm.Sign(WithIdentity(i))
			req.NoError(err)

			token, err = jwt.Parse(signed)
			req.NoError(err)

			ift := IdentityFromToken(token)

			req.Equal(i.Identity(), ift.Identity())
			req.Equal(i.Roles(), ift.Roles())
		})

	}
}
