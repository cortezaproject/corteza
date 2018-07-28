package auth

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	configuration struct {
		jwtSecret string
		jwtExpiry int64
	}
)

var config configuration

func (c configuration) validate() error {
	if c.jwtSecret == "" {
		return errors.New("JWT Secret not set for AUTH")
	}

	return nil
}

// Flags should be called from main to register flags
func Flags() {
	flag.StringVar(&config.jwtSecret, "auth-jwt-secret", "", "JWT Secret")
	flag.Int64Var(&config.jwtExpiry, "auth-jwt-expiry", 3600, "JWT Expiration in minutes")
}
