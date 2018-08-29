package rest

import (
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/rest/handlers"
	"github.com/crusttech/crust/sam/service"
	"github.com/go-chi/chi"
)

func MountRoutes(jwtAuth auth.TokenEncoder) func(chi.Router) {
	// Initialize services
	var (
		channelSvc      = service.Channel()
		messageSvc      = service.Message()
		organisationSvc = service.Organisation()
		teamSvc         = service.Team()
		userSvc         = service.User()
	)

	// Initialize handers & controllers.
	return func(r chi.Router) {
		handlers.NewAuth(Auth{}.New(userSvc, jwtAuth)).MountRoutes(r)

		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthenticationMiddlewareValidOnly)

			handlers.NewChannel(Channel{}.New(channelSvc)).MountRoutes(r)
			handlers.NewMessage(Message{}.New(messageSvc)).MountRoutes(r)
			handlers.NewOrganisation(Organisation{}.New(organisationSvc)).MountRoutes(r)
			handlers.NewTeam(Team{}.New(teamSvc)).MountRoutes(r)
			handlers.NewUser(User{}.New(userSvc, messageSvc)).MountRoutes(r)
		})
	}
}
