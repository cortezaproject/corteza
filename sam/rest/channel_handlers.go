package rest

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

func (ch *ChannelHandlers) List(w http.ResponseWriter, r *http.Request) {
	params := ChannelListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.List(params) })
}
func (ch *ChannelHandlers) Create(w http.ResponseWriter, r *http.Request) {
	params := ChannelCreateRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Create(params) })
}
func (ch *ChannelHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := ChannelEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Edit(params) })
}
func (ch *ChannelHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := ChannelReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Read(params) })
}
func (ch *ChannelHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	params := ChannelDeleteRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Delete(params) })
}
