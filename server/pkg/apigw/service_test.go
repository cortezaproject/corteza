package apigw

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/apigw/registry"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	st "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	// overriding types.MockHandler with only
	// the merge function
	mockExistingHandler struct {
		*types.MockHandler
		merge func(params []byte) (types.Handler, error)
	}
)

func Test_serviceLoadRoutes(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
	)

	mockStorer := &types.MockStorer{
		R: func(c context.Context, arf st.ApigwRouteFilter) (s st.ApigwRouteSet, f st.ApigwRouteFilter, err error) {
			s = st.ApigwRouteSet{
				{ID: 1, Endpoint: "/endpoint", Method: "GET", Enabled: true, Group: 0},
				{ID: 2, Endpoint: "/endpoint2", Method: "POST", Enabled: true, Group: 0},
			}
			return
		},
	}

	service := &apigw{
		storer: mockStorer,
	}

	r, err := service.loadRoutes(ctx)

	req.NoError(err)
	req.Len(r, 2)
}

func Test_serviceLoadFunctions(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
	)

	mockStorer := &types.MockStorer{
		F: func(c context.Context, aff st.ApigwFilterFilter) (s st.ApigwFilterSet, f st.ApigwFilterFilter, err error) {
			s = st.ApigwFilterSet{
				{ID: 1, Route: 1},
				{ID: 2, Route: 2},
			}
			return
		},
	}

	service := &apigw{
		storer: mockStorer,
	}

	r, err := service.loadFilters(ctx, 2)

	req.NoError(err)
	req.Len(r, 2)
}

func Test_serviceInit(t *testing.T) {
	type (
		tf struct {
			name   string
			expLen int
			st     types.MockStorer
			reg    map[string]types.Handler
		}
	)

	var (
		tcc = []tf{
			{
				name: "could not register 1 function for route",
				st: types.MockStorer{
					R: func(c context.Context, arf st.ApigwRouteFilter) (s st.ApigwRouteSet, f st.ApigwRouteFilter, err error) {
						s = st.ApigwRouteSet{
							{ID: 1, Endpoint: "/endpoint", Method: "GET", Enabled: true, Group: 0},
						}
						return
					},
					F: func(c context.Context, aff st.ApigwFilterFilter) (s st.ApigwFilterSet, f st.ApigwFilterFilter, err error) {
						s = st.ApigwFilterSet{
							{ID: 1, Route: 1, Ref: "testExistingFilter"},
							{ID: 2, Route: 1, Ref: "testNotExistingFunction"},
						}
						return
					},
				},
				reg:    map[string]types.Handler{"testExistingFilter": &types.MockHandler{}},
				expLen: 1,
			},
			{
				name: "successful register of 2 functions for route",
				st: types.MockStorer{
					R: func(c context.Context, arf st.ApigwRouteFilter) (s st.ApigwRouteSet, f st.ApigwRouteFilter, err error) {
						s = st.ApigwRouteSet{
							{ID: 1, Endpoint: "/endpoint", Method: "GET", Enabled: true, Group: 0},
						}
						return
					},
					F: func(c context.Context, aff st.ApigwFilterFilter) (s st.ApigwFilterSet, f st.ApigwFilterFilter, err error) {
						s = st.ApigwFilterSet{
							{ID: 1, Route: 1, Ref: "testExistingFilter"},
							{ID: 2, Route: 1, Ref: "testExistingFilter"},
						}
						return
					},
				},
				reg:    map[string]types.Handler{"testExistingFilter": &types.MockHandler{}},
				expLen: 2,
			},
			{
				name: "could not merge params for function",
				st: types.MockStorer{
					R: func(c context.Context, arf st.ApigwRouteFilter) (s st.ApigwRouteSet, f st.ApigwRouteFilter, err error) {
						s = st.ApigwRouteSet{
							{ID: 1, Endpoint: "/endpoint", Method: "GET", Enabled: true, Group: 0},
						}
						return
					},
					F: func(c context.Context, aff st.ApigwFilterFilter) (s st.ApigwFilterSet, f st.ApigwFilterFilter, err error) {
						s = st.ApigwFilterSet{
							{ID: 1, Route: 1, Ref: "testExistingFilter", Params: st.ApigwFilterParams{}},
						}
						return
					},
				},
				reg: map[string]types.Handler{"testExistingFilter": &mockExistingHandler{
					MockHandler: &types.MockHandler{},
					merge: func(params []byte) (types.Handler, error) {
						return nil, fmt.Errorf("testttt")
					},
				}},
				expLen: 0,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
				ctx = context.Background()
			)

			reg := registry.NewRegistry(types.Config{})

			for hn, h := range tc.reg {
				reg.Add(hn, h)
			}

			service := &apigw{
				log:    zap.NewNop(),
				storer: tc.st,
				reg:    reg,
			}

			rr, err := service.loadRoutes(ctx)
			req.NoError(err)

			service.Init(ctx, rr...)

			req.NotEmpty(service.routes)
		})
	}

}

func (h mockExistingHandler) Merge(params []byte, cfg types.Config) (types.Handler, error) {
	return h.merge(params)
}

func Test_serviceAppendRoutes(t *testing.T) {
	var (
		req = require.New(t)

		routes = func(rr ...*route) []*route {
			return rr
		}

		r1 = &route{
			method:   "GET",
			endpoint: "/test",
		}
		r2 = &route{
			method:   "POST",
			endpoint: "/test",
		}
		r3 = &route{
			method:   "PUT",
			endpoint: "/test",
		}
		r4 = &route{
			method:   "DELETE",
			endpoint: "/test",
		}
		r5 = &route{
			endpoint: "GET",
			method:   "/test2",
		}

		tests = []struct {
			name     string
			ag       apigw
			routes   []*route
			expected []*route
		}{
			{
				name: "add new",
				ag: apigw{
					routes: routes(r1, r2, r3),
				},
				routes:   routes(r4),
				expected: routes(r1, r2, r3, r4),
			},
			{
				name: "no duplicates",
				ag: apigw{
					routes: routes(r1, r2, r3, r4, r5),
				},
				routes:   routes(r1, r5),
				expected: routes(r1, r2, r3, r4, r5),
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ag.AppendRoutes(tt.routes...)
			req.Equal(tt.expected, tt.ag.routes)
		})
	}
}
