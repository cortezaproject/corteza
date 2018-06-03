package sam

import (
	"net/http"

	"github.com/pkg/errors"
)

// Channel edit request parameters
type ChannelEditRequest struct {
	id    uint64
	name  string
	topic string
}

func (ChannelEditRequest) new() *ChannelEditRequest {
	return &ChannelEditRequest{}
}

func (c *ChannelEditRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: ChannelEditRequest.Fill")
}

var _ RequestFiller = ChannelEditRequest{}.new()

// Channel remove request parameters
type ChannelRemoveRequest struct {
	id uint64
}

func (ChannelRemoveRequest) new() *ChannelRemoveRequest {
	return &ChannelRemoveRequest{}
}

func (c *ChannelRemoveRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: ChannelRemoveRequest.Fill")
}

var _ RequestFiller = ChannelRemoveRequest{}.new()

// Channel read request parameters
type ChannelReadRequest struct {
	id uint64
}

func (ChannelReadRequest) new() *ChannelReadRequest {
	return &ChannelReadRequest{}
}

func (c *ChannelReadRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: ChannelReadRequest.Fill")
}

var _ RequestFiller = ChannelReadRequest{}.new()

// Channel search request parameters
type ChannelSearchRequest struct {
	query string
}

func (ChannelSearchRequest) new() *ChannelSearchRequest {
	return &ChannelSearchRequest{}
}

func (c *ChannelSearchRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: ChannelSearchRequest.Fill")
}

var _ RequestFiller = ChannelSearchRequest{}.new()

// Channel archive request parameters
type ChannelArchiveRequest struct {
	id uint64
}

func (ChannelArchiveRequest) new() *ChannelArchiveRequest {
	return &ChannelArchiveRequest{}
}

func (c *ChannelArchiveRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: ChannelArchiveRequest.Fill")
}

var _ RequestFiller = ChannelArchiveRequest{}.new()

// Channel move request parameters
type ChannelMoveRequest struct {
	id uint64
}

func (ChannelMoveRequest) new() *ChannelMoveRequest {
	return &ChannelMoveRequest{}
}

func (c *ChannelMoveRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: ChannelMoveRequest.Fill")
}

var _ RequestFiller = ChannelMoveRequest{}.new()

// Channel merge request parameters
type ChannelMergeRequest struct {
	destination uint64
	source      uint64
}

func (ChannelMergeRequest) new() *ChannelMergeRequest {
	return &ChannelMergeRequest{}
}

func (c *ChannelMergeRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: ChannelMergeRequest.Fill")
}

var _ RequestFiller = ChannelMergeRequest{}.new()
