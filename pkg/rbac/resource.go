package rbac

import (
	"strconv"
	"strings"
)

type (
	Resource string
)

const (
	resourceDelimiter = ':'
	resourceWildcard  = '*'
)

func (r Resource) append(suffix string) Resource {
	if !r.IsAppendable() {
		panic("cannot append to non appendable resource '" + r.String() + "'")
	}

	return Resource(r.String() + suffix)
}

// Resource to satisfty interfaces and ease development
func (r Resource) RBACResource() Resource {
	return r
}

// DynamicRoles satisfies Resourcable interface when Resource is
// used directly
func (r Resource) DynamicRoles(i uint64) []uint64 {
	return nil
}

func (r Resource) AppendID(ID uint64) Resource {
	return r.append(strconv.FormatUint(ID, 10))
}

func (r Resource) AppendWildcard() Resource {
	return r.TrimID().append(string(resourceWildcard))
}

// Trims off wildcard/id from resource
func (r Resource) TrimID() Resource {
	s := r.String()
	p := strings.LastIndexByte(s, resourceDelimiter)
	if p > 0 {
		return Resource(s[0 : p+1])
	}

	return r
}

// GetID returns the identifier for this resource
func (r Resource) GetID() (id uint64, err error) {
	s := r.String()
	p := strings.LastIndexByte(s, resourceDelimiter)
	if p > 0 {
		id, err = strconv.ParseUint(s[p+1:], 10, 64)
		return id, err
	}

	return 0, nil
}

// IsAppendable checks if Resource has trailing resource delimiter
func (r Resource) IsAppendable() bool {
	return strings.IndexByte(r.String(), resourceDelimiter) > -1
}

// IsValid does basic resource validation
func (r Resource) IsValid() bool {
	return len(r) > 0 && r[len(r)-1] != resourceDelimiter
}

// IsServiceLevel checks for resource delimiters - service level resources do not have it
func (r Resource) GetService() Resource {
	s := r.String()
	p := strings.IndexByte(s, resourceDelimiter)
	if p > 0 {
		return Resource(s[0:p])
	}

	return r
}

// HasWildcard checks if resource has wildcard char at the end
func (r Resource) HasWildcard() bool {
	return len(r) > 0 && r[len(r)-1] == resourceWildcard
}

func (r Resource) String() string {
	return string(r)
}
