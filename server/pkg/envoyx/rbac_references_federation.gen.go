package envoyx

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza/server/federation/types"
)

// FederationNodeRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func FederationNodeRbacReferences(node string) (res *Ref, pp []*Ref, err error) {
	if node != "*" {
		res = &Ref{ResourceType: types.NodeResourceType, Identifiers: MakeIdentifiers(node)}
	}

	return
}

// FederationExposedModuleRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func FederationExposedModuleRbacReferences(nodeID string, exposedModule string) (res *Ref, pp []*Ref, err error) {
	if nodeID != "*" {
		pp = append(pp, &Ref{ResourceType: types.NodeResourceType, Identifiers: MakeIdentifiers(nodeID)})
	}
	if exposedModule != "*" {
		res = &Ref{ResourceType: types.ExposedModuleResourceType, Identifiers: MakeIdentifiers(exposedModule)}
	}

	return
}

// FederationSharedModuleRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func FederationSharedModuleRbacReferences(nodeID string, sharedModule string) (res *Ref, pp []*Ref, err error) {
	if nodeID != "*" {
		pp = append(pp, &Ref{ResourceType: types.NodeResourceType, Identifiers: MakeIdentifiers(nodeID)})
	}
	if sharedModule != "*" {
		res = &Ref{ResourceType: types.SharedModuleResourceType, Identifiers: MakeIdentifiers(sharedModule)}
	}

	return
}
