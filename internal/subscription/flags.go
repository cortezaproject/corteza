package subscription

import (
	"github.com/crusttech/crust/internal/config"
)

type (
	localFlags struct {
		subscription *config.Subscription
	}
)

var flags *localFlags

// Flags matches signature for main()
func Flags(prefix ...string) {
	new(localFlags).Init(prefix...)
}

func (f *localFlags) Validate() error {
	if err := f.subscription.Validate(); err != nil {
		return err
	}
	return nil
}

func (f *localFlags) Init(prefix ...string) *localFlags {
	if flags != nil {
		return flags
	}
	flags = &localFlags{
		new(config.Subscription).Init(prefix...),
	}
	return flags
}
