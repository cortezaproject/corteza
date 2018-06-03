package sam

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (mh *MessageHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := messageEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Edit(params) })
}
func (mh *MessageHandlers) Attach(w http.ResponseWriter, r *http.Request) {
	params := messageAttachRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Attach(params) })
}
func (mh *MessageHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := messageRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Remove(params) })
}
func (mh *MessageHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := messageReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Read(params) })
}
func (mh *MessageHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := messageSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Search(params) })
}
func (mh *MessageHandlers) Pin(w http.ResponseWriter, r *http.Request) {
	params := messagePinRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Pin(params) })
}
func (mh *MessageHandlers) Flag(w http.ResponseWriter, r *http.Request) {
	params := messageFlagRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return mh.Message.Flag(params) })
}
