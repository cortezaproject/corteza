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
		// Setting JWT secret to random string to prevent security accidents...
		Secret: string(rand.Bytes(32)),

		Expiry: time.Hour * 24 * 30,
	}

	fill(o, pfix)

	return
}
