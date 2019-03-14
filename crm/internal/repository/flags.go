package repository

import (
	_ "github.com/crusttech/crust/internal/config"
)

type (
	Flags struct {
		// No config yet
	}
)

var flags *Flags

func (f *Flags) Validate() error {
	return nil
}

func (f *Flags) Init(prefix ...string) *Flags {
	if flags != nil {
		return flags
	}
	flags = &Flags{}
	return flags
}
