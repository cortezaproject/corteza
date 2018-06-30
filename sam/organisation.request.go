package sam

import (
	"github.com/go-chi/chi"
	"net/http"
)

var _ = chi.URLParam

// Organisation edit request parameters
type organisationEditRequest struct {
	id   uint64
	name string
}

func (organisationEditRequest) new() *organisationEditRequest {
	return &organisationEditRequest{}
}

func (o *organisationEditRequest) Fill(r *http.Request) error {
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
	return nil
}

var _ RequestFiller = organisationEditRequest{}.new()

// Organisation remove request parameters
type organisationRemoveRequest struct {
	id uint64
}

func (organisationRemoveRequest) new() *organisationRemoveRequest {
	return &organisationRemoveRequest{}
}

func (o *organisationRemoveRequest) Fill(r *http.Request) error {
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
	return nil
}

var _ RequestFiller = organisationRemoveRequest{}.new()

// Organisation read request parameters
type organisationReadRequest struct {
	id uint64
}

func (organisationReadRequest) new() *organisationReadRequest {
	return &organisationReadRequest{}
}

func (o *organisationReadRequest) Fill(r *http.Request) error {
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
	return nil
}

var _ RequestFiller = organisationReadRequest{}.new()

// Organisation search request parameters
type organisationSearchRequest struct {
	query string
}

func (organisationSearchRequest) new() *organisationSearchRequest {
	return &organisationSearchRequest{}
}

func (o *organisationSearchRequest) Fill(r *http.Request) error {
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
	return nil
}

var _ RequestFiller = organisationSearchRequest{}.new()

// Organisation archive request parameters
type organisationArchiveRequest struct {
	id uint64
}

func (organisationArchiveRequest) new() *organisationArchiveRequest {
	return &organisationArchiveRequest{}
}

func (o *organisationArchiveRequest) Fill(r *http.Request) error {
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
	return nil
}

var _ RequestFiller = organisationArchiveRequest{}.new()
