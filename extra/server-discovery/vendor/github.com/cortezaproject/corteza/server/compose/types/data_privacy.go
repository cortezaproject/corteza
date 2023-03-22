package types

import (
	"github.com/cortezaproject/corteza/server/pkg/filter"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
)

type (
	PrivacyModule struct {
		Module    PrivacyModuleMeta    `json:"module"`
		Namespace PrivacyNamespaceMeta `json:"namespace"`

		ConnectionID uint64                  `json:"-"`
		Connection   *sysTypes.DalConnection `json:"connection"`
	}

	PrivacyModuleMeta struct {
		ID     uint64         `json:"moduleID,string"`
		Name   string         `json:"name"`
		Handle string         `json:"handle"`
		Fields ModuleFieldSet `json:"fields"`
	}

	PrivacyNamespaceMeta struct {
		ID   uint64 `json:"namespaceID,string"`
		Slug string `json:"slug"`
		Name string `json:"name"`
	}

	PrivacyModuleFilter struct {
		NamespaceID  uint64   `json:"-"`
		ConnectionID []uint64 `json:"connectionID,string"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)
