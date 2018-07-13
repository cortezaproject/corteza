package sam

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `team.go`, `team.util.go` or `team_test.go` to
	implement your API calls, helper functions and tests. The file `team.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (th *TeamHandlers) List(w http.ResponseWriter, r *http.Request) {
	params := teamListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.List(params) })
}
func (th *TeamHandlers) Create(w http.ResponseWriter, r *http.Request) {
	params := teamCreateRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Create(params) })
}
func (th *TeamHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := teamEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Edit(params) })
}
func (th *TeamHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := teamReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Read(params) })
}
func (th *TeamHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := teamRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Remove(params) })
}
func (th *TeamHandlers) Archive(w http.ResponseWriter, r *http.Request) {
	params := teamArchiveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Archive(params) })
}
func (th *TeamHandlers) Move(w http.ResponseWriter, r *http.Request) {
	params := teamMoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Move(params) })
}
func (th *TeamHandlers) Merge(w http.ResponseWriter, r *http.Request) {
	params := teamMergeRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Merge(params) })
}
