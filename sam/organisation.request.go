package sam

import (
	"net/http"

	"github.com/pkg/errors"
)

// Organisation edit request parameters
type OrganisationEditRequest struct {
	id   uint64
	name string
}

func (OrganisationEditRequest) new() *OrganisationEditRequest {
	return &OrganisationEditRequest{}
}

func (o *OrganisationEditRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	o.id = parseUInt64(post["id"])

	o.name = post["name"]
	return errors.New("Not implemented: OrganisationEditRequest.Fill")
}

var _ RequestFiller = OrganisationEditRequest{}.new()

// Organisation remove request parameters
type OrganisationRemoveRequest struct {
	id uint64
}

func (OrganisationRemoveRequest) new() *OrganisationRemoveRequest {
	return &OrganisationRemoveRequest{}
}

func (o *OrganisationRemoveRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	o.id = parseUInt64(get["id"])
	return errors.New("Not implemented: OrganisationRemoveRequest.Fill")
}

var _ RequestFiller = OrganisationRemoveRequest{}.new()

// Organisation read request parameters
type OrganisationReadRequest struct {
	id uint64
}

func (OrganisationReadRequest) new() *OrganisationReadRequest {
	return &OrganisationReadRequest{}
}

func (o *OrganisationReadRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	o.id = parseUInt64(get["id"])
	return errors.New("Not implemented: OrganisationReadRequest.Fill")
}

var _ RequestFiller = OrganisationReadRequest{}.new()

// Organisation search request parameters
type OrganisationSearchRequest struct {
	query string
}

func (OrganisationSearchRequest) new() *OrganisationSearchRequest {
	return &OrganisationSearchRequest{}
}

func (o *OrganisationSearchRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	o.query = get["query"]
	return errors.New("Not implemented: OrganisationSearchRequest.Fill")
}

var _ RequestFiller = OrganisationSearchRequest{}.new()

// Organisation archive request parameters
type OrganisationArchiveRequest struct {
	id uint64
}

func (OrganisationArchiveRequest) new() *OrganisationArchiveRequest {
	return &OrganisationArchiveRequest{}
}

func (o *OrganisationArchiveRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	o.id = parseUInt64(post["id"])
	return errors.New("Not implemented: OrganisationArchiveRequest.Fill")
}

var _ RequestFiller = OrganisationArchiveRequest{}.new()
