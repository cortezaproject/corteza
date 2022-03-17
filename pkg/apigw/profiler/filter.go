package profiler

import (
	"encoding/base64"

	"github.com/cortezaproject/corteza-server/system/types"
)

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

func encodeRoutePath(p string) string {
	return base64.URLEncoding.EncodeToString([]byte(p))
}

func decodeRoutePath(p string) (s string, err error) {
	b, err := base64.URLEncoding.DecodeString(p)
	s = string(b)

	return
}
