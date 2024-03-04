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
	// Component struct serves as a virtual resource type for the system component
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

// RbacResource returns string representation of RBAC resource for Application by calling ApplicationRbacResource fn
//
// RBAC resource is in the corteza::system:application/... format
//
// This function is auto-generated
func (r Application) RbacResource() string {
	return ApplicationRbacResource(r.ID)
}

// ApplicationRbacResource returns string representation of RBAC resource for Application
//
// RBAC resource is in the corteza::system:application/... format
//
// This function is auto-generated
func ApplicationRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{ApplicationResourceType}
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

	out := fmt.Sprintf(ApplicationRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func ApplicationRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ApigwRoute by calling ApigwRouteRbacResource fn
//
// RBAC resource is in the corteza::system:apigw-route/... format
//
// This function is auto-generated
func (r ApigwRoute) RbacResource() string {
	return ApigwRouteRbacResource(r.ID)
}

// ApigwRouteRbacResource returns string representation of RBAC resource for ApigwRoute
//
// RBAC resource is in the corteza::system:apigw-route/... format
//
// This function is auto-generated
func ApigwRouteRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{ApigwRouteResourceType}
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

	out := fmt.Sprintf(ApigwRouteRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func ApigwRouteRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for AuthClient by calling AuthClientRbacResource fn
//
// RBAC resource is in the corteza::system:auth-client/... format
//
// This function is auto-generated
func (r AuthClient) RbacResource() string {
	return AuthClientRbacResource(r.ID)
}

// AuthClientRbacResource returns string representation of RBAC resource for AuthClient
//
// RBAC resource is in the corteza::system:auth-client/... format
//
// This function is auto-generated
func AuthClientRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{AuthClientResourceType}
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

	out := fmt.Sprintf(AuthClientRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func AuthClientRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for DataPrivacyRequest by calling DataPrivacyRequestRbacResource fn
//
// RBAC resource is in the corteza::system:data-privacy-request/... format
//
// This function is auto-generated
func (r DataPrivacyRequest) RbacResource() string {
	return DataPrivacyRequestRbacResource(r.ID)
}

// DataPrivacyRequestRbacResource returns string representation of RBAC resource for DataPrivacyRequest
//
// RBAC resource is in the corteza::system:data-privacy-request/... format
//
// This function is auto-generated
func DataPrivacyRequestRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{DataPrivacyRequestResourceType}
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

	out := fmt.Sprintf(DataPrivacyRequestRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func DataPrivacyRequestRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Queue by calling QueueRbacResource fn
//
// RBAC resource is in the corteza::system:queue/... format
//
// This function is auto-generated
func (r Queue) RbacResource() string {
	return QueueRbacResource(r.ID)
}

// QueueRbacResource returns string representation of RBAC resource for Queue
//
// RBAC resource is in the corteza::system:queue/... format
//
// This function is auto-generated
func QueueRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{QueueResourceType}
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

	out := fmt.Sprintf(QueueRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func QueueRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Report by calling ReportRbacResource fn
//
// RBAC resource is in the corteza::system:report/... format
//
// This function is auto-generated
func (r Report) RbacResource() string {
	return ReportRbacResource(r.ID)
}

// ReportRbacResource returns string representation of RBAC resource for Report
//
// RBAC resource is in the corteza::system:report/... format
//
// This function is auto-generated
func ReportRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{ReportResourceType}
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

	out := fmt.Sprintf(ReportRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func ReportRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Role by calling RoleRbacResource fn
//
// RBAC resource is in the corteza::system:role/... format
//
// This function is auto-generated
func (r Role) RbacResource() string {
	return RoleRbacResource(r.ID)
}

// RoleRbacResource returns string representation of RBAC resource for Role
//
// RBAC resource is in the corteza::system:role/... format
//
// This function is auto-generated
func RoleRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{RoleResourceType}
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

	out := fmt.Sprintf(RoleRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func RoleRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Template by calling TemplateRbacResource fn
//
// RBAC resource is in the corteza::system:template/... format
//
// This function is auto-generated
func (r Template) RbacResource() string {
	return TemplateRbacResource(r.ID)
}

// TemplateRbacResource returns string representation of RBAC resource for Template
//
// RBAC resource is in the corteza::system:template/... format
//
// This function is auto-generated
func TemplateRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{TemplateResourceType}
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

	out := fmt.Sprintf(TemplateRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func TemplateRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for User by calling UserRbacResource fn
//
// RBAC resource is in the corteza::system:user/... format
//
// This function is auto-generated
func (r User) RbacResource() string {
	return UserRbacResource(r.ID)
}

// UserRbacResource returns string representation of RBAC resource for User
//
// RBAC resource is in the corteza::system:user/... format
//
// This function is auto-generated
func UserRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{UserResourceType}
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

	out := fmt.Sprintf(UserRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func UserRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for DalConnection by calling DalConnectionRbacResource fn
//
// RBAC resource is in the corteza::system:dal-connection/... format
//
// This function is auto-generated
func (r DalConnection) RbacResource() string {
	return DalConnectionRbacResource(r.ID)
}

// DalConnectionRbacResource returns string representation of RBAC resource for DalConnection
//
// RBAC resource is in the corteza::system:dal-connection/... format
//
// This function is auto-generated
func DalConnectionRbacResource(id uint64) string {
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, id)
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{}{DalConnectionResourceType}
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

	out := fmt.Sprintf(DalConnectionRbacResourceTpl(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, id)

	return out

}

func DalConnectionRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza::system/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza::system/ format
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
