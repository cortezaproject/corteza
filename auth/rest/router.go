package rest

import (
	"github.com/crusttech/crust/auth/rest/handlers"
	"github.com/crusttech/crust/auth/service"
	"github.com/crusttech/crust/auth/types"
	"github.com/go-chi/chi"
)

func MountRoutes(jwtAuth types.TokenEncoder) func(chi.Router) {
	var userSvc = service.User()

	// Initialize handers & controllers.
	return func(r chi.Router) {
		handlers.NewAuth(Auth{}.New(userSvc, jwtAuth)).MountRoutes(r)
	}
}
