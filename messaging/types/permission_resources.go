package types

import (
	"github.com/crusttech/crust/internal/rules"
)

const PermissionResource = rules.Resource("messaging")
const ChannelPermissionResource = rules.Resource("messaging:channel:")
const WebhookPermissionResource = rules.Resource("messaging:webhook:")
