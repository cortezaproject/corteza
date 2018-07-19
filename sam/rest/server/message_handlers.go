package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `message.go`, `message.util.go` or `message_test.go` to
	implement your API calls, helper functions and tests. The file `message.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (mh *MessageHandlers) Create(w http.ResponseWriter, r *http.Request) {
	params := MessageCreateRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Create(r.Context(), params) })
}
func (mh *MessageHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := MessageEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Edit(r.Context(), params) })
}
func (mh *MessageHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	params := MessageDeleteRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Delete(r.Context(), params) })
}
func (mh *MessageHandlers) Attach(w http.ResponseWriter, r *http.Request) {
	params := MessageAttachRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Attach(r.Context(), params) })
}
func (mh *MessageHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := MessageSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Search(r.Context(), params) })
}
func (mh *MessageHandlers) Pin(w http.ResponseWriter, r *http.Request) {
	params := MessagePinRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Pin(r.Context(), params) })
}
func (mh *MessageHandlers) Unpin(w http.ResponseWriter, r *http.Request) {
	params := MessageUnpinRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Unpin(r.Context(), params) })
}
func (mh *MessageHandlers) Flag(w http.ResponseWriter, r *http.Request) {
	params := MessageFlagRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Flag(r.Context(), params) })
}
func (mh *MessageHandlers) Deflag(w http.ResponseWriter, r *http.Request) {
	params := MessageDeflagRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Deflag(r.Context(), params) })
}
func (mh *MessageHandlers) React(w http.ResponseWriter, r *http.Request) {
	params := MessageReactRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.React(r.Context(), params) })
}
func (mh *MessageHandlers) Unreact(w http.ResponseWriter, r *http.Request) {
	params := MessageUnreactRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Unreact(r.Context(), params) })
}
