package external

import (
	"net/http"

	"github.com/cortezaproject/corteza-server/auth/saml"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/crewjam/saml/samlsp"
	"github.com/markbates/goth/gothic"
)

type (
	ExternalAuthHandler interface {
		BeginUserAuth(http.ResponseWriter, *http.Request)
		CompleteUserAuth(http.ResponseWriter, *http.Request) (u *types.ExternalAuthUser, err error)
	}

	externalSamlAuthHandler struct {
		service saml.SamlSPService
	}

	externalDefaultAuthHandler struct{}
)

func (eh *externalDefaultAuthHandler) BeginUserAuth(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (eh *externalDefaultAuthHandler) CompleteUserAuth(w http.ResponseWriter, r *http.Request) (u *types.ExternalAuthUser, err error) {
	gu, err := gothic.CompleteUserAuth(w, r)
	u = &types.ExternalAuthUser{gu}

	return
}

func (eh *externalSamlAuthHandler) BeginUserAuth(w http.ResponseWriter, r *http.Request) {
	_, err := eh.service.Handler().Session.GetSession(r)

	if err == samlsp.ErrNoSession {
		eh.service.Handler().HandleStartAuthFlow(w, r)
	}
}

func (eh *externalSamlAuthHandler) CompleteUserAuth(w http.ResponseWriter, r *http.Request) (u *types.ExternalAuthUser, err error) {
	var (
		session samlsp.Session
	)

	if session, err = eh.service.Handler().Session.GetSession(r); err != nil {
		return
	}

	if session != nil {
		sess := (session.(samlsp.JWTSessionClaims))

		if sess.StandardClaims.Valid() != nil {
			return nil, sess.StandardClaims.Valid()
		}

		u = &types.ExternalAuthUser{}
		u.Provider = "saml"

		// get identifier for use with Corteza (email)
		u.Email = eh.service.GuessIdentifier(sess.Attributes)

		// try to get email from jwt claims
		if u.Email == "" {
			u.Email = sess.StandardClaims.Subject
			u.Name = sess.StandardClaims.Subject
			u.NickName = sess.StandardClaims.Subject
		} else {
			u.Name = sess.Attributes.Get(eh.service.IDPUserMeta.Name)
			u.NickName = sess.Attributes.Get(eh.service.IDPUserMeta.Handle)
		}

		u.UserID = sess.Attributes.Get("SessionIndex")
	}

	return
}

func NewSamlExternalHandler(s saml.SamlSPService) *externalSamlAuthHandler {
	return &externalSamlAuthHandler{
		service: s,
	}
}

func NewDefaultExternalHandler() *externalDefaultAuthHandler {
	return &externalDefaultAuthHandler{}
}
