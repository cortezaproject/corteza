package profiler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"time"

	h "github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/system/types"
)

const (
	// default fallback on amount of items
	FILTER_NUM_ITEMS = 20

	// default fallback on amount of aggregated items
	FILTER_NUM_AGG_ITEMS = 10
)

var sortAggFields = []string{"path", "count", "size_min", "size_max", "size_avg", "time_min", "time_max", "time_avg"}
var sortRouteFields = []string{"time_start", "time_finish", "time_duration"}

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

	Sort struct {
		Hit    string
		Path   string
		Size   uint64
		Before string
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

func (p *Profiler) Dump(s Sort) Hits {
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

func StartHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// add some info to context
		next.ServeHTTP(rw, r)
	})
}

func FinishHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r)
	})
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

func SortAggregation(list *types.ApigwProfilerAggregationSet, filter *types.ApigwProfilerFilter) {
	for _, ff := range sortAggFields {
		fe := filter.Sort.Get(ff)

		if fe == nil {
			continue
		}

		if filter.Sort.Get(ff).Descending {
			sort.Sort(sort.Reverse(getSortType(ff, list)))
			break
		}

		sort.Sort(getSortType(ff, list))
		break
	}
}

func FilterAggregation(list *types.ApigwProfilerAggregationSet, filter *types.ApigwProfilerFilter) {
	var (
		dec string = ""
		i   uint   = 0
		b          = filter.Before == ""
	)

	if filter.Limit == 0 {
		filter.Limit = FILTER_NUM_AGG_ITEMS
	}

	dec, _ = decodeRoutePath(filter.Before)

	*list, _ = list.Filter(func(apa *types.ApigwProfilerAggregation) (bool, error) {
		// after a specific hit and inside the limits
		if b && i < filter.Limit {
			i++
			filter.Next = encodeRoutePath(apa.Path)
			return true, nil
		}

		// after the specific hit check
		if dec != "" && b == false {
			b = apa.Path == dec
		}

		return false, nil
	})

	return
}

func FilterHits(list *types.ApigwProfilerHitSet, filter *types.ApigwProfilerFilter) {
	var (
		i uint = 0
		b      = filter.Before == ""
	)

	if filter.Limit == 0 {
		filter.Limit = FILTER_NUM_ITEMS
	}

	*list, _ = list.Filter(func(aph *types.ApigwProfilerHit) (bool, error) {
		// after a specific hit and inside the limits
		if b && i < filter.Limit {
			i++
			filter.Next = aph.ID
			return true, nil
		}

		// after the specific hit check
		if filter.Before != "" && b == false {
			b = aph.ID == filter.Before
		}

		return false, nil
	})

	return
}

func SortHits(list *types.ApigwProfilerHitSet, filter *types.ApigwProfilerFilter) {
	for _, ff := range sortRouteFields {
		fe := filter.Sort.Get(ff)

		if fe == nil {
			continue
		}

		if filter.Sort.Get(ff).Descending {
			sort.Sort(sort.Reverse(getSortTypeHit(ff, list)))
			break
		}

		sort.Sort(getSortTypeHit(ff, list))
		break
	}
}

func getSortType(s string, list *types.ApigwProfilerAggregationSet) sort.Interface {
	switch s {
	case "path":
		return types.ByPath(*list)
	case "count":
		return types.ByCount(*list)
	case "size_min":
		return types.BySizeMin(*list)
	case "size_max":
		return types.BySizeMax(*list)
	case "size_avg":
		return types.BySizeAvg(*list)
	case "time_min":
		return types.ByTimeMin(*list)
	case "time_max":
		return types.ByTimeMax(*list)
	case "time_avg":
		return types.ByTimeAvg(*list)
	default:
		return types.ByCount(*list)
	}
}

func getSortTypeHit(s string, list *types.ApigwProfilerHitSet) sort.Interface {
	switch s {
	case "time_start":
		return types.BySTime(*list)
	case "time_finish":
		return types.ByFTime(*list)
	case "time_duration":
		return types.ByDuration(*list)
	default:
		return types.BySTime(*list)
	}
}

func encodeRoutePath(p string) string {
	return base64.URLEncoding.EncodeToString([]byte(p))
}

func decodeRoutePath(p string) (s string, err error) {
	b, err := base64.URLEncoding.DecodeString(p)
	s = string(b)

	return
}
