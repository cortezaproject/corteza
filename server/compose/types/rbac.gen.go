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
	// Component struct serves as a virtual resource type for the compose component
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

// RbacResource returns string representation of RBAC resource for Chart by calling ChartRbacResource fn
//
// RBAC resource is in the corteza::compose:chart/... format
//
// This function is auto-generated
func (r Chart) RbacResource() string {
	return ChartRbacResource(r.NamespaceID, r.ID)
}

// ChartRbacResource returns string representation of RBAC resource for Chart
//
// RBAC resource is in the corteza::compose:chart/... format
//
// This function is auto-generated
func ChartRbacResource(namespaceID uint64, id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{ChartResourceType}
	if namespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(namespaceID, 10))
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

	out := fmt.Sprintf(ChartRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, namespaceID, id)

	return out

}

func ChartRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Module by calling ModuleRbacResource fn
//
// RBAC resource is in the corteza::compose:module/... format
//
// This function is auto-generated
func (r Module) RbacResource() string {
	return ModuleRbacResource(r.NamespaceID, r.ID)
}

// ModuleRbacResource returns string representation of RBAC resource for Module
//
// RBAC resource is in the corteza::compose:module/... format
//
// This function is auto-generated
func ModuleRbacResource(namespaceID uint64, id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{ModuleResourceType}
	if namespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(namespaceID, 10))
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

	out := fmt.Sprintf(ModuleRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, namespaceID, id)

	return out

}

func ModuleRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for ModuleField by calling ModuleFieldRbacResource fn
//
// RBAC resource is in the corteza::compose:module-field/... format
//
// This function is auto-generated
func (r ModuleField) RbacResource() string {
	return ModuleFieldRbacResource(r.NamespaceID, r.ModuleID, r.ID)
}

// ModuleFieldRbacResource returns string representation of RBAC resource for ModuleField
//
// RBAC resource is in the corteza::compose:module-field/... format
//
// This function is auto-generated
func ModuleFieldRbacResource(namespaceID uint64, moduleID uint64, id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, moduleID, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{ModuleFieldResourceType}
	if namespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(namespaceID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if moduleID != 0 {
		cpts = append(cpts, strconv.FormatUint(moduleID, 10))
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

	out := fmt.Sprintf(ModuleFieldRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, namespaceID, moduleID, id)

	return out

}

func ModuleFieldRbacResourceTpl() string {
	return "%s/%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Namespace by calling NamespaceRbacResource fn
//
// RBAC resource is in the corteza::compose:namespace/... format
//
// This function is auto-generated
func (r Namespace) RbacResource() string {
	return NamespaceRbacResource(r.ID)
}

// NamespaceRbacResource returns string representation of RBAC resource for Namespace
//
// RBAC resource is in the corteza::compose:namespace/... format
//
// This function is auto-generated
func NamespaceRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{NamespaceResourceType}
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

	out := fmt.Sprintf(NamespaceRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func NamespaceRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Page by calling PageRbacResource fn
//
// RBAC resource is in the corteza::compose:page/... format
//
// This function is auto-generated
func (r Page) RbacResource() string {
	return PageRbacResource(r.NamespaceID, r.ID)
}

// PageRbacResource returns string representation of RBAC resource for Page
//
// RBAC resource is in the corteza::compose:page/... format
//
// This function is auto-generated
func PageRbacResource(namespaceID uint64, id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{PageResourceType}
	if namespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(namespaceID, 10))
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

	out := fmt.Sprintf(PageRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, namespaceID, id)

	return out

}

func PageRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for PageLayout by calling PageLayoutRbacResource fn
//
// RBAC resource is in the corteza::compose:page-layout/... format
//
// This function is auto-generated
func (r PageLayout) RbacResource() string {
	return PageLayoutRbacResource(r.NamespaceID, r.PageID, r.ID)
}

// PageLayoutRbacResource returns string representation of RBAC resource for PageLayout
//
// RBAC resource is in the corteza::compose:page-layout/... format
//
// This function is auto-generated
func PageLayoutRbacResource(namespaceID uint64, pageID uint64, id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, pageID, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{PageLayoutResourceType}
	if namespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(namespaceID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if pageID != 0 {
		cpts = append(cpts, strconv.FormatUint(pageID, 10))
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

	out := fmt.Sprintf(PageLayoutRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, namespaceID, pageID, id)

	return out

}

func PageLayoutRbacResourceTpl() string {
	return "%s/%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Record by calling RecordRbacResource fn
//
// RBAC resource is in the corteza::compose:record/... format
//
// This function is auto-generated
func (r Record) RbacResource() string {
	return RecordRbacResource(r.NamespaceID, r.ModuleID, r.ID)
}

// RecordRbacResource returns string representation of RBAC resource for Record
//
// RBAC resource is in the corteza::compose:record/... format
//
// This function is auto-generated
func RecordRbacResource(namespaceID uint64, moduleID uint64, id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, moduleID, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{RecordResourceType}
	if namespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(namespaceID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if moduleID != 0 {
		cpts = append(cpts, strconv.FormatUint(moduleID, 10))
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

	out := fmt.Sprintf(RecordRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, namespaceID, moduleID, id)

	return out

}

func RecordRbacResourceTpl() string {
	return "%s/%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza::compose/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza::compose/ format
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
