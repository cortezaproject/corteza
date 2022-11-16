package rest

import (
	"context"
	"github.com/cortezaproject/corteza/server/automation/rest/request"
	"github.com/cortezaproject/corteza/server/automation/service"
)

type (
	EventTypes struct {
		reg interface {
			Types() []string
		}
	}

	eventTypePayload struct {
		Set []eventTypeDef `json:"set"`
	}

	eventTypeDef struct {
		ResourceType string                   `json:"resourceType"`
		EventType    string                   `json:"eventType"`
		Properties   []eventTypePropertyDef   `json:"properties"`
		Constraints  []eventTypeConstraintDef `json:"constraints"`
	}

	eventTypePropertyDef struct {
		Name      string `json:"name"`
		Type      string `json:"type"`
		Immutable bool   `json:"immutable"`
	}

	eventTypeConstraintDef struct {
		Name string `json:"name"`
	}
)

func (EventTypes) New() *EventTypes {
	ctrl := &EventTypes{reg: service.Registry()}
	return ctrl
}

func (ctrl EventTypes) List(_ context.Context, _ *request.EventTypesList) (interface{}, error) {
	return eventTypePayload{Set: getEventTypeDefinitions()}, nil
}
