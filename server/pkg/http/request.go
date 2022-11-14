package http

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type (
	Request struct {
		*http.Request
		Body io.Reader
	}

	// The BufferedReader behaves exactly like a bytes.Reader, with the exception
	// when the last block is read, it automatically rewinds the internal pointer to the start,
	// so effectively, the content can be read again without calling Seek() externally.
	BufferedReader struct {
		s        []byte
		i        int64 // current reading index
		prevRune int   // index of previous rune; or < 0
	}
)

func NewRequest(r *http.Request) (rr *Request, err error) {
	rs, err := NewBufferedReader(r.Body)

	if err != nil {
		return
	}

	rr = &Request{r, rs}
	return
}

// NewBufferedReader copies original data to the
// BufferedReader
func NewBufferedReader(rr io.Reader) (bb *BufferedReader, err error) {
	var (
		buf = &bytes.Buffer{}
	)

	bb = &BufferedReader{}

	_, err = io.Copy(buf, rr)

	if err != nil {
		return
	}

	return &BufferedReader{
		s:        buf.Bytes(),
		i:        0,
		prevRune: -1,
	}, nil
}

func (r *BufferedReader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		n = 0
		err = io.EOF
		r.Seek(0, io.SeekStart)
		return
	}
	r.prevRune = -1
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

func (r *BufferedReader) Seek(offset int64, whence int) (int64, error) {
	r.prevRune = -1
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(r.s)) + offset
	default:
		return 0, errors.New("bytes.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("bytes.Reader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
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
