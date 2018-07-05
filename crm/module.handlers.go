package crm

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (mh *ModuleHandlers) List(w http.ResponseWriter, r *http.Request) {
	params := moduleListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.List(params) })
}
func (mh *ModuleHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := moduleEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.Edit(params) })
}
func (mh *ModuleHandlers) ContentList(w http.ResponseWriter, r *http.Request) {
	params := moduleContentListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.ContentList(params) })
}
func (mh *ModuleHandlers) ContentEdit(w http.ResponseWriter, r *http.Request) {
	params := moduleContentEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.ContentEdit(params) })
}
func (mh *ModuleHandlers) ContentDelete(w http.ResponseWriter, r *http.Request) {
	params := moduleContentDeleteRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Module.ContentDelete(params) })
}
