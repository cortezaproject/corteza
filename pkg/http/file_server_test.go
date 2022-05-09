package http

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	//go:embed test_data/file_server
	testAssets embed.FS
)

func TestSPA(t *testing.T) {
	testAssetsSub, err := fs.Sub(testAssets, "test_data/file_server")
	if err != nil {
		t.Errorf("failed to dive into subdirectory: %v", err)
	}

	testDirectAssets := os.DirFS("test_data/file_server")

	var (
		cases = []struct {
			name string
			fs   fs.FS
			cc   []configurator
			url  string
			rsp  string
		}{
			{
				name: "no special config",
				fs:   testAssetsSub,
				url:  "index.html",
				rsp:  "index html\n",
			},

			{
				name: "empty path",
				fs:   testAssetsSub,
				cc:   []configurator{Fallbacks("index.html")},
				url:  "",
				rsp:  "index html\n",
			},
			{
				name: "slash",
				fs:   testAssetsSub,
				cc:   []configurator{Fallbacks("index.html")},
				url:  "/",
				rsp:  "index html\n",
			},

			{
				name: "prefix, no slash",
				fs:   testAssetsSub,
				cc:   []configurator{UrlPrefix("/my-spa"), Fallbacks("index.html")},
				url:  "/my-spa",
				rsp:  "index html\n",
			},

			{
				name: "sub dir with url prefix",
				fs:   testAssetsSub,
				cc:   []configurator{UrlPrefix("/my-spa"), Fallbacks("index.html")},
				url:  "/my-spa/index.html",
				rsp:  "index html\n",
			},
			{
				name: "sub dir with url prefix, none existing file",
				fs:   testAssetsSub,
				cc:   []configurator{UrlPrefix("/my-spa"), Fallbacks("index.html")},
				url:  "/not-here",
				rsp:  "index html\n",
			},
			{
				name: "sub dir with url prefix, dir",
				fs:   testAssetsSub,
				cc:   []configurator{UrlPrefix("/my-spa"), Fallbacks("index.html")},
				url:  "/my-spa/sub",
				rsp:  "index html\n",
			},
			{
				name: "sub dir with url prefix, sub index",
				fs:   testAssetsSub,
				cc:   []configurator{UrlPrefix("/my-spa"), Fallbacks("index.html")},
				url:  "/my-spa/sub/index.html",
				rsp:  "sub index html\n",
			},
			{
				name: "sub dir with url prefix, sub file",
				fs:   testAssetsSub,
				cc:   []configurator{UrlPrefix("/my-spa"), Fallbacks("index.html")},
				url:  "/my-spa/sub/test.html",
				rsp:  "sub test html\n",
			},
		}
	)

	// another set of test cases for direct fs
	total := len(cases)
	for c := 0; c < total; c++ {
		direct := cases[c]
		cases[c].name += "; embedded"
		direct.name += "; direct"
		direct.fs = testDirectAssets

		cases = append(cases, direct)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				handler http.Handler
				req     = require.New(t)
				w       = httptest.NewRecorder()
				r, err  = http.NewRequest(http.MethodGet, c.url, nil)
			)

			handler, err = FileServer(c.fs, c.cc...)

			req.NoError(err)
			handler.ServeHTTP(w, r)

			req.Equal(c.rsp, w.Body.String())
			req.Equal(http.StatusOK, w.Result().StatusCode)
		})
	}

}
