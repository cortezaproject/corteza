package profiler

import (
	"encoding/base64"
	"fmt"
	"time"

	h "github.com/cortezaproject/corteza-server/pkg/http"
)

type (
	Hit struct {
		ID     string
		Status int
		Route  uint64

		R *h.Request

		Ts *time.Time
		Tf *time.Time
		D  *time.Duration
	}

	Hits map[string][]*Hit
)

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
