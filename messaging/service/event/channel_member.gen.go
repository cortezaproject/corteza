package event

// This file is auto-generated.
//
// YAML event definitions:
//   messaging/service/event/events.yaml
//
// Regenerate with:
//   go run codegen/v2/events.go --service messaging
//

import (
	"github.com/cortezaproject/corteza-server/messaging/types"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// channelMemberBase
	//
	// This type is auto-generated.
	channelMemberBase struct {
		member  *types.ChannelMember
		channel *types.Channel
		invoker auth.Identifiable
	}

	// channelMemberBeforeJoin
	//
	// This type is auto-generated.
	channelMemberBeforeJoin struct {
		*channelMemberBase
	}

	// channelMemberBeforePart
	//
	// This type is auto-generated.
	channelMemberBeforePart struct {
		*channelMemberBase
	}

	// channelMemberBeforeAdd
	//
	// This type is auto-generated.
	channelMemberBeforeAdd struct {
		*channelMemberBase
	}

	// channelMemberBeforeRemove
	//
	// This type is auto-generated.
	channelMemberBeforeRemove struct {
		*channelMemberBase
	}

	// channelMemberAfterJoin
	//
	// This type is auto-generated.
	channelMemberAfterJoin struct {
		*channelMemberBase
	}

	// channelMemberAfterPart
	//
	// This type is auto-generated.
	channelMemberAfterPart struct {
		*channelMemberBase
	}

	// channelMemberAfterAdd
	//
	// This type is auto-generated.
	channelMemberAfterAdd struct {
		*channelMemberBase
	}

	// channelMemberAfterRemove
	//
	// This type is auto-generated.
	channelMemberAfterRemove struct {
		*channelMemberBase
	}
)

// ResourceType returns "messaging:channel:member"
//
// This function is auto-generated.
func (channelMemberBase) ResourceType() string {
	return "messaging:channel:member"
}

// EventType on channelMemberBeforeJoin returns "beforeJoin"
//
// This function is auto-generated.
func (channelMemberBeforeJoin) EventType() string {
	return "beforeJoin"
}

// EventType on channelMemberBeforePart returns "beforePart"
//
// This function is auto-generated.
func (channelMemberBeforePart) EventType() string {
	return "beforePart"
}

// EventType on channelMemberBeforeAdd returns "beforeAdd"
//
// This function is auto-generated.
func (channelMemberBeforeAdd) EventType() string {
	return "beforeAdd"
}

// EventType on channelMemberBeforeRemove returns "beforeRemove"
//
// This function is auto-generated.
func (channelMemberBeforeRemove) EventType() string {
	return "beforeRemove"
}

// EventType on channelMemberAfterJoin returns "afterJoin"
//
// This function is auto-generated.
func (channelMemberAfterJoin) EventType() string {
	return "afterJoin"
}

// EventType on channelMemberAfterPart returns "afterPart"
//
// This function is auto-generated.
func (channelMemberAfterPart) EventType() string {
	return "afterPart"
}

// EventType on channelMemberAfterAdd returns "afterAdd"
//
// This function is auto-generated.
func (channelMemberAfterAdd) EventType() string {
	return "afterAdd"
}

// EventType on channelMemberAfterRemove returns "afterRemove"
//
// This function is auto-generated.
func (channelMemberAfterRemove) EventType() string {
	return "afterRemove"
}

// ChannelMemberBeforeJoin creates beforeJoin for messaging:channel:member resource
//
// This function is auto-generated.
func ChannelMemberBeforeJoin(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberBeforeJoin {
	return &channelMemberBeforeJoin{
		channelMemberBase: &channelMemberBase{
			member:  argMember,
			channel: argChannel,
		},
	}
}

// ChannelMemberBeforePart creates beforePart for messaging:channel:member resource
//
// This function is auto-generated.
func ChannelMemberBeforePart(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberBeforePart {
	return &channelMemberBeforePart{
		channelMemberBase: &channelMemberBase{
			member:  argMember,
			channel: argChannel,
		},
	}
}

// ChannelMemberBeforeAdd creates beforeAdd for messaging:channel:member resource
//
// This function is auto-generated.
func ChannelMemberBeforeAdd(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberBeforeAdd {
	return &channelMemberBeforeAdd{
		channelMemberBase: &channelMemberBase{
			member:  argMember,
			channel: argChannel,
		},
	}
}

// ChannelMemberBeforeRemove creates beforeRemove for messaging:channel:member resource
//
// This function is auto-generated.
func ChannelMemberBeforeRemove(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberBeforeRemove {
	return &channelMemberBeforeRemove{
		channelMemberBase: &channelMemberBase{
			member:  argMember,
			channel: argChannel,
		},
	}
}

// ChannelMemberAfterJoin creates afterJoin for messaging:channel:member resource
//
// This function is auto-generated.
func ChannelMemberAfterJoin(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberAfterJoin {
	return &channelMemberAfterJoin{
		channelMemberBase: &channelMemberBase{
			member:  argMember,
			channel: argChannel,
		},
	}
}

// ChannelMemberAfterPart creates afterPart for messaging:channel:member resource
//
// This function is auto-generated.
func ChannelMemberAfterPart(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberAfterPart {
	return &channelMemberAfterPart{
		channelMemberBase: &channelMemberBase{
			member:  argMember,
			channel: argChannel,
		},
	}
}

// ChannelMemberAfterAdd creates afterAdd for messaging:channel:member resource
//
// This function is auto-generated.
func ChannelMemberAfterAdd(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberAfterAdd {
	return &channelMemberAfterAdd{
		channelMemberBase: &channelMemberBase{
			member:  argMember,
			channel: argChannel,
		},
	}
}

// ChannelMemberAfterRemove creates afterRemove for messaging:channel:member resource
//
// This function is auto-generated.
func ChannelMemberAfterRemove(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberAfterRemove {
	return &channelMemberAfterRemove{
		channelMemberBase: &channelMemberBase{
			member:  argMember,
			channel: argChannel,
		},
	}
}

// SetMember sets new member value
//
// This function is auto-generated.
func (res *channelMemberBase) SetMember(argMember *types.ChannelMember) {
	res.member = argMember
}

// Member returns member
//
// This function is auto-generated.
func (res channelMemberBase) Member() *types.ChannelMember {
	return res.member
}

// SetChannel sets new channel value
//
// This function is auto-generated.
func (res *channelMemberBase) SetChannel(argChannel *types.Channel) {
	res.channel = argChannel
}

// Channel returns channel
//
// This function is auto-generated.
func (res channelMemberBase) Channel() *types.Channel {
	return res.channel
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *channelMemberBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res channelMemberBase) Invoker() auth.Identifiable {
	return res.invoker
}
