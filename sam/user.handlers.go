package sam

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (uh *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	params := UserLoginRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return uh.User.Login(params) })
}
func (uh *UserHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := UserSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return uh.User.Search(params) })
}
