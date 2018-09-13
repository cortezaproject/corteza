package repository

import (
	"github.com/crusttech/crust/internal/config"
)

type (
	Flags struct {
		PubSub    *config.PubSub
		Websocket *config.Websocket
	}
)

var flags *Flags

func (f *Flags) Validate() error {
	if flags == nil {
		return ErrConfigError.New()
	}
	if err := f.PubSub.Validate(); err != nil {
		return err
	}
	if err := f.Websocket.Validate(); err != nil {
		return err
	}
	return nil
}

func (f *Flags) Init(prefix ...string) *Flags {
	if flags != nil {
		return flags
	}
	flags = &Flags{
		new(config.PubSub).Init(prefix...),
		new(config.Websocket).Init(prefix...),
	}
	return flags
}
