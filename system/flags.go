package service

import (
	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/config"
)

type (
	appFlags struct {
		smtp    *config.SMTP
		http    *config.HTTP
		monitor *config.Monitor
		db      *config.Database
		oidc    *config.OIDC
		social  *config.Social
	}
)

var flags *appFlags

func (c *appFlags) Validate() error {
	if c == nil {
		return errors.New("AUTH flags are not initialized, need to call Flags() or FullFlags()")
	}
	if err := c.http.Validate(); err != nil {
		return err
	}
	if err := c.smtp.Validate(); err != nil {
		return err
	}
	if err := c.monitor.Validate(); err != nil {
		return err
	}
	if err := c.db.Validate(); err != nil {
		return err
	}
	if err := c.oidc.Validate(); err != nil {
		return err
	}
	//if err := c.jwt.Validate(); err != nil {
	//	return err
	//}
	return nil
}

func Flags(prefix ...string) {
	if flags != nil {
		return
	}
	if len(prefix) == 0 {
		panic("auth.Flags() needs prefix on first call")
	}
	flags = &appFlags{
		new(config.SMTP).Init(prefix...),
		new(config.HTTP).Init(prefix...),
		new(config.Monitor).Init(prefix...),
		new(config.Database).Init(prefix...),
		new(config.OIDC).Init(prefix...),
		new(config.Social).Init(prefix...),
	}
}
