package types

import (
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	h "github.com/cortezaproject/corteza/server/pkg/http"
)

type (
	ApigwProfilerHit struct {
		ID string `json:"ID"`

		Body    []byte    `json:"body"`
		Request h.Request `json:"request"`
		Route   uint64    `json:"route,string"`
		Status  int       `json:"http_status_code,string"`

		Ts *time.Time     `json:"time_start"`
		Tf *time.Time     `json:"time_finish"`
		D  *time.Duration `json:"-"`
		Dr float64        `json:"time_duration"`
	}

	ApigwProfilerAggregation struct {
		Path  string  `json:"path"`
		Count uint64  `json:"count"`
		Smin  int64   `json:"size_min"`
		Smax  int64   `json:"size_max"`
		Savg  float64 `json:"size_avg"`
		Tmin  float64 `json:"time_min"`
		Tmax  float64 `json:"time_max"`
		Tavg  float64 `json:"time_avg"`
	}

	ApigwProfilerFilter struct {
		Hit    string `json:"hit,omitempty"`
		Path   string `json:"path,omitempty"`
		Before string `json:"before,omitempty"`
		Next   string `json:"next,omitempty"`

		filter.Sorting
		filter.Paging
	}
)

// sorting methods
type (
	ByPath    ApigwProfilerAggregationSet
	ByCount   ApigwProfilerAggregationSet
	BySizeMin ApigwProfilerAggregationSet
	BySizeMax ApigwProfilerAggregationSet
	BySizeAvg ApigwProfilerAggregationSet
	ByTimeMin ApigwProfilerAggregationSet
	ByTimeMax ApigwProfilerAggregationSet
	ByTimeAvg ApigwProfilerAggregationSet

	BySTime         ApigwProfilerHitSet
	ByFTime         ApigwProfilerHitSet
	ByDuration      ApigwProfilerHitSet
	ByContentLength ApigwProfilerHitSet
	ByStatus        ApigwProfilerHitSet
	ByMethod        ApigwProfilerHitSet
)

func (h ByPath) Len() int {
	return len(h)
}

func (h ByPath) Less(i, j int) bool {
	return h[i].Path < h[j].Path
}

func (h ByPath) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h ByCount) Len() int {
	return len(h)
}

func (h ByCount) Less(i, j int) bool {
	// make sure to sort via path also
	// if the counts are the same
	if h[i].Count == h[j].Count {
		return h[i].Path < h[j].Path
	}

	return h[i].Count < h[j].Count
}

func (h ByCount) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h BySizeMin) Len() int {
	return len(h)
}

func (h BySizeMin) Less(i, j int) bool {
	// make sure to sort via path also
	// if the counts are the same
	if h[i].Smin == h[j].Smin {
		return h[i].Path < h[j].Path
	}

	return h[i].Smin < h[j].Smin
}

func (h BySizeMin) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h BySizeMax) Len() int {
	return len(h)
}

func (h BySizeMax) Less(i, j int) bool {
	// make sure to sort via path also
	// if the counts are the same
	if h[i].Smax == h[j].Smax {
		return h[i].Path < h[j].Path
	}

	return h[i].Smax < h[j].Smax
}

func (h BySizeMax) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h BySizeAvg) Len() int {
	return len(h)
}

func (h BySizeAvg) Less(i, j int) bool {
	// make sure to sort via path also
	// if the counts are the same
	if h[i].Savg == h[j].Savg {
		return h[i].Path < h[j].Path
	}

	return h[i].Savg < h[j].Savg
}

func (h BySizeAvg) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h ByTimeAvg) Len() int {
	return len(h)
}

func (h ByTimeAvg) Less(i, j int) bool {
	// make sure to sort via path also
	// if the counts are the same
	if h[i].Tavg == h[j].Tavg {
		return h[i].Path < h[j].Path
	}

	return h[i].Tavg < h[j].Tavg
}

func (h ByTimeAvg) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h ByTimeMax) Len() int {
	return len(h)
}

func (h ByTimeMax) Less(i, j int) bool {
	// make sure to sort via path also
	// if the counts are the same
	if h[i].Tmax == h[j].Tmax {
		return h[i].Path < h[j].Path
	}

	return h[i].Tmax < h[j].Tmax
}

func (h ByTimeMax) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h ByTimeMin) Len() int {
	return len(h)
}

func (h ByTimeMin) Less(i, j int) bool {
	// make sure to sort via path also
	// if the counts are the same
	if h[i].Tmin == h[j].Tmin {
		return h[i].Path < h[j].Path
	}

	return h[i].Tmin < h[j].Tmin
}

func (h ByTimeMin) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

//
// Sorting hits
//
func (h BySTime) Len() int {
	return len(h)
}

func (h BySTime) Less(i, j int) bool {
	return h[j].Ts.After(*h[i].Ts)
}

func (h BySTime) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h ByFTime) Len() int {
	return len(h)
}

func (h ByFTime) Less(i, j int) bool {
	return h[j].Tf.After(*h[i].Tf)
}

func (h ByFTime) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h ByDuration) Len() int {
	return len(h)
}

func (h ByDuration) Less(i, j int) bool {
	return h[i].D.Microseconds() > h[j].D.Microseconds()
}

func (h ByDuration) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h ByContentLength) Len() int {
	return len(h)
}

func (h ByContentLength) Less(i, j int) bool {
	// there could be many requests with equal amount of data
	// so for paging we need a secondary sort
	if h[i].Request.ContentLength == h[j].Request.ContentLength {
		return h[j].Ts.After(*h[i].Ts)
	}
	return h[i].Request.ContentLength < h[j].Request.ContentLength
}

func (h ByContentLength) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h ByStatus) Len() int {
	return len(h)
}

func (h ByStatus) Less(i, j int) bool {
	// most will have the same status
	// so for paging we need a secondary sort
	if h[i].Status == h[j].Status {
		return h[j].Ts.After(*h[i].Ts)
	}
	return h[i].Status < h[j].Status
}

func (h ByStatus) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

func (h ByMethod) Len() int {
	return len(h)
}

func (h ByMethod) Less(i, j int) bool {
	// most will have the same status
	// so for paging we need a secondary sort
	if h[i].Request.Method == h[j].Request.Method {
		return h[j].Ts.After(*h[i].Ts)
	}
	return h[i].Request.Method < h[j].Request.Method
}

func (h ByMethod) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}
