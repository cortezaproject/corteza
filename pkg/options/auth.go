package options

import (
	"github.com/cortezaproject/corteza-server/pkg/rand"
)

func (o *AuthOpt) Defaults() {

	if o.Secret == "" {
		o.Secret = string(rand.Bytes(32))
	}
}
