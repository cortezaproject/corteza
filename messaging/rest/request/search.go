package request

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `search.go`, `search.util.go` or `search_test.go` to
	implement your API calls, helper functions and tests. The file `search.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"io"
	"strings"

	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var _ = chi.URLParam
var _ = multipart.FileHeader{}

// Search messages request parameters
type SearchMessages struct {
	ChannelID       []uint64 `json:",string"`
	AfterMessageID  uint64   `json:",string"`
	BeforeMessageID uint64   `json:",string"`
	FromMessageID   uint64   `json:",string"`
	ToMessageID     uint64   `json:",string"`
	ThreadID        []uint64 `json:",string"`
	UserID          []uint64 `json:",string"`
	Type            []string
	PinnedOnly      bool
	BookmarkedOnly  bool
	Limit           uint
	Query           string
}

func NewSearchMessages() *SearchMessages {
	return &SearchMessages{}
}

func (sReq *SearchMessages) Fill(r *http.Request) (err error) {
	if strings.ToLower(r.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(r.Body).Decode(sReq)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

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

	if val, ok := urlQuery["channelID[]"]; ok {
		sReq.ChannelID = parseUInt64A(val)
	} else if val, ok = urlQuery["channelID"]; ok {
		sReq.ChannelID = parseUInt64A(val)
	}

	if val, ok := get["afterMessageID"]; ok {

		sReq.AfterMessageID = parseUInt64(val)
	}
	if val, ok := get["beforeMessageID"]; ok {

		sReq.BeforeMessageID = parseUInt64(val)
	}
	if val, ok := get["fromMessageID"]; ok {

		sReq.FromMessageID = parseUInt64(val)
	}
	if val, ok := get["toMessageID"]; ok {

		sReq.ToMessageID = parseUInt64(val)
	}

	if val, ok := urlQuery["threadID[]"]; ok {
		sReq.ThreadID = parseUInt64A(val)
	} else if val, ok = urlQuery["threadID"]; ok {
		sReq.ThreadID = parseUInt64A(val)
	}

	if val, ok := urlQuery["userID[]"]; ok {
		sReq.UserID = parseUInt64A(val)
	} else if val, ok = urlQuery["userID"]; ok {
		sReq.UserID = parseUInt64A(val)
	}

	if val, ok := get["pinnedOnly"]; ok {

		sReq.PinnedOnly = parseBool(val)
	}
	if val, ok := get["bookmarkedOnly"]; ok {

		sReq.BookmarkedOnly = parseBool(val)
	}
	if val, ok := get["limit"]; ok {

		sReq.Limit = parseUint(val)
	}
	if val, ok := get["query"]; ok {

		sReq.Query = val
	}

	return err
}

var _ RequestFiller = NewSearchMessages()
