package rbac

import (
	"github.com/crusttech/crust/internal/config"
)

var flags *config.RBAC

func Flags(prefix ...string) {
	flags = new(config.RBAC).Init(prefix...)
}
