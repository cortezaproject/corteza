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

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}
