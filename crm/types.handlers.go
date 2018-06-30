package crm

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (th *TypesHandlers) List(w http.ResponseWriter, r *http.Request) {
	params := typesListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Types.List(params) })
}
func (th *TypesHandlers) Type(w http.ResponseWriter, r *http.Request) {
	params := typesTypeRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return th.Types.Type(params) })
}
