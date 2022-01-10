package handlers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwk"

	"github.com/cortezaproject/corteza-server/auth/oauth2"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	systemService "github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	oauth2def "github.com/go-oauth2/oauth2/v4"
	oauth2errors "github.com/go-oauth2/oauth2/v4/errors"
	"go.uber.org/zap"
)

// oauth2 flow authorize step
//
// OA2 server internals first run user check (see SetUserAuthorizationHandler lambda)
// to ensure user is authenticated;
func (h AuthHandlers) oauth2Authorize(req *request.AuthReq) (err error) {
	if form := request.GetOAuth2AuthParams(req.Session); form != nil {
		req.Request.Form = form
		h.Log.Debug("restarting oauth2 authorization flow", zap.Any("params", req.Request.Form))
	} else {
		h.Log.Debug("starting new oauth2 authorization flow", zap.Any("params", req.Request.Form))

	}

	request.SetOauth2AuthParams(req.Session, nil)

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
		request.SetOauth2Client(req.Session, client)
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
		client = request.GetOauth2Client(req.Session)
	)

	if client == nil {
		return fmt.Errorf("flow broken; client missing")
	}

	if err = client.Verify(); err != nil {
		return
	}

	if !h.canAuthorizeClient(req.Context(), req.Client) {
		h.Log.Error("user's roles do not allow authorization of this client", zap.Uint64("ID", req.Client.ID), zap.String("handle", req.Client.Handle))
		request.SetOauth2Client(req.Session, nil)
		request.SetOauth2AuthParams(req.Session, nil)
		req.RedirectTo = GetLinks().Profile

		t := translator(req, "auth")
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "danger",
			Text: t("oauth2-authorize-client.alerts.denied", "client", req.Client.Meta.Name),
		})
		return nil
	}

	if !req.AuthUser.User.EmailConfirmed {
		req.Data["invalidUser"] = true
		req.Data["disabled"] = true
	} else if client.Trusted {
		h.Log.Debug("pre-authorized trusted oauth2 client")

		// Client is trusted, no need to show this screen
		// move forward and authorize oauth2 request
		request.SetOauth2ClientAuthorized(req.Session, true)
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
		request.SetOauth2Client(req.Session, nil)
		if _, allow := req.Request.Form["allow"]; allow {
			request.SetOauth2ClientAuthorized(req.Session, true)
			req.RedirectTo = GetLinks().OAuth2Authorize
			return
		}
	}

	// handle deny client action from authorize-client form
	//
	// This occurs when user pressed "DENY" button on authorize-client form
	// Remove all and redirect to profile
	//
	request.SetOauth2AuthParams(req.Session, nil)
	req.RedirectTo = GetLinks().Profile
	t := translator(req, "auth")
	req.NewAlerts = append(req.NewAlerts, request.Alert{
		Type: "warning",
		Text: t("oauth2-authorize-client.alerts.denied", "client", req.Client.Meta.Name),
	})
	return
}

// Verifies weather current user can authorize this client or not
func (h AuthHandlers) canAuthorizeClient(ctx context.Context, c *types.AuthClient) bool {
	return systemService.DefaultAccessControl.CanAuthorizeAuthClient(ctx, c)
}

func (h AuthHandlers) oauth2Token(req *request.AuthReq) (err error) {
	// Cleanup
	request.SetOauth2ClientAuthorized(req.Session, false)

	req.Status = -1

	client, err := h.loadRequestedClient(req)

	if err != nil {
		return h.tokenError(req.Response, err)
	}

	return h.handleTokenRequest(req, client)
}

func (h AuthHandlers) oauth2Info(w http.ResponseWriter, r *http.Request) {
	ti, err := h.OAuth2.ValidationBearerToken(r)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			// Gracefully handle request cancellation
			//
			// This happens during the login procedure when browser is
			// sent through a couple of quick redirects that can
			// terminate requests
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var (
			data   = make(map[string]interface{})
			code   int
			header http.Header
		)

		if errors.Is(err, oauth2errors.ErrInvalidAccessToken) {
			code = http.StatusForbidden
			data["error"] = err.Error()
		} else {
			data, code, header = h.OAuth2.GetErrorData(err)
		}

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

	if state, has := req.Request.Form["state"]; has && len(state) > 0 {
		params.Set("state", state[0])
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

	req.Request = r

	return h.handleTokenRequest(req, h.DefaultClient)
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

		if client = request.GetOauth2Client(req.Session); client != nil {
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

func (h AuthHandlers) handleTokenRequest(req *request.AuthReq, client *types.AuthClient) error {
	req.Status = -1

	var (
		r   = req.Request
		w   = req.Response
		ctx = req.Context()
	)

	req.Status = -1

	if err := client.Verify(); err != nil {
		return h.tokenError(w, fmt.Errorf("invalid client: %w", err))
	}

	// add client to context so we can reach it from client store via context.Value() fn
	// this way we work around the limitations we have with the oauth2 lib.
	ctx = context.WithValue(ctx, &oauth2.ContextClientStore{}, client)
	r = req.Request.Clone(ctx)

	gt, tgr, err := h.OAuth2.ValidationTokenRequest(r)
	if err != nil {
		return h.tokenError(w, err)
	}

	if gt == oauth2def.ClientCredentials {
		// Authenticated with client credentials!
		//
		// We'll use info from client security
		if client.Security == nil || client.Security.ImpersonateUser == 0 {
			return h.tokenError(w, errors.Internal("auth client security configuration invalid"))
		}

		tgr.UserID = strings.Join(append(
			[]string{fmt.Sprintf("%d", client.Security.ImpersonateUser)},
			client.Security.ForcedRoles...,
		), " ")
	}

	ti, err := h.OAuth2.GetAccessToken(ctx, gt, tgr)
	if err != nil {
		return h.tokenError(w, err)
	}

	return token(w, h.OAuth2.GetTokenData(ti), nil)
}

func (h AuthHandlers) tokenError(w http.ResponseWriter, err error) error {
	data, statusCode, header := h.OAuth2.GetErrorData(err)
	return token(w, data, header, statusCode)
}

func (h AuthHandlers) oauth2PublicKeys(w http.ResponseWriter, r *http.Request) {
	// handle error
	handleErr := func(code int, err error) {
		var (
			data   = make(map[string]interface{})
			header http.Header
		)

		data["error"] = err.Error()

		_ = header.Write(w)
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(data)
	}

	// @todo determine weather to get it form user and save it in settings
	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		handleErr(http.StatusBadRequest, fmt.Errorf("failed to generate new RSA privatre key: %w", err))
		return
	}

	key, err := jwk.New(raw)
	if err != nil {
		handleErr(http.StatusInternalServerError, fmt.Errorf("failed to create symmetric key: %w", err))
		return
	}

	publicKey, err := key.PublicKey()
	if err != nil {
		handleErr(http.StatusInternalServerError, fmt.Errorf("failed to create public key: %w", err))
		return
	}

	_ = json.NewEncoder(w).Encode(publicKey)
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
// map is filled with username (handle), email and name
func Profile(ctx context.Context, ti oauth2def.TokenInfo, data map[string]interface{}) error {
	if !auth.CheckScope(ti.GetScope(), "profile") {
		return nil
	}

	userID, roles := auth.ExtractFromSubClaim(ti.GetUserID())
	if userID == 0 {
		return fmt.Errorf("invalid user ID in 'sub' claim")
	}

	user, err := systemService.DefaultUser.FindByID(
		// inject ad-hoc identity into context so that user service is aware who is
		// doing the lookup
		auth.SetIdentityToContext(ctx, auth.Authenticated(userID, roles...)),
		userID,
	)

	if err != nil {
		return err
	}

	data["handle"] = user.Handle
	data["name"] = user.Name
	data["email"] = user.Email

	if user.Meta != nil && user.Meta.PreferredLanguage != "" {
		data["preferred_language"] = user.Meta.PreferredLanguage
	}

	return nil
}

func token(w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	for key := range header {
		w.Header().Set(key, header.Get(key))
	}

	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}

	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
