package profiler

import (
	"net/http"
	"sync"
	"time"

	h "github.com/cortezaproject/corteza/server/pkg/http"
)

type (
	Profiler struct {
		mux sync.RWMutex
		l   Hits
	}
)

func New() *Profiler {
	return &Profiler{l: make(Hits)}
}

func (p *Profiler) Hit(r *h.Request) (h *Hit) {
	var (
		n = time.Now()
	)

	h = &Hit{"", http.StatusOK, 0, r, &n, nil, nil}
	h.generateID()

	return
}

func (p *Profiler) Push(h *Hit) (id string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if h.Tf == nil {
		n := time.Now()
		d := n.Sub(*h.Ts)

		h.Tf = &n
		h.D = &d
	}

	h.generateID()

	id = p.id(h.R)
	p.l[id] = append(p.l[id], h)

	return
}

func (p *Profiler) Hits(s Sort) Hits {
	p.mux.RLock()
	defer p.mux.RUnlock()

	ll := p.l.Filter(func(k string, v *Hit) bool {
		var b bool = true

		if s.Path != "" && v.R.URL.Path != s.Path {
			b = false
		}

		if s.Hit != "" && v.ID != s.Hit {
			b = false
		}

		return b
	})

	return ll
}

func (p *Profiler) Purge(f *PurgeFilter) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if f.RouteID == 0 {
		p.l = make(Hits, 0)
		return
	}

	p.l = p.l.Filter(func(k string, v *Hit) bool {
		return v.Route != f.RouteID
	})
}

func (p *Profiler) id(r *h.Request) string {
	return r.URL.Path
}
