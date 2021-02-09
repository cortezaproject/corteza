package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/auth/oauth2"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/session"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	systemService "github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	oauth2def "github.com/go-oauth2/oauth2/v4"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// oauth2 flow authorize step
//
// OA2 server internals first run user check (see SetUserAuthorizationHandler lambda)
// to ensure user is authenticated;
func (h AuthHandlers) oauth2Authorize(req *request.AuthReq) (err error) {
	if form := session.GetOAuth2AuthParams(req.Session); form != nil {
		req.Request.Form = form
		h.Log.Debug("restarting oauth2 authorization flow", zap.Any("params", req.Request.Form))
	} else {
		h.Log.Debug("starting new oauth2 authorization flow", zap.Any("params", req.Request.Form))

	}

	session.SetOauth2AuthParams(req.Session, nil)

	var (
		ctx    context.Context
		client *types.AuthClient
	)

	if client, err = h.loadRequestedClient(req); err != nil {
		return err
	}

	// add client to context so we can reach it from client store via context.Value() fn
	//
	// this way we work around the limitations we have with the oauth2 lib.
	ctx = context.WithValue(req.Context(), &oauth2.ContextClientStore{}, client)

	if client != nil {
		// No client validation is done at this point;
		// first, see if user is able to authenticate.
		session.SetOauth2Client(req.Session, client)
	}

	// set to -1 to make sure that wrapping request handler
	// does not send status code!
	req.Status = -1

	// handle authorize request with extended context that now holds the loaded client!
	// we do this
	err = h.OAuth2.HandleAuthorizeRequest(req.Response, req.Request.Clone(ctx))
	if err != nil {
		req.Status = http.StatusInternalServerError
		req.Template = TmplInternalError
		req.Data["error"] = err
	}

	return nil
}

func (h AuthHandlers) oauth2AuthorizeClient(req *request.AuthReq) (err error) {
	var (
		client = session.GetOauth2Client(req.Session)
	)

	if client == nil {
		return fmt.Errorf("flow broken; client missing")
	}

	if err = client.Verify(); err != nil {
		return
	}

	if !h.canAuthorizeClient(req.Context(), req.Client) {
		h.Log.Error("user's roles do not allow authorization of this client", zap.Uint64("ID", req.Client.ID), zap.String("handle", req.Client.Handle))
		session.SetOauth2Client(req.Session, nil)
		session.SetOauth2AuthParams(req.Session, nil)
		req.RedirectTo = GetLinks().Profile
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "danger",
			Text: fmt.Sprintf("Can not authorize '%s', no permissions.", req.Client),
		})
		return nil
	}

	if !req.User.EmailConfirmed {
		req.Data["invalidUser"] = template.HTML(fmt.Sprintf(
			`Can not continue with unauthorized email, 
			visit <a href="%s">your profile</a> and resolve the issue.`,
			GetLinks().Profile,
		))

		req.Data["disabled"] = true
	} else if client.Trusted {
		h.Log.Debug("pre-authorized trusted oauth2 client")

		// Client is trusted, no need to show this screen
		// move forward and authorize oauth2 request
		session.SetOauth2ClientAuthorized(req.Session, true)
		req.RedirectTo = GetLinks().OAuth2Authorize
		return nil
	}

	h.Log.Debug("showing oauth2 client auth form")

	req.Template = TmplOAuth2AuthorizeClient
	return nil
}

func (h AuthHandlers) oauth2AuthorizeClientProc(req *request.AuthReq) (err error) {
	// permissions are already check ed in the  oauth2AuthorizeClient fn,
	// just making sure
	if h.canAuthorizeClient(req.Context(), req.Client) {
		session.SetOauth2Client(req.Session, nil)
		if _, allow := req.Request.Form["allow"]; allow {
			session.SetOauth2ClientAuthorized(req.Session, true)
			req.RedirectTo = GetLinks().OAuth2Authorize
			return
		}
	}

	// handle deny client action from authorize-client form
	//
	// This occurs when user pressed "DENY" button on authorize-client form
	// Remove all and redirect to profile
	//
	session.SetOauth2AuthParams(req.Session, nil)
	req.RedirectTo = GetLinks().Profile
	req.NewAlerts = append(req.NewAlerts, request.Alert{
		Type: "warning",
		Text: "Access for client denied",
	})
	return
}

// Verifies weather current user can authorize this client or not
func (h AuthHandlers) canAuthorizeClient(ctx context.Context, c *types.AuthClient) bool {
	return systemService.DefaultAccessControl.CanAuthorizeAuthClient(ctx, c)
}

func (h AuthHandlers) oauth2Token(req *request.AuthReq) (err error) {
	// Cleanup
	session.SetOauth2ClientAuthorized(req.Session, false)

	req.Status = -1

	client, err := h.loadRequestedClient(req)

	if err != nil {
		return
	}

	if err = client.Verify(); err != nil {
		return fmt.Errorf("invalid client: %w", err)
	} else {
		// add client to context so we can reach it from client store via context.Value() fn
		//
		// this way we work around the limitations we have with the oauth2 lib.
		r := req.Request.Clone(context.WithValue(req.Context(), &oauth2.ContextClientStore{}, client))

		// handle token request with extended context that now holds client!
		err = h.OAuth2.HandleTokenRequest(req.Response, r)
	}

	return
}

func (h AuthHandlers) oauth2Info(w http.ResponseWriter, r *http.Request) {
	ti, err := h.OAuth2.ValidationBearerToken(r)
	if err != nil {
		data, code, header := h.OAuth2.GetErrorData(err)

		_ = header.Write(w)
		w.WriteHeader(code)

		data["active"] = false
		_ = json.NewEncoder(w).Encode(data)

		return
	}

	data := map[string]interface{}{
		"active":    true,
		"scope":     ti.GetScope(),
		"client_id": ti.GetClientID(),
		"exp":       int64(ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"aud":       ti.GetClientID(),
	}

	SubSplit(ti, data)
	if err = Profile(r.Context(), ti, data); err != nil {
		h.Log.Error("failed to add profile data", zap.Error(err))
	}

	_ = json.NewEncoder(w).Encode(data)
}

// oauth2authorizeDefaultClient acts as a proxy for default client
//
// Responsibilities:
//   - handles parameterless request to initialize authorization code flow
//   - accepts redirect_uri via query string that's used to build the oauth2 authorization URL
// (for the rest of the flow, see oauth2authorizeDefaultClientProc)
func (h AuthHandlers) oauth2authorizeDefaultClient(req *request.AuthReq) (err error) {
	if err = h.verifyDefaultClient(); err != nil {
		return
	}

	var (
		params = url.Values{}
	)

	params.Set("client_id", fmt.Sprintf("%d", h.DefaultClient.ID))
	params.Set("response_type", "code")
	params.Set("response_mode", "query")

	if redir, has := req.Request.Form["redirect_uri"]; has && len(redir) > 0 {
		params.Set("redirect_uri", redir[0])
	}

	if scope, has := req.Request.Form["scope"]; has && len(scope) > 0 {
		params.Set("scope", scope[0])
	}

	req.RedirectTo = GetLinks().OAuth2Authorize + "?" + params.Encode()
	return
}

// oauth2authorizeDefaultClient acts as a proxy for default client
//
// Responsibilities:
//   - handles exchange of authorization-code for token
//   - handles issuing of new access token requests
// (for the first part of the flow, see oauth2authorizeDefaultClient)
func (h AuthHandlers) oauth2authorizeDefaultClientProc(req *request.AuthReq) (err error) {
	if err = h.verifyDefaultClient(); err != nil {
		return
	}

	var (
		// extend context and set default client for oauth2server internals
		ctx = context.WithValue(req.Context(), &oauth2.ContextClientStore{}, h.DefaultClient)

		// Clone of the initial request
		// that we'll use for token request validation
		r = req.Request.Clone(ctx)
	)

	if _, has := r.Form["code"]; has {
		r.Form.Set("grant_type", oauth2def.AuthorizationCode.String())
	} else if _, has := r.Form["refresh_token"]; has {
		r.Form.Set("grant_type", oauth2def.Refreshing.String())
	} else {
		// make sure we do not get surprised, remove grant_type param
		// and let oauth2 internals invalidate the request
		r.Form.Del("grant_type")
	}

	// Add id and secret from the default client
	r.SetBasicAuth(
		strconv.FormatUint(h.DefaultClient.ID, 10),
		h.DefaultClient.Secret,
	)

	req.Status = -1
	return h.OAuth2.HandleTokenRequest(req.Response, r)
}

func (h AuthHandlers) verifyDefaultClient() error {
	if h.DefaultClient == nil {
		return fmt.Errorf("no default client configured")
	}

	if err := h.DefaultClient.Verify(); err != nil {
		return fmt.Errorf("invalid client: %w", err)
	}

	return nil
}

// loads client from the request params and verifies other request params against client settings
func (h AuthHandlers) loadRequestedClient(req *request.AuthReq) (client *types.AuthClient, err error) {
	return client, func() (err error) {
		var (
			id       string
			clientID uint64
			found    bool
		)

		if id, _, found = req.Request.BasicAuth(); !found {
			if _, found = req.Request.Form["client_id"]; !found {
				return
			} else {
				id = req.Request.Form.Get("client_id")
			}
		}

		h.Log.Debug("loading client", zap.String("info", id))

		if clientID, err = strconv.ParseUint(id, 10, 64); err != nil {
			return errors.InvalidData("failed to parse client ID from params: %v", err)

		} else if clientID == 0 {
			return errors.InvalidData("invalid client ID")
		}

		if client = session.GetOauth2Client(req.Session); client != nil {
			h.Log.Debug("client loaded from session", zap.Uint64("ID", client.ID))

			// ensure that session holds the right client and
			// not some leftover from a previous flow
			if client.ID != clientID {
				h.Log.Debug("stale client found in session")

				// cleanup leftovers
				client = nil
			} else {
				return
			}
		}

		client, err = h.ClientService.LookupByID(req.Context(), clientID)
		if err != nil {
			return fmt.Errorf("invalid client: %w", err)
		}

		h.Log.Debug("client loaded from store", zap.Uint64("ID", client.ID))
		return
	}()
}

func SubSplit(ti oauth2def.TokenInfo, data map[string]interface{}) {
	userIdWithRoles := strings.SplitN(ti.GetUserID(), " ", 2)
	data["sub"] = userIdWithRoles[0]
	if len(userIdWithRoles) > 1 {
		data["roles"] = userIdWithRoles[1]
	}
}

// Profile fills map with user's data
//
// If scope supports it (contains "profile") user is loaded and
// map is sfilled with username (handle), email and name
func Profile(ctx context.Context, ti oauth2def.TokenInfo, data map[string]interface{}) error {
	if !auth.CheckScope(ti.GetScope(), "profile") {
		return nil
	}

	userID := auth.ExtractUserIDFromSubClaim(ti.GetUserID())
	if userID == 0 {
		return fmt.Errorf("invalid user ID in 'sub' claim")
	}

	user, err := systemService.DefaultUser.FindByID(
		// inject ad-hoc identity into context so that user service is aware who is
		// doing the lookup
		auth.SetIdentityToContext(ctx, auth.NewIdentity(userID)),
		userID,
	)

	if err != nil {
		return err
	}

	data["handle"] = user.Handle
	data["name"] = user.Name
	data["email"] = user.Email

	return nil
}
