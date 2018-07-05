package sam

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type ChannelHandlers struct {
	*Channel
}

func (ChannelHandlers) new() *ChannelHandlers {
	return &ChannelHandlers{
		Channel{}.new(),
	}
}

// Internal API interface
type ChannelAPI interface {
	Edit(*channelEditRequest) (interface{}, error)
	Remove(*channelRemoveRequest) (interface{}, error)
	Read(*channelReadRequest) (interface{}, error)
	Search(*channelSearchRequest) (interface{}, error)
	Archive(*channelArchiveRequest) (interface{}, error)
	Move(*channelMoveRequest) (interface{}, error)
}

// HTTP API interface
type ChannelHandlersAPI interface {
	Edit(http.ResponseWriter, *http.Request)
	Remove(http.ResponseWriter, *http.Request)
	Read(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
	Archive(http.ResponseWriter, *http.Request)
	Move(http.ResponseWriter, *http.Request)
}

// Compile time check to see if we implement the interfaces
var _ ChannelHandlersAPI = &ChannelHandlers{}
var _ ChannelAPI = &Channel{}
