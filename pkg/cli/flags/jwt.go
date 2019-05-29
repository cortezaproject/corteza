package flags

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/internal/rand"
)

type (
	JWTOpt struct {
		Secret string
		Expiry time.Duration
	}
)

func JWT(cmd *cobra.Command) (o *JWTOpt) {
	o = &JWTOpt{}

	// Setting JWT secret to random string to prevent security accidents...
	BindString(cmd, &o.Secret,
		"auth-jwt-secret", string(rand.Bytes(32)),
		"JWT Secret")

	BindDuration(cmd, &o.Expiry,
		"auth-jwt-expiry", time.Hour*24*30,
		"JWT Expiration")

	return
}
