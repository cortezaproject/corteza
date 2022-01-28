package http

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type (
	Request struct {
		*http.Request
		Body io.Reader
	}

	BufferedReader struct {
		buffer []byte
	}
)

// NewBufferedReader creates a new reader from readcloser
func NewBufferedReader(r io.ReadCloser) (b *BufferedReader, err error) {
	var bb []byte

	if bb, err = io.ReadAll(r); err != nil {
		return
	}

	b = &BufferedReader{bb}
	return
}

// NewRequest creates a new Request with the buffered ready body
func NewRequest(r *http.Request) (rr *Request, err error) {
	rs, err := NewBufferedReader(r.Body)

	if err != nil {
		return
	}

	rr = &Request{r, rs}
	return
}

func (bb *BufferedReader) Read(p []byte) (n int, err error) {
	if len(bb.buffer) <= n {
		err = io.EOF
		return
	}

	if c := cap(p); c > 0 {
		for n < c {
			if len(bb.buffer) <= n {
				err = io.EOF
				break
			}

			p[n] = bb.buffer[n]
			n++
		}
	}

	return
}

func (bb *Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Method        string
		URL           *url.URL
		Header        http.Header
		ContentLength int64
		Host          string
		Form          url.Values
		PostForm      url.Values
		MultipartForm *multipart.Form
		RemoteAddr    string
		RequestURI    string
	}{
		Method:        bb.Method,
		URL:           bb.URL,
		Header:        bb.Header,
		ContentLength: bb.ContentLength,
		Host:          bb.Host,
		Form:          bb.Form,
		PostForm:      bb.PostForm,
		MultipartForm: bb.MultipartForm,
		RemoteAddr:    bb.RemoteAddr,
		RequestURI:    bb.RequestURI,
	})
}
