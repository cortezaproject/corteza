package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `field.go`, `field.util.go` or `field_test.go` to
	implement your API calls, helper functions and tests. The file `field.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (fh *FieldHandlers) List(w http.ResponseWriter, r *http.Request) {
	params := FieldListRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return fh.Field.List(r.Context(), params) })
}
func (fh *FieldHandlers) Type(w http.ResponseWriter, r *http.Request) {
	params := FieldTypeRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return fh.Field.Type(r.Context(), params) })
}
