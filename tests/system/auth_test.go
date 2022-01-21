package system

import (
	"context"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestAuthImpersonate(t *testing.T) {

	var (
		h     = newHelper(t)
		ctx   = context.Background()
		user  = h.createUserWithEmail(h.randEmail())
		input = &struct {
			UserID uint64 `json:",string"`
		}{
			UserID: user.ID,
		}
	)

	helpers.AllowMe(h, types.UserRbacResource(user.ID), "impersonate")

	var res struct {
		Response struct {
			JWT string `json:"jwt"`
		} `json:"response"`
	}
	h.apiInit().
		Post("/auth/impersonate").
		Header("Accept", "application/json").
		JSON(helpers.JSON(input)).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		End().
		JSON(&res)

	// make sure response has JWT token
	signedToken := res.Response.JWT
	h.a.NotEmpty(signedToken)

	token, err := jwt.Parse([]byte(signedToken))
	h.a.Nil(err)

	at, err := service.DefaultStore.LookupAuthOa2tokenByAccess(ctx, token.JwtID())
	h.a.Nil(err)
	h.a.NotNil(at)
	h.a.Equal(at.Access, token.JwtID())
}
