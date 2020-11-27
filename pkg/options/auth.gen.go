package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/auth.yaml

import (
	"time"
)

type (
	AuthOpt struct {
		Secret string        `env:"AUTH_JWT_SECRET"`
		Expiry time.Duration `env:"AUTH_JWT_EXPIRY"`
	}
)

// Auth initializes and returns a AuthOpt with default values
func Auth() (o *AuthOpt) {
	o = &AuthOpt{
		Expiry: time.Hour * 24 * 30,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Auth) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
