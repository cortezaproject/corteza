package rest

import (
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/sam/rest/handlers"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/store"
	"github.com/go-chi/chi"
	"log"
)

func MountRoutes(jwtAuth types.TokenEncoder) func(chi.Router) {
	// Initialize services
	fs, err := store.New("/tmp/crust/messages")
	if err != nil {
		log.Fatalf("Failed to initialize stor: %v", err)
	}

	var (
		channelSvc      = service.Channel()
		attachmentSvc   = service.Attachment(fs)
		messageSvc      = service.Message(attachmentSvc)
		organisationSvc = service.Organisation()
		teamSvc         = service.Team()
		userSvc         = service.User()
	)

	var (
		channel      = Channel{}.New(channelSvc, attachmentSvc)
		message      = Message{}.New(messageSvc)
		organisation = Organisation{}.New(organisationSvc)
		team         = Team{}.New(teamSvc)
		user         = User{}.New(userSvc, messageSvc)
		attachment   = Attachment{}.New(attachmentSvc)
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
			handlers.NewAttachment(attachment).MountRoutes(r)
		})
	}
}
