package options

import (
	"time"

	"github.com/cortezaproject/corteza-server/internal/rand"
)

type (
	JWTOpt struct {
		Secret string        `env:"AUTH_JWT_SECRET"`
		Expiry time.Duration `env:"AUTH_JWT_EXPIRY"`
	}
)

func JWT(pfix string) (o *JWTOpt) {
	o = &JWTOpt{
		Expiry: time.Hour * 24 * 30,
	}

	fill(o, pfix)

	// Setting JWT secret to random string to prevent security accidents...
	if o.Secret == "" {
		o.Secret = string(rand.Bytes(32))
	}

	return
}
