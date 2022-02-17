package profiler

// import (
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/cortezaproject/corteza-server/pkg/filter"
// 	h "github.com/cortezaproject/corteza-server/pkg/http"
// 	"github.com/cortezaproject/corteza-server/system/types"
// 	"github.com/davecgh/go-spew/spew"
// )

// func Test_Profiler(t *testing.T) {
// 	// var (
// 	// 	p = Profiler{
// 	// 		l: make(map[string]*Hit),
// 	// 	}
// 	// 	req = require.New(t)
// 	// )

// 	// // in goes the h.Request
// 	// rr, err := http.NewRequest("POST", "/foo", strings.NewReader(`foo`))

// 	// req.NoError(err)

// 	// hh, err := h.NewRequest(rr)
// 	// req.NoError(err)

// 	// // need to create an internal profiling struct to hold the request?
// 	// p.Push(hh)

// 	// spew.Dump(p.l)

// 	t.Fail()
// }

// func Test_Profiler2(t *testing.T) {
// 	// types:
// 	//  + list of hits, aggregated by endpoint (ie /parse/js)
// 	//  - list of hits for a specific endpoint
// 	//  - list of hits for a specific registered route

// 	now := time.Now()
// 	then := time.Date(2022, time.March, 1, 1, 1, 1, 0, time.UTC)
// 	later := time.Date(2022, time.March, 1, 1, 1, 1, 30, time.UTC)

// 	pp := New()
// 	hr, _ := h.NewRequest(httptest.NewRequest("POST", "/foo", strings.NewReader(`foo`)))
// 	hr2, _ := h.NewRequest(httptest.NewRequest("GET", "/sometotherpath", strings.NewReader(`foo`)))
// 	pp.Push(&Hit{R: hr, Ts: &then})
// 	pp.Push(&Hit{R: hr, Ts: &then})
// 	pp.Push(&Hit{R: hr, Ts: &now})
// 	pp.Push(&Hit{R: hr, Ts: &now})
// 	pp.Push(&Hit{R: hr2, Ts: &now})
// 	pp.Push(&Hit{R: hr2, Ts: &later})
// 	pp.Push(&Hit{R: hr2, Ts: &later})
// 	pp.Push(&Hit{R: hr2, Ts: &later})
// 	pp.Push(&Hit{R: hr2, Ts: &later})
// 	pp.Push(&Hit{R: hr2, Ts: &later})
// 	pp.Push(&Hit{R: hr2, Ts: &later})
// 	pp.Push(&Hit{R: hr2, Ts: &later})

// 	var err error
// 	f := types.ApigwProfilerFilter{}
// 	if f.Sorting, err = filter.NewSorting("count DESC"); err != nil {
// 		spew.Dump(err)
// 	}

// 	list := pp.Dump(Sort{})

// 	// list of aggregations
// 	//  - keep showing, just refresh

// 	// list := pp.Dump(Sort{Before: &later, Size: 3})
// 	// spew.Dump(list)

// 	// list = list.Filter(func(k string, v *Hit) bool {
// 	// 	if v.R.URL.Path == "/sometotherpath" {
// 	// 		return true
// 	// 	}

// 	// 	return false
// 	// })

// 	// spew.Dump("LIST", list)

// 	// for p, v := range list {
// 	// 	for _, vv := range v {
// 	// 		spew.Dump(fmt.Sprintf("Path: %s, S: %s, F: %s", p, vv.Ts, vv.Ts))
// 	// 	}
// 	// }

// 	var (
// 		r                = make(types.ApigwProfilerAggregationSet, 0)
// 		tsum, tmin, tmax time.Duration
// 		ssum, smin, smax int64
// 		i                uint64 = 1
// 	)

// 	for p, v := range list {
// 		tmin, tmax, tsum = time.Hour, 0, 0
// 		smin, smax, ssum = 0, 0, 0

// 		i = 0

// 		for _, vv := range v {
// 			var (
// 				d = vv.Tf.Sub(*vv.Ts)
// 				s = vv.R.ContentLength
// 			)

// 			if d < tmin {
// 				tmin = d
// 			}

// 			if d > tmax {
// 				tmax = d
// 			}

// 			if s < smin {
// 				smin = s
// 			}

// 			if s > smax {
// 				smax = s
// 			}

// 			tsum += d
// 			ssum += s
// 			i++
// 		}

// 		spew.Dump("TSUM", tsum.Seconds())

// 		r = append(r, &types.ApigwProfilerAggregation{
// 			Path:  p,
// 			Count: i,
// 			Tmin:  tmin,
// 			Tmax:  tmax,
// 			Tavg:  time.Duration(int64(tsum.Seconds()/float64(i))) * time.Second,
// 			Smin:  smin,
// 			Smax:  smax,
// 			Savg:  float64(ssum) / float64(i),
// 		})

// 	}

// 	spew.Dump(r)

// 	SortAggregation(&r, &f)

// 	spew.Dump(r)

// 	t.Fail()
// }
