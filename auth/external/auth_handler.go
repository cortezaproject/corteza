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
		s := (session.(samlsp.JWTSessionClaims)).GetAttributes()

		u = &types.ExternalAuthUser{}
		u.Provider = "saml"

		// get identifier for use with Corteza (email)
		u.Email = eh.service.GuessIdentifier(s)
		u.Name = s.Get(eh.service.IDPUserMeta.Name)
		u.NickName = s.Get(eh.service.IDPUserMeta.Handle)
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
