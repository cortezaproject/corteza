package options

import (
	"time"

	"github.com/cortezaproject/corteza-server/internal/rand"
)

type (
	JWTOpt struct {
		Secret string
		Expiry time.Duration
	}
)

func JWT(pfix string) (o *JWTOpt) {
	o = &JWTOpt{
		Secret: EnvString(pfix, "AUTH_JWT_SECRET", string(rand.Bytes(32))),
		// Setting JWT secret to random string to prevent security accidents...
		Expiry: EnvDuration(pfix, "AUTH_JWT_EXPIRY", time.Hour*24*30),
	}

	return
}
