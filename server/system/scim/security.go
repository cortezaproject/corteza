package scim

import (
	"context"
	"net/http"

	"github.com/crusttech/human/server/pkg/auth"
)

type (
	getSecurityContextFn func(r *http.Request) context.Context
)

// Set service user to request's identity
func getSecurityContext(r *http.Request) context.Context {
	return auth.SetIdentityToContext(r.Context(), auth.ServiceUser())
}
