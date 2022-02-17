package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	h "github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/pkg/errors"
)

type (
	ApigwRoute struct {
		ID       uint64         `json:"routeID,string"`
		Endpoint string         `json:"endpoint"`
		Method   string         `json:"method"`
		Enabled  bool           `json:"enabled"`
		Group    uint64         `json:"group,string"`
		Meta     ApigwRouteMeta `json:"meta"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

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

	ApigwRouteMeta struct {
		Debug bool `json:"debug"`
		Async bool `json:"async"`
	}

	ApigwRouteFilter struct {
		Route string `json:"route"`
		Group string `json:"group"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*ApigwRoute) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
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

func (cc *ApigwRouteMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*cc = ApigwRouteMeta{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, cc); err != nil {
			return errors.Wrapf(err, "cannot scan '%v' into ApigwRouteMeta", string(b))
		}
	}

	return nil
}

func (cc ApigwRouteMeta) Value() (driver.Value, error) {
	return json.Marshal(cc)
}

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

	BySTime    ApigwProfilerHitSet
	ByFTime    ApigwProfilerHitSet
	ByDuration ApigwProfilerHitSet
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
	return h[i].Tmin < h[j].Tmin
}

func (h ByTimeMin) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	return
}

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
