package auth

import (
	"github.com/crusttech/crust/internal/config"
)

type (
	localFlags struct {
		jwt *config.JWT
	}
)

var flags *localFlags

// Flags matches signature for main()
func Flags(prefix ...string) {
	new(localFlags).Init(prefix...)
}

func (f *localFlags) Validate() error {
	if flags == nil {
		return ErrConfigError.New()
	}
	if err := f.jwt.Validate(); err != nil {
		return err
	}
	return nil
}

func (f *localFlags) Init(prefix ...string) *localFlags {
	if flags != nil {
		return flags
	}
	flags = &localFlags{
		new(config.JWT).Init(prefix...),
	}
	return flags
}
