package options

import (
	"crypto/md5"
	"fmt"
)

func (o *AuthOpt) Defaults() {
	if o.Secret == "" {
		// if JWT secret is empty generate it from virtualhost/hostname and DB_DSN value.
		// this will keep the secret the same through restarts
		o.Secret = EnvString("DB_DSN", "memory")
		// pick one of the env that holds hostname
		o.Secret += EnvString("HOSTNAME", "localhost")

		o.Secret = fmt.Sprintf("%x", md5.Sum([]byte(o.Secret)))
	}
}
