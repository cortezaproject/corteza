package messagebus

import (
	"strconv"
)

const (
	// QueueResourceType
	// Using system as component (even if it is currently parked under pkg/messagebus)
	// and queue as name (even if is currently named queue-settings
	QueueResourceType = "corteza::system:queue"
)

// RbacResource
//
// Not following the naming pattern here (queue settings is queue)
func (r QueueSettings) RbacResource() string {
	return QueueRbacResource(r.ID)
}

// QueueRbacResource returns string representation of RBAC resource for Queue
func QueueRbacResource(id uint64) string {
	out := QueueResourceType
	if id != 0 {
		out += "/" + strconv.FormatUint(id, 10)
	} else {
		out += "/*"
	}

	return out

}
