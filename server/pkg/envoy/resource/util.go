package resource

import (
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
)

func firstOkString(ss ...string) string {
	for _, s := range ss {
		if s != "" && s != "0" {
			return s
		}
	}
	return ""
}

// Taken (and modified) from compose/service/values/sanitizer.go
func toTime(v string) *time.Time {
	ff := []string{
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC850,
		time.RFC822Z,
		time.RFC822,
		time.RubyDate,
		time.UnixDate,
		time.ANSIC,
		"2006/_1/_2 15:04:05",
		"2006/_1/_2 15:04",
	}

	for _, f := range ff {
		parsed, err := time.Parse(f, v)
		if err == nil {
			return &parsed
		}
	}

	return nil
}

// Ref builders

func MakeNamespaceRef(id uint64, handle, name string) *Ref {
	return makeGenericRef(types.NamespaceResourceType, id, handle, name)
}

func MakeModuleRef(id uint64, handle, name string) *Ref {
	return makeGenericRef(types.ModuleResourceType, id, handle, name)
}

func MakePageRef(id uint64, handle, title string) *Ref {
	return makeGenericRef(types.PageResourceType, id, handle, title)
}

func MakeRoleRef(id uint64, handle, name string) *Ref {
	return makeGenericRef(systemTypes.RoleResourceType, id, handle, name)
}

func makeGenericRef(t string, id uint64, ii ...string) *Ref {
	args := make([]interface{}, len(ii)+1)
	args[0] = id
	for i := range ii {
		args[i+1] = ii[i]
	}

	aux := &Ref{
		ResourceType: t,
		Identifiers:  identifiers(args...),
	}

	if t == "" || len(aux.Identifiers) == 0 {
		return nil
	}

	return aux
}
