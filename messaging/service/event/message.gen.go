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
	// messageBase
	//
	// This type is auto-generated.
	messageBase struct {
		message    *types.Message
		oldMessage *types.Message
		channel    *types.Channel
		invoker    auth.Identifiable
	}

	// messageOnManual
	//
	// This type is auto-generated.
	messageOnManual struct {
		*messageBase
	}

	// messageBeforeCreate
	//
	// This type is auto-generated.
	messageBeforeCreate struct {
		*messageBase
	}

	// messageBeforeUpdate
	//
	// This type is auto-generated.
	messageBeforeUpdate struct {
		*messageBase
	}

	// messageBeforeDelete
	//
	// This type is auto-generated.
	messageBeforeDelete struct {
		*messageBase
	}

	// messageAfterCreate
	//
	// This type is auto-generated.
	messageAfterCreate struct {
		*messageBase
	}

	// messageAfterUpdate
	//
	// This type is auto-generated.
	messageAfterUpdate struct {
		*messageBase
	}

	// messageAfterDelete
	//
	// This type is auto-generated.
	messageAfterDelete struct {
		*messageBase
	}
)

// ResourceType returns "messaging:message"
//
// This function is auto-generated.
func (messageBase) ResourceType() string {
	return "messaging:message"
}

// EventType on messageOnManual returns "onManual"
//
// This function is auto-generated.
func (messageOnManual) EventType() string {
	return "onManual"
}

// EventType on messageBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (messageBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on messageBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (messageBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on messageBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (messageBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on messageAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (messageAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on messageAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (messageAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on messageAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (messageAfterDelete) EventType() string {
	return "afterDelete"
}

// MessageOnManual creates onManual for messaging:message resource
//
// This function is auto-generated.
func MessageOnManual(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageOnManual {
	return &messageOnManual{
		messageBase: &messageBase{
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageBeforeCreate creates beforeCreate for messaging:message resource
//
// This function is auto-generated.
func MessageBeforeCreate(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageBeforeCreate {
	return &messageBeforeCreate{
		messageBase: &messageBase{
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageBeforeUpdate creates beforeUpdate for messaging:message resource
//
// This function is auto-generated.
func MessageBeforeUpdate(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageBeforeUpdate {
	return &messageBeforeUpdate{
		messageBase: &messageBase{
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageBeforeDelete creates beforeDelete for messaging:message resource
//
// This function is auto-generated.
func MessageBeforeDelete(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageBeforeDelete {
	return &messageBeforeDelete{
		messageBase: &messageBase{
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageAfterCreate creates afterCreate for messaging:message resource
//
// This function is auto-generated.
func MessageAfterCreate(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageAfterCreate {
	return &messageAfterCreate{
		messageBase: &messageBase{
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageAfterUpdate creates afterUpdate for messaging:message resource
//
// This function is auto-generated.
func MessageAfterUpdate(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageAfterUpdate {
	return &messageAfterUpdate{
		messageBase: &messageBase{
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageAfterDelete creates afterDelete for messaging:message resource
//
// This function is auto-generated.
func MessageAfterDelete(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageAfterDelete {
	return &messageAfterDelete{
		messageBase: &messageBase{
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// SetMessage sets new message value
//
// This function is auto-generated.
func (res *messageBase) SetMessage(argMessage *types.Message) {
	res.message = argMessage
}

// Message returns message
//
// This function is auto-generated.
func (res messageBase) Message() *types.Message {
	return res.message
}

// OldMessage returns oldMessage
//
// This function is auto-generated.
func (res messageBase) OldMessage() *types.Message {
	return res.oldMessage
}

// SetChannel sets new channel value
//
// This function is auto-generated.
func (res *messageBase) SetChannel(argChannel *types.Channel) {
	res.channel = argChannel
}

// Channel returns channel
//
// This function is auto-generated.
func (res messageBase) Channel() *types.Channel {
	return res.channel
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *messageBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res messageBase) Invoker() auth.Identifiable {
	return res.invoker
}
