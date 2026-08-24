package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/crusttech/human/server/auth/request"
	"github.com/crusttech/human/server/system/service"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// Generic OAuth2 authorization-code "Connect" for catalog connections. authorize
// redirects the logged-in user to the provider's consent screen; callback validates
// the single-use state (held in the session, CSRF guard), exchanges the code and
// links the delegated token to the configured connection. Provider-agnostic — one
// flow for every provider.

const connectionOAuth2SessionKey = "oauth2_connect"

type connectionOAuth2State struct {
	CCID     uint64 `json:"ccID"`
	State    string `json:"state"`
	Verifier string `json:"verifier"`
}

func (h *AuthHandlers) connectionOAuth2Authorize(req *request.AuthReq) error {
	if req.AuthUser == nil || req.AuthUser.User == nil {
		return fmt.Errorf("authentication required")
	}

	ccID, err := strconv.ParseUint(req.Request.URL.Query().Get("configuredConnectionID"), 10, 64)
	if err != nil || ccID == 0 {
		return fmt.Errorf("invalid configuredConnectionID")
	}

	state, err := randomToken()
	if err != nil {
		return err
	}
	verifier := oauth2.GenerateVerifier()

	url, err := service.DefaultConfiguredConnection.OAuth2AuthorizeURL(req.Context(), ccID, state, verifier)
	if err != nil {
		return err
	}

	buf, _ := json.Marshal(connectionOAuth2State{CCID: ccID, State: state, Verifier: verifier})
	req.Session.Values[connectionOAuth2SessionKey] = string(buf)

	req.RedirectTo = url
	return nil
}

func (h *AuthHandlers) connectionOAuth2Callback(req *request.AuthReq) error {
	if req.AuthUser == nil || req.AuthUser.User == nil {
		return fmt.Errorf("authentication required")
	}

	// Single-use: consume the pending state regardless of outcome.
	raw, _ := req.Session.Values[connectionOAuth2SessionKey].(string)
	delete(req.Session.Values, connectionOAuth2SessionKey)

	var st connectionOAuth2State
	if raw == "" || json.Unmarshal([]byte(raw), &st) != nil {
		return fmt.Errorf("no pending oauth2 connection")
	}

	q := req.Request.URL.Query()
	if s := q.Get("state"); s == "" || s != st.State {
		return fmt.Errorf("invalid oauth2 state")
	}

	if q.Get("error") != "" {
		req.RedirectTo = connectionOAuth2Done(false)
		return nil
	}

	code := q.Get("code")
	if code == "" {
		return fmt.Errorf("missing authorization code")
	}

	if _, err := service.DefaultConfiguredConnection.OAuth2Complete(req.Context(), st.CCID, code, st.Verifier, req.AuthUser.User.ID); err != nil {
		h.Log.Error("oauth2 connection callback failed", zap.Error(err))
		req.RedirectTo = connectionOAuth2Done(false)
		return nil
	}

	req.RedirectTo = connectionOAuth2Done(true)
	return nil
}

// connectionOAuth2Done is where the callback lands the browser. The Connect flow is
// run in a popup; the frontend detects the flag, notifies the opener and closes.
func connectionOAuth2Done(ok bool) string {
	if ok {
		return "/?oauth2Connected=1"
	}
	return "/?oauth2Connected=0"
}

func randomToken() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
