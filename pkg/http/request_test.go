package http

import (
	"io"
	h "net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_requestReadMultiple(t *testing.T) {
	var req = require.New(t)
	r, _ := h.NewRequest("POST", "/foo", strings.NewReader(`foo body`))

	rs, err := NewBufferedReader(r.Body)

	req.NoError(err)
	req.Equal(`foo body`, must(io.ReadAll(rs)))
	req.Equal(`foo body`, must(io.ReadAll(rs)))
}

func Test_requestReadMultipleNoBody(t *testing.T) {
	var req = require.New(t)
	r, _ := h.NewRequest("POST", "/foo", h.NoBody)

	rs, err := NewBufferedReader(r.Body)

	req.NoError(err)
	req.Equal(``, must(io.ReadAll(rs)))
	req.Equal(``, must(io.ReadAll(rs)))
}

func must(b []byte, e error) string {
	if e != nil {
		panic(e)
	}

	return string(b)
}
