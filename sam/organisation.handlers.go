package sam

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (oh *OrganisationHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := organisationEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Edit(params) })
}
func (oh *OrganisationHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := organisationRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Remove(params) })
}
func (oh *OrganisationHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := organisationReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Read(params) })
}
func (oh *OrganisationHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := organisationSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Search(params) })
}
func (oh *OrganisationHandlers) Archive(w http.ResponseWriter, r *http.Request) {
	params := organisationArchiveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Archive(params) })
}
