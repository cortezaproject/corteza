package crm

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `types.go`, `types.util.go` or `types_test.go` to
	implement your API calls, helper functions and tests. The file `types.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

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
