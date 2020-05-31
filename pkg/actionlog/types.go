package actionlog

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// Any additional data
	// that can be packed with the raised audit event
	Meta map[string]interface{}

	// Severity determinants event severity level
	Severity uint8

	// Standardized data structure for audit log events
	Action struct {
		// Timestamp of the raised event
		Timestamp time.Time `json:"timestamp"`

		// Origin of the action (rest-api, cli, grpc, system)
		RequestOrigin string `json:"requestOrigin"`

		// Request ID
		RequestID string `json:"requestID"`

		// This can contain a series of IP addresses (when proxied)
		// https://en.wikipedia.org/wiki/X-Forwarded-For#Format
		ActorIPAddr string `json:"actorIPAddr"`

		// ID of the user (if not anonymous)
		ActorID uint64 `json:"actorID,string"`

		// Resource
		Resource string `json:"resource"`

		// Type of action
		Action string `json:"action"`

		// Type of error
		Error string `json:"error"`

		// Action severity
		Severity Severity `json:"severity"`

		// Description of the event
		Description string `json:"description"`

		// Meta data, resource specific values
		Meta Meta `json:"meta"`
	}

	Filter struct {
		From     *time.Time `json:"from"`
		To       *time.Time `json:"to"`
		ActorID  []uint64   `json:"actorID"`
		Resource string     `json:"resource"`
		Action   string     `json:"action"`

		// @todo pending implementation
		// Query   string     `json:"query"`

		// Standard paging fields & helpers
		rh.PageFilter
	}

	loggableMetaValue interface {
		ActionLogMetaValue() (interface{}, bool)
	}
)

const (
	// Not using log/syslog LOG_* constants as they are only
	// available outside windows env.
	Emergency Severity = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Info
	Debug
)

func (a *Action) LoggableAction() *Action { return a }

// Set value as key to meta, or skip it if value is empty (and omitempty is set)
func (m Meta) Set(key string, in interface{}, omitempty bool) {
	// if value has fn() ActionLogMetaValue that can
	// reformat value and (maybe) omit empty values, call that first
	if lmv, is := in.(loggableMetaValue); is {
		if v, empty := lmv.ActionLogMetaValue(); !omitempty || !empty {
			m[key] = v
		}

		return
	}

	if !omitempty {
		// Nothing special,
		// assign value and quit
		m[key] = in
		return
	}

	// for the rest, we need to determine what kind of

	if str, is := in.(string); is {
		if !omitempty || len(str) > 0 {
			m[key] = str
		}

		return
	}

	if bb, is := in.([]byte); is {
		if !omitempty || len(bb) > 0 {
			m[key] = bb
		}

		return
	}

	if b, is := in.(bool); is {
		m[key] = b
		return
	}

	// cast to (int|uint|float)64
	num := func(n interface{}) interface{} {
		switch n := n.(type) {
		case int:
			return int64(n)
		case int8:
			return int64(n)
		case int16:
			return int64(n)
		case int32:
			return int64(n)
		case int64:
			return n
		case uint:
			return uint64(n)
		case uintptr:
			return uint64(n)
		case uint8:
			return uint64(n)
		case uint16:
			return uint64(n)
		case uint32:
			return uint64(n)
		case uint64:
			return n
		case float32:
			return float64(n)
		case float64:
			return n
		}
		return n
	}(in)

	switch num := num.(type) {
	case nil:
	case uint64:
		if !omitempty || num > 0 {
			m[key] = in
		}
		return

	case int64:
		if !omitempty || num > 0 {
			m[key] = in
		}
		return

	case float64:
		if !omitempty || num > 0 {
			m[key] = in
		}

		return
	}

	// for the rest (slices, etc..)
	// just set the value
	m[key] = in
}
