package websocket

import (
	"github.com/go-chi/chi/v5"
)

// MountRoutes initialize route for websocket
// No middleware used, since anyone can open connection and
// send first message with valid JWT token,
// If it's valid then we keep the connection open or close it
func (ws *server) MountRoutes(r chi.Router) {
	// Initialize handlers & controllers.
	r.Get("/", ws.Open)
}
