package scim

import (
	"net/http"
	"regexp"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type (
	Config struct {
		ExternalIdAsPrimary bool
		ExternalIdValidator *regexp.Regexp
	}
)

var (
	log = zap.NewNop()
)

func Guard(opt options.SCIMOpt) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// temp authorization mechanism so we do not have to
		// pre-create users and generate their auth tokens
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authPrefix := "Bearer "
			authHeader := r.Header.Get("Authorization")
			if (len(authPrefix)+len(opt.Secret)) == len(authHeader) && opt.Secret == authHeader[len(authPrefix):] {
				// all good, auth header matches the secret
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "Unauthorized", http.StatusForbidden)

		})
	}
}

func Routes(r chi.Router, cfg Config) {
	r.Route("/Users", func(r chi.Router) {
		uh := &usersHandler{
			externalIdAsPrimary: cfg.ExternalIdAsPrimary,
			externalIdValidator: cfg.ExternalIdValidator,

			svc:     service.DefaultUser,
			passSvc: service.DefaultAuth,
			sec:     getSecurityContext,
		}

		r.Get("/{id}", uh.get)
		r.Post("/", uh.create)
		r.Put("/{id}", uh.replace)
		r.Delete("/{id}", uh.delete)
	})

	r.Route("/Groups", func(r chi.Router) {
		gh := &groupsHandler{
			externalIdAsPrimary: cfg.ExternalIdAsPrimary,
			externalIdValidator: cfg.ExternalIdValidator,

			svc:     service.DefaultRole,
			userSvc: service.DefaultUser,
			sec:     getSecurityContext,
		}

		r.Get("/{id}", gh.get)
		r.Post("/", gh.create)
		r.Put("/{id}", gh.replace)
		r.Patch("/{id}", gh.patch)
		r.Delete("/{id}", gh.delete)
	})
}
