package types

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	st "github.com/cortezaproject/corteza/server/system/types"
)

type (
	MockExecer struct {
		Exec_ func(context.Context, *Scp) (err error)
		Type_ func() FilterKind
	}

	MockErrorHandler struct {
		Handler_ ErrorHandlerFunc
	}

	MockHandler struct {
		Foo      string `json:"foo"`
		Handler_ HandlerFunc
	}

	MockStorer struct {
		F func(context.Context, st.ApigwFilterFilter) (st.ApigwFilterSet, st.ApigwFilterFilter, error)
		R func(context.Context, st.ApigwRouteFilter) (st.ApigwRouteSet, st.ApigwRouteFilter, error)
	}

	MockRoundTripper func(*http.Request) (*http.Response, error)
)

func (h MockHandler) New(cfg Config) Handler {
	return MockHandler{}
}

func (h MockHandler) Enabled() bool {
	return true
}

func (h MockHandler) String() string {
	return "MockHandler"
}

func (h MockHandler) Handler() HandlerFunc {
	return h.Handler_
}

func (h MockHandler) Merge(params []byte, cfg Config) (Handler, error) {
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

func (h MockErrorHandler) Handler() ErrorHandlerFunc {
	return h.Handler_
}

func (mrt MockRoundTripper) RoundTrip(rq *http.Request) (r *http.Response, err error) {
	return mrt(rq)
}
