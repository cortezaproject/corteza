package config

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	JWT struct {
		Secret     string
		Expiry     int64
		DebugToken bool
	}
)

func (c *JWT) Validate() error {
	if c == nil {
		return nil
	}
	if c.Secret == "" {
		return errors.New("JWT Secret not set for AUTH")
	}
	return nil
}

func (c *JWT) Init(prefix ...string) *JWT {
	flag.StringVar(&c.Secret, "auth-jwt-secret", "", "JWT Secret")
	flag.Int64Var(&c.Expiry, "auth-jwt-expiry", 3600, "JWT Expiration in minutes")
	flag.BoolVar(&c.DebugToken, "auth-jwt-debug", false, "Generate debug JWT key")
	return c
}
