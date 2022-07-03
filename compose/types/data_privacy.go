package types

import (
	"github.com/cortezaproject/corteza-server/pkg/filter"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	PrivacyModule struct {
		ID     uint64 `json:"moduleID,string"`
		Name   string `json:"name"`
		Handle string `json:"handle"`
		Owner  bool   `json:"owner"`

		ConnectionID uint64 `json:"-"`
		Connection   *sysTypes.DalConnection
	}

	PrivacyModuleFilter struct {
		NamespaceID  uint64   `json:"namespaceID,string"`
		ConnectionID []uint64 `json:"connectionID,string"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)
