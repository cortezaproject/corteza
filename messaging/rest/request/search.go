package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
)

type (
	// Internal API interface
	SearchMessages struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// ChannelID GET parameter
		//
		// Filter by channels
		ChannelID []string

		// AfterMessageID GET parameter
		//
		// ID of the first message in the list (exclusive)
		AfterMessageID uint64 `json:",string"`

		// BeforeMessageID GET parameter
		//
		// ID of the last message in the list (exclusive)
		BeforeMessageID uint64 `json:",string"`

		// FromMessageID GET parameter
		//
		// ID of the first message in the list (inclusive)
		FromMessageID uint64 `json:",string"`

		// ToMessageID GET parameter
		//
		// ID of the last message the list (inclusive)
		ToMessageID uint64 `json:",string"`

		// ThreadID GET parameter
		//
		// Filter by thread message ID
		ThreadID []string

		// UserID GET parameter
		//
		// Filter by one or more user
		UserID []string

		// Type GET parameter
		//
		// Filter by message type (text, inlineImage, attachment, ...)
		Type []string

		// PinnedOnly GET parameter
		//
		// Return only pinned messages
		PinnedOnly bool

		// BookmarkedOnly GET parameter
		//
		// Only bookmarked messages
		BookmarkedOnly bool

		// Limit GET parameter
		//
		// Max number of messages
		Limit uint
	}

	SearchThreads struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// ChannelID GET parameter
		//
		// Filter by channels
		ChannelID []string

		// Limit GET parameter
		//
		// Max number of messages
		Limit uint
	}
)

// NewSearchMessages request
func NewSearchMessages() *SearchMessages {
	return &SearchMessages{}
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":           r.Query,
		"channelID":       r.ChannelID,
		"afterMessageID":  r.AfterMessageID,
		"beforeMessageID": r.BeforeMessageID,
		"fromMessageID":   r.FromMessageID,
		"toMessageID":     r.ToMessageID,
		"threadID":        r.ThreadID,
		"userID":          r.UserID,
		"type":            r.Type,
		"pinnedOnly":      r.PinnedOnly,
		"bookmarkedOnly":  r.BookmarkedOnly,
		"limit":           r.Limit,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetChannelID() []string {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetAfterMessageID() uint64 {
	return r.AfterMessageID
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetBeforeMessageID() uint64 {
	return r.BeforeMessageID
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetFromMessageID() uint64 {
	return r.FromMessageID
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetToMessageID() uint64 {
	return r.ToMessageID
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetThreadID() []string {
	return r.ThreadID
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetUserID() []string {
	return r.UserID
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetType() []string {
	return r.Type
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetPinnedOnly() bool {
	return r.PinnedOnly
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetBookmarkedOnly() bool {
	return r.BookmarkedOnly
}

// Auditable returns all auditable/loggable parameters
func (r SearchMessages) GetLimit() uint {
	return r.Limit
}

// Fill processes request and fills internal variables
func (r *SearchMessages) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["channelID[]"]; ok {
			r.ChannelID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["channelID"]; ok {
			r.ChannelID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["afterMessageID"]; ok && len(val) > 0 {
			r.AfterMessageID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["beforeMessageID"]; ok && len(val) > 0 {
			r.BeforeMessageID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["fromMessageID"]; ok && len(val) > 0 {
			r.FromMessageID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["toMessageID"]; ok && len(val) > 0 {
			r.ToMessageID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["threadID[]"]; ok {
			r.ThreadID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["threadID"]; ok {
			r.ThreadID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["userID[]"]; ok {
			r.UserID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["userID"]; ok {
			r.UserID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["type[]"]; ok {
			r.Type, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["type"]; ok {
			r.Type, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["pinnedOnly"]; ok && len(val) > 0 {
			r.PinnedOnly, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["bookmarkedOnly"]; ok && len(val) > 0 {
			r.BookmarkedOnly, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewSearchThreads request
func NewSearchThreads() *SearchThreads {
	return &SearchThreads{}
}

// Auditable returns all auditable/loggable parameters
func (r SearchThreads) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":     r.Query,
		"channelID": r.ChannelID,
		"limit":     r.Limit,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SearchThreads) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r SearchThreads) GetChannelID() []string {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r SearchThreads) GetLimit() uint {
	return r.Limit
}

// Fill processes request and fills internal variables
func (r *SearchThreads) Fill(req *http.Request) (err error) {
	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["channelID[]"]; ok {
			r.ChannelID, err = val, nil
			if err != nil {
				return err
			}
		} else if val, ok := tmp["channelID"]; ok {
			r.ChannelID, err = val, nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
