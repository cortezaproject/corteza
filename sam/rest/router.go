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

	var (
		channel      = Channel{}.New(channelSvc)
		message      = Message{}.New(messageSvc)
		organisation = Organisation{}.New(organisationSvc)
		team         = Team{}.New(teamSvc)
		user         = User{}.New(userSvc, messageSvc)
	)

	// Initialize handers & controllers.
	return func(r chi.Router) {
		handlers.NewAuth(Auth{}.New(userSvc, jwtAuth)).MountRoutes(r)

		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthenticationMiddlewareValidOnly)

			handlers.NewChannel(channel).MountRoutes(r)
			handlers.NewMessage(message).MountRoutes(r)
			handlers.NewOrganisation(organisation).MountRoutes(r)
			handlers.NewTeam(team).MountRoutes(r)
			handlers.NewUser(user).MountRoutes(r)
		})
	}
}
