package auth

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	jwtFlags struct {
		secret     string
		expiry     int64
		debugToken bool
	}
)

func (c *jwtFlags) validate() error {
	if c.secret == "" {
		return errors.New("JWT Secret not set for AUTH")
	}
	return nil
}

func (c *jwtFlags) flags(prefix ...string) {
	flag.StringVar(&c.secret, "auth-jwt-secret", "", "JWT Secret")
	flag.Int64Var(&c.expiry, "auth-jwt-expiry", 3600, "JWT Expiration in minutes")
	flag.BoolVar(&c.debugToken, "auth-jwt-debug", false, "Generate debug JWT key")
}
