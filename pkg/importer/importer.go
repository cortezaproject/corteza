package importer

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	Interface interface {
		CastSet(set interface{}) error
		Cast(handle string, set interface{}) error
	}

	PermissionImporter interface {
		CastSet(string, string, interface{}) error
		CastResourcesSet(string, interface{}) error
		UpdateResources(base string, handle string, ID uint64)
		UpdateRoles(handle string, ID uint64)
		Store(context.Context, permissions.ImportKeeper) error
	}
)

func NormalizeHandle(in string) string {
	return in
}

func IsValidHandle(in string) bool {
	return true
}
