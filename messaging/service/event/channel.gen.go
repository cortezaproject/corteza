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
	// channelBase
	//
	// This type is auto-generated.
	channelBase struct {
		channel    *types.Channel
		oldChannel *types.Channel
		invoker    auth.Identifiable
	}

	// channelOnManual
	//
	// This type is auto-generated.
	channelOnManual struct {
		*channelBase
	}

	// channelBeforeCreate
	//
	// This type is auto-generated.
	channelBeforeCreate struct {
		*channelBase
	}

	// channelBeforeUpdate
	//
	// This type is auto-generated.
	channelBeforeUpdate struct {
		*channelBase
	}

	// channelBeforeDelete
	//
	// This type is auto-generated.
	channelBeforeDelete struct {
		*channelBase
	}

	// channelAfterCreate
	//
	// This type is auto-generated.
	channelAfterCreate struct {
		*channelBase
	}

	// channelAfterUpdate
	//
	// This type is auto-generated.
	channelAfterUpdate struct {
		*channelBase
	}

	// channelAfterDelete
	//
	// This type is auto-generated.
	channelAfterDelete struct {
		*channelBase
	}
)

// ResourceType returns "messaging:channel"
//
// This function is auto-generated.
func (channelBase) ResourceType() string {
	return "messaging:channel"
}

// EventType on channelOnManual returns "onManual"
//
// This function is auto-generated.
func (channelOnManual) EventType() string {
	return "onManual"
}

// EventType on channelBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (channelBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on channelBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (channelBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on channelBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (channelBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on channelAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (channelAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on channelAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (channelAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on channelAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (channelAfterDelete) EventType() string {
	return "afterDelete"
}

// ChannelOnManual creates onManual for messaging:channel resource
//
// This function is auto-generated.
func ChannelOnManual(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelOnManual {
	return &channelOnManual{
		channelBase: &channelBase{
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelBeforeCreate creates beforeCreate for messaging:channel resource
//
// This function is auto-generated.
func ChannelBeforeCreate(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelBeforeCreate {
	return &channelBeforeCreate{
		channelBase: &channelBase{
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelBeforeUpdate creates beforeUpdate for messaging:channel resource
//
// This function is auto-generated.
func ChannelBeforeUpdate(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelBeforeUpdate {
	return &channelBeforeUpdate{
		channelBase: &channelBase{
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelBeforeDelete creates beforeDelete for messaging:channel resource
//
// This function is auto-generated.
func ChannelBeforeDelete(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelBeforeDelete {
	return &channelBeforeDelete{
		channelBase: &channelBase{
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelAfterCreate creates afterCreate for messaging:channel resource
//
// This function is auto-generated.
func ChannelAfterCreate(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelAfterCreate {
	return &channelAfterCreate{
		channelBase: &channelBase{
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelAfterUpdate creates afterUpdate for messaging:channel resource
//
// This function is auto-generated.
func ChannelAfterUpdate(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelAfterUpdate {
	return &channelAfterUpdate{
		channelBase: &channelBase{
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelAfterDelete creates afterDelete for messaging:channel resource
//
// This function is auto-generated.
func ChannelAfterDelete(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelAfterDelete {
	return &channelAfterDelete{
		channelBase: &channelBase{
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// SetChannel sets new channel value
//
// This function is auto-generated.
func (res *channelBase) SetChannel(argChannel *types.Channel) {
	res.channel = argChannel
}

// Channel returns channel
//
// This function is auto-generated.
func (res channelBase) Channel() *types.Channel {
	return res.channel
}

// OldChannel returns oldChannel
//
// This function is auto-generated.
func (res channelBase) OldChannel() *types.Channel {
	return res.oldChannel
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *channelBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res channelBase) Invoker() auth.Identifiable {
	return res.invoker
}
