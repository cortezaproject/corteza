package rest

import (
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/crm/rest/server"
	"github.com/go-chi/chi"
)

type (
	authTokenEncoder interface {
		Encode(identity auth.Identifiable) string
	}
)

func MountRoutes(jwtAuth authTokenEncoder) func(chi.Router) {
	// @todo pass jwtAuth to auth handlers (signUp) for JWT generation

	// Initialize handers & controllers.
	return func(r chi.Router) {
		r.Use(auth.AuthenticationMiddlewareValidOnly)

		(&server.FieldHandlers{
			Field: Field{}.New(),
		}).MountRoutes(r)

		(&server.ModuleHandlers{
			Module: Module{}.New(),
		}).MountRoutes(r)
	}
}
