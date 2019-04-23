package service

import (
	"context"

	"github.com/go-chi/chi"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/middleware"
	"github.com/crusttech/crust/internal/routes"
	"github.com/crusttech/crust/system/rest"
)

func Routes(ctx context.Context, th auth.TokenHandler) *chi.Mux {
	r := chi.NewRouter()
	middleware.Mount(ctx, r, flags.http)
	MountRoutes(ctx, r, th)
	routes.Print(r)
	middleware.MountSystemRoutes(ctx, r, flags.http)
	return r
}

func MountRoutes(ctx context.Context, r chi.Router, th auth.TokenHandler) {
	// Only protect application routes with JWT
	r.Group(func(r chi.Router) {
		r.Use(th.Verifier(), th.Authenticator())
		mountRoutes(r, flags.http, rest.MountRoutes(th))
	})
}

func mountRoutes(r chi.Router, opts *config.HTTP, mounts ...func(r chi.Router)) {
	for _, mount := range mounts {
		mount(r)
	}
}
