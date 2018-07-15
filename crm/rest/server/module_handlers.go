package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (mh *ModuleHandlers) List(w http.ResponseWriter, r *http.Request) {
	params := ModuleListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.List(params) })
}
func (mh *ModuleHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := ModuleEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.Edit(params) })
}
func (mh *ModuleHandlers) ContentList(w http.ResponseWriter, r *http.Request) {
	params := ModuleContentListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.ContentList(params) })
}
func (mh *ModuleHandlers) ContentEdit(w http.ResponseWriter, r *http.Request) {
	params := ModuleContentEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.ContentEdit(params) })
}
func (mh *ModuleHandlers) ContentDelete(w http.ResponseWriter, r *http.Request) {
	params := ModuleContentDeleteRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.ContentDelete(params) })
}
