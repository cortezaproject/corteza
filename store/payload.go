package store

import (
	"github.com/cortezaproject/corteza-server/pkg/slice"
)

type (
	// Payload servers as a generic interface between incoming structs scheduled
	// for update or insertion and the db layer
	Payload map[string]interface{}
)

// Only returns new payload with subset fields (if empty slice is send, it returns entire payload)
func (p Payload) Only(ff ...string) Payload {
	if len(ff) == 0 {
		return p
	}

	o := Payload{}
	for _, k := range ff {
		if v, has := p[k]; has {
			o[k] = v
		}
	} //

	return o
}

// Skip returns new payload with subset fields (if empty slice is send, it returns entire payload)
func (p Payload) Skip(ff ...string) Payload {
	if len(ff) == 0 {
		return p
	}

	m := slice.ToStringBoolMap(ff)

	o := Payload{}
	for k, v := range p {
		if !m[k] {
			o[k] = v
		}
	}

	return o
}

// Alias adds alias (<alias>.<key>) to all keys
func (p Payload) Alias(alias string) Payload {
	o := Payload{}
	for k, v := range p {
		o[alias+"."+k] = v
	}

	return o
}
