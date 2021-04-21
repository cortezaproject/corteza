package types

import (
	"fmt"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"time"
)

type (
	// ResourceAction determinants resource action event
	ResourceAction string

	// ResourceActivity Standardized data structure for resource activity
	ResourceActivity struct {
		ID uint64 `json:"activityID,string"`

		// ResourceID of the corteza resource
		ResourceID uint64 `json:"resourceID,string"`

		// ResourceType
		ResourceType string `json:"resourceType"`

		// ResourceAction Type of action
		ResourceAction string `json:"resourceAction"`

		// Timestamp of the raised event
		Timestamp time.Time `json:"timestamp"`
	}

	ResourceActivityFilter struct {
		FromTimestamp *time.Time `json:"from"`
		ToTimestamp   *time.Time `json:"to"`

		Limit uint `json:"limit"`

		// Check fn is called by store backend for each resource found function can
		// modify the activity and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*ResourceActivity) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	ResDecoder interface {
		EventType() string
		ResourceType() string
	}
	userDecoder interface {
		ResDecoder
		User() *systemTypes.User
	}
	nsDecoder interface {
		ResDecoder
		Namespace() *composeTypes.Namespace
	}
	mDecoder interface {
		ResDecoder
		Module() *composeTypes.Module
	}
	recDecoder interface {
		ResDecoder
		Record() *composeTypes.Record
	}
)

const (
	AfterCreate ResourceAction = "afterCreate"
	AfterUpdate ResourceAction = "afterUpdate"
	AfterDelete ResourceAction = "afterDelete"
)

func (s ResourceAction) String() string {
	switch s {
	case AfterCreate:
		return "create"
	case AfterUpdate:
		return "update"
	case AfterDelete:
		return "delete"
	}

	return ""
}

func CastToResourceActivity(dec ResDecoder) (a *ResourceActivity, err error) {
	a = &ResourceActivity{
		ID:             id.Next(),
		Timestamp:      time.Now(),
		ResourceType:   dec.ResourceType(),
		ResourceAction: ResourceAction(dec.EventType()).String(),
	}

	setResourceID := func(ID uint64) {
		a.ResourceID = ID
	}

	switch a.ResourceType {
	case "system:user": // @todo system/service/service.go#134
		if v, ok := dec.(userDecoder); ok {
			if v.User() != nil {
				setResourceID(v.User().ID)
			}
		}
	case (composeTypes.Namespace{}).LabelResourceKind():
		if v, ok := dec.(nsDecoder); ok {
			if v.Namespace() != nil {
				setResourceID(v.Namespace().ID)
			}
		}
	case (composeTypes.Module{}).LabelResourceKind():
		if v, ok := dec.(mDecoder); ok {
			if v.Module() != nil {
				setResourceID(v.Module().ID)
			}
		}
	case (composeTypes.Record{}).LabelResourceKind():
		if v, ok := dec.(recDecoder); ok {
			if v.Record() != nil {
				setResourceID(v.Record().ID)
			}
		}
	default:
		return a, fmt.Errorf("unsupported resource type")
	}

	return
}
