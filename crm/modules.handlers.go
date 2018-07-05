package crm

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (mh *ModulesHandlers) List(w http.ResponseWriter, r *http.Request) {
	params := modulesListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Modules.List(params) })
}
func (mh *ModulesHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := modulesEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Modules.Edit(params) })
}
func (mh *ModulesHandlers) ContentList(w http.ResponseWriter, r *http.Request) {
	params := modulesContentListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Modules.ContentList(params) })
}
func (mh *ModulesHandlers) ContentEdit(w http.ResponseWriter, r *http.Request) {
	params := modulesContentEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Modules.ContentEdit(params) })
}
func (mh *ModulesHandlers) ContentDelete(w http.ResponseWriter, r *http.Request) {
	params := modulesContentDeleteRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Modules.ContentDelete(params) })
}
