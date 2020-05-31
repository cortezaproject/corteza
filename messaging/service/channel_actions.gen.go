package service

// This file is auto-generated from messaging/service/channel_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
)

type (
	channelActionProps struct {
		channel  *types.Channel
		changed  *types.Channel
		filter   *types.ChannelFilter
		flag     string
		memberID uint64
	}

	channelAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *channelActionProps
	}

	channelError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *channelActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setChannel updates channelActionProps's channel
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *channelActionProps) setChannel(channel *types.Channel) *channelActionProps {
	p.channel = channel
	return p
}

// setChanged updates channelActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *channelActionProps) setChanged(changed *types.Channel) *channelActionProps {
	p.changed = changed
	return p
}

// setFilter updates channelActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *channelActionProps) setFilter(filter *types.ChannelFilter) *channelActionProps {
	p.filter = filter
	return p
}

// setFlag updates channelActionProps's flag
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *channelActionProps) setFlag(flag string) *channelActionProps {
	p.flag = flag
	return p
}

// setMemberID updates channelActionProps's memberID
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *channelActionProps) setMemberID(memberID uint64) *channelActionProps {
	p.memberID = memberID
	return p
}

// serialize converts channelActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p channelActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.channel != nil {
		m.Set("channel.name", p.channel.Name, true)
		m.Set("channel.topic", p.channel.Topic, true)
		m.Set("channel.type", p.channel.Type, true)
		m.Set("channel.ID", p.channel.ID, true)
	}
	if p.changed != nil {
		m.Set("changed.name", p.changed.Name, true)
		m.Set("changed.topic", p.changed.Topic, true)
		m.Set("changed.type", p.changed.Type, true)
		m.Set("changed.ID", p.changed.ID, true)
		m.Set("changed.meta", p.changed.Meta, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.channelID", p.filter.ChannelID, true)
		m.Set("filter.currentUserID", p.filter.CurrentUserID, true)
		m.Set("filter.includeDeleted", p.filter.IncludeDeleted, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}
	m.Set("flag", p.flag, true)
	m.Set("memberID", p.memberID, true)

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p channelActionProps) tr(in string, err error) string {
	var (
		pairs = []string{"{err}"}
		// first non-empty string
		fns = func(ii ...interface{}) string {
			for _, i := range ii {
				if s := fmt.Sprintf("%v", i); len(s) > 0 {
					return s
				}
			}

			return ""
		}
	)

	if err != nil {
		for {
			// Unwrap errors
			ue := errors.Unwrap(err)
			if ue == nil {
				break
			}

			err = ue
		}

		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}

	if p.channel != nil {
		// replacement for "{channel}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{channel}",
			fns(
				p.channel.Name,
				p.channel.Topic,
				p.channel.Type,
				p.channel.ID,
			),
		)
		pairs = append(pairs, "{channel.name}", fns(p.channel.Name))
		pairs = append(pairs, "{channel.topic}", fns(p.channel.Topic))
		pairs = append(pairs, "{channel.type}", fns(p.channel.Type))
		pairs = append(pairs, "{channel.ID}", fns(p.channel.ID))
	}

	if p.changed != nil {
		// replacement for "{changed}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{changed}",
			fns(
				p.changed.Name,
				p.changed.Topic,
				p.changed.Type,
				p.changed.ID,
				p.changed.Meta,
			),
		)
		pairs = append(pairs, "{changed.name}", fns(p.changed.Name))
		pairs = append(pairs, "{changed.topic}", fns(p.changed.Topic))
		pairs = append(pairs, "{changed.type}", fns(p.changed.Type))
		pairs = append(pairs, "{changed.ID}", fns(p.changed.ID))
		pairs = append(pairs, "{changed.meta}", fns(p.changed.Meta))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.ChannelID,
				p.filter.CurrentUserID,
				p.filter.IncludeDeleted,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.channelID}", fns(p.filter.ChannelID))
		pairs = append(pairs, "{filter.currentUserID}", fns(p.filter.CurrentUserID))
		pairs = append(pairs, "{filter.includeDeleted}", fns(p.filter.IncludeDeleted))
		pairs = append(pairs, "{filter.sort}", fns(p.filter.Sort))
	}
	pairs = append(pairs, "{flag}", fns(p.flag))
	pairs = append(pairs, "{memberID}", fns(p.memberID))
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
//
func (a *channelAction) String() string {
	var props = &channelActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *channelAction) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error methods

// String returns loggable description as string
//
// It falls back to message if log is not set
//
// This function is auto-generated.
//
func (e *channelError) String() string {
	var props = &channelActionProps{}

	if e.props != nil {
		props = e.props
	}

	if e.wrap != nil && !strings.Contains(e.log, "{err}") {
		// Suffix error log with {err} to ensure
		// we log the cause for this error
		e.log += ": {err}"
	}

	return props.tr(e.log, e.wrap)
}

// Error satisfies
//
// This function is auto-generated.
//
func (e *channelError) Error() string {
	var props = &channelActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *channelError) Is(Resource error) bool {
	t, ok := Resource.(*channelError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps channelError around another error
//
// This function is auto-generated.
//
func (e *channelError) Wrap(err error) *channelError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *channelError) Unwrap() error {
	return e.wrap
}

func (e *channelError) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Error:       e.Error(),
		Meta:        e.props.serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// ChannelActionCreate returns "messaging:channel.create" error
//
// This function is auto-generated.
//
func ChannelActionCreate(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "create",
		log:       "created {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionUpdate returns "messaging:channel.update" error
//
// This function is auto-generated.
//
func ChannelActionUpdate(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "update",
		log:       "updated {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionDelete returns "messaging:channel.delete" error
//
// This function is auto-generated.
//
func ChannelActionDelete(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "delete",
		log:       "deleted {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionUndelete returns "messaging:channel.undelete" error
//
// This function is auto-generated.
//
func ChannelActionUndelete(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "undelete",
		log:       "undeleted {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionArchive returns "messaging:channel.archive" error
//
// This function is auto-generated.
//
func ChannelActionArchive(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "archive",
		log:       "archived {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionUnarchive returns "messaging:channel.unarchive" error
//
// This function is auto-generated.
//
func ChannelActionUnarchive(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "unarchive",
		log:       "unarchived {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionSetFlag returns "messaging:channel.setFlag" error
//
// This function is auto-generated.
//
func ChannelActionSetFlag(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "setFlag",
		log:       "set flag {flag} on {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionInviteMember returns "messaging:channel.inviteMember" error
//
// This function is auto-generated.
//
func ChannelActionInviteMember(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "inviteMember",
		log:       "member {memberID} invited to {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionRemoveMember returns "messaging:channel.removeMember" error
//
// This function is auto-generated.
//
func ChannelActionRemoveMember(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "removeMember",
		log:       "member {memberID} removed from {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionAddMember returns "messaging:channel.addMember" error
//
// This function is auto-generated.
//
func ChannelActionAddMember(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "addMember",
		log:       "member {memberID} added to {channel}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// ChannelErrGeneric returns "messaging:channel.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrGeneric(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotFound returns "messaging:channel.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ChannelErrNotFound(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notFound",
		action:    "error",
		message:   "channel does not exist",
		log:       "channel does not exist",
		severity:  actionlog.Warning,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrInvalidID returns "messaging:channel.invalidID" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrInvalidID(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrInvalidType returns "messaging:channel.invalidType" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrInvalidType(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "invalidType",
		action:    "error",
		message:   "invalid type",
		log:       "invalid type",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNameLength returns "messaging:channel.nameLength" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNameLength(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "nameLength",
		action:    "error",
		message:   "name too long",
		log:       "name too long",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNameEmpty returns "messaging:channel.nameEmpty" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNameEmpty(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "nameEmpty",
		action:    "error",
		message:   "name not set",
		log:       "name not set",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrTopicLength returns "messaging:channel.topicLength" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrTopicLength(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "topicLength",
		action:    "error",
		message:   "topic too long",
		log:       "topic too long",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrAlreadyDeleted returns "messaging:channel.alreadyDeleted" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrAlreadyDeleted(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "alreadyDeleted",
		action:    "error",
		message:   "channel already deleted",
		log:       "channel already deleted",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotDeleted returns "messaging:channel.notDeleted" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotDeleted(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notDeleted",
		action:    "error",
		message:   "channel is not deleted",
		log:       "channel is not deleted",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrAlreadyArchived returns "messaging:channel.alreadyArchived" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrAlreadyArchived(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "alreadyArchived",
		action:    "error",
		message:   "channel already archived",
		log:       "channel already archived",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotArchived returns "messaging:channel.notArchived" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotArchived(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notArchived",
		action:    "error",
		message:   "channel is not archived",
		log:       "channel is not archived",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotMember returns "messaging:channel.notMember" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotMember(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notMember",
		action:    "error",
		message:   "not a member of this channel",
		log:       "not a member of this channel",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrUnableToManageGroupMembers returns "messaging:channel.unableToManageGroupMembers" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrUnableToManageGroupMembers(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "unableToManageGroupMembers",
		action:    "error",
		message:   "channel already deleted",
		log:       "channel already deleted",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotAllowedToRead returns "messaging:channel.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToRead(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this channel",
		log:       "could not read {channel}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotAllowedToListChannels returns "messaging:channel.notAllowedToListChannels" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToListChannels(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notAllowedToListChannels",
		action:    "error",
		message:   "not allowed to list this channels",
		log:       "could not list channels; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotAllowedToCreate returns "messaging:channel.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToCreate(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create channels",
		log:       "could not create channels; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotAllowedToUpdate returns "messaging:channel.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToUpdate(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this channel",
		log:       "could not update {channel}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotAllowedToJoin returns "messaging:channel.notAllowedToJoin" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToJoin(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notAllowedToJoin",
		action:    "error",
		message:   "not allowed to join this channel",
		log:       "could not join {channel}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotAllowedToPart returns "messaging:channel.notAllowedToPart" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToPart(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notAllowedToPart",
		action:    "error",
		message:   "not allowed to part this channel",
		log:       "could not part {channel}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotAllowedToDelete returns "messaging:channel.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToDelete(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this channel",
		log:       "could not delete {channel}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotAllowedToUndelete returns "messaging:channel.notAllowedToUndelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToUndelete(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this channel",
		log:       "could not undelete {channel}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ChannelErrNotAllowedToManageMembers returns "messaging:channel.notAllowedToManageMembers" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToManageMembers(props ...*channelActionProps) *channelError {
	var e = &channelError{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		error:     "notAllowedToManageMembers",
		action:    "error",
		message:   "not allowed to manage channel members",
		log:       "could not manage channel members; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *channelActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// context is used to enrich audit log entry with current user info, request ID, IP address...
// props are collected action/error properties
// action (optional) fn will be used to construct channelAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc channel) recordAction(ctx context.Context, props *channelActionProps, action func(...*channelActionProps) *channelAction, err error) error {
	var (
		ok bool

		// Return error
		retError *channelError

		// Recorder error
		recError *channelError
	)

	if err != nil {
		if retError, ok = err.(*channelError); !ok {
			// got non-channel error, wrap it with ChannelErrGeneric
			retError = ChannelErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use ChannelErrGeneric for recording too
			// because it can hold more info
			recError = retError
		} else if retError != nil {
			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}
			// start with copy of return error for recording
			// this will be updated with tha root cause as we try and
			// unwrap the error
			recError = retError

			// find the original recError for this error
			// for the purpose of logging
			var unwrappedError error = retError
			for {
				if unwrappedError = errors.Unwrap(unwrappedError); unwrappedError == nil {
					// nothing wrapped
					break
				}

				// update recError ONLY of wrapped error is of type channelError
				if unwrappedSinkError, ok := unwrappedError.(*channelError); ok {
					recError = unwrappedSinkError
				}
			}

			if retError.props == nil {
				// set props on returning error if empty
				retError.props = props
			}

			if recError.props == nil {
				// set props on recording error if empty
				recError.props = props
			}
		}
	}

	if svc.actionlog != nil {
		if retError != nil {
			// failed action, log error
			svc.actionlog.Record(ctx, recError)
		} else if action != nil {
			// successful
			svc.actionlog.Record(ctx, action(props))
		}
	}

	if err == nil {
		// retError not an interface and that WILL (!!) cause issues
		// with nil check (== nil) when it is not explicitly returned
		return nil
	}

	return retError
}
