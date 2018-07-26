package rest

import (
	"context"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/service"
	"github.com/go-chi/chi"
)

type (
	suspender interface {
		Suspend(context.Context, uint64) error
		Unsuspend(context.Context, uint64) error
	}

	archiver interface {
		Archive(context.Context, uint64) error
		Unarchive(context.Context, uint64) error
	}

	deleter interface {
		Delete(context.Context, uint64) error
	}
)

func MountRoutes(jwtAuth authTokenEncoder) func(chi.Router) {
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
		(&server.AuthHandlers{
			Auth: (&Auth{}).New(userSvc, jwtAuth),
		}).MountRoutes(r)

		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthenticationMiddlewareValidOnly)

			(&server.ChannelHandlers{
				Channel: (&Channel{}).New(channelSvc),
			}).MountRoutes(r)

			(&server.MessageHandlers{
				Message: (&Message{}).New(messageSvc),
			}).MountRoutes(r)

			(&server.OrganisationHandlers{
				Organisation: (&Organisation{}).New(organisationSvc),
			}).MountRoutes(r)

			(&server.TeamHandlers{
				Team: (&Team{}).New(teamSvc),
			}).MountRoutes(r)

			(&server.UserHandlers{
				User: (&User{}).New(userSvc),
			}).MountRoutes(r)
		})
	}
}
