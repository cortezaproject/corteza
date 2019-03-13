package service

import (
	"context"

	"github.com/go-chi/chi"

	"github.com/crusttech/crust/crm/rest"
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/middleware"
	"github.com/crusttech/crust/internal/routes"
)

func Routes(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()
	middleware.Mount(ctx, r, flags.http)
	MountRoutes(ctx, r)
	routes.Print(r)
	middleware.MountSystemRoutes(ctx, r, flags.http)
	return r
}

func MountRoutes(ctx context.Context, r chi.Router) {
	// Only protect application routes with JWT
	r.Group(func(r chi.Router) {
		r.Use(jwtVerifier, jwtAuthenticator)
		mountRoutes(r, flags.http, rest.MountRoutes())
	})
}

func mountRoutes(r chi.Router, opts *config.HTTP, mounts ...func(r chi.Router)) {
	for _, mount := range mounts {
		mount(r)
	}
}
