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

func (c *configuration) validate() error {
	if c == nil {
		return errors.New("CRM config is not initialized, need to call Flags() or FullFlags()")
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
	config = &configuration{
		jwt: new(jwtFlags).flags(prefix...),
	}
}

func FullFlags(prefix ...string) {
	if config != nil {
		return
	}
	if len(prefix) == 0 {
		panic("crm.Flags() needs prefix on first call")
	}
	config = &configuration{
		new(httpFlags).flags(prefix...),
		new(dbFlags).flags(prefix...),
		new(jwtFlags).flags(prefix...),
	}
}
