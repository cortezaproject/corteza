package service

import (
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"math"
	"sort"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/apigw"
	"github.com/cortezaproject/corteza-server/pkg/apigw/profiler"

	"github.com/cortezaproject/corteza-server/system/types"
)

var (
	sortAggFields   = []string{"path", "count", "size_min", "size_max", "size_avg", "time_min", "time_max", "time_avg"}
	sortRouteFields = []string{"time_start", "time_finish", "time_duration", "content_length", "http_status_code", "http_method"}
)

const (
	// default fallback on amount of items
	FILTER_NUM_ITEMS = 20

	// default fallback on amount of aggregated items
	FILTER_NUM_AGG_ITEMS = 10
)

type (
	apigwProfiler struct{}
)

func Profiler() *apigwProfiler {
	return &apigwProfiler{}
}

// HitsAggregated fetches a list of hits from integration gateway profiler
func (svc *apigwProfiler) Hits(ctx context.Context, filter types.ApigwProfilerFilter) (r types.ApigwProfilerHitSet, f types.ApigwProfilerFilter, err error) {

	f = filter
	r = make(types.ApigwProfilerHitSet, 0)

	uDec, err := base64.URLEncoding.DecodeString(filter.Path)

	if err != nil {
		return
	}

	filter.Path = string(uDec)

	if filter.Path == "" && filter.Hit == "" {
		err = errors.New("fetching all hits (no route and hit specified) not supported")
		return
	}

	var sorting = profiler.Sort{
		Hit:  filter.Hit,
		Path: filter.Path,
	}

	var (
		list = apigw.Service().Profiler().Hits(sorting)
	)

	var pp = ""

	for k, _ := range list {
		if filter.Hit != "" {
			pp = k
			break
		}

		if filter.Path != "" && k == filter.Path {
			pp = k
			break
		}
	}

	if pp == "" {
		return
	}

	for _, h := range list[pp] {
		hh := &types.ApigwProfilerHit{
			ID: h.ID,

			Route:   h.Route,
			Status:  h.Status,
			Request: *h.R,

			Ts: h.Ts,
			Tf: h.Tf,
			D:  h.D,
			Dr: float64(h.D.Microseconds()) / 1000,
		}

		// fetch body only on hit details
		if filter.Hit != "" {
			hh.Body, _ = ioutil.ReadAll(hh.Request.Body)
		}

		r = append(r, hh)
	}

	// sort
	sortHits(&r, &f)

	// filter sorted
	filterHits(&r, &f)

	return
}

// HitsAggregated fetches a list of hits from integration gateway profiler
// and aggregates them with assigned filters
func (svc *apigwProfiler) HitsAggregated(ctx context.Context, filter types.ApigwProfilerFilter) (r types.ApigwProfilerAggregationSet, f types.ApigwProfilerFilter, err error) {
	f = filter
	r = make(types.ApigwProfilerAggregationSet, 0)

	uDec, err := base64.URLEncoding.DecodeString(filter.Path)

	if err != nil {
		return
	}

	filter.Path = string(uDec)

	var (
		list = apigw.Service().Profiler().Hits(profiler.Sort{
			Path: filter.Path,
		})

		tsum, tmin, tmax time.Duration
		ssum, smin, smax int64
		i                uint64 = 1
	)

	for p, v := range list {
		tmin, tmax, tsum = time.Hour, 0, 0
		smin, smax, ssum = math.MaxInt64, 0, 0

		i = 0

		for _, vv := range v {
			var (
				d = *vv.D
				s = vv.R.ContentLength
			)

			if d < tmin {
				tmin = d
			}

			if d > tmax {
				tmax = d
			}

			if s < smin {
				smin = s
			}

			if s > smax {
				smax = s
			}

			tsum += d
			ssum += s
			i++
		}

		r = append(r, &types.ApigwProfilerAggregation{
			Path:  p,
			Count: i,
			Tmin:  float64(tmin.Microseconds()) / 1000,
			Tmax:  float64(tmax.Microseconds()) / 1000,
			Tavg:  float64(tsum.Microseconds()) / float64(i) / 1000,
			Smin:  smin,
			Smax:  smax,
			Savg:  float64(ssum) / float64(i),
		})
	}

	// sort first, by primary and secondary fields
	// primary is the one chosen
	// secondary is path
	sortAggregation(&r, &f)

	// filter
	filterAggregation(&r, &f)

	return
}

func sortAggregation(list *types.ApigwProfilerAggregationSet, filter *types.ApigwProfilerFilter) {
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

func sortHits(list *types.ApigwProfilerHitSet, filter *types.ApigwProfilerFilter) {
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

func filterAggregation(list *types.ApigwProfilerAggregationSet, filter *types.ApigwProfilerFilter) {
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

func filterHits(list *types.ApigwProfilerHitSet, filter *types.ApigwProfilerFilter) {
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

func encodeRoutePath(p string) string {
	return base64.URLEncoding.EncodeToString([]byte(p))
}

func decodeRoutePath(p string) (s string, err error) {
	b, err := base64.URLEncoding.DecodeString(p)
	s = string(b)

	return
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
	case "content_length":
		return types.ByContentLength(*list)
	case "http_status_code":
		return types.ByStatus(*list)
	case "http_method":
		return types.ByMethod(*list)
	default:
		return types.BySTime(*list)
	}
}
