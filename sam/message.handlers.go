package sam

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (mh *MessageHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := MessageEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Edit(params) })
}
func (mh *MessageHandlers) Attach(w http.ResponseWriter, r *http.Request) {
	params := MessageAttachRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Attach(params) })
}
func (mh *MessageHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := MessageRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Remove(params) })
}
func (mh *MessageHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := MessageReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Read(params) })
}
func (mh *MessageHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := MessageSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Search(params) })
}
func (mh *MessageHandlers) Pin(w http.ResponseWriter, r *http.Request) {
	params := MessagePinRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Pin(params) })
}
func (mh *MessageHandlers) Flag(w http.ResponseWriter, r *http.Request) {
	params := MessageFlagRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Flag(params) })
}
