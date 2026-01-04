package oauth2

import (
	"fmt"
	"net/http"
	"net/url"

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
			prompt = r.Form.Get("prompt")
		)

		// Handle prompt=none (OIDC silent authentication)
		// Per OIDC Core 1.0 Section 3.1.2.1, when prompt=none the Authorization Server
		// MUST NOT display any authentication or consent UI. If the user is not already
		// authenticated or consent is required, an error must be returned.
		if prompt == "none" {
			redirectURI := r.Form.Get("redirect_uri")
			state := r.Form.Get("state")

			if au == nil {
				// User not logged in - return login_required error
				redirectWithOIDCError(w, r, redirectURI, state, "login_required", "User is not authenticated")
				return "", nil
			}

			// Check if client authorization is needed
			// Trusted clients don't require explicit consent
			if client != nil && !client.Trusted && !request.IsOauth2ClientAuthorized(ses) {
				// User logged in but consent is required for non-trusted client
				redirectWithOIDCError(w, r, redirectURI, state, "interaction_required", "User consent is required")
				return "", nil
			}

			// User is authenticated and either:
			// - Client is trusted (no consent needed)
			// - Client was previously authorized in this session
			// Mark client as authorized and continue with the flow
			if client != nil {
				request.SetOauth2ClientAuthorized(ses, true)
				sm.Save(w, r)
			}
		}

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

// redirectWithOIDCError redirects to the client's redirect_uri with an OIDC error response
// Per OIDC Core 1.0 Section 3.1.2.6, errors are returned as query parameters
func redirectWithOIDCError(w http.ResponseWriter, r *http.Request, redirectURI, state, errorCode, errorDesc string) {
	u, err := url.Parse(redirectURI)
	if err != nil {
		// If redirect_uri is invalid, we can't redirect - return a simple error page
		http.Error(w, "Invalid redirect_uri", http.StatusBadRequest)
		return
	}

	q := u.Query()
	q.Set("error", errorCode)
	q.Set("error_description", errorDesc)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusFound)
}
