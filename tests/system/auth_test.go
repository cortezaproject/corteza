package system

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"net/http"
	"testing"

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
	jwt := res.Response.JWT
	h.a.Greater(len(jwt), 0)

	at, err := service.DefaultStore.LookupAuthOa2tokenByAccess(ctx, jwt)
	h.a.Nil(err)
	h.a.NotNil(at)
	h.a.Greater(len(at.Access), 0)
	h.a.Equal(at.Access, jwt)
}
