package oauth2

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/auth/session"
	"github.com/go-oauth2/oauth2/v4/server"
	"net/http"
)

func NewUserAuthorizer(sm *session.Manager, loginURL, clientAuthURL string) server.UserAuthorizationHandler {
	return func(w http.ResponseWriter, r *http.Request) (identity string, err error) {
		var (
			ses    = sm.Get(r)
			user   = session.GetUser(ses)
			client = session.GetOauth2Client(ses)
		)

		// temporary break oauth2 flow by redirecting to
		// login form and ask user to authenticate
		session.SetOauth2AuthParams(ses, r.Form)

		// make sure session is saved!
		sm.Save(w, r)

		// @todo harden security by enforcing login
		//       for each new authorization flow
		if user == nil {
			// user is currently not logged-in;
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		} else {
			if !session.IsOauth2ClientAuthorized(ses) || client == nil {
				// user logged in but we need to re-authenticate the client
				http.Redirect(w, r, clientAuthURL, http.StatusSeeOther)
				return
			}
		}

		var roles = session.GetRoleMemberships(ses)
		if client.Security != nil {
			// filter user's roles with client security settings
			roles = client.Security.ProcessRoles(roles...)
		}

		// User authenticated, client authorized!
		// remove authorization values from session
		session.SetOauth2AuthParams(ses, nil)
		session.SetOauth2Client(ses, nil)
		session.SetOauth2ClientAuthorized(ses, false)

		// make sure session is saved!
		sm.Save(w, r)

		return UserIDSerializer(user.ID, roles...), nil
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
