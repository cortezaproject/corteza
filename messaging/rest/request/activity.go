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
	ActivitySend struct {
		// ChannelID POST parameter
		//
		// Channel ID, if set, activity will be send only to subscribed users
		ChannelID uint64 `json:",string"`

		// MessageID POST parameter
		//
		// Message ID, if set, channelID must be set as well
		MessageID uint64 `json:",string"`

		// Kind POST parameter
		//
		// Arbitrary string
		Kind string
	}
)

// NewActivitySend request
func NewActivitySend() *ActivitySend {
	return &ActivitySend{}
}

// Auditable returns all auditable/loggable parameters
func (r ActivitySend) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"channelID": r.ChannelID,
		"messageID": r.MessageID,
		"kind":      r.Kind,
	}
}

// Auditable returns all auditable/loggable parameters
func (r ActivitySend) GetChannelID() uint64 {
	return r.ChannelID
}

// Auditable returns all auditable/loggable parameters
func (r ActivitySend) GetMessageID() uint64 {
	return r.MessageID
}

// Auditable returns all auditable/loggable parameters
func (r ActivitySend) GetKind() string {
	return r.Kind
}

// Fill processes request and fills internal variables
func (r *ActivitySend) Fill(req *http.Request) (err error) {
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
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["channelID"]; ok && len(val) > 0 {
			r.ChannelID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["messageID"]; ok && len(val) > 0 {
			r.MessageID, err = payload.ParseUint64(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["kind"]; ok && len(val) > 0 {
			r.Kind, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}
