package service

import (
	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/config"
)

type (
	appFlags struct {
		http    *config.HTTP
		monitor *config.Monitor
		db      *config.Database
	}
)

var flags *appFlags

func (c *appFlags) Validate() error {
	if c == nil {
		return errors.New("CRM flags are not initialized, need to call Flags()")
	}
	if err := c.http.Validate(); err != nil {
		return err
	}
	if err := c.monitor.Validate(); err != nil {
		return err
	}
	if err := c.db.Validate(); err != nil {
		return err
	}
	return nil
}

func Flags(prefix ...string) {
	if flags != nil {
		return
	}
	if len(prefix) == 0 {
		panic("crm.Flags() needs prefix on first call")
	}

	flags = &appFlags{
		new(config.HTTP).Init(prefix...),
		new(config.Monitor).Init(prefix...),
		new(config.Database).Init(prefix...),
	}
}
