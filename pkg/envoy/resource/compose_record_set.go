package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeRecordSet struct {
		*base

		Walk Walker
	}

	Walker func(r *types.Record) error
)

// @todo add record provider
func NewComposeRecordSet() *ComposeRecordSet {
	r := &ComposeRecordSet{base: &base{}}
	r.SetResourceType("compose:record")

	return r
}
