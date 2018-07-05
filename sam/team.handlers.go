package sam

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (th *TeamHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := teamEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Edit(params) })
}
func (th *TeamHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := teamRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Remove(params) })
}
func (th *TeamHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := teamReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Read(params) })
}
func (th *TeamHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := teamSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Search(params) })
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
