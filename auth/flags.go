package auth

import (
	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/config"
)

type (
	appFlags struct {
		http *config.HTTP
		db   *config.Database
		jwt  *config.JWT
		oidc *config.OIDC
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
	if err := c.db.Validate(); err != nil {
		return err
	}
	if err := c.jwt.Validate(); err != nil {
		return err
	}
	if err := c.oidc.Validate(); err != nil {
		return err
	}
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
		new(config.HTTP).Init(prefix...),
		new(config.Database).Init(prefix...),
		new(config.JWT).Init(prefix...),
		new(config.OIDC).Init(prefix...),
	}
}
