package actionlog

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/sql"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	// Any additional data
	// that can be packed with the raised audit event
	Meta map[string]interface{}

	// Severity determinants event severity level
	Severity uint8

	// Standardized data structure for audit log events
	Action struct {
		ID uint64 `json:"actionID,string"`

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

		// Error, if any
		Error string `json:"error"`

		// Action severity
		Severity Severity `json:"severity"`

		// Description of the event
		Description string `json:"description"`

		// Meta data, resource specific values
		Meta Meta `json:"meta"`
	}

	Filter struct {
		FromTimestamp *time.Time `json:"from"`
		ToTimestamp   *time.Time `json:"to"`

		BeforeActionID uint64 `json:"beforeActionID"`

		ActorID  []uint64 `json:"actorID"`
		Origin   string   `json:"origin"`
		Resource string   `json:"resource"`
		Action   string   `json:"action"`
		Limit    uint     `json:"limit"`

		// Standard helpers for sorting
		filter.Sorting
	}

	loggableMetaValue interface {
		ActionLogMetaValue() (interface{}, bool)
	}
)

// Severity constants
//
// not using log/syslog LOG_* constants as they are only
// available outside windows env.
const (
	Emergency Severity = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Info
	Debug

	ActionResourceType = "corteza::generic:action"
)

func (a *Action) ToAction() *Action { return a }

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

	{
		// properly encode big numbers
		// before storing them as JSON

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
				return int64(n)
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
				return uint64(n)
			case float32:
				return float64(n)
			case float64:
				return n
			}
			return n
		}(in)

		switch num := num.(type) {
		case int64:
			if !omitempty || num > 0 {
				m[key] = strconv.FormatInt(num, 10)
			}
			return
		case uint64:
			if !omitempty || num > 0 {
				m[key] = strconv.FormatUint(num, 10)
			}
			return

		case float64:
			if !omitempty || num > 0 {
				m[key] = in
			}

			return
		}
	}

	// for the rest (string, slices, etc..)
	// just set the value
	m[key] = in
}

func (s Severity) String() string {
	switch s {
	case Emergency:
		return "emergency"
	case Alert:
		return "alert"
	case Critical:
		return "critical"
	case Error:
		return "err"
	case Warning:
		return "warning"
	case Notice:
		return "notice"
	case Info:
		return "info"
	case Debug:
		return "debug"
	}

	return ""
}

func NewSeverity(s string) Severity {
	switch s {
	case "emergency":
		return Emergency
	case "alert":
		return Alert
	case "critical":
		return Critical
	case "err":
		return Error
	case "warning":
		return Warning
	case "notice":
		return Notice
	case "info":
		return Info
	}

	return Debug
}

func (m *Meta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m Meta) Value() (driver.Value, error) { return json.Marshal(m) }
