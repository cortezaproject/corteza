package apigw

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	mockExecer struct {
		exec func(context.Context, *scp) (err error)
	}

	mockErrorExecer struct {
		exec func(context.Context, *scp, error)
	}

	mockHandler struct {
		Foo string `json:"foo"`
	}

	mockStorer struct {
		f func(context.Context, types.ApigwFunctionFilter) (types.ApigwFunctionSet, types.ApigwFunctionFilter, error)
		r func(context.Context, types.ApigwRouteFilter) (types.ApigwRouteSet, types.ApigwRouteFilter, error)
	}
)

func (h mockHandler) String() string {
	return "mockHandler"
}

func (h mockHandler) Exec(_ context.Context, _ *scp) error {
	panic("not implemented") // TODO: Implement
}

func (h mockHandler) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h)
	return h, err
}

func (h mockHandler) Meta() functionMeta {
	return functionMeta{
		Name: "return mocked function metadata",
	}
}

func (td mockStorer) SearchApigwRoutes(ctx context.Context, f types.ApigwRouteFilter) (s types.ApigwRouteSet, ff types.ApigwRouteFilter, err error) {
	return td.r(ctx, f)
}

func (td mockStorer) SearchApigwFunctions(ctx context.Context, f types.ApigwFunctionFilter) (s types.ApigwFunctionSet, ff types.ApigwFunctionFilter, err error) {
	return td.f(ctx, f)
}

func (me mockExecer) String() string {
	return "mockExecer"
}

func (me mockExecer) Exec(ctx context.Context, s *scp) (err error) {
	return me.exec(ctx, s)
}

func (me mockErrorExecer) Exec(ctx context.Context, s *scp, e error) {
	me.exec(ctx, s, e)
}
