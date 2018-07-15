package rest

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `websocket.go`, `websocket.util.go` or `websocket_test.go` to
	implement your API calls, helper functions and tests. The file `websocket.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (wh *WebsocketHandlers) Client(w http.ResponseWriter, r *http.Request) {
	params := WebsocketClientRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return wh.Websocket.Client(params) })
}
