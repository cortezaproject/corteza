package discovery

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

	Filter struct {
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

	// @fixme too many interface; too many Type assertions in CastToResourceActivity
	resDecoder interface {
		EventType() string
		ResourceType() string
		User() *systemTypes.User
		Namespace() *composeTypes.Namespace
		Module() *composeTypes.Module
		Record() *composeTypes.Record
	}
	systemDecoder interface {
		EventType() string
		ResourceType() string
		User() *systemTypes.User
	}
	composeDecoder interface {
		EventType() string
		ResourceType() string
		Namespace() *composeTypes.Namespace
		Module() *composeTypes.Module
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

func CastToResourceActivity(dec interface{}) (a *ResourceActivity, err error) {
	a = &ResourceActivity{
		ID:             id.Next(),
		Timestamp:      time.Now(),
		ResourceType:   dec.(resDecoder).ResourceType(),
		ResourceAction: ResourceAction(dec.(resDecoder).EventType()).String(),
	}
	setResourceID := func(ID uint64) {
		a.ResourceID = ID
	}

	switch a.ResourceType {
	case "system:user": // @fixme
		setResourceID(dec.(systemDecoder).User().ID)
	case (composeTypes.Namespace{}).LabelResourceKind():
		setResourceID(dec.(composeDecoder).Namespace().ID)
	case (composeTypes.Module{}).LabelResourceKind():
		setResourceID(dec.(composeDecoder).Module().ID)
	case (composeTypes.Record{}).LabelResourceKind():
		setResourceID(dec.(composeDecoder).Record().ID)
	default:
		return a, fmt.Errorf("unsupported resource type")
	}

	return
}
