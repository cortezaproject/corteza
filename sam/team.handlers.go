package sam

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (th *TeamHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := TeamEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Edit(params) })
}
func (th *TeamHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := TeamRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Remove(params) })
}
func (th *TeamHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := TeamReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Read(params) })
}
func (th *TeamHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := TeamSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Search(params) })
}
func (th *TeamHandlers) Archive(w http.ResponseWriter, r *http.Request) {
	params := TeamArchiveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Archive(params) })
}
func (th *TeamHandlers) Move(w http.ResponseWriter, r *http.Request) {
	params := TeamMoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Team.Move(params) })
}
