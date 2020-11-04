package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

const (
	COMPOSE_RECORD_SET_RESOURCE_TYPE = "ComposeRecordSet"
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
	r.SetResourceType(COMPOSE_RECORD_SET_RESOURCE_TYPE)

	return r
}
