package websocket

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func Watch(ctx context.Context) {
	events := repository.Events()

	go func() {
		defer sentry.Recover()

		for {
			if err := eq.feedSessions(ctx, events, store); err != nil {
				if err == context.Canceled {
					return
				}

				logger.Default().Error("session event feed error", zap.Error(err))
			}
		}
	}()
	eq.store(ctx, events)
}

func (ws Websocket) ApiServerRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/websocket", func(r chi.Router) {
			r.Get("/", ws.Open)
		})
	})
}
