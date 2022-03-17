package profiler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	h "github.com/cortezaproject/corteza-server/pkg/http"
)

const (
	// default fallback on amount of items
	FILTER_NUM_ITEMS = 20

	// default fallback on amount of aggregated items
	FILTER_NUM_AGG_ITEMS = 10
)

type (
	Hits map[string][]*Hit

	Profiler struct {
		l Hits
	}

	Hit struct {
		ID     string
		Status int
		Route  uint64

		R *h.Request

		Ts *time.Time
		Tf *time.Time
		D  *time.Duration
	}

	CtxHit []*Stage

	Stage struct {
		Name string
		Ts   *time.Time
		Tf   *time.Time
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

func (h *Hit) generateID() {
	h.ID = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s_%d", h.R.URL.Path, h.Ts.UnixNano())))
}

func (s Hits) Filter(fn func(k string, v *Hit) bool) Hits {
	ss := Hits{}

	for k, v := range s {
		for _, vv := range v {
			if !fn(k, vv) {
				continue
			}

			ss[k] = append(ss[k], vv)
		}
	}

	return ss
}

func (h Hits) Len() int {
	return h.Len()
}

func (h Hits) Less(i, j int) bool {
	return false
}

func (h Hits) Swap(i, j int) {
	return
}
