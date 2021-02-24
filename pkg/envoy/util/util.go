package util

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/id"
)

var (
	// wrapper around NextID that will aid service testing
	NextID = func() uint64 {
		return id.Next()
	}

	// wrapper around time.Now() that will aid testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}
)
