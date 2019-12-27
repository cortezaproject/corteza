package types

import (
	"io/ioutil"
	"net/http"
)

type (
	SinkRequest struct {
		RequestURL string
		Header     http.Header

		// RawBody will be base64 encoded!
		// (might contain binary data)
		RawBody []byte `json:"rawBody,string"`
	}
)

func NewSinkRequest(r *http.Request) (sr *SinkRequest, err error) {
	sr = &SinkRequest{
		Header:     r.Header,
		RequestURL: r.RequestURI,
	}

	if sr.RawBody, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, err
	}

	return sr, err
}
