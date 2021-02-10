package event

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// messaging/service/event/events.yaml

import (
	"encoding/json"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

// dummy placing to simplify import generation logic
var _ = json.NewEncoder

type (

	// messagingBase
	//
	// This type is auto-generated.
	messagingBase struct {
		immutable bool
		invoker   auth.Identifiable
	}

	// messagingOnManual
	//
	// This type is auto-generated.
	messagingOnManual struct {
		*messagingBase
	}

	// messagingOnInterval
	//
	// This type is auto-generated.
	messagingOnInterval struct {
		*messagingBase
	}

	// messagingOnTimestamp
	//
	// This type is auto-generated.
	messagingOnTimestamp struct {
		*messagingBase
	}

	// channelBase
	//
	// This type is auto-generated.
	channelBase struct {
		immutable  bool
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

	// channelMemberBase
	//
	// This type is auto-generated.
	channelMemberBase struct {
		immutable bool
		member    *types.ChannelMember
		channel   *types.Channel
		invoker   auth.Identifiable
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

	// commandBase
	//
	// This type is auto-generated.
	commandBase struct {
		immutable bool
		command   *types.Command
		channel   *types.Channel
		invoker   auth.Identifiable
	}

	// commandOnInvoke
	//
	// This type is auto-generated.
	commandOnInvoke struct {
		*commandBase
	}

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

// ResourceType returns "messaging"
//
// This function is auto-generated.
func (messagingBase) ResourceType() string {
	return "messaging"
}

// EventType on messagingOnManual returns "onManual"
//
// This function is auto-generated.
func (messagingOnManual) EventType() string {
	return "onManual"
}

// EventType on messagingOnInterval returns "onInterval"
//
// This function is auto-generated.
func (messagingOnInterval) EventType() string {
	return "onInterval"
}

// EventType on messagingOnTimestamp returns "onTimestamp"
//
// This function is auto-generated.
func (messagingOnTimestamp) EventType() string {
	return "onTimestamp"
}

// MessagingOnManual creates onManual for messaging resource
//
// This function is auto-generated.
func MessagingOnManual() *messagingOnManual {
	return &messagingOnManual{
		messagingBase: &messagingBase{
			immutable: false,
		},
	}
}

// MessagingOnManualImmutable creates onManual for messaging resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessagingOnManualImmutable() *messagingOnManual {
	return &messagingOnManual{
		messagingBase: &messagingBase{
			immutable: true,
		},
	}
}

// MessagingOnInterval creates onInterval for messaging resource
//
// This function is auto-generated.
func MessagingOnInterval() *messagingOnInterval {
	return &messagingOnInterval{
		messagingBase: &messagingBase{
			immutable: false,
		},
	}
}

// MessagingOnIntervalImmutable creates onInterval for messaging resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessagingOnIntervalImmutable() *messagingOnInterval {
	return &messagingOnInterval{
		messagingBase: &messagingBase{
			immutable: true,
		},
	}
}

// MessagingOnTimestamp creates onTimestamp for messaging resource
//
// This function is auto-generated.
func MessagingOnTimestamp() *messagingOnTimestamp {
	return &messagingOnTimestamp{
		messagingBase: &messagingBase{
			immutable: false,
		},
	}
}

// MessagingOnTimestampImmutable creates onTimestamp for messaging resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MessagingOnTimestampImmutable() *messagingOnTimestamp {
	return &messagingOnTimestamp{
		messagingBase: &messagingBase{
			immutable: true,
		},
	}
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *messagingBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res messagingBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res messagingBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res messagingBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *messagingBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
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

func (res *messagingBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

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
			immutable:  false,
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelOnManualImmutable creates onManual for messaging:channel resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelOnManualImmutable(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelOnManual {
	return &channelOnManual{
		channelBase: &channelBase{
			immutable:  true,
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
			immutable:  false,
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelBeforeCreateImmutable creates beforeCreate for messaging:channel resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelBeforeCreateImmutable(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelBeforeCreate {
	return &channelBeforeCreate{
		channelBase: &channelBase{
			immutable:  true,
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
			immutable:  false,
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelBeforeUpdateImmutable creates beforeUpdate for messaging:channel resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelBeforeUpdateImmutable(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelBeforeUpdate {
	return &channelBeforeUpdate{
		channelBase: &channelBase{
			immutable:  true,
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
			immutable:  false,
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelBeforeDeleteImmutable creates beforeDelete for messaging:channel resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelBeforeDeleteImmutable(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelBeforeDelete {
	return &channelBeforeDelete{
		channelBase: &channelBase{
			immutable:  true,
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
			immutable:  false,
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelAfterCreateImmutable creates afterCreate for messaging:channel resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelAfterCreateImmutable(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelAfterCreate {
	return &channelAfterCreate{
		channelBase: &channelBase{
			immutable:  true,
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
			immutable:  false,
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelAfterUpdateImmutable creates afterUpdate for messaging:channel resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelAfterUpdateImmutable(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelAfterUpdate {
	return &channelAfterUpdate{
		channelBase: &channelBase{
			immutable:  true,
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
			immutable:  false,
			channel:    argChannel,
			oldChannel: argOldChannel,
		},
	}
}

// ChannelAfterDeleteImmutable creates afterDelete for messaging:channel resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelAfterDeleteImmutable(
	argChannel *types.Channel,
	argOldChannel *types.Channel,
) *channelAfterDelete {
	return &channelAfterDelete{
		channelBase: &channelBase{
			immutable:  true,
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res channelBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["channel"], err = json.Marshal(res.channel); err != nil {
		return nil, err
	}

	if args["oldChannel"], err = json.Marshal(res.oldChannel); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res channelBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for *types.Channel

	// Could not found expression-type counterpart for *types.Channel

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *channelBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.channel != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.channel); err != nil {
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

	// Do not decode oldChannel; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}

func (res *channelBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for *types.Channel
	// oldChannel marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

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
			immutable: false,
			member:    argMember,
			channel:   argChannel,
		},
	}
}

// ChannelMemberBeforeJoinImmutable creates beforeJoin for messaging:channel:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelMemberBeforeJoinImmutable(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberBeforeJoin {
	return &channelMemberBeforeJoin{
		channelMemberBase: &channelMemberBase{
			immutable: true,
			member:    argMember,
			channel:   argChannel,
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
			immutable: false,
			member:    argMember,
			channel:   argChannel,
		},
	}
}

// ChannelMemberBeforePartImmutable creates beforePart for messaging:channel:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelMemberBeforePartImmutable(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberBeforePart {
	return &channelMemberBeforePart{
		channelMemberBase: &channelMemberBase{
			immutable: true,
			member:    argMember,
			channel:   argChannel,
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
			immutable: false,
			member:    argMember,
			channel:   argChannel,
		},
	}
}

// ChannelMemberBeforeAddImmutable creates beforeAdd for messaging:channel:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelMemberBeforeAddImmutable(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberBeforeAdd {
	return &channelMemberBeforeAdd{
		channelMemberBase: &channelMemberBase{
			immutable: true,
			member:    argMember,
			channel:   argChannel,
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
			immutable: false,
			member:    argMember,
			channel:   argChannel,
		},
	}
}

// ChannelMemberBeforeRemoveImmutable creates beforeRemove for messaging:channel:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelMemberBeforeRemoveImmutable(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberBeforeRemove {
	return &channelMemberBeforeRemove{
		channelMemberBase: &channelMemberBase{
			immutable: true,
			member:    argMember,
			channel:   argChannel,
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
			immutable: false,
			member:    argMember,
			channel:   argChannel,
		},
	}
}

// ChannelMemberAfterJoinImmutable creates afterJoin for messaging:channel:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelMemberAfterJoinImmutable(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberAfterJoin {
	return &channelMemberAfterJoin{
		channelMemberBase: &channelMemberBase{
			immutable: true,
			member:    argMember,
			channel:   argChannel,
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
			immutable: false,
			member:    argMember,
			channel:   argChannel,
		},
	}
}

// ChannelMemberAfterPartImmutable creates afterPart for messaging:channel:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelMemberAfterPartImmutable(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberAfterPart {
	return &channelMemberAfterPart{
		channelMemberBase: &channelMemberBase{
			immutable: true,
			member:    argMember,
			channel:   argChannel,
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
			immutable: false,
			member:    argMember,
			channel:   argChannel,
		},
	}
}

// ChannelMemberAfterAddImmutable creates afterAdd for messaging:channel:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelMemberAfterAddImmutable(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberAfterAdd {
	return &channelMemberAfterAdd{
		channelMemberBase: &channelMemberBase{
			immutable: true,
			member:    argMember,
			channel:   argChannel,
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
			immutable: false,
			member:    argMember,
			channel:   argChannel,
		},
	}
}

// ChannelMemberAfterRemoveImmutable creates afterRemove for messaging:channel:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ChannelMemberAfterRemoveImmutable(
	argMember *types.ChannelMember,
	argChannel *types.Channel,
) *channelMemberAfterRemove {
	return &channelMemberAfterRemove{
		channelMemberBase: &channelMemberBase{
			immutable: true,
			member:    argMember,
			channel:   argChannel,
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res channelMemberBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["member"], err = json.Marshal(res.member); err != nil {
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

// Encode internal data to be passed as event params & arguments to workflow
func (res channelMemberBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for *types.ChannelMember

	// Could not found expression-type counterpart for *types.Channel

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *channelMemberBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.member != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.member); err != nil {
				return
			}
		}
	}

	if res.member != nil {
		if r, ok := results["member"]; ok {
			if err = json.Unmarshal(r, res.member); err != nil {
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

func (res *channelMemberBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for *types.ChannelMember
	// Could not find expression-type counterpart for *types.Channel
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

// ResourceType returns "messaging:command"
//
// This function is auto-generated.
func (commandBase) ResourceType() string {
	return "messaging:command"
}

// EventType on commandOnInvoke returns "onInvoke"
//
// This function is auto-generated.
func (commandOnInvoke) EventType() string {
	return "onInvoke"
}

// CommandOnInvoke creates onInvoke for messaging:command resource
//
// This function is auto-generated.
func CommandOnInvoke(
	argCommand *types.Command,
	argChannel *types.Channel,
) *commandOnInvoke {
	return &commandOnInvoke{
		commandBase: &commandBase{
			immutable: false,
			command:   argCommand,
			channel:   argChannel,
		},
	}
}

// CommandOnInvokeImmutable creates onInvoke for messaging:command resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func CommandOnInvokeImmutable(
	argCommand *types.Command,
	argChannel *types.Channel,
) *commandOnInvoke {
	return &commandOnInvoke{
		commandBase: &commandBase{
			immutable: true,
			command:   argCommand,
			channel:   argChannel,
		},
	}
}

// Command returns command
//
// This function is auto-generated.
func (res commandBase) Command() *types.Command {
	return res.command
}

// Channel returns channel
//
// This function is auto-generated.
func (res commandBase) Channel() *types.Channel {
	return res.channel
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *commandBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res commandBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res commandBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["command"], err = json.Marshal(res.command); err != nil {
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

// Encode internal data to be passed as event params & arguments to workflow
func (res commandBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for *types.Command

	// Could not found expression-type counterpart for *types.Channel

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *commandBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.command != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.command); err != nil {
				return
			}
		}
	}

	// Do not decode command; marked as immutable

	// Do not decode channel; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}

func (res *commandBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// command marked as immutable
	// channel marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

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

// Encode internal data to be passed as event params & arguments to workflow
func (res messageBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for *types.Message

	// Could not found expression-type counterpart for *types.Message

	// Could not found expression-type counterpart for *types.Channel

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
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

	// Do not decode oldMessage; marked as immutable

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

func (res *messageBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for *types.Message
	// oldMessage marked as immutable
	// Could not find expression-type counterpart for *types.Channel
	// Could not find expression-type counterpart for auth.Identifiable

	return
}
