package store

import (
	"github.com/crusttech/human/server/pkg/envoy/resource"
	"github.com/crusttech/human/server/system/types"
)

type (
	setting struct {
		cfg *EncoderConfig

		res *resource.Setting
		st  *types.SettingValue

		ux *userIndex
	}
)
