package store

import (
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	setting struct {
		cfg *EncoderConfig

		res *resource.Setting
		st  *types.SettingValue

		ux *userIndex
	}
)
