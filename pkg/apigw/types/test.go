package types

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	st "github.com/cortezaproject/corteza-server/system/types"
)

type (
	MockExecer struct {
		Exec_ func(context.Context, *Scp) (err error)
		Type_ func() FilterKind
	}

	MockErrorExecer struct {
		Exec_ func(context.Context, *Scp, error)
	}

	MockHandler struct {
		Foo string `json:"foo"`
	}

	MockStorer struct {
		F func(context.Context, st.ApigwFilterFilter) (st.ApigwFilterSet, st.ApigwFilterFilter, error)
		R func(context.Context, st.ApigwRouteFilter) (st.ApigwRouteSet, st.ApigwRouteFilter, error)
	}

	MockRoundTripper func(*http.Request) (*http.Response, error)
)

func (h MockHandler) String() string {
	return "MockHandler"
}

func (h MockHandler) Type() FilterKind {
	return PreFilter
}

func (h MockHandler) Weight() int {
	return 0
}

func (h MockHandler) Exec(_ context.Context, _ *Scp) error {
	panic("not implemented") // TODO: Implement
}

func (h MockHandler) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h)
	return h, err
}

func (h MockHandler) Meta() FilterMeta {
	return FilterMeta{
		Name: "return mocked filter metadata",
	}
}

func (td MockStorer) SearchApigwRoutes(ctx context.Context, f st.ApigwRouteFilter) (s st.ApigwRouteSet, ff st.ApigwRouteFilter, err error) {
	return td.R(ctx, f)
}

func (td MockStorer) SearchApigwFilters(ctx context.Context, f st.ApigwFilterFilter) (s st.ApigwFilterSet, ff st.ApigwFilterFilter, err error) {
	return td.F(ctx, f)
}

func (me MockExecer) String() string {
	return "MockExecer"
}

func (h MockExecer) Type() FilterKind {
	return PreFilter
}

func (h MockExecer) Weight() int {
	return 0
}

func (me MockExecer) Exec(ctx context.Context, s *Scp) (err error) {
	return me.Exec_(ctx, s)
}

func (me MockErrorExecer) Exec(ctx context.Context, s *Scp, e error) {
	me.Exec_(ctx, s, e)
}

func (mrt MockRoundTripper) RoundTrip(rq *http.Request) (r *http.Response, err error) {
	return mrt(rq)
}
