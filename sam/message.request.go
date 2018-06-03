package sam

import (
	"net/http"

	"github.com/pkg/errors"
)

// Message edit request parameters
type MessageEditRequest struct {
	id         uint64
	channel_id uint64
	contents   string
}

func (MessageEditRequest) new() *MessageEditRequest {
	return &MessageEditRequest{}
}

func (m *MessageEditRequest) Fill(r *http.Request) error {
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

	m.id = parseUInt64(post["id"])

	m.channel_id = parseUInt64(post["channel_id"])

	m.contents = post["contents"]
	return errors.New("Not implemented: MessageEditRequest.Fill")
}

var _ RequestFiller = MessageEditRequest{}.new()

// Message attach request parameters
type MessageAttachRequest struct {
}

func (MessageAttachRequest) new() *MessageAttachRequest {
	return &MessageAttachRequest{}
}

func (m *MessageAttachRequest) Fill(r *http.Request) error {
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
	return errors.New("Not implemented: MessageAttachRequest.Fill")
}

var _ RequestFiller = MessageAttachRequest{}.new()

// Message remove request parameters
type MessageRemoveRequest struct {
	id uint64
}

func (MessageRemoveRequest) new() *MessageRemoveRequest {
	return &MessageRemoveRequest{}
}

func (m *MessageRemoveRequest) Fill(r *http.Request) error {
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

	m.id = parseUInt64(get["id"])
	return errors.New("Not implemented: MessageRemoveRequest.Fill")
}

var _ RequestFiller = MessageRemoveRequest{}.new()

// Message read request parameters
type MessageReadRequest struct {
	channel_id uint64
}

func (MessageReadRequest) new() *MessageReadRequest {
	return &MessageReadRequest{}
}

func (m *MessageReadRequest) Fill(r *http.Request) error {
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

	m.channel_id = parseUInt64(post["channel_id"])
	return errors.New("Not implemented: MessageReadRequest.Fill")
}

var _ RequestFiller = MessageReadRequest{}.new()

// Message search request parameters
type MessageSearchRequest struct {
	query        string
	message_type string
}

func (MessageSearchRequest) new() *MessageSearchRequest {
	return &MessageSearchRequest{}
}

func (m *MessageSearchRequest) Fill(r *http.Request) error {
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

	m.query = get["query"]

	m.message_type = get["message_type"]
	return errors.New("Not implemented: MessageSearchRequest.Fill")
}

var _ RequestFiller = MessageSearchRequest{}.new()

// Message pin request parameters
type MessagePinRequest struct {
	id uint64
}

func (MessagePinRequest) new() *MessagePinRequest {
	return &MessagePinRequest{}
}

func (m *MessagePinRequest) Fill(r *http.Request) error {
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

	m.id = parseUInt64(post["id"])
	return errors.New("Not implemented: MessagePinRequest.Fill")
}

var _ RequestFiller = MessagePinRequest{}.new()

// Message flag request parameters
type MessageFlagRequest struct {
	id uint64
}

func (MessageFlagRequest) new() *MessageFlagRequest {
	return &MessageFlagRequest{}
}

func (m *MessageFlagRequest) Fill(r *http.Request) error {
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

	m.id = parseUInt64(post["id"])
	return errors.New("Not implemented: MessageFlagRequest.Fill")
}

var _ RequestFiller = MessageFlagRequest{}.new()
