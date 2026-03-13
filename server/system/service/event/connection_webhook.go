package event

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	connectionWebhookEvent struct {
		resourceType           string
		connectionID           uint64
		configuredConnectionID uint64
		eventType              string
		vars                   map[string]any
	}
)

// ConnectionWebhookEvent constructs a new webhook event.
// resourceType is scoped to the parent Connection so one trigger covers all configured instances;
// configuredConnectionID identifies which instance fired and is passed as a var.
func ConnectionWebhookEvent(connectionID, configuredConnectionID uint64, eventType string, vars map[string]any) *connectionWebhookEvent {
	return &connectionWebhookEvent{
		resourceType:           fmt.Sprintf("corteza::system:connection-webhook/%d", connectionID),
		connectionID:           connectionID,
		configuredConnectionID: configuredConnectionID,
		eventType:              eventType,
		vars:                   vars,
	}
}

// ResourceType is scoped to the parent Connection, not the ConfiguredConnection.
// This means one trigger registration covers all configured instances.
func (e *connectionWebhookEvent) ResourceType() string {
	return e.resourceType
}

func (e *connectionWebhookEvent) EventType() string {
	return e.eventType
}

// Match supports constraints on "connectionID", "configuredConnectionID", "eventType", and any payload field.
func (e *connectionWebhookEvent) Match(c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "connectionID":
		return c.Match(fmt.Sprintf("%d", e.connectionID))
	case "configuredConnectionID":
		return c.Match(fmt.Sprintf("%d", e.configuredConnectionID))
	case "eventType":
		return c.Match(e.eventType)
	}

	if v, ok := e.vars[c.Name()]; ok {
		return c.Match(fmt.Sprintf("%v", v))
	}

	return true
}

// Vars returns the payload variables carried by this event.
func (e *connectionWebhookEvent) Vars() map[string]any {
	return e.vars
}

// EncodeVars encodes the event payload into expr.Vars so ng-automation
// can access webhook variables in the execution scope.
func (e *connectionWebhookEvent) EncodeVars() (out *expr.Vars, err error) {
	out = &expr.Vars{}
	for k, v := range e.vars {
		if err = out.Set(k, v); err != nil {
			return
		}
	}
	return
}
