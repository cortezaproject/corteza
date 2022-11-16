package profiler

import (
	"net/http"
	"time"

	h "github.com/cortezaproject/corteza/server/pkg/http"
)

type (
	Profiler struct {
		l Hits
	}
)

func New() *Profiler {
	return &Profiler{make(Hits)}
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

func (p *Profiler) id(r *h.Request) string {
	return r.URL.Path
}
