package apigw

import (
	"context"
	"errors"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	// overriding mockHandler with only
	// the merge function
	mockExistingHandler struct {
		*mockHandler
		merge func(params []byte) (Handler, error)
	}
)

func Test_serviceLoadRoutes(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
	)

	mockStorer := &mockStorer{
		r: func(c context.Context, arf types.ApigwRouteFilter) (s types.ApigwRouteSet, f types.ApigwRouteFilter, err error) {
			s = types.ApigwRouteSet{
				{ID: 1, Endpoint: "/endpoint", Method: "GET", Debug: false, Enabled: true, Group: 0},
				{ID: 2, Endpoint: "/endpoint2", Method: "POST", Debug: false, Enabled: true, Group: 0},
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

	mockStorer := &mockStorer{
		f: func(c context.Context, aff types.ApigwFunctionFilter) (s types.ApigwFunctionSet, f types.ApigwFunctionFilter, err error) {
			s = types.ApigwFunctionSet{
				{ID: 1, Route: 1},
				{ID: 2, Route: 2},
			}
			return
		},
	}

	service := &apigw{

		storer: mockStorer,
	}

	r, err := service.loadFunctions(ctx, 2)

	req.NoError(err)
	req.Len(r, 2)
}

func Test_serviceInit(t *testing.T) {
	type (
		tf struct {
			name   string
			expLen int
			st     mockStorer
			reg    *registry
		}
	)

	var (
		tcc = []tf{
			{
				name: "could not register 1 function for route",
				st: mockStorer{
					r: func(c context.Context, arf types.ApigwRouteFilter) (s types.ApigwRouteSet, f types.ApigwRouteFilter, err error) {
						s = types.ApigwRouteSet{
							{ID: 1, Endpoint: "/endpoint", Method: "GET", Debug: false, Enabled: true, Group: 0},
						}
						return
					},
					f: func(c context.Context, aff types.ApigwFunctionFilter) (s types.ApigwFunctionSet, f types.ApigwFunctionFilter, err error) {
						s = types.ApigwFunctionSet{
							{ID: 1, Route: 1, Ref: "testExistingFunction"},
							{ID: 2, Route: 1, Ref: "testNotExistingFunction"},
						}
						return
					},
				},
				reg: &registry{
					h: map[string]Handler{"testExistingFunction": &mockHandler{}},
				},
				expLen: 1,
			},
			{
				name: "successful register of 2 functions for route",
				st: mockStorer{
					r: func(c context.Context, arf types.ApigwRouteFilter) (s types.ApigwRouteSet, f types.ApigwRouteFilter, err error) {
						s = types.ApigwRouteSet{
							{ID: 1, Endpoint: "/endpoint", Method: "GET", Debug: false, Enabled: true, Group: 0},
						}
						return
					},
					f: func(c context.Context, aff types.ApigwFunctionFilter) (s types.ApigwFunctionSet, f types.ApigwFunctionFilter, err error) {
						s = types.ApigwFunctionSet{
							{ID: 1, Route: 1, Ref: "testExistingFunction"},
							{ID: 2, Route: 1, Ref: "testExistingFunction"},
						}
						return
					},
				},
				reg: &registry{
					h: map[string]Handler{"testExistingFunction": &mockHandler{}},
				},
				expLen: 2,
			},
			{
				name: "could not merge params for function",
				st: mockStorer{
					r: func(c context.Context, arf types.ApigwRouteFilter) (s types.ApigwRouteSet, f types.ApigwRouteFilter, err error) {
						s = types.ApigwRouteSet{
							{ID: 1, Endpoint: "/endpoint", Method: "GET", Debug: false, Enabled: true, Group: 0},
						}
						return
					},
					f: func(c context.Context, aff types.ApigwFunctionFilter) (s types.ApigwFunctionSet, f types.ApigwFunctionFilter, err error) {
						s = types.ApigwFunctionSet{
							{ID: 1, Route: 1, Ref: "testExistingFunction", Params: types.ApigwFuncParams{}},
						}
						return
					},
				},
				reg: &registry{
					h: map[string]Handler{
						"testExistingFunction": &mockExistingHandler{
							merge: func(params []byte) (Handler, error) {
								return nil, errors.New("testttt")
							},
						},
					},
				},
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

			service := &apigw{
				log:    zap.NewNop(),
				storer: tc.st,
				reg:    tc.reg,
			}

			rr, err := service.loadRoutes(ctx)
			req.NoError(err)

			service.Init(ctx, rr...)

			req.NotEmpty(service.routes)
			req.Len(service.routes[0].pipe.w, tc.expLen)
		})
	}

}

func (h mockExistingHandler) Merge(params []byte) (Handler, error) {
	return h.merge(params)
}
