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

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/spf13/cast"

	"github.com/go-chi/jwtauth"
	oauth2errors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/cortezaproject/corteza/server/auth/request"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	systemService "github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	oauth2def "github.com/go-oauth2/oauth2/v4"
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
		client *types.AuthClient
	)

	if client, err = h.loadRequestedClient(req); err != nil {
		return err
	}

	if client != nil {
		// No client validation is done at this point;
		// first, see if user is able to authenticate.
		request.SetOauth2Client(req.Session, client)

		// ensure we're dealing with client ID in case someone used handle
		req.Request.Form.Set("client_id", strconv.FormatUint(client.ID, 10))
	}

	// set to -1 to make sure wrapping request handler
	// does not send status code!
	req.Status = -1

	err = h.OAuth2.HandleAuthorizeRequest(req.Response, req.Request)
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

// oauth2Info handler validates token and responds with decoded claims
func (h AuthHandlers) oauth2Info(w http.ResponseWriter, r *http.Request) {
	var (
		jt     jwt.Token
		claims map[string]interface{}
	)

	err := func() (err error) {
		if jt, claims, err = jwtauth.FromContext(r.Context()); err != nil {
			return
		}

		if err = auth.TokenIssuer.Validate(r.Context(), jt); err != nil {
			return
		}

		return nil
	}()

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

	_ = json.NewEncoder(w).Encode(claims)
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
		// Clone of the initial request
		// that we'll use for token request validation
		r = req.Request.Clone(req.Context())
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
			id    string
			found bool
		)

		if id, _, found = req.Request.BasicAuth(); !found {
			if _, found = req.Request.Form["client_id"]; !found {
				return
			} else {
				id = req.Request.Form.Get("client_id")
			}
		}

		client, err = h.ClientService.Lookup(req.Context(), id)
		if err != nil {
			return fmt.Errorf("invalid client: %w", err)
		}

		h.Log.Debug("client loaded from store", zap.Uint64("ID", client.ID))
		return
	}()
}

func (h AuthHandlers) handleTokenRequest(req *request.AuthReq, client *types.AuthClient) error {
	var (
		r   = req.Request
		w   = req.Response
		ctx = req.Context()

		user *types.User
	)

	req.Status = -1

	if err := client.Verify(); err != nil {
		return h.tokenError(w, fmt.Errorf("invalid client: %w", err))
	}

	r = req.Request.Clone(ctx)

	gt, tgr, err := h.OAuth2.ValidationTokenRequest(r)
	if err != nil {
		return h.tokenError(w, err)
	}

	// Manually set the impersonating user & roles when client credentials
	if gt == oauth2def.ClientCredentials {
		tgr.UserID = strings.Join(
			append([]string{strconv.FormatUint(client.Security.ImpersonateUser, 10)}, client.Security.ForcedRoles...),
			" ",
		)
	}

	ti, err := h.OAuth2.GetAccessToken(ctx, gt, tgr)
	if err != nil {
		return h.tokenError(w, err)
	}

	suCtx := auth.SetIdentityToContext(ctx, auth.ServiceUser())

	switch gt {
	case oauth2def.ClientCredentials:
		// Authenticated with client credentials!

		// First, validate client's security settings
		if client.Security == nil || client.Security.ImpersonateUser == 0 {
			return h.tokenError(w, errors.Internal("auth client security configuration invalid"))
		}

		// Load the user
		if user, err = h.UserService.FindByAny(suCtx, client.Security.ImpersonateUser); err != nil {
			return h.tokenError(w, fmt.Errorf("could not generate token for impersonated user: %v", err))
		}

	case oauth2def.AuthorizationCode, oauth2def.Refreshing:
		userID := ti.GetUserID()
		if i := strings.Index(ti.GetUserID(), " "); i > 0 {
			// userID field from the token could contain encoded roles
			// @todo investigate if role-encoding into user-id field is still needed?
			userID = userID[:i]
		}

		if req.AuthUser != nil && req.AuthUser.User != nil && req.AuthUser.User.ID == cast.ToUint64(userID) {
			user = req.AuthUser.User
		} else if user, err = h.UserService.FindByAny(suCtx, userID); err != nil {
			return h.tokenError(w, fmt.Errorf("could not generate token: %v", err))
		}

	default:
		return fmt.Errorf("unsupported oauth2 grant type: %v", gt)
	}

	var (
		signed []byte
		scope  = strings.Split(ti.GetScope(), " ")
	)

	// Here set roles to signed
	signed, err = auth.TokenIssuer.Sign(
		auth.WithAccessToken(ti.GetAccess()),
		auth.WithIdentity(user),
		func(tr *auth.TokenRequest) error {
			// Calculate user's roles
			roles := user.Roles()
			if client.Security != nil {
				roles = auth.ApplyRoleSecurity(
					payload.ParseUint64s(client.Security.PermittedRoles),
					payload.ParseUint64s(client.Security.ProhibitedRoles),
					payload.ParseUint64s(client.Security.ForcedRoles),
					roles...,
				)
			}
			tr.Roles = roles
			return nil
		},
		auth.WithClientID(client.ID),
		auth.WithScope(scope...),
	)

	if err != nil {
		return h.tokenError(w, err)
	}

	// modify token info with signed JWT
	// this will be sent back to the user
	ti.SetAccess(string(signed))

	response := h.OAuth2.GetTokenData(ti)

	// in case client is configured with "openid" scope,
	// we'll add "id_token" with all required (by OIDC) details encoded
	if strings.Contains(client.Scope, "openid") {
		var idToken []byte
		if idToken, err = generateIdToken(user, client, ti, h.Opt.BaseURL); err != nil {
			return h.tokenError(w, err)
		}
		response["id_token"] = string(idToken)
	}

	return writeResponse(w, response, nil)
}

func (h AuthHandlers) tokenError(w http.ResponseWriter, err error) error {
	data, statusCode, header := h.OAuth2.GetErrorData(err)
	return writeResponse(w, data, header, statusCode)
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

// Generates ID token that is part of OIDC flow for doing corteza-to-corteza auth
func generateIdToken(user *types.User, client *types.AuthClient, ti oauth2def.TokenInfo, baseURL string) (_ []byte, err error) {
	token := jwt.New()
	if err = token.Set(jwt.IssuerKey, baseURL); err != nil {
		return
	}

	// we do not know what the admin used for client key value
	// on the receiving end, so we'll encode both,
	// client's ID, and it's handle
	aud := []string{strconv.FormatUint(client.ID, 10)}
	if len(client.Handle) > 0 {
		aud = append(aud, client.Handle)
	}

	if err = token.Set("aud", aud); err != nil {
		return
	}
	if err = token.Set("user_id", strconv.FormatUint(user.ID, 10)); err != nil {
		return
	}
	if err = token.Set("email", user.Email); err != nil {
		return
	}
	if err = token.Set(jwt.ExpirationKey, now().Add(ti.GetAccessExpiresIn()).Unix()); err != nil {
		return
	}

	return jwt.Sign(token, jwa.HS512, []byte(client.Secret))
}

func writeResponse(w http.ResponseWriter, data map[string]interface{}, header http.Header, statusCode ...int) error {
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
