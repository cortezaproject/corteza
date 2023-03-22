package types

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type (
	SinkRequest struct {
		Method string      `json:"method"`
		Path   string      `json:"path"`
		Host   string      `json:"host"`
		Header http.Header `json:"header"`
		Query  url.Values  `json:"query"`

		Username   string `json:"username"`
		Password   string `json:"password"`
		RemoteAddr string `json:"remoteAddress"`

		// Make sure to set content-type to
		// application/octet-stream or application/x-www-form-urlencoded
		// to fill up PostForm
		PostForm url.Values `json:"postForm"`

		// RawBody will be base64 encoded!
		// (might contain binary data)
		Body []byte `json:"rawBody,string"`
	}

	SinkResponse struct {
		Status int         `json:"status"`
		Header http.Header `json:"header"`
		Body   interface{} `json:"body,string"`
	}
)

func NewSinkRequest(r *http.Request, body io.Reader) (sr *SinkRequest, err error) {
	sr = &SinkRequest{
		Method:     r.Method,
		Header:     r.Header,
		Host:       r.Host,
		Path:       r.URL.Path,
		Query:      r.URL.Query(),
		RemoteAddr: r.RemoteAddr,
	}

	sr.Username, sr.Password, _ = r.BasicAuth()

	if body != nil {
		if err = r.ParseForm(); err != nil {
			return nil, err
		} else {
			sr.PostForm = r.PostForm
		}

		if sr.Body, err = ioutil.ReadAll(body); err != nil {
			return nil, err
		}
	}

	return sr, err
}
