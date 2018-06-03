package sam

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (ch *ChannelHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := ChannelEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Edit(params) })
}
func (ch *ChannelHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := ChannelRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Remove(params) })
}
func (ch *ChannelHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := ChannelReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Read(params) })
}
func (ch *ChannelHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := ChannelSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Search(params) })
}
func (ch *ChannelHandlers) Archive(w http.ResponseWriter, r *http.Request) {
	params := ChannelArchiveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Archive(params) })
}
func (ch *ChannelHandlers) Move(w http.ResponseWriter, r *http.Request) {
	params := ChannelMoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Move(params) })
}
func (ch *ChannelHandlers) Merge(w http.ResponseWriter, r *http.Request) {
	params := ChannelMergeRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return ch.Channel.Merge(params) })
}
