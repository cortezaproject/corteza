package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// messaging/service/channel_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"strings"
	"time"
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

	channelLogMetaKey   struct{}
	channelPropsMetaKey struct{}
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

// Serialize converts channelActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p channelActionProps) Serialize() actionlog.Meta {
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
func (p channelActionProps) Format(in string, err error) string {
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

	return props.Format(a.log, nil)
}

func (e *channelAction) ToAction() *actionlog.Action {
	return &actionlog.Action{
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.Serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// ChannelActionCreate returns "messaging:channel.create" action
//
// This function is auto-generated.
//
func ChannelActionCreate(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "create",
		log:       "created {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionUpdate returns "messaging:channel.update" action
//
// This function is auto-generated.
//
func ChannelActionUpdate(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "update",
		log:       "updated {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionDelete returns "messaging:channel.delete" action
//
// This function is auto-generated.
//
func ChannelActionDelete(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "delete",
		log:       "deleted {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionUndelete returns "messaging:channel.undelete" action
//
// This function is auto-generated.
//
func ChannelActionUndelete(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "undelete",
		log:       "undeleted {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionArchive returns "messaging:channel.archive" action
//
// This function is auto-generated.
//
func ChannelActionArchive(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "archive",
		log:       "archived {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionUnarchive returns "messaging:channel.unarchive" action
//
// This function is auto-generated.
//
func ChannelActionUnarchive(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "unarchive",
		log:       "unarchived {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionSetFlag returns "messaging:channel.setFlag" action
//
// This function is auto-generated.
//
func ChannelActionSetFlag(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "setFlag",
		log:       "set flag {flag} on {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionInviteMember returns "messaging:channel.inviteMember" action
//
// This function is auto-generated.
//
func ChannelActionInviteMember(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "inviteMember",
		log:       "member {memberID} invited to {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionRemoveMember returns "messaging:channel.removeMember" action
//
// This function is auto-generated.
//
func ChannelActionRemoveMember(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "removeMember",
		log:       "member {memberID} removed from {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChannelActionAddMember returns "messaging:channel.addMember" action
//
// This function is auto-generated.
//
func ChannelActionAddMember(props ...*channelActionProps) *channelAction {
	a := &channelAction{
		timestamp: time.Now(),
		resource:  "messaging:channel",
		action:    "addMember",
		log:       "member {memberID} added to {channel}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// ChannelErrGeneric returns "messaging:channel.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrGeneric(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "{err}"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotFound returns "messaging:channel.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotFound(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("channel does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrInvalidID returns "messaging:channel.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrInvalidID(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrInvalidType returns "messaging:channel.invalidType" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrInvalidType(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid type", nil),

		errors.Meta("type", "invalidType"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNameLength returns "messaging:channel.nameLength" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNameLength(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("name too long", nil),

		errors.Meta("type", "nameLength"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNameEmpty returns "messaging:channel.nameEmpty" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNameEmpty(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("name not set", nil),

		errors.Meta("type", "nameEmpty"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrTopicLength returns "messaging:channel.topicLength" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrTopicLength(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("topic too long", nil),

		errors.Meta("type", "topicLength"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrAlreadyDeleted returns "messaging:channel.alreadyDeleted" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrAlreadyDeleted(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("channel already deleted", nil),

		errors.Meta("type", "alreadyDeleted"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotDeleted returns "messaging:channel.notDeleted" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotDeleted(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("channel is not deleted", nil),

		errors.Meta("type", "notDeleted"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrAlreadyArchived returns "messaging:channel.alreadyArchived" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrAlreadyArchived(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("channel already archived", nil),

		errors.Meta("type", "alreadyArchived"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotArchived returns "messaging:channel.notArchived" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotArchived(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("channel is not archived", nil),

		errors.Meta("type", "notArchived"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotMember returns "messaging:channel.notMember" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotMember(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not a member of this channel", nil),

		errors.Meta("type", "notMember"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrUnableToManageGroupMembers returns "messaging:channel.unableToManageGroupMembers" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrUnableToManageGroupMembers(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("channel already deleted", nil),

		errors.Meta("type", "unableToManageGroupMembers"),
		errors.Meta("resource", "messaging:channel"),

		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotAllowedToRead returns "messaging:channel.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToRead(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this channel", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "could not read {channel}; insufficient permissions"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotAllowedToListChannels returns "messaging:channel.notAllowedToListChannels" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToListChannels(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list this channels", nil),

		errors.Meta("type", "notAllowedToListChannels"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "could not list channels; insufficient permissions"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotAllowedToCreate returns "messaging:channel.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToCreate(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create channels", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "could not create channels; insufficient permissions"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotAllowedToUpdate returns "messaging:channel.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToUpdate(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this channel", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "could not update {channel}; insufficient permissions"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotAllowedToJoin returns "messaging:channel.notAllowedToJoin" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToJoin(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to join this channel", nil),

		errors.Meta("type", "notAllowedToJoin"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "could not join {channel}; insufficient permissions"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotAllowedToPart returns "messaging:channel.notAllowedToPart" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToPart(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to part this channel", nil),

		errors.Meta("type", "notAllowedToPart"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "could not part {channel}; insufficient permissions"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotAllowedToDelete returns "messaging:channel.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToDelete(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this channel", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "could not delete {channel}; insufficient permissions"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotAllowedToUndelete returns "messaging:channel.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToUndelete(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this channel", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "could not undelete {channel}; insufficient permissions"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChannelErrNotAllowedToManageMembers returns "messaging:channel.notAllowedToManageMembers" as *errors.Error
//
//
// This function is auto-generated.
//
func ChannelErrNotAllowedToManageMembers(mm ...*channelActionProps) *errors.Error {
	var p = &channelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage channel members", nil),

		errors.Meta("type", "notAllowedToManageMembers"),
		errors.Meta("resource", "messaging:channel"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(channelLogMetaKey{}, "could not manage channel members; insufficient permissions"),
		errors.Meta(channelPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// It will wrap unrecognized/internal errors with generic errors.
//
// This function is auto-generated.
//
func (svc channel) recordAction(ctx context.Context, props *channelActionProps, actionFn func(...*channelActionProps) *channelAction, err error) error {
	if svc.actionlog == nil || actionFn == nil {
		// action log disabled or no action fn passed, return error as-is
		return err
	} else if err == nil {
		// action completed w/o error, record it
		svc.actionlog.Record(ctx, actionFn(props).ToAction())
		return nil
	}

	a := actionFn(props).ToAction()

	// Extracting error information and recording it as action
	a.Error = err.Error()

	switch c := err.(type) {
	case *errors.Error:
		m := c.Meta()

		a.Error = err.Error()
		a.Severity = actionlog.Severity(m.AsInt("severity"))
		a.Description = props.Format(m.AsString(channelLogMetaKey{}), err)

		if p, has := m[channelPropsMetaKey{}]; has {
			a.Meta = p.(*channelActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
