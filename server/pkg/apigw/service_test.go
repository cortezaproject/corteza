package apigw

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/apigw/registry"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	mockExistingHandler struct {
		*types.MockHandler
		merge func(params []byte) (types.Handler, error)
	}

	mockRouteServicer struct {
		lr  func(context.Context, string, string) ([]*types.Route, error)
		lrs func(context.Context) ([]*types.Route, error)
	}

	mockFilterServicer struct {
		lf func(context.Context, uint64) ([]*types.RouteFilter, error)
	}
)

func Test_serviceLoadRoutes(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
	)

	routeServicer := mockRouteServicer{
		lr: func(ctx context.Context, method, endpoint string) ([]*types.Route, error) {
			return []*types.Route{
				{ID: 1, Endpoint: "/endpoint", Method: "GET"},
				{ID: 2, Endpoint: "/endpoint2", Method: "POST"},
			}, nil
		},
		lrs: func(ctx context.Context) ([]*types.Route, error) {
			return []*types.Route{
				{ID: 1, Endpoint: "/endpoint", Method: "GET"},
				{ID: 2, Endpoint: "/endpoint2", Method: "POST"},
			}, nil
		},
	}

	service := &apigw{rs: routeServicer}

	r, err := service.loadRoutes(ctx)

	req.NoError(err)
	req.Len(r, 2)
}

func Test_serviceLoadFilters(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
	)

	filterServicer := mockFilterServicer{
		lf: func(ctx context.Context, routeID uint64) (s []*types.RouteFilter, err error) {
			s = []*types.RouteFilter{
				{ID: 1, Route: 1},
				{ID: 2, Route: 2},
			}
			return
		},
	}

	service := &apigw{
		fs: filterServicer,
	}

	r, err := service.fs.LoadFilters(ctx, 2)

	req.NoError(err)
	req.Len(r, 2)
}

func Test_serviceInit(t *testing.T) {
	type (
		tf struct {
			name   string
			expLen int
			reg    map[string]types.Handler

			rs mockRouteServicer
			fs mockFilterServicer
		}
	)

	var (
		tcc = []tf{
			{
				name: "could not register 1 function for route",
				rs: mockRouteServicer{
					lrs: func(ctx context.Context) ([]*types.Route, error) {
						return []*types.Route{
							{ID: 1, Endpoint: "/endpoint", Method: "GET"},
						}, nil
					},
				},
				fs: mockFilterServicer{
					lf: func(ctx context.Context, routeID uint64) (s []*types.RouteFilter, err error) {
						s = []*types.RouteFilter{
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
				rs: mockRouteServicer{
					lrs: func(ctx context.Context) ([]*types.Route, error) {
						return []*types.Route{
							{ID: 1, Endpoint: "/endpoint", Method: "GET"},
						}, nil
					},
				},
				fs: mockFilterServicer{
					lf: func(ctx context.Context, routeID uint64) (s []*types.RouteFilter, err error) {
						s = []*types.RouteFilter{
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
				rs: mockRouteServicer{
					lrs: func(ctx context.Context) ([]*types.Route, error) {
						return []*types.Route{
							{ID: 1, Endpoint: "/endpoint", Method: "GET"},
						}, nil
					},
				},
				fs: mockFilterServicer{
					lf: func(ctx context.Context, routeID uint64) (s []*types.RouteFilter, err error) {
						s = []*types.RouteFilter{
							{ID: 1, Route: 1, Ref: "testExistingFilter", Params: types.RouteFilterParams{}},
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
				log: zap.NewNop(),
				reg: reg,

				rs: tc.rs,
				fs: tc.fs,
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
			Route: types.Route{
				Method:   "GET",
				Endpoint: "/test",
			},
		}
		r2 = &route{
			Route: types.Route{
				Method:   "POST",
				Endpoint: "/test",
			},
		}
		r3 = &route{

			Route: types.Route{
				Method:   "PUT",
				Endpoint: "/test",
			},
		}
		r4 = &route{
			Route: types.Route{
				Method:   "DELETE",
				Endpoint: "/test",
			},
		}
		r5 = &route{
			Route: types.Route{
				Method:   "GET",
				Endpoint: "/test2",
			},
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

func (mrs mockRouteServicer) LoadRoute(ctx context.Context, method, endpoint string) ([]*types.Route, error) {
	return mrs.lr(ctx, method, endpoint)
}

func (mrs mockRouteServicer) LoadRoutes(ctx context.Context) ([]*types.Route, error) {
	return mrs.lrs(ctx)
}

func (mfs mockFilterServicer) LoadFilters(ctx context.Context, route uint64) ([]*types.RouteFilter, error) {
	return mfs.lf(ctx, route)
}
