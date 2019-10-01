package types

import (
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

const MessagingPermissionResource = permissions.Resource("messaging")
const ChannelPermissionResource = permissions.Resource("messaging:channel:")
const WebhookPermissionResource = permissions.Resource("messaging:webhook:")
