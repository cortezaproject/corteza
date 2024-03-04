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
	// Component struct serves as a virtual resource type for the automation component
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

// RbacResource returns string representation of RBAC resource for Workflow by calling WorkflowRbacResource fn
//
// RBAC resource is in the corteza::automation:workflow/... format
//
// This function is auto-generated
func (r Workflow) RbacResource() string {
	return WorkflowRbacResource(r.ID)
}

// WorkflowRbacResource returns string representation of RBAC resource for Workflow
//
// RBAC resource is in the corteza::automation:workflow/... format
//
// This function is auto-generated
func WorkflowRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{WorkflowResourceType}
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

	out := fmt.Sprintf(WorkflowRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func WorkflowRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza::automation/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza::automation/ format
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
