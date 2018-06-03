package sam

import (
	"net/http"
)

// Channel edit request parameters
type channelEditRequest struct {
	id    uint64
	name  string
	topic string
}

func (channelEditRequest) new() *channelEditRequest {
	return &channelEditRequest{}
}

func (c *channelEditRequest) Fill(r *http.Request) error {
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

	c.id = parseUInt64(post["id"])

	c.name = post["name"]

	c.topic = post["topic"]
	return nil
}

var _ RequestFiller = channelEditRequest{}.new()

// Channel remove request parameters
type channelRemoveRequest struct {
	id uint64
}

func (channelRemoveRequest) new() *channelRemoveRequest {
	return &channelRemoveRequest{}
}

func (c *channelRemoveRequest) Fill(r *http.Request) error {
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

	c.id = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = channelRemoveRequest{}.new()

// Channel read request parameters
type channelReadRequest struct {
	id uint64
}

func (channelReadRequest) new() *channelReadRequest {
	return &channelReadRequest{}
}

func (c *channelReadRequest) Fill(r *http.Request) error {
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

	c.id = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = channelReadRequest{}.new()

// Channel search request parameters
type channelSearchRequest struct {
	query string
}

func (channelSearchRequest) new() *channelSearchRequest {
	return &channelSearchRequest{}
}

func (c *channelSearchRequest) Fill(r *http.Request) error {
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

	c.query = get["query"]
	return nil
}

var _ RequestFiller = channelSearchRequest{}.new()

// Channel archive request parameters
type channelArchiveRequest struct {
	id uint64
}

func (channelArchiveRequest) new() *channelArchiveRequest {
	return &channelArchiveRequest{}
}

func (c *channelArchiveRequest) Fill(r *http.Request) error {
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

	c.id = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = channelArchiveRequest{}.new()

// Channel move request parameters
type channelMoveRequest struct {
	id uint64
}

func (channelMoveRequest) new() *channelMoveRequest {
	return &channelMoveRequest{}
}

func (c *channelMoveRequest) Fill(r *http.Request) error {
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

	c.id = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = channelMoveRequest{}.new()

// Channel merge request parameters
type channelMergeRequest struct {
	destination uint64
	source      uint64
}

func (channelMergeRequest) new() *channelMergeRequest {
	return &channelMergeRequest{}
}

func (c *channelMergeRequest) Fill(r *http.Request) error {
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

	c.destination = parseUInt64(post["destination"])

	c.source = parseUInt64(post["source"])
	return nil
}

var _ RequestFiller = channelMergeRequest{}.new()
