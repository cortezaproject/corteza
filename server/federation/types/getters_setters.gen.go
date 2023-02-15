package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza/server/pkg/cast2"
)

func (r Node) GetID() uint64 { return r.ID }

func (r *Node) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "authToken", "AuthToken":
		return r.AuthToken, nil
	case "baseURL", "BaseURL":
		return r.BaseURL, nil
	case "contact", "Contact":
		return r.Contact, nil
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "id", "ID":
		return r.ID, nil
	case "name", "Name":
		return r.Name, nil
	case "pairToken", "PairToken":
		return r.PairToken, nil
	case "sharedNodeID", "SharedNodeID":
		return r.SharedNodeID, nil
	case "status", "Status":
		return r.Status, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Node) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "authToken", "AuthToken":
		return cast2.String(value, &r.AuthToken)
	case "baseURL", "BaseURL":
		return cast2.String(value, &r.BaseURL)
	case "contact", "Contact":
		return cast2.String(value, &r.Contact)
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "pairToken", "PairToken":
		return cast2.String(value, &r.PairToken)
	case "sharedNodeID", "SharedNodeID":
		return cast2.Uint64(value, &r.SharedNodeID)
	case "status", "Status":
		return cast2.String(value, &r.Status)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r NodeSync) GetID() uint64 {
	// The resource does not define an ID field
	return 0
}

func (r *NodeSync) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "moduleID", "ModuleID":
		return r.ModuleID, nil
	case "nodeID", "NodeID":
		return r.NodeID, nil
	case "syncStatus", "SyncStatus":
		return r.SyncStatus, nil
	case "syncType", "SyncType":
		return r.SyncType, nil
	case "timeOfAction", "TimeOfAction":
		return r.TimeOfAction, nil

	}
	return nil, nil
}

func (r *NodeSync) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "moduleID", "ModuleID":
		return cast2.Uint64(value, &r.ModuleID)
	case "nodeID", "NodeID":
		return cast2.Uint64(value, &r.NodeID)
	case "syncStatus", "SyncStatus":
		return cast2.String(value, &r.SyncStatus)
	case "syncType", "SyncType":
		return cast2.String(value, &r.SyncType)
	case "timeOfAction", "TimeOfAction":
		return cast2.Time(value, &r.TimeOfAction)

	}
	return nil
}

func (r ExposedModule) GetID() uint64 { return r.ID }

func (r *ExposedModule) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "composeModuleID", "ComposeModuleID":
		return r.ComposeModuleID, nil
	case "composeNamespaceID", "ComposeNamespaceID":
		return r.ComposeNamespaceID, nil
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "name", "Name":
		return r.Name, nil
	case "nodeID", "NodeID":
		return r.NodeID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ExposedModule) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "composeModuleID", "ComposeModuleID":
		return cast2.Uint64(value, &r.ComposeModuleID)
	case "composeNamespaceID", "ComposeNamespaceID":
		return cast2.Uint64(value, &r.ComposeNamespaceID)
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "nodeID", "NodeID":
		return cast2.Uint64(value, &r.NodeID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r SharedModule) GetID() uint64 { return r.ID }

func (r *SharedModule) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "externalFederationModuleID", "ExternalFederationModuleID":
		return r.ExternalFederationModuleID, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "name", "Name":
		return r.Name, nil
	case "nodeID", "NodeID":
		return r.NodeID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *SharedModule) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "externalFederationModuleID", "ExternalFederationModuleID":
		return cast2.Uint64(value, &r.ExternalFederationModuleID)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "nodeID", "NodeID":
		return cast2.Uint64(value, &r.NodeID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r ModuleMapping) GetID() uint64 {
	// The resource does not define an ID field
	return 0
}

func (r *ModuleMapping) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "composeModuleID", "ComposeModuleID":
		return r.ComposeModuleID, nil
	case "composeNamespaceID", "ComposeNamespaceID":
		return r.ComposeNamespaceID, nil
	case "federationModuleID", "FederationModuleID":
		return r.FederationModuleID, nil
	case "nodeID", "NodeID":
		return r.NodeID, nil

	}
	return nil, nil
}

func (r *ModuleMapping) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "composeModuleID", "ComposeModuleID":
		return cast2.Uint64(value, &r.ComposeModuleID)
	case "composeNamespaceID", "ComposeNamespaceID":
		return cast2.Uint64(value, &r.ComposeNamespaceID)
	case "federationModuleID", "FederationModuleID":
		return cast2.Uint64(value, &r.FederationModuleID)
	case "nodeID", "NodeID":
		return cast2.Uint64(value, &r.NodeID)

	}
	return nil
}
