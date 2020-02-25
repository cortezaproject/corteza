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
	"encoding/json"

	"github.com/cortezaproject/corteza-server/messaging/types"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// messageBase
	//
	// This type is auto-generated.
	messageBase struct {
		immutable  bool
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
			immutable:  false,
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageOnManualImmutable creates onManual for messaging:message resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessageOnManualImmutable(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageOnManual {
	return &messageOnManual{
		messageBase: &messageBase{
			immutable:  true,
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
			immutable:  false,
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageBeforeCreateImmutable creates beforeCreate for messaging:message resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessageBeforeCreateImmutable(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageBeforeCreate {
	return &messageBeforeCreate{
		messageBase: &messageBase{
			immutable:  true,
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
			immutable:  false,
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageBeforeUpdateImmutable creates beforeUpdate for messaging:message resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessageBeforeUpdateImmutable(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageBeforeUpdate {
	return &messageBeforeUpdate{
		messageBase: &messageBase{
			immutable:  true,
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
			immutable:  false,
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageBeforeDeleteImmutable creates beforeDelete for messaging:message resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessageBeforeDeleteImmutable(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageBeforeDelete {
	return &messageBeforeDelete{
		messageBase: &messageBase{
			immutable:  true,
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
			immutable:  false,
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageAfterCreateImmutable creates afterCreate for messaging:message resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessageAfterCreateImmutable(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageAfterCreate {
	return &messageAfterCreate{
		messageBase: &messageBase{
			immutable:  true,
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
			immutable:  false,
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageAfterUpdateImmutable creates afterUpdate for messaging:message resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessageAfterUpdateImmutable(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageAfterUpdate {
	return &messageAfterUpdate{
		messageBase: &messageBase{
			immutable:  true,
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
			immutable:  false,
			message:    argMessage,
			oldMessage: argOldMessage,
			channel:    argChannel,
		},
	}
}

// MessageAfterDeleteImmutable creates afterDelete for messaging:message resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessageAfterDeleteImmutable(
	argMessage *types.Message,
	argOldMessage *types.Message,
	argChannel *types.Channel,
) *messageAfterDelete {
	return &messageAfterDelete{
		messageBase: &messageBase{
			immutable:  true,
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res messageBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["message"], err = json.Marshal(res.message); err != nil {
		return nil, err
	}

	if args["oldMessage"], err = json.Marshal(res.oldMessage); err != nil {
		return nil, err
	}

	if args["channel"], err = json.Marshal(res.channel); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Decode return values from Corredor script into struct props
func (res *messageBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.message != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.message); err != nil {
				return
			}
		}
	}

	if res.message != nil {
		if r, ok := results["message"]; ok {
			if err = json.Unmarshal(r, res.message); err != nil {
				return
			}
		}
	}

	if res.channel != nil {
		if r, ok := results["channel"]; ok {
			if err = json.Unmarshal(r, res.channel); err != nil {
				return
			}
		}
	}

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}
