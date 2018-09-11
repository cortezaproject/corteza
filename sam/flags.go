package sam

import (
	"github.com/pkg/errors"

	"github.com/crusttech/crust/config"
	"github.com/crusttech/crust/sam/repository"
)

type (
	appFlags struct {
		http       *config.HTTP
		db         *config.Database
		repository *repository.Flags
	}
)

var flags *appFlags

func (c *appFlags) Validate() error {
	if c == nil {
		return errors.New("SAM flags are not initialized, need to call Flags()")
	}
	if err := c.http.Validate(); err != nil {
		return err
	}
	if err := c.db.Validate(); err != nil {
		return err
	}
	if err := c.repository.Validate(); err != nil {
		return err
	}
	return nil
}

func Flags(prefix ...string) {
	if flags != nil {
		return
	}
	if len(prefix) == 0 {
		panic("sam.Flags() needs prefix on first call")
	}

	flags = &appFlags{
		new(config.HTTP).Init(prefix...),
		new(config.Database).Init(prefix...),
		new(repository.Flags).Init(prefix...),
	}
}
