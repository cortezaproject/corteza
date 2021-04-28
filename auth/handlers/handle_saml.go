package handlers

import (
	"net/http"

	"github.com/cortezaproject/corteza-server/auth/external"
	"go.uber.org/zap"
)

func (h AuthHandlers) samlInit(w http.ResponseWriter, r *http.Request) {
	r = copyProviderToContext(r)
	h.Log.Info("starting saml authentication flow")

	ex := external.NewSamlExternalHandler(h.SamlSPService)
	beginUserAuth(w, r, ex)

	if user, err := completeUserAuth(w, r, ex); err != nil {
		h.Log.Error("failed to complete user auth", zap.Error(err))
		h.handleFailedExternalAuth(w, r, err)
	} else {
		h.handleSuccessfulExternalAuth(w, r, *user)
	}
}
