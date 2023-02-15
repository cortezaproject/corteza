package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza/server/pkg/cast2"
)

func (r Workflow) GetID() uint64 { return r.ID }

func (r *Workflow) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "enabled", "Enabled":
		return r.Enabled, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "keepSessions", "KeepSessions":
		return r.KeepSessions, nil
	case "ownedBy", "OwnedBy":
		return r.OwnedBy, nil
	case "runAs", "RunAs":
		return r.RunAs, nil
	case "trace", "Trace":
		return r.Trace, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Workflow) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "enabled", "Enabled":
		return cast2.Bool(value, &r.Enabled)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "keepSessions", "KeepSessions":
		return cast2.Int(value, &r.KeepSessions)
	case "ownedBy", "OwnedBy":
		return cast2.Uint64(value, &r.OwnedBy)
	case "runAs", "RunAs":
		return cast2.Uint64(value, &r.RunAs)
	case "trace", "Trace":
		return cast2.Bool(value, &r.Trace)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r Session) GetID() uint64 { return r.ID }

func (r *Session) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "completedAt", "CompletedAt":
		return r.CompletedAt, nil
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "error", "Error":
		return r.Error, nil
	case "eventType", "EventType":
		return r.EventType, nil
	case "id", "ID":
		return r.ID, nil
	case "purgeAt", "PurgeAt":
		return r.PurgeAt, nil
	case "resourceType", "ResourceType":
		return r.ResourceType, nil
	case "suspendedAt", "SuspendedAt":
		return r.SuspendedAt, nil
	case "workflowID", "WorkflowID":
		return r.WorkflowID, nil

	}
	return nil, nil
}

func (r *Session) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "completedAt", "CompletedAt":
		return cast2.TimePtr(value, &r.CompletedAt)
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "error", "Error":
		return cast2.String(value, &r.Error)
	case "eventType", "EventType":
		return cast2.String(value, &r.EventType)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "purgeAt", "PurgeAt":
		return cast2.TimePtr(value, &r.PurgeAt)
	case "resourceType", "ResourceType":
		return cast2.String(value, &r.ResourceType)
	case "suspendedAt", "SuspendedAt":
		return cast2.TimePtr(value, &r.SuspendedAt)
	case "workflowID", "WorkflowID":
		return cast2.Uint64(value, &r.WorkflowID)

	}
	return nil
}

func (r Trigger) GetID() uint64 { return r.ID }

func (r *Trigger) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "enabled", "Enabled":
		return r.Enabled, nil
	case "eventType", "EventType":
		return r.EventType, nil
	case "id", "ID":
		return r.ID, nil
	case "ownedBy", "OwnedBy":
		return r.OwnedBy, nil
	case "resourceType", "ResourceType":
		return r.ResourceType, nil
	case "stepID", "StepID":
		return r.StepID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil
	case "workflowID", "WorkflowID":
		return r.WorkflowID, nil

	}
	return nil, nil
}

func (r *Trigger) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "enabled", "Enabled":
		return cast2.Bool(value, &r.Enabled)
	case "eventType", "EventType":
		return cast2.String(value, &r.EventType)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "ownedBy", "OwnedBy":
		return cast2.Uint64(value, &r.OwnedBy)
	case "resourceType", "ResourceType":
		return cast2.String(value, &r.ResourceType)
	case "stepID", "StepID":
		return cast2.Uint64(value, &r.StepID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)
	case "workflowID", "WorkflowID":
		return cast2.Uint64(value, &r.WorkflowID)

	}
	return nil
}
