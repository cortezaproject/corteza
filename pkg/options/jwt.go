package options

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/rand"
)

type (
	AuthOpt struct {
		Secret string        `env:"AUTH_JWT_SECRET"`
		Expiry time.Duration `env:"AUTH_JWT_EXPIRY"`
	}
)

func Auth() (o *AuthOpt) {
	o = &AuthOpt{
		Expiry: time.Hour * 24 * 30,
	}

	fill(o)

	// Setting JWT secret to random string to prevent security accidents...
	//
	// @todo check if this is a monolith system
	//       on microservice setup we can not afford to autogenerate secret:
	//       each subsystem will get it's own
	if o.Secret == "" {
		o.Secret = string(rand.Bytes(32))
	}

	return
}
