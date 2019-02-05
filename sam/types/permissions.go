package types

import (
	"github.com/crusttech/crust/system/types"
)

type (
	ResourceProvider interface {
		Resource() types.Resource
	}
)

// These entities create resources in RBAC
var _ ResourceProvider = &Organisation{}
var _ ResourceProvider = &Team{}
var _ ResourceProvider = &Channel{}
