package types

import (
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

const MessagingRBACResource = rbac.Resource("messaging")
const ChannelRBACResource = rbac.Resource("messaging:channel:")
