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
	// commandBase
	//
	// This type is auto-generated.
	commandBase struct {
		command *types.Command
		channel *types.Channel
		invoker auth.Identifiable
	}

	// commandOnInvoke
	//
	// This type is auto-generated.
	commandOnInvoke struct {
		*commandBase
	}
)

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
			command: argCommand,
			channel: argChannel,
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
