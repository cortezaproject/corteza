package oauth2

import (
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza/server/auth/request"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/go-oauth2/oauth2/v4/server"
)

func NewUserAuthorizer(sm *request.SessionManager, loginURL, clientAuthURL string) server.UserAuthorizationHandler {
	return func(w http.ResponseWriter, r *http.Request) (identity string, err error) {
		var (
			ses    = sm.Get(r)
			au     = request.GetAuthUser(ses)
			client = request.GetOauth2Client(ses)
		)

		// temporary break oauth2 flow by redirecting to
		// login form and ask user to authenticate
		request.SetOauth2AuthParams(ses, r.Form)

		// make sure session is saved!
		sm.Save(w, r)

		// @todo harden security by enforcing login
		//       for each new authorization flow
		if au == nil {
			// user is currently not logged-in;
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		} else {
			if !request.IsOauth2ClientAuthorized(ses) || client == nil {
				// user logged in but we need to re-authenticate the client
				http.Redirect(w, r, clientAuthURL, http.StatusSeeOther)
				return
			}
		}

		roles := au.User.Roles()
		if client.Security != nil {
			// filter user's roles with client security settings
			roles = internalAuth.ApplyRoleSecurity(
				payload.ParseUint64s(client.Security.PermittedRoles),
				payload.ParseUint64s(client.Security.ProhibitedRoles),
				payload.ParseUint64s(client.Security.ForcedRoles),
				roles...,
			)
		}

		// User authenticated, client authorized!
		// remove authorization values from session
		request.SetOauth2AuthParams(ses, nil)
		request.SetOauth2Client(ses, nil)
		request.SetOauth2ClientAuthorized(ses, false)

		// make sure session is saved!
		sm.Save(w, r)

		return UserIDSerializer(au.User.ID, roles...), nil
	}
}

// UserIDSerializer packs user ID and IDs of all roles
// into space delimited string
//
// Main reason to do this is to simplify JWT claims encoding;
// we do not have  access to user's membership info from there)
func UserIDSerializer(userID uint64, rr ...uint64) string {
	identity := fmt.Sprintf("%d", userID)
	for _, r := range rr {
		identity += fmt.Sprintf(" %d", r)
	}

	return identity
}
