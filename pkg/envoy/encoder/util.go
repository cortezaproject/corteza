package encoder

import "github.com/cortezaproject/corteza-server/pkg/id"

var (
	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)
