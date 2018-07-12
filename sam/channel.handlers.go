package sam

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `channel.go`, `channel.util.go` or `channel_test.go` to
	implement your API calls, helper functions and tests. The file `channel.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (ch *ChannelHandlers) Create(w http.ResponseWriter, r *http.Request) {
	params := channelCreateRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Create(params) })
}
func (ch *ChannelHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := channelEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Edit(params) })
}
func (ch *ChannelHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := channelRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Remove(params) })
}
func (ch *ChannelHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := channelReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Read(params) })
}
func (ch *ChannelHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := channelSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Search(params) })
}
