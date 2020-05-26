package actionlog

import (
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"time"
)

type (
	// Any additional data
	// that can be packed with the raised audit event
	Meta map[string]string

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
