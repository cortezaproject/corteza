package types

import (
	"encoding/json"
	"fmt"
	"time"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/jmoiron/sqlx/types"
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

		// Meta of the related resources
		Meta types.JSONText `json:"meta"`
	}

	ResourceActivityFilter struct {
		FromTimestamp *time.Time `json:"from"`
		ToTimestamp   *time.Time `json:"to"`

		// Check fn is called by store backend for each resource found function can
		// modify the activity and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*ResourceActivity) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	ResourceActivityMeta struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		ModuleID    uint64 `json:"moduleID,string"`
	}

	ResDecoder interface {
		EventType() string
		ResourceType() string
	}
	userDecoder interface {
		ResDecoder
		User() *systemTypes.User
		OldUser() *systemTypes.User
	}
	nsDecoder interface {
		ResDecoder
		Namespace() *composeTypes.Namespace
		OldNamespace() *composeTypes.Namespace
	}
	mDecoder interface {
		ResDecoder
		Module() *composeTypes.Module
		OldModule() *composeTypes.Module
	}
	recDecoder interface {
		ResDecoder
		Record() *composeTypes.Record
		OldRecord() *composeTypes.Record
	}
)

const (
	AfterCreate ResourceAction = "afterCreate"
	AfterUpdate ResourceAction = "afterUpdate"
	AfterDelete ResourceAction = "afterDelete"

	ResourceActivityResourceType = "corteza::generic:resource-activity"
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
	setMeta := func(nsID, mID uint64) error {
		var meta ResourceActivityMeta
		if nsID > 0 {
			meta.NamespaceID = nsID
		}
		if mID > 0 {
			meta.ModuleID = mID
		}

		a.Meta, err = json.Marshal(meta)
		if err != nil {
			return err
		}

		return nil
	}

	switch a.ResourceType {
	case "system:user": // @todo system/service/service.go#134
		if v, ok := dec.(userDecoder); ok {
			user := v.User()
			// fallback to OldUser for afterDelete event
			if user == nil {
				user = v.OldUser()
			}
			if user != nil {
				setResourceID(user.ID)
			}
		}
	case (composeTypes.Namespace{}).LabelResourceKind():
		if v, ok := dec.(nsDecoder); ok {
			ns := v.Namespace()
			// fallback to OldNamespace for afterDelete event
			if ns == nil {
				ns = v.OldNamespace()
			}
			if ns != nil {
				setResourceID(ns.ID)
			}
		}
	case (composeTypes.Module{}).LabelResourceKind():
		if v, ok := dec.(mDecoder); ok {
			mod := v.Module()
			// fallback to OldModule for afterDelete event
			if mod == nil {
				mod = v.OldModule()
			}
			if mod != nil {
				setResourceID(mod.ID)
				err = setMeta(mod.NamespaceID, 0)
				if err != nil {
					return
				}
			}
		}
	case "compose:record":
		if v, ok := dec.(recDecoder); ok {
			rec := v.Record()
			// fallback to OldRecord for afterDelete event
			if rec == nil {
				rec = v.OldRecord()
			}
			if rec != nil {
				setResourceID(rec.ID)
				err = setMeta(rec.NamespaceID, rec.ModuleID)
				if err != nil {
					return
				}
			}
		}
	default:
		return a, fmt.Errorf("unsupported resource type")
	}

	return
}
