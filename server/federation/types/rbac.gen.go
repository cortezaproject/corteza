package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/ds"
	"strconv"
)

type (
	// Component struct serves as a virtual resource type for the federation component
	//
	// This struct is auto-generated
	Component struct{}

	indexWrapper struct {
		resource string
		counter  uint
	}
)

var (
	_ = fmt.Printf
	_ = strconv.FormatUint
)

var (
	resourceIndex = ds.Trie[uint64, *indexWrapper]()
)

var (
	resourceIndexMaxSize = 1000
)

// RbacResource returns string representation of RBAC resource for Node by calling NodeRbacResource fn
//
// RBAC resource is in the corteza::federation:node/... format
//
// This function is auto-generated
func (r Node) RbacResource() string {
	return NodeRbacResource(r.ID)
}

// NodeRbacResource returns string representation of RBAC resource for Node
//
// RBAC resource is in the corteza::federation:node/... format
//
// This function is auto-generated
func NodeRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{NodeResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	// Remove the least used ones
	// @todo for now just rebuild the index, later do this properly
	if resourceIndex.Size+1 > resourceIndexMaxSize {
		resourceIndex = ds.Trie[uint64, *indexWrapper]()
	}

	out := fmt.Sprintf(NodeRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func NodeRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ExposedModule by calling ExposedModuleRbacResource fn
//
// RBAC resource is in the corteza::federation:exposed-module/... format
//
// This function is auto-generated
func (r ExposedModule) RbacResource() string {
	return ExposedModuleRbacResource(r.NodeID, r.ID)
}

// ExposedModuleRbacResource returns string representation of RBAC resource for ExposedModule
//
// RBAC resource is in the corteza::federation:exposed-module/... format
//
// This function is auto-generated
func ExposedModuleRbacResource(nodeID uint64, id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, nodeID, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{ExposedModuleResourceType}
	if nodeID != 0 {
		cpts = append(cpts, strconv.FormatUint(nodeID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	// Remove the least used ones
	// @todo for now just rebuild the index, later do this properly
	if resourceIndex.Size+1 > resourceIndexMaxSize {
		resourceIndex = ds.Trie[uint64, *indexWrapper]()
	}

	out := fmt.Sprintf(ExposedModuleRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, nodeID, id)

	return out

}

func ExposedModuleRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for SharedModule by calling SharedModuleRbacResource fn
//
// RBAC resource is in the corteza::federation:shared-module/... format
//
// This function is auto-generated
func (r SharedModule) RbacResource() string {
	return SharedModuleRbacResource(r.NodeID, r.ID)
}

// SharedModuleRbacResource returns string representation of RBAC resource for SharedModule
//
// RBAC resource is in the corteza::federation:shared-module/... format
//
// This function is auto-generated
func SharedModuleRbacResource(nodeID uint64, id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, nodeID, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{SharedModuleResourceType}
	if nodeID != 0 {
		cpts = append(cpts, strconv.FormatUint(nodeID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	// Remove the least used ones
	// @todo for now just rebuild the index, later do this properly
	if resourceIndex.Size+1 > resourceIndexMaxSize {
		resourceIndex = ds.Trie[uint64, *indexWrapper]()
	}

	out := fmt.Sprintf(SharedModuleRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, nodeID, id)

	return out

}

func SharedModuleRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza::federation/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza::federation/ format
//
// This function is auto-generated
func ComponentRbacResource() string {
	return ComponentResourceType + "/"

}

func ComponentRbacResourceTpl() string {
	return "%s"
}

func merge(a, b *indexWrapper) *indexWrapper {
	a.counter += b.counter
	return a
}
