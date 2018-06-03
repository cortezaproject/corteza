package sam

import (
	"net/http"

	"github.com/pkg/errors"
)

// Team edit request parameters
type TeamEditRequest struct {
	id      uint64
	name    string
	members []uint64
}

func (TeamEditRequest) new() *TeamEditRequest {
	return &TeamEditRequest{}
}

func (t *TeamEditRequest) Fill(r *http.Request) error {
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

	t.id = parseUInt64(post["id"])

	t.name = post["name"]
	return errors.New("Not implemented: TeamEditRequest.Fill")
}

var _ RequestFiller = TeamEditRequest{}.new()

// Team remove request parameters
type TeamRemoveRequest struct {
	id uint64
}

func (TeamRemoveRequest) new() *TeamRemoveRequest {
	return &TeamRemoveRequest{}
}

func (t *TeamRemoveRequest) Fill(r *http.Request) error {
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

	t.id = parseUInt64(get["id"])
	return errors.New("Not implemented: TeamRemoveRequest.Fill")
}

var _ RequestFiller = TeamRemoveRequest{}.new()

// Team read request parameters
type TeamReadRequest struct {
	id uint64
}

func (TeamReadRequest) new() *TeamReadRequest {
	return &TeamReadRequest{}
}

func (t *TeamReadRequest) Fill(r *http.Request) error {
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

	t.id = parseUInt64(get["id"])
	return errors.New("Not implemented: TeamReadRequest.Fill")
}

var _ RequestFiller = TeamReadRequest{}.new()

// Team search request parameters
type TeamSearchRequest struct {
	query string
}

func (TeamSearchRequest) new() *TeamSearchRequest {
	return &TeamSearchRequest{}
}

func (t *TeamSearchRequest) Fill(r *http.Request) error {
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

	t.query = get["query"]
	return errors.New("Not implemented: TeamSearchRequest.Fill")
}

var _ RequestFiller = TeamSearchRequest{}.new()

// Team archive request parameters
type TeamArchiveRequest struct {
	id uint64
}

func (TeamArchiveRequest) new() *TeamArchiveRequest {
	return &TeamArchiveRequest{}
}

func (t *TeamArchiveRequest) Fill(r *http.Request) error {
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

	t.id = parseUInt64(post["id"])
	return errors.New("Not implemented: TeamArchiveRequest.Fill")
}

var _ RequestFiller = TeamArchiveRequest{}.new()

// Team move request parameters
type TeamMoveRequest struct {
	id              uint64
	organisation_id uint64
}

func (TeamMoveRequest) new() *TeamMoveRequest {
	return &TeamMoveRequest{}
}

func (t *TeamMoveRequest) Fill(r *http.Request) error {
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

	t.id = parseUInt64(post["id"])

	t.organisation_id = parseUInt64(post["organisation_id"])
	return errors.New("Not implemented: TeamMoveRequest.Fill")
}

var _ RequestFiller = TeamMoveRequest{}.new()
