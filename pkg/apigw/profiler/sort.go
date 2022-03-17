package profiler

import (
	"sort"

	"github.com/cortezaproject/corteza-server/system/types"
)

var (
	sortAggFields   = []string{"path", "count", "size_min", "size_max", "size_avg", "time_min", "time_max", "time_avg"}
	sortRouteFields = []string{"time_start", "time_finish", "time_duration"}
)

type (
	Sort struct {
		Hit    string
		Path   string
		Size   uint64
		Before string
	}
)

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
