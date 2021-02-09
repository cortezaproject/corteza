package request

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

type (
	// auth context simplifies auth request & response handling
	AuthReq struct {
		// HTTP request sent
		Request *http.Request

		Response http.ResponseWriter

		// Sessions
		Session *sessions.Session

		// Loaded user (from session)
		User *types.User

		// Current client (when in oauth2 flow)
		Client *types.AuthClient

		// Redirect to
		RedirectTo string

		// Template to render
		Template string

		// Data to render with the template
		Data map[string]interface{}

		// handling flash alerts of all types
		//
		// should not be used for form errors; store them into sep. session keys
		PrevAlerts, NewAlerts []Alert

		// HTTP status to send
		Status int
	}

	Alert struct {
		// primary, secondary, danger, warning...
		Type string
		Text string
		Html template.URL
	}

	// ExtraReqInfo serves as transport struct for additional
	// request information we want to store with the oauth2 token
	//
	// There is effortless way to extend token info that is created inside go-oauth2 lib
	// so we'll attach this struct to (request's) context with middleware (see MountHttpRoutes)
	// and unpack from context when token is created in CortezaTokenStore.Create()
	//
	// ExtraReqInfo struct also serves as context value key!
	ExtraReqInfo struct {
		RemoteAddr string
		UserAgent  string
	}
)

func (req AuthReq) Context() context.Context { return req.Request.Context() }

func (req *AuthReq) SetInternalError(err error) bool {
	if err == nil {
		return false
	}

	req.Status = http.StatusInternalServerError
	req.Data["error"] = err
	return true
}

func (req *AuthReq) PopAlerts() []Alert {
	val, has := req.Session.Values["alerts"]
	if !has {
		return nil
	}

	delete(req.Session.Values, "alerts")

	return val.([]Alert)
}

func (req *AuthReq) SetAlerts(aa ...Alert) {
	if len(aa) == 0 {
		return
	}

	req.Session.Values["alerts"] = aa
}

// retrives key-value from session, stored under request-uri key
func (req *AuthReq) GetKV() map[string]string {
	val, has := req.Session.Values["KV:"+req.Request.RequestURI]
	if !has {
		return nil
	}

	return val.(map[string]string)
}

// sets key-value value to session under request-uri key
func (req *AuthReq) SetKV(val map[string]string) {
	if val == nil {
		delete(req.Session.Values, "KV:"+req.Request.RequestURI)
	} else {
		req.Session.Values["KV:"+req.Request.RequestURI] = val
	}
}

func ExtraReqInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ExtraReqInfo{}, ExtraReqInfo{
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
		})))
	})
}

func GetExtraReqInfoFromContext(ctx context.Context) ExtraReqInfo {
	eti := ctx.Value(ExtraReqInfo{})
	if eti != nil {
		return eti.(ExtraReqInfo)
	} else {
		return ExtraReqInfo{}
	}
}
