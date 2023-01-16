package profiler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	h "github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/stretchr/testify/require"
)

const (
	day = time.Hour * 24
)

func Test_ApigwProfiler_newHit(t *testing.T) {
	var (
		p   = New()
		req = require.New(t)

		rr, err = h.NewRequest(httptest.NewRequest("POST", "/foo", strings.NewReader(`foo`)))
		hit     = p.Hit(rr)
	)

	req.NoError(err)
	req.Equal(hit.Status, http.StatusOK)
	req.NotNil(hit.Ts)
	req.NotNil(hit.ID)
}

func Test_ApigwProfiler_push(t *testing.T) {
	var (
		p       = New()
		req     = require.New(t)
		rr, err = h.NewRequest(httptest.NewRequest("POST", "/foo", strings.NewReader(`foo`)))

		hit = p.Hit(rr)
		id  = p.Push(hit)
	)

	_, found := p.l[id]

	req.NoError(err)
	req.NotEmpty(id)
	req.Len(p.l, 1)
	req.True(found)
}

func Test_ApigwProfiler_filterPath(t *testing.T) {
	var (
		pp  = New()
		req = require.New(t)

		now   = time.Date(2022, time.March, 1, 1, 1, 1, 0, time.UTC)
		then  = now.Add(-1 * day)
		later = now.Add(day)

		hr, _  = h.NewRequest(httptest.NewRequest("POST", "/foo", strings.NewReader(`foo`)))
		hr2, _ = h.NewRequest(httptest.NewRequest("GET", "/bar", strings.NewReader(`foo`)))
		hr3, _ = h.NewRequest(httptest.NewRequest("GET", "/baz", strings.NewReader(`foo`)))
	)

	pp.Push(&Hit{R: hr, Ts: &then})
	pp.Push(&Hit{R: hr, Ts: &then})
	pp.Push(&Hit{R: hr2, Ts: &now})
	pp.Push(&Hit{R: hr3, Ts: &later})

	list := pp.Hits(Sort{
		Path: hr3.RequestURI,
	})

	_, found := list[hr3.RequestURI]

	req.True(found)
	req.Len(list[hr3.RequestURI], 1)
}

func Test_ApigwProfiler_filterHit(t *testing.T) {
	var (
		pp  = New()
		req = require.New(t)

		now   = time.Date(2022, time.March, 1, 1, 1, 1, 0, time.UTC)
		then  = now.Add(-1 * day)
		later = now.Add(day)

		hr, _  = h.NewRequest(httptest.NewRequest("POST", "/foo", strings.NewReader(`foo`)))
		hr2, _ = h.NewRequest(httptest.NewRequest("GET", "/bar", strings.NewReader(`foo`)))
		hr3, _ = h.NewRequest(httptest.NewRequest("GET", "/baz", strings.NewReader(`foo`)))
	)

	pp.Push(&Hit{R: hr, Ts: &then})
	pp.Push(&Hit{R: hr, Ts: &now})
	pp.Push(&Hit{R: hr2, Ts: &now})
	pp.Push(&Hit{R: hr2, Ts: &later})
	pp.Push(&Hit{R: hr2, Ts: &later})

	h := pp.Hit(hr3)
	h.Ts = &later

	id := pp.Push(h)

	list := pp.Hits(Sort{
		Hit: h.ID,
	})

	_, found := list[id]

	req.True(found)
	req.Len(list[id], 1)
}

func Test_ApigwProfiler_purgeAll(t *testing.T) {
	var (
		pp  = New()
		req = require.New(t)

		now   = time.Date(2022, time.March, 1, 1, 1, 1, 0, time.UTC)
		then  = now.Add(-1 * day)
		later = now.Add(day)

		hr, _  = h.NewRequest(httptest.NewRequest("POST", "/foo", strings.NewReader(`foo`)))
		hr2, _ = h.NewRequest(httptest.NewRequest("GET", "/bar", strings.NewReader(`foo`)))
	)

	pp.Push(&Hit{R: hr, Ts: &then})
	pp.Push(&Hit{R: hr, Ts: &now})
	pp.Push(&Hit{R: hr2, Ts: &now})
	pp.Push(&Hit{R: hr2, Ts: &later})
	pp.Push(&Hit{R: hr2, Ts: &later})

	req.Len(pp.l, 2)
	pp.Purge(&PurgeFilter{})
	req.Len(pp.l, 0)
}

func Test_ApigwProfiler_purgeRoute(t *testing.T) {
	var (
		pp  = New()
		req = require.New(t)

		now   = time.Date(2022, time.March, 1, 1, 1, 1, 0, time.UTC)
		then  = now.Add(-1 * day)
		later = now.Add(day)

		hr, _  = h.NewRequest(httptest.NewRequest("POST", "/foo", strings.NewReader(`foo`)))
		hr2, _ = h.NewRequest(httptest.NewRequest("GET", "/bar", strings.NewReader(`foo`)))
	)

	pp.Push(&Hit{Route: 1, R: hr, Ts: &then})
	pp.Push(&Hit{Route: 1, R: hr, Ts: &now})
	pp.Push(&Hit{Route: 2, R: hr2, Ts: &now})
	pp.Push(&Hit{Route: 2, R: hr2, Ts: &later})
	pp.Push(&Hit{Route: 2, R: hr2, Ts: &later})

	req.Len(pp.l, 2)

	pp.Purge(&PurgeFilter{RouteID: 2})

	req.Len(pp.l, 1)
}
