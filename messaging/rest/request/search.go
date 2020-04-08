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

// SearchMessages request parameters
type SearchMessages struct {
	hasChannelID bool
	rawChannelID []string
	ChannelID    []string

	hasAfterMessageID bool
	rawAfterMessageID string
	AfterMessageID    uint64 `json:",string"`

	hasBeforeMessageID bool
	rawBeforeMessageID string
	BeforeMessageID    uint64 `json:",string"`

	hasFromMessageID bool
	rawFromMessageID string
	FromMessageID    uint64 `json:",string"`

	hasToMessageID bool
	rawToMessageID string
	ToMessageID    uint64 `json:",string"`

	hasThreadID bool
	rawThreadID []string
	ThreadID    []string

	hasUserID bool
	rawUserID []string
	UserID    []string

	hasType bool
	rawType []string
	Type    []string

	hasPinnedOnly bool
	rawPinnedOnly string
	PinnedOnly    bool

	hasBookmarkedOnly bool
	rawBookmarkedOnly string
	BookmarkedOnly    bool

	hasLimit bool
	rawLimit string
	Limit    uint

	hasQuery bool
	rawQuery string
	Query    string
}

// NewSearchMessages request
func NewSearchMessages() *SearchMessages {
	return &SearchMessages{}
}

// Auditable returns all auditable/loggable parameters
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

// Fill processes request and fills internal variables
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

	if val, ok := urlQuery["channelID[]"]; ok {
		r.hasChannelID = true
		r.rawChannelID = val
		r.ChannelID = parseStrings(val)
	} else if val, ok = urlQuery["channelID"]; ok {
		r.hasChannelID = true
		r.rawChannelID = val
		r.ChannelID = parseStrings(val)
	}

	if val, ok := get["afterMessageID"]; ok {
		r.hasAfterMessageID = true
		r.rawAfterMessageID = val
		r.AfterMessageID = parseUInt64(val)
	}
	if val, ok := get["beforeMessageID"]; ok {
		r.hasBeforeMessageID = true
		r.rawBeforeMessageID = val
		r.BeforeMessageID = parseUInt64(val)
	}
	if val, ok := get["fromMessageID"]; ok {
		r.hasFromMessageID = true
		r.rawFromMessageID = val
		r.FromMessageID = parseUInt64(val)
	}
	if val, ok := get["toMessageID"]; ok {
		r.hasToMessageID = true
		r.rawToMessageID = val
		r.ToMessageID = parseUInt64(val)
	}

	if val, ok := urlQuery["threadID[]"]; ok {
		r.hasThreadID = true
		r.rawThreadID = val
		r.ThreadID = parseStrings(val)
	} else if val, ok = urlQuery["threadID"]; ok {
		r.hasThreadID = true
		r.rawThreadID = val
		r.ThreadID = parseStrings(val)
	}

	if val, ok := urlQuery["userID[]"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseStrings(val)
	} else if val, ok = urlQuery["userID"]; ok {
		r.hasUserID = true
		r.rawUserID = val
		r.UserID = parseStrings(val)
	}

	if val, ok := urlQuery["type[]"]; ok {
		r.hasType = true
		r.rawType = val
		r.Type = parseStrings(val)
	} else if val, ok = urlQuery["type"]; ok {
		r.hasType = true
		r.rawType = val
		r.Type = parseStrings(val)
	}

	if val, ok := get["pinnedOnly"]; ok {
		r.hasPinnedOnly = true
		r.rawPinnedOnly = val
		r.PinnedOnly = parseBool(val)
	}
	if val, ok := get["bookmarkedOnly"]; ok {
		r.hasBookmarkedOnly = true
		r.rawBookmarkedOnly = val
		r.BookmarkedOnly = parseBool(val)
	}
	if val, ok := get["limit"]; ok {
		r.hasLimit = true
		r.rawLimit = val
		r.Limit = parseUint(val)
	}
	if val, ok := get["query"]; ok {
		r.hasQuery = true
		r.rawQuery = val
		r.Query = val
	}

	return err
}

var _ RequestFiller = NewSearchMessages()

// SearchThreads request parameters
type SearchThreads struct {
	hasChannelID bool
	rawChannelID []string
	ChannelID    []string

	hasLimit bool
	rawLimit string
	Limit    uint

	hasQuery bool
	rawQuery string
	Query    string
}

// NewSearchThreads request
func NewSearchThreads() *SearchThreads {
	return &SearchThreads{}
}

// Auditable returns all auditable/loggable parameters
func (r SearchThreads) Auditable() map[string]interface{} {
	var out = map[string]interface{}{}

	out["channelID"] = r.ChannelID
	out["limit"] = r.Limit
	out["query"] = r.Query

	return out
}

// Fill processes request and fills internal variables
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

	if val, ok := urlQuery["channelID[]"]; ok {
		r.hasChannelID = true
		r.rawChannelID = val
		r.ChannelID = parseStrings(val)
	} else if val, ok = urlQuery["channelID"]; ok {
		r.hasChannelID = true
		r.rawChannelID = val
		r.ChannelID = parseStrings(val)
	}

	if val, ok := get["limit"]; ok {
		r.hasLimit = true
		r.rawLimit = val
		r.Limit = parseUint(val)
	}
	if val, ok := get["query"]; ok {
		r.hasQuery = true
		r.rawQuery = val
		r.Query = val
	}

	return err
}

var _ RequestFiller = NewSearchThreads()

// HasChannelID returns true if channelID was set
func (r *SearchMessages) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *SearchMessages) RawChannelID() []string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *SearchMessages) GetChannelID() []string {
	return r.ChannelID
}

// HasAfterMessageID returns true if afterMessageID was set
func (r *SearchMessages) HasAfterMessageID() bool {
	return r.hasAfterMessageID
}

// RawAfterMessageID returns raw value of afterMessageID parameter
func (r *SearchMessages) RawAfterMessageID() string {
	return r.rawAfterMessageID
}

// GetAfterMessageID returns casted value of  afterMessageID parameter
func (r *SearchMessages) GetAfterMessageID() uint64 {
	return r.AfterMessageID
}

// HasBeforeMessageID returns true if beforeMessageID was set
func (r *SearchMessages) HasBeforeMessageID() bool {
	return r.hasBeforeMessageID
}

// RawBeforeMessageID returns raw value of beforeMessageID parameter
func (r *SearchMessages) RawBeforeMessageID() string {
	return r.rawBeforeMessageID
}

// GetBeforeMessageID returns casted value of  beforeMessageID parameter
func (r *SearchMessages) GetBeforeMessageID() uint64 {
	return r.BeforeMessageID
}

// HasFromMessageID returns true if fromMessageID was set
func (r *SearchMessages) HasFromMessageID() bool {
	return r.hasFromMessageID
}

// RawFromMessageID returns raw value of fromMessageID parameter
func (r *SearchMessages) RawFromMessageID() string {
	return r.rawFromMessageID
}

// GetFromMessageID returns casted value of  fromMessageID parameter
func (r *SearchMessages) GetFromMessageID() uint64 {
	return r.FromMessageID
}

// HasToMessageID returns true if toMessageID was set
func (r *SearchMessages) HasToMessageID() bool {
	return r.hasToMessageID
}

// RawToMessageID returns raw value of toMessageID parameter
func (r *SearchMessages) RawToMessageID() string {
	return r.rawToMessageID
}

// GetToMessageID returns casted value of  toMessageID parameter
func (r *SearchMessages) GetToMessageID() uint64 {
	return r.ToMessageID
}

// HasThreadID returns true if threadID was set
func (r *SearchMessages) HasThreadID() bool {
	return r.hasThreadID
}

// RawThreadID returns raw value of threadID parameter
func (r *SearchMessages) RawThreadID() []string {
	return r.rawThreadID
}

// GetThreadID returns casted value of  threadID parameter
func (r *SearchMessages) GetThreadID() []string {
	return r.ThreadID
}

// HasUserID returns true if userID was set
func (r *SearchMessages) HasUserID() bool {
	return r.hasUserID
}

// RawUserID returns raw value of userID parameter
func (r *SearchMessages) RawUserID() []string {
	return r.rawUserID
}

// GetUserID returns casted value of  userID parameter
func (r *SearchMessages) GetUserID() []string {
	return r.UserID
}

// HasType returns true if type was set
func (r *SearchMessages) HasType() bool {
	return r.hasType
}

// RawType returns raw value of type parameter
func (r *SearchMessages) RawType() []string {
	return r.rawType
}

// GetType returns casted value of  type parameter
func (r *SearchMessages) GetType() []string {
	return r.Type
}

// HasPinnedOnly returns true if pinnedOnly was set
func (r *SearchMessages) HasPinnedOnly() bool {
	return r.hasPinnedOnly
}

// RawPinnedOnly returns raw value of pinnedOnly parameter
func (r *SearchMessages) RawPinnedOnly() string {
	return r.rawPinnedOnly
}

// GetPinnedOnly returns casted value of  pinnedOnly parameter
func (r *SearchMessages) GetPinnedOnly() bool {
	return r.PinnedOnly
}

// HasBookmarkedOnly returns true if bookmarkedOnly was set
func (r *SearchMessages) HasBookmarkedOnly() bool {
	return r.hasBookmarkedOnly
}

// RawBookmarkedOnly returns raw value of bookmarkedOnly parameter
func (r *SearchMessages) RawBookmarkedOnly() string {
	return r.rawBookmarkedOnly
}

// GetBookmarkedOnly returns casted value of  bookmarkedOnly parameter
func (r *SearchMessages) GetBookmarkedOnly() bool {
	return r.BookmarkedOnly
}

// HasLimit returns true if limit was set
func (r *SearchMessages) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *SearchMessages) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *SearchMessages) GetLimit() uint {
	return r.Limit
}

// HasQuery returns true if query was set
func (r *SearchMessages) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *SearchMessages) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *SearchMessages) GetQuery() string {
	return r.Query
}

// HasChannelID returns true if channelID was set
func (r *SearchThreads) HasChannelID() bool {
	return r.hasChannelID
}

// RawChannelID returns raw value of channelID parameter
func (r *SearchThreads) RawChannelID() []string {
	return r.rawChannelID
}

// GetChannelID returns casted value of  channelID parameter
func (r *SearchThreads) GetChannelID() []string {
	return r.ChannelID
}

// HasLimit returns true if limit was set
func (r *SearchThreads) HasLimit() bool {
	return r.hasLimit
}

// RawLimit returns raw value of limit parameter
func (r *SearchThreads) RawLimit() string {
	return r.rawLimit
}

// GetLimit returns casted value of  limit parameter
func (r *SearchThreads) GetLimit() uint {
	return r.Limit
}

// HasQuery returns true if query was set
func (r *SearchThreads) HasQuery() bool {
	return r.hasQuery
}

// RawQuery returns raw value of query parameter
func (r *SearchThreads) RawQuery() string {
	return r.rawQuery
}

// GetQuery returns casted value of  query parameter
func (r *SearchThreads) GetQuery() string {
	return r.Query
}
