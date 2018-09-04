package auth

import (
	"github.com/pkg/errors"
)

type (
	configuration struct {
		http *httpFlags
		db   *dbFlags
		jwt  *jwtFlags
	}
)

var config *configuration

func (configuration) New() *configuration {
	return &configuration{
		new(httpFlags),
		new(dbFlags),
		new(jwtFlags),
	}
}

func (c *configuration) validate() error {
	if c == nil {
		return errors.New("CRM config is not initialized, need to call Flags()")
	}
	if err := c.http.validate(); err != nil {
		return err
	}
	if err := c.db.validate(); err != nil {
		return err
	}
	if err := c.jwt.validate(); err != nil {
		return err
	}
	return nil
}

func Flags(prefix ...string) {
	if config != nil {
		return
	}
	if len(prefix) == 0 {
		panic("crm.Flags() needs prefix on first call")
	}
	config = configuration{}.New()
	config.http.flags(prefix...)
	config.db.flags(prefix...)
	config.jwt.flags(prefix...)
}
