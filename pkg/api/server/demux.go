package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/atomic"
)

type (
	// demux (demultiplexer) routes request to one of the underlying routers
	// according to current state
	demux struct {
		state   *atomic.Uint32
		routers map[uint32]chi.Router
	}
)

var _ http.Handler = &demux{}

func Demux(state uint32, r chi.Router) *demux {
	return &demux{
		state:   atomic.NewUint32(state),
		routers: map[uint32]chi.Router{state: r},
	}
}

func (d *demux) State(s uint32) {
	d.state.Store(s)
}

func (d *demux) Router(s uint32, r chi.Router) {
	d.routers[s] = r
}

func (d *demux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer panicRecovery(r.Context(), w)

	var (
		state          = d.state.Load()
		router, exists = d.routers[state]
	)

	if !exists {
		_, _ = fmt.Fprintf(w, "unconfigured request demultiplexor state (%d)", state)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	router.ServeHTTP(w, r)

}
