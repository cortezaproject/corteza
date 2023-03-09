package types

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
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

func (h MockErrorHandler) Handler() ErrorHandlerFunc {
	return h.Handler_
}

func (mrt MockRoundTripper) RoundTrip(rq *http.Request) (r *http.Response, err error) {
	return mrt(rq)
}
