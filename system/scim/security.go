package scim

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"net/http"
)

type (
	getSecurityContextFn func(r *http.Request) context.Context
)

// All actions are in security context of a superuser for now
//
func getSecurityContext(r *http.Request) context.Context {
	return auth.SetSuperUserContext(r.Context())
}
