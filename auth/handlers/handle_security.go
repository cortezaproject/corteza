package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
)

func (h *AuthHandlers) security(req *request.AuthReq) error {
	req.Template = TmplSecurity
	return nil
}
