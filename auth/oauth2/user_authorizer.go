package oauth2

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/go-oauth2/oauth2/v4/server"
	"net/http"
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

		var roles = request.GetRoleMemberships(ses)
		if client.Security != nil {
			// filter user's roles with client security settings
			roles = client.Security.ProcessRoles(roles...)
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
