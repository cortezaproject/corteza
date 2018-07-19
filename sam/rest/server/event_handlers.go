package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `event.go`, `event.util.go` or `event_test.go` to
	implement your API calls, helper functions and tests. The file `event.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (eh *EventHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := EventEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return eh.Event.Edit(r.Context(), params) })
}
func (eh *EventHandlers) Attach(w http.ResponseWriter, r *http.Request) {
	params := EventAttachRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return eh.Event.Attach(r.Context(), params) })
}
func (eh *EventHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := EventRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return eh.Event.Remove(r.Context(), params) })
}
func (eh *EventHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := EventReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return eh.Event.Read(r.Context(), params) })
}
func (eh *EventHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := EventSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return eh.Event.Search(r.Context(), params) })
}
func (eh *EventHandlers) Pin(w http.ResponseWriter, r *http.Request) {
	params := EventPinRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return eh.Event.Pin(r.Context(), params) })
}
func (eh *EventHandlers) Flag(w http.ResponseWriter, r *http.Request) {
	params := EventFlagRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return eh.Event.Flag(r.Context(), params) })
}
