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
	ChannelID       []string
	AfterMessageID  uint64 `json:",string"`
	BeforeMessageID uint64 `json:",string"`
	FromMessageID   uint64 `json:",string"`
	ToMessageID     uint64 `json:",string"`
	ThreadID        []string
	UserID          []string
	Type            []string
	PinnedOnly      bool
	BookmarkedOnly  bool
	Limit           uint
	Query           string
}

func NewSearchMessages() *SearchMessages {
	return &SearchMessages{}
}

func (r SearchMessages) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["afterMessageID"] = r.AfterMessageID
	out["beforeMessageID"] = r.BeforeMessageID
	out["fromMessageID"] = r.FromMessageID
	out["toMessageID"] = r.ToMessageID
	out["threadID"] = r.ThreadID
	out["userID"] = r.UserID
	out["type"] = r.Type
	out["pinnedOnly"] = r.PinnedOnly
	out["bookmarkedOnly"] = r.BookmarkedOnly
	out["limit"] = r.Limit
	out["query"] = r.Query

	return out
}

func (r *SearchMessages) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := get["afterMessageID"]; ok {
		r.AfterMessageID = parseUInt64(val)
	}
	if val, ok := get["beforeMessageID"]; ok {
		r.BeforeMessageID = parseUInt64(val)
	}
	if val, ok := get["fromMessageID"]; ok {
		r.FromMessageID = parseUInt64(val)
	}
	if val, ok := get["toMessageID"]; ok {
		r.ToMessageID = parseUInt64(val)
	}
	if val, ok := get["pinnedOnly"]; ok {
		r.PinnedOnly = parseBool(val)
	}
	if val, ok := get["bookmarkedOnly"]; ok {
		r.BookmarkedOnly = parseBool(val)
	}
	if val, ok := get["limit"]; ok {
		r.Limit = parseUint(val)
	}
	if val, ok := get["query"]; ok {
		r.Query = val
	}

	return err
}

var _ RequestFiller = NewSearchMessages()

// Search threads request parameters
type SearchThreads struct {
	ChannelID []string
	Limit     uint
	Query     string
}

func NewSearchThreads() *SearchThreads {
	return &SearchThreads{}
}

func (r SearchThreads) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["limit"] = r.Limit
	out["query"] = r.Query

	return out
}

func (r *SearchThreads) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return errors.Wrap(err, "error parsing http request body")
		}
	}

	if err = req.ParseForm(); err != nil {
		return err
	}

	get := map[string]string{}
	post := map[string]string{}
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := req.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	if val, ok := get["limit"]; ok {
		r.Limit = parseUint(val)
	}
	if val, ok := get["query"]; ok {
		r.Query = val
	}

	return err
}

var _ RequestFiller = NewSearchThreads()
